// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/felipepimentel/daggerverse/python/internal/dagger"
)

// Common error messages
const (
	errBuildContainer = "failed to build container"
	errBuildTestEnv   = "failed to build test environment"
	errGetVersion     = "error getting version"
)

// Default configurations
const (
	// DefaultPythonVersion is the default Python version to use
	DefaultPythonVersion = "3.12-slim"
	// DefaultGitEmail is the default Git email for commits
	DefaultGitEmail = "github-actions[bot]@users.noreply.github.com"
	// DefaultGitName is the default Git username for commits
	DefaultGitName = "github-actions[bot]"
)

// Container configuration
const (
	containerWorkdir = "/src"
	registryURLFmt   = "ttl.sh/python-pipeline-%s"
)

// PythonPipeline orchestrates Python project workflows using Poetry and PyPI.
// It provides a complete CI/CD pipeline for Python projects, including testing,
// building, and publishing to PyPI.
type PythonPipeline struct {
	// pythonVersion specifies the Python version to use
	pythonVersion string
	// gitEmail is used for Git configuration
	gitEmail string
	// gitName is used for Git configuration
	gitName string
}

// Option configures a PythonPipeline
type Option func(*PythonPipeline)

// WithPythonVersion sets the Python version
func WithPythonVersion(version string) Option {
	return func(p *PythonPipeline) {
		p.pythonVersion = version
	}
}

// WithGitConfig sets the Git configuration
func WithGitConfig(email, name string) Option {
	return func(p *PythonPipeline) {
		p.gitEmail = email
		p.gitName = name
	}
}

// New creates a new instance of PythonPipeline with the provided options.
// If no options are provided, default values will be used.
func New() *PythonPipeline {
	p := &PythonPipeline{
		pythonVersion: DefaultPythonVersion,
		gitEmail:      DefaultGitEmail,
		gitName:       DefaultGitName,
	}
	return p
}

// Publish builds, tests, and publishes the Python package to PyPI and the container image.
// It returns the address of the published container or an error if any step fails.
func (m *PythonPipeline) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
	fmt.Println("Starting publish process...")

	// Run tests first
	if _, err := m.Test(ctx, source); err != nil {
		return "", fmt.Errorf("tests failed: %w", err)
	}

	fmt.Println("Getting version from versioner module...")
	version, err := dag.Versioner().BumpVersion(ctx, source, true)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errGetVersion, err)
	}

	fmt.Println("Using version:", version)

	// Build container with all dependencies
	baseContainer := m.Build(ctx, source)
	if baseContainer == nil {
		return "", errors.New(errBuildContainer)
	}

	// Configure Git in container
	_ = baseContainer.
		WithEnvVariable("GIT_AUTHOR_EMAIL", m.gitEmail).
		WithEnvVariable("GIT_AUTHOR_NAME", m.gitName).
		WithEnvVariable("GIT_COMMITTER_EMAIL", m.gitEmail).
		WithEnvVariable("GIT_COMMITTER_NAME", m.gitName)

	// Use Poetry module for dependency management
	poetryContainer := dag.Container().From(fmt.Sprintf("python:%s", m.pythonVersion))
	poetryContainer = poetryContainer.
		WithDirectory(containerWorkdir, source).
		WithWorkdir(containerWorkdir).
		WithExec([]string{"pip", "install", "poetry"}).
		WithExec([]string{"poetry", "install", "--with", "dev"})

	// Use PyPI module for publishing
	pypiContainer := poetryContainer.
		WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token).
		WithEnvVariable("VERSION", version).
		WithExec([]string{"poetry", "build"}).
		WithExec([]string{"poetry", "publish", "--username", "__token__", "--no-interaction"})

	// Run Ruff checks
	lintContainer := dag.Container().
		From("astral/ruff").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"ruff", "check", "."})

	if _, err := lintContainer.Stdout(ctx); err != nil {
		return "", fmt.Errorf("ruff check failed: %w", err)
	}

	fmt.Println("Package version", version, "published successfully to PyPI")

	// Publish container
	fmt.Println("Publishing container...")
	address, err := pypiContainer.Publish(ctx, fmt.Sprintf(registryURLFmt, version))
	if err != nil {
		return "", fmt.Errorf("failed to publish container: %w", err)
	}

	fmt.Println("Container published successfully to:", address)
	return address, nil
}

// Build creates a container with all dependencies installed and configured.
// It returns the configured container or nil if the build fails.
func (m *PythonPipeline) Build(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("python:%s", m.pythonVersion)).
		WithDirectory(containerWorkdir, source).
		WithWorkdir(containerWorkdir).
		WithExec([]string{"pip", "install", "poetry"}).
		WithExec([]string{"poetry", "install", "--with", "dev"})
}

// Test runs all quality checks and returns the combined test output.
// It returns an error if any check fails.
func (m *PythonPipeline) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	container := m.Build(ctx, source)
	if container == nil {
		return "", errors.New(errBuildTestEnv)
	}

	// Run tests with Poetry
	if _, err := container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx); err != nil {
		return "", fmt.Errorf("pytest failed: %w", err)
	}

	// Run Ruff checks
	lintContainer := dag.Container().
		From("astral/ruff").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"ruff", "check", "."})

	if _, err := lintContainer.Stdout(ctx); err != nil {
		return "", fmt.Errorf("ruff check failed: %w", err)
	}

	return "All tests and checks passed successfully!", nil
}

// BuildEnv creates a development environment with all dependencies installed.
// It returns the configured container.
func (m *PythonPipeline) BuildEnv(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return m.Build(ctx, source)
}

