package main

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger/python/internal/dagger"
)

// PyPIConfig holds PyPI deployment configuration
type PyPIConfig struct {
	// Registry URL (default: https://upload.pypi.org/legacy/)
	Registry string
	// Token for authentication
	Token *dagger.Secret
	// Skip existing versions (default: false)
	SkipExisting bool
	// Allow dirty versions (default: false)
	AllowDirty bool
}

// Python represents a Python module with Poetry support
type Python struct {
	// PythonVersion specifies the Python version to use (default: "3.12")
	PythonVersion string
	// PackagePath specifies the path to the package within the source (default: ".")
	PackagePath string
	// PyPIConfig holds the PyPI deployment configuration
	PyPIConfig *PyPIConfig
}

// WithPythonVersion sets the Python version to use
func (m *Python) WithPythonVersion(version string) *Python {
	m.PythonVersion = version
	return m
}

// WithPackagePath sets the package path within the source
func (m *Python) WithPackagePath(path string) *Python {
	m.PackagePath = path
	return m
}

// WithPyPIConfig sets the PyPI deployment configuration
func (m *Python) WithPyPIConfig(config *PyPIConfig) *Python {
	m.PyPIConfig = config
	return m
}

// getBaseImage returns the Python base image with the configured version
func (m *Python) getBaseImage() string {
	version := m.PythonVersion
	if version == "" {
		version = "3.12"
	}
	return fmt.Sprintf("python:%s-slim", version)
}

// getWorkdir returns the working directory path
func (m *Python) getWorkdir(basePath string) string {
	if m.PackagePath == "" {
		return basePath
	}
	return filepath.Join(basePath, m.PackagePath)
}

// getPyPIRegistry returns the configured PyPI registry URL with a default
func (m *Python) getPyPIRegistry() string {
	if m.PyPIConfig == nil || m.PyPIConfig.Registry == "" {
		return "https://upload.pypi.org/legacy/"
	}
	return m.PyPIConfig.Registry
}

// Publish builds, tests and publishes the Python package to a registry
func (m *Python) Publish(ctx context.Context, source *dagger.Directory) (string, error) {
	// Run tests before publishing
	if _, err := m.Test(ctx, source); err != nil {
		return "", fmt.Errorf("tests failed: %w", err)
	}

	// Build the package
	container := m.Build(source)

	// Configure Poetry for publishing
	container = container.WithExec([]string{
		"poetry", "config",
		"repositories.pypi.url", m.getPyPIRegistry(),
	})

	// Add authentication if token is provided
	if m.PyPIConfig != nil && m.PyPIConfig.Token != nil {
		container = container.
			WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", m.PyPIConfig.Token)
	}

	// Prepare publish command
	publishCmd := []string{"poetry", "publish", "--build", "--no-interaction"}
	
	// Add optional flags based on configuration
	if m.PyPIConfig != nil {
		if m.PyPIConfig.SkipExisting {
			publishCmd = append(publishCmd, "--skip-existing")
		}
		if m.PyPIConfig.AllowDirty {
			publishCmd = append(publishCmd, "--allow-dirty")
		}
	}

	// Execute publish command
	return container.WithExec(publishCmd).Stdout(ctx)
}

// Build creates a Python package using Poetry
func (m *Python) Build(source *dagger.Directory) *dagger.Container {
	container := m.BuildEnv(source).
		WithExec([]string{
			"poetry", "build",
			"--no-interaction",
		})
	
	return container.WithDirectory("/dist", container.Directory("/app/dist"))
}

// Test runs the test suite using pytest with coverage reporting
func (m *Python) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BuildEnv(source).
		WithExec([]string{
			"poetry", "run", "pytest",
			"--verbose",
			"--color=yes",
			fmt.Sprintf("--cov=%s", m.PackagePath),
			"--cov-report=xml",
			"--cov-report=term",
			"--cov-report=html:coverage_html",
			"--no-cov-on-fail",
		}).
		Stdout(ctx)
}

// BuildEnv prepares a Python development environment with Poetry
func (m *Python) BuildEnv(source *dagger.Directory) *dagger.Container {
	poetryCache := dag.CacheVolume("poetry-cache")
	pipCache := dag.CacheVolume("pip-cache")
	
	workdir := m.getWorkdir("/app")
	
	return dag.Container().
		From(m.getBaseImage()).
		WithDirectory("/app", source).
		WithMountedCache("/root/.cache/pypoetry", poetryCache).
		WithMountedCache("/root/.cache/pip", pipCache).
		WithWorkdir(workdir).
		WithExec([]string{
			"pip", "install",
			"--no-cache-dir",
			"--upgrade",
			"pip",
			"poetry",
		}).
		WithExec([]string{
			"poetry", "config",
			"virtualenvs.in-project", "true",
		}).
		WithExec([]string{
			"poetry", "install",
			"--no-interaction",
			"--no-root",
			"--with", "dev",
		})
}
