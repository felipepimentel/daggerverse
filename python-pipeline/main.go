// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"

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

		// Create a backup of the original pyproject.toml
		container = container.WithExec([]string{"cp", "pyproject.toml", "pyproject.toml.bak"})

		// Add semantic-release config if it doesn't exist
		container = container.WithExec([]string{"bash", "-c", `
if ! grep -q "\[tool.semantic_release\]" pyproject.toml; then
	cat >> pyproject.toml << 'EOF'

[tool.semantic_release]
version_variables = ["pyproject.toml:version"]
commit_author = "github-actions[bot] <github-actions[bot]@users.noreply.github.com>"
commit_parser = "angular"
branch = "main"
upload_to_pypi = true
build_command = "poetry build"
repository = "pepperpy-core"
repository_owner = "felipepimentel"
EOF
fi`})

		// Set GH_TOKEN environment variable for semantic-release
		container = container.WithEnvVariable("GH_TOKEN", "ghp_your_token_here")

		// Run semantic-release version to determine and update version
		_, err = container.WithExec([]string{
			"semantic-release",
			"version",
		}).Stdout(ctx)
		if err != nil {
			// Restore original pyproject.toml if semantic-release fails
			container = container.WithExec([]string{"mv", "pyproject.toml.bak", "pyproject.toml"})
			return fmt.Errorf("error running semantic-release version: %v", err)
		}

		// Clean up backup file
		container = container.WithExec([]string{"rm", "-f", "pyproject.toml.bak"})

		// Run semantic-release publish to handle publishing
		_, err = container.WithExec([]string{
			"semantic-release",
			"publish",
		}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error running semantic-release publish: %v", err)
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