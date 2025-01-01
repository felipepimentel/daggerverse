// Package main provides functionality for managing Python projects with Poetry.
package main

import (
	"context"
	"fmt"

	"dagger/python-poetry/internal/dagger"
)

// PythonPoetry handles Python project management with Poetry.
type PythonPoetry struct{}

// New creates a new instance of PythonPoetry.
func New(ctx context.Context) (*PythonPoetry, error) {
	return &PythonPoetry{}, nil
}

// Install installs project dependencies using Poetry.
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// Returns:
// - *dagger.Directory: The directory with installed dependencies
// - error: Any error that occurred during installation
func (m *PythonPoetry) Install(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	client := dagger.Connect()

	container := client.Container().
		From("python:3.12-alpine").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"}).
		WithExec([]string{"poetry", "config", "virtualenvs.create", "false"}).
		WithExec([]string{"poetry", "install", "--no-interaction"})

	return container.Directory("/src"), nil
}

// Build builds the Python package using Poetry.
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// Returns:
// - *dagger.Directory: The directory containing the built package
// - error: Any error that occurred during build
func (m *PythonPoetry) Build(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	client := dagger.Connect()

	container := client.Container().
		From("python:3.12-alpine").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"}).
		WithExec([]string{"poetry", "build"})

	return container.Directory("/src/dist"), nil
}

// Test runs tests using Poetry.
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// Returns:
// - string: The test output
// - error: Any error that occurred during testing
func (m *PythonPoetry) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	client := dagger.Connect()

	container := client.Container().
		From("python:3.12-alpine").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"}).
		WithExec([]string{"poetry", "config", "virtualenvs.create", "false"}).
		WithExec([]string{"poetry", "install", "--no-interaction"})

	output, err := container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error running tests: %v", err)
	}

	return output, nil
}

// Lock updates the poetry.lock file.
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// Returns:
// - *dagger.Directory: The directory containing the updated lock file
// - error: Any error that occurred during lock update
func (m *PythonPoetry) Lock(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	client := dagger.Connect()

	container := client.Container().
		From("python:3.12-alpine").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"}).
		WithExec([]string{"poetry", "lock", "--no-update"})

	return container.Directory("/src"), nil
}

// Update updates dependencies to their latest versions.
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// Returns:
// - *dagger.Directory: The directory with updated dependencies
// - error: Any error that occurred during update
func (m *PythonPoetry) Update(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	client := dagger.Connect()

	container := client.Container().
		From("python:3.12-alpine").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"}).
		WithExec([]string{"poetry", "update"})

	return container.Directory("/src"), nil
} 