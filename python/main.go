package main

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger/python/internal/dagger"
)

// Python represents a Python module with Poetry support
type Python struct {
	// PythonVersion specifies the Python version to use (default: "3.12")
	PythonVersion string
	// PackagePath specifies the path to the package within the source (default: ".")
	PackagePath string
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

// Publish builds, tests and publishes the Python package to a registry
func (m *Python) Publish(ctx context.Context, source *dagger.Directory, registry string) (string, error) {
	// Run tests before publishing
	if _, err := m.Test(ctx, source); err != nil {
		return "", fmt.Errorf("tests failed: %w", err)
	}

	// Build the package
	build := m.Build(source)

	// If no registry specified, use TestPyPI as default
	if registry == "" {
		registry = "https://test.pypi.org/legacy/"
	}

	// Publish to the specified registry
	return build.
		WithEnvVariable("POETRY_REPOSITORIES_PYPI_URL", registry).
		WithExec([]string{
			"poetry", "publish",
			"--build",
			"--no-interaction",
			"--skip-existing",
		}).
		Stdout(ctx)
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
