// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
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

// Python orchestrates Python project workflows using Poetry and PyPI.
// It provides a complete CI/CD pipeline for Python projects, including testing,
// building, and publishing to PyPI.
type Python struct {
	// pythonVersion specifies the Python version to use
	pythonVersion string
	// gitEmail is used for Git configuration
	gitEmail string
	// gitName is used for Git configuration
	gitName string
}

// Option configures a Python
type Option func(*Python)

// WithPythonVersion sets the Python version
func WithPythonVersion(version string) Option {
	return func(p *Python) {
		p.pythonVersion = version
	}
}

// WithGitConfig sets the Git configuration
func WithGitConfig(email, name string) Option {
	return func(p *Python) {
		p.gitEmail = email
		p.gitName = name
	}
}

// New creates a new instance of Python with the provided options.
// If no options are provided, default values will be used.
func New() *Python {
	p := &Python{
		pythonVersion: DefaultPythonVersion,
		gitEmail:      DefaultGitEmail,
		gitName:       DefaultGitName,
	}
	return p
}

// Publish builds, tests, and publishes the Python package to PyPI and the container image.
// It returns the address of the published container or an error if any step fails.
func (m *Python) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
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

	// Build package using Poetry module
	buildDir := dag.Poetry().Build(source)

	// Use PyPI module for publishing
	address, err := dag.Container().
		From(fmt.Sprintf("python:%s", m.pythonVersion)).
		WithDirectory(containerWorkdir, buildDir).
		WithWorkdir(containerWorkdir).
		WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token).
		WithEnvVariable("VERSION", version).
		WithExec([]string{"poetry", "publish", "--username", "__token__", "--no-interaction"}).
		Publish(ctx, fmt.Sprintf(registryURLFmt, version))

	if err != nil {
		return "", fmt.Errorf("failed to publish container: %w", err)
	}

	fmt.Println("Container published successfully to:", address)
	return address, nil
}

// Build creates a container with all dependencies installed and configured.
// It returns the configured container or nil if the build fails.
func (m *Python) Build(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("python:%s", m.pythonVersion)).
		WithDirectory(containerWorkdir, dag.Poetry().Install(source)).
		WithWorkdir(containerWorkdir)
}

// Test runs all quality checks and returns the combined test output.
// It returns an error if any check fails.
func (m *Python) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	// Run tests using Poetry module
	testOutput, err := dag.Poetry().Test(ctx, source)
	if err != nil {
		return "", fmt.Errorf("poetry test failed: %w", err)
	}

	// Run Ruff checks
	if err := dag.Ruff().Lint(source).Assert(ctx); err != nil {
		return "", fmt.Errorf("ruff check failed: %w", err)
	}

	return fmt.Sprintf("Test output:\n%s\nAll tests and checks passed successfully!", testOutput), nil
}

// BuildEnv creates a development environment with all dependencies installed.
// It returns the configured container.
func (m *Python) BuildEnv(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return m.Build(ctx, source)
}

