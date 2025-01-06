// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/felipepimentel/daggerverse/python-pipeline/internal/dagger"
)

// PythonPipeline orchestrates Python project workflows using Poetry and PyPI.
type PythonPipeline struct{}

// New creates a new instance of PythonPipeline.
func New() *PythonPipeline {
	return &PythonPipeline{}
}

// CICD runs the complete CI/CD pipeline for a Python project.
func (m *PythonPipeline) CICD(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	// Retrieve version from environment variable
	version := os.Getenv("VERSION")
	if version == "" {
		return fmt.Errorf("VERSION environment variable is not set")
	}

	// Setup Python container with Poetry
	container := dag.Container().
		From("python:3.12-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "curl", "ca-certificates"}).
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Configure git
	container = container.
		WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"})

	// Ensure pyproject.toml is correctly configured
	container = container.WithExec([]string{"bash", "-c", `
		if ! grep -q "[tool.semantic_release]" pyproject.toml; then
			echo "Adding semantic-release configuration to pyproject.toml"
			cat >> pyproject.toml <<EOL

[tool.semantic_release]
version_variables = ["pyproject.toml:version"]
EOL
		fi
	`})

	// Install dependencies
	container = container.WithExec([]string{"poetry", "install", "--no-interaction"})

	// Run tests
	_, err := container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running tests: %v", err)
	}

	// Build and publish package
	container = container.WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token)
	_, err = container.WithExec([]string{"poetry", "publish", "--no-interaction"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error publishing to PyPI: %v", err)
	}

	return nil
}
