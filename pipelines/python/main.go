// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipepimentel/daggerverse/pipelines/python/internal/dagger"
)

// Error messages for common failures.
const (
	errBuildContainer = "failed to build container"
	errBuildTestEnv   = "failed to build test environment"
	errGetVersion     = "error getting version"
	errPublish        = "failed to publish container"
	errPoetryTest     = "poetry test failed"
	errRuffCheck      = "ruff check failed"
	errPypiPublish    = "failed to publish package to PyPI"
	errBuild          = "failed to build package"
)

// Log messages for progress tracking.
const (
	logStartPublish    = "Starting publish process..."
	logStartTests      = "Running tests..."
	logStartLint       = "Running linting checks..."
	logStartBuild      = "🏗️  Building package..."
	logStartPyPI       = "📦 Publishing to PyPI..."
	logStartContainer  = "Publishing container..."
	logSuccessTests    = "All tests passed successfully!"
	logSuccessLint     = "All linting checks passed!"
	logSuccessPyPI     = "✅ Package published to PyPI successfully!"
	logSuccessVersion  = "Using version: %s"
	logSuccessPublish  = "Container published successfully to: %s"
)

// Python configuration defaults.
const (
	// DefaultPythonVersion is the default Python version to use.
	DefaultPythonVersion = "3.12-alpine"
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
	// dockerUsername is used for Docker Hub authentication.
	dockerUsername string
	// dockerPassword is used for Docker Hub authentication.
	dockerPassword *dagger.Secret
	// skipTests indicates whether to skip running tests
	skipTests bool
	// skipLint indicates whether to skip running linting checks
	skipLint bool
	// githubToken is used for GitHub authentication
	// +private
	githubToken *dagger.Secret
}

// New creates a new instance of Python with the provided configuration.
func New(
	// Python version to use
	// +optional
	// +default="3.12-alpine"
	pythonVersion string,
	// Git email for commits
	// +optional
	// +default="github-actions[bot]@users.noreply.github.com"
	gitEmail string,
	// Git username for commits
	// +optional
	// +default="github-actions[bot]"
	gitName string,
	// Docker Hub username
	// +optional
	dockerUsername string,
	// Docker Hub password
	// +optional
	dockerPassword *dagger.Secret,
	// Skip running tests
	// +optional
	// +default=false
	skipTests bool,
	// Skip running linting checks
	// +optional
	// +default=false
	skipLint bool,
	// GitHub token for authentication
	// +optional
	githubToken *dagger.Secret,
) *Python {
	if pythonVersion == "" {
		pythonVersion = DefaultPythonVersion
	}
	if gitEmail == "" {
		gitEmail = DefaultGitEmail
	}
	if gitName == "" {
		gitName = DefaultGitName
	}

	return &Python{
		pythonVersion:   pythonVersion,
		gitEmail:        gitEmail,
		gitName:         gitName,
		dockerUsername:  dockerUsername,
		dockerPassword:  dockerPassword,
		skipTests:       skipTests,
		skipLint:        skipLint,
		githubToken:     githubToken,
	}
}

// Publish builds and publishes a Python package to PyPI
func (m *Python) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	// Create base container with git and poetry
	container := dag.Container().
		From("python:3.12-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry", "python-semantic-release", "tomli"})

	// Configure git for semantic-release
	container = container.
		WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"})

	// Mount source code
	container = container.WithMountedDirectory("/src", source).WithWorkdir("/src")

	// Check if repository is shallow
	isShallow, err := container.WithExec([]string{"git", "rev-parse", "--is-shallow-repository"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if repository is shallow: %w", err)
	}

	// If repository is shallow, unshallow it
	if strings.TrimSpace(isShallow) == "true" {
		container = container.WithExec([]string{"git", "fetch", "--unshallow", "--tags"})
	} else {
		container = container.WithExec([]string{"git", "fetch", "--tags"})
	}

	// Get the next version using semantic-release
	version, err := container.
		WithExec([]string{"semantic-release", "version", "--no-commit", "--no-tag"}).
		Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to determine version: %w", err)
	}
	version = strings.TrimSpace(version)

	// Update version in pyproject.toml
	container = container.WithExec([]string{"poetry", "version", version})

	// Build the package
	dist := dag.Poetry().BuildWithVersion(source, version)

	// Publish to PyPI
	err = dag.Pypi().
		Publish(ctx, dist, token)
	if err != nil {
		return fmt.Errorf("failed to publish package: %w", err)
	}

	// Create and push tag if GitHub token is provided
	if m.githubToken != nil {
		container = container.WithSecretVariable("GITHUB_TOKEN", m.githubToken)
		_, err = container.WithExec([]string{"semantic-release", "publish"}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("failed to create and push tag: %w", err)
		}
	}

	return nil
}

// Build creates a container with all dependencies installed and configured.
// It returns the configured container or nil if the build fails.
func (p *Python) Build(ctx context.Context, source *dagger.Directory) *dagger.Container {
	container := dag.Container()
	
	// Add Docker Hub authentication if credentials are provided
	if p.dockerUsername != "" && p.dockerPassword != nil {
		container = container.WithRegistryAuth("docker.io", p.dockerUsername, p.dockerPassword)
	}
	
	return container.
		From(fmt.Sprintf("python:%s", p.pythonVersion)).
		WithDirectory(containerWorkdir, dag.Poetry().Install(source)).
		WithWorkdir(containerWorkdir)
}

// Test runs all quality checks and returns the combined test output.
// It returns an error if any check fails.
func (p *Python) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	var testOutput string
	var err error

	if !p.skipTests {
		fmt.Println(logStartTests)
		// Run tests using Poetry module
		testOutput, err = dag.Poetry().Test(ctx, source)
		if err != nil {
			return "", fmt.Errorf("%s: %w", errPoetryTest, err)
		}
		fmt.Println(logSuccessTests)
	}

	// Run linting checks if not skipped
	if !p.skipLint {
		if err := p.Lint(ctx, source); err != nil {
			return "", err
		}
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

