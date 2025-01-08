// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/python/internal/dagger"
)

// Error messages for common failures.
const (
	errBuildContainer = "failed to build container"
	errBuildTestEnv   = "failed to build test environment"
	errGetVersion     = "error getting version"
	errPublish        = "failed to publish container"
	errPoetryTest     = "poetry test failed"
	errRuffCheck      = "ruff check failed"
	errPypiPublish    = "failed to publish to PyPI"
)

// Log messages for progress tracking.
const (
	logStartPublish    = "Starting publish process..."
	logStartTests      = "Running tests..."
	logStartLint       = "Running linting checks..."
	logStartBuild      = "Building package..."
	logStartPyPI       = "Publishing to PyPI..."
	logStartContainer  = "Publishing container..."
	logSuccessTests    = "All tests passed successfully!"
	logSuccessLint     = "All linting checks passed!"
	logSuccessPyPI     = "Package published successfully to PyPI"
	logSuccessVersion  = "Using version: %s"
	logSuccessPublish  = "Container published successfully to: %s"
)

// Python configuration defaults.
const (
	// DefaultPythonVersion is the default Python version to use.
	DefaultPythonVersion = "3.12-slim"
)

// Git configuration defaults.
const (
	// DefaultGitEmail is the default Git email for commits.
	DefaultGitEmail = "github-actions[bot]@users.noreply.github.com"
	// DefaultGitName is the default Git username for commits.
	DefaultGitName = "github-actions[bot]"
)

// Container configuration constants.
const (
	// containerWorkdir is the working directory inside the container.
	containerWorkdir = "/src"
	// registryURLFmt is the format string for the container registry URL.
	registryURLFmt = "ttl.sh/python-pipeline-%s"
)

// Python orchestrates Python project workflows using Poetry and PyPI.
// It provides a complete CI/CD pipeline for Python projects, including testing,
// building, and publishing to PyPI.
type Python struct {
	// pythonVersion specifies the Python version to use.
	pythonVersion string
	// gitEmail is used for Git configuration.
	gitEmail string
	// gitName is used for Git configuration.
	gitName string
}

// Option configures a Python instance.
type Option func(*Python)

// WithPythonVersion sets the Python version.
func WithPythonVersion(version string) Option {
	return func(p *Python) {
		p.pythonVersion = version
	}
}

// WithGitConfig sets the Git configuration.
func WithGitConfig(email, name string) Option {
	return func(p *Python) {
		p.gitEmail = email
		p.gitName = name
	}
}

// New creates a new instance of Python with the provided options.
// If no options are provided, default values will be used.
func New(opts ...Option) *Python {
	p := &Python{
		pythonVersion: DefaultPythonVersion,
		gitEmail:      DefaultGitEmail,
		gitName:       DefaultGitName,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Publish builds, tests, and publishes the Python package to PyPI and the container image.
// It returns the address of the published container or an error if any step fails.
func (p *Python) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
	fmt.Println(logStartPublish)

	// Run tests first
	if _, err := p.Test(ctx, source); err != nil {
		return "", fmt.Errorf("tests failed: %w", err)
	}

	// Get version from versioner module
	version, err := dag.Versioner().BumpVersion(ctx, source, true)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errGetVersion, err)
	}

	fmt.Printf(logSuccessVersion+"\n", version)

	fmt.Println(logStartBuild)
	// Build package using Poetry module
	buildDir := dag.Poetry().Build(source)

	fmt.Println(logStartPyPI)
	// Publish to PyPI using the pypi module
	if err := dag.Pypi().Publish(ctx, buildDir, token); err != nil {
		return "", fmt.Errorf("%s: %w", errPypiPublish, err)
	}
	fmt.Println(logSuccessPyPI)

	fmt.Println(logStartContainer)
	// Publish container
	address, err := dag.Container().
		From(fmt.Sprintf("python:%s", p.pythonVersion)).
		WithDirectory(containerWorkdir, buildDir).
		WithWorkdir(containerWorkdir).
		Publish(ctx, fmt.Sprintf(registryURLFmt, version))

	if err != nil {
		return "", fmt.Errorf("%s: %w", errPublish, err)
	}

	fmt.Printf(logSuccessPublish+"\n", address)
	return address, nil
}

// Build creates a container with all dependencies installed and configured.
// It returns the configured container or nil if the build fails.
func (p *Python) Build(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("python:%s", p.pythonVersion)).
		WithDirectory(containerWorkdir, dag.Poetry().Install(source)).
		WithWorkdir(containerWorkdir)
}

// Test runs all quality checks and returns the combined test output.
// It returns an error if any check fails.
func (p *Python) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	fmt.Println(logStartTests)

	// Run tests using Poetry module
	testOutput, err := dag.Poetry().Test(ctx, source)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errPoetryTest, err)
	}

	fmt.Println(logSuccessTests)

	// Run linting checks
	if err := p.Lint(ctx, source); err != nil {
		return "", err
	}

	return fmt.Sprintf("Test output:\n%s", testOutput), nil
}

// Lint runs code quality checks using Ruff.
// It returns an error if any check fails.
func (p *Python) Lint(ctx context.Context, source *dagger.Directory) error {
	fmt.Println(logStartLint)

	if err := dag.Ruff().Lint(source).Assert(ctx); err != nil {
		return fmt.Errorf("%s: %w", errRuffCheck, err)
	}

	fmt.Println(logSuccessLint)
	return nil
}

// BuildEnv creates a development environment with all dependencies installed.
// It returns the configured container.
func (p *Python) BuildEnv(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return p.Build(ctx, source)
}

