// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipepimentel/daggerverse/python-pipeline/internal/dagger"
)

// PythonPipeline orchestrates Python project workflows using Poetry and PyPI.
type PythonPipeline struct{}

// New creates a new instance of PythonPipeline.
func New() *PythonPipeline {
	return &PythonPipeline{}
}

// CICD runs the complete CI/CD pipeline for a Python project.
// This includes:
// 1. Installing dependencies
// 2. Running tests
// 3. Running linting (if configured)
// 4. Building the package
// 5. Publishing to PyPI (if token is provided)
//
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// - token: Optional PyPI token for publishing. If provided, the package will be published
//
// Returns:
// - error: Any error that occurred during the process
func (m *PythonPipeline) CICD(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	// Setup Python container with Poetry
	container := dag.Container().
		From("python:3.12-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git"}).
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Configure git
	container = container.
		WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"})

	// Install dependencies
	container = container.WithExec([]string{"poetry", "install", "--no-interaction"})

	// Run tests
	_, err := container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running tests: %v", err)
	}

	// Run black check
	_, err = container.WithExec([]string{"poetry", "run", "black", ".", "--check"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running black check: %v", err)
	}

	// Run ruff check
	_, err = container.WithExec([]string{"poetry", "run", "ruff", "check", "."}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running ruff check: %v", err)
	}

	// If token is provided, publish to PyPI
	if token != nil {
		container = container.WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token)

		// Install python-semantic-release
		container = container.WithExec([]string{"pip", "install", "python-semantic-release"})

		// Configure git
		container = container.
			WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
			WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"})

		// Run semantic release to determine version and create release
		version, err := container.WithExec([]string{
			"semantic-release",
			"version",
			"--print",
		}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error running semantic-release: %v", err)
		}

		// Clean version string (remove newline)
		version = strings.TrimSpace(version)

		// Get current version
		currentVersion, err := container.WithExec([]string{"poetry", "version", "--short"}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error getting current version: %v", err)
		}
		currentVersion = strings.TrimSpace(currentVersion)

		// Only update version if it's different
		if version != currentVersion {
			// Update version in pyproject.toml using poetry
			_, err = container.WithExec([]string{"poetry", "version", version}).Stdout(ctx)
			if err != nil {
				return fmt.Errorf("error updating version in pyproject.toml: %v", err)
			}

			// Check if there are changes to commit
			status, err := container.WithExec([]string{"git", "status", "--porcelain"}).Stdout(ctx)
			if err != nil {
				return fmt.Errorf("error checking git status: %v", err)
			}

			// Only commit if there are changes
			if strings.TrimSpace(status) != "" {
				container = container.
					WithExec([]string{"git", "add", "pyproject.toml"}).
					WithExec([]string{"git", "commit", "-m", fmt.Sprintf("chore(release): bump version to %s [skip ci]", version)}).
					WithExec([]string{"git", "push", "origin", "HEAD"})
			}
		}

		// Build and publish using poetry
		container = container.WithExec([]string{"poetry", "build"})
		_, err = container.WithExec([]string{"poetry", "publish", "--no-interaction"}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error publishing to PyPI: %v", err)
		}
	}

	return nil
}

// BuildAndPublish builds a Python package and publishes it to PyPI.
// The process includes:
// 1. Installing dependencies with Poetry
// 2. Running tests
// 3. Building the package
// 4. Publishing to PyPI
//
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// - token: The PyPI authentication token
//
// Returns:
// - error: Any error that occurred during the process
func (m *PythonPipeline) BuildAndPublish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	// Setup Python container with Poetry
	container := dag.Container().
		From("python:3.12-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Install dependencies
	container = container.WithExec([]string{"poetry", "install", "--no-interaction"})

	// Run tests
	_, err := container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running tests: %v", err)
	}

	// Build package
	container = container.WithExec([]string{"poetry", "build"})

	// Configure PyPI credentials and publish
	container = container.WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token)
	_, err = container.WithExec([]string{"poetry", "publish", "--no-interaction"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error publishing to PyPI: %v", err)
	}

	return nil
}

// UpdateDependencies updates project dependencies and lock file.
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
//
// Returns:
// - *dagger.Directory: The directory with updated dependencies
// - error: Any error that occurred during the update
func (m *PythonPipeline) UpdateDependencies(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	// Setup Python container with Poetry
	container := dag.Container().
		From("python:3.12-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Update dependencies
	container = container.WithExec([]string{"poetry", "update", "--no-interaction"})

	// Export the updated directory
	return container.Directory("/src"), nil
} 