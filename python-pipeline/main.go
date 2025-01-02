// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/python-pipeline/internal/dagger"
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
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Install dependencies
	container = container.WithExec([]string{"poetry", "install", "--no-interaction"})

	// Run tests
	_, err := container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running tests: %v", err)
	}

	// Run linting if configured (looking for .pylintrc or setup.cfg)
	if _, err := container.WithExec([]string{"test", "-f", ".pylintrc"}).Stdout(ctx); err == nil {
		_, err := container.WithExec([]string{"poetry", "run", "pylint", "."}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error running linting: %v", err)
		}
	}

	// Build package
	container = container.WithExec([]string{"poetry", "build"})

	// If token is provided, publish to PyPI
	if token != nil {
		container = container.WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token)
		
		// Check if version already exists
		version, err := container.WithExec([]string{"poetry", "version", "--short"}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error getting package version: %v", err)
		}

		// Try to get package info from PyPI
		checkCmd := fmt.Sprintf("pip install %s==%s 2>/dev/null || echo 'Version not found'", 
			"pepperpy-core", // TODO: Get package name dynamically
			version)
		out, err := container.WithExec([]string{"sh", "-c", checkCmd}).Stdout(ctx)
		if err == nil && !strings.Contains(out, "Version not found") {
			return fmt.Errorf("version %s already exists on PyPI", version)
		}

		// Publish if version doesn't exist
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