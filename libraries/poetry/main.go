// Package main provides functionality for managing Python projects with Poetry.
package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/libraries/poetry/internal/dagger"
)

// Poetry handles Python project management with Poetry.
type Poetry struct {
	// Add base container configuration
	baseImage string
	client    *dagger.Client
}

// New creates a new instance of Poetry.
func New(ctx context.Context) (*Poetry, error) {
	client := dagger.Connect()
	return &Poetry{
		baseImage: "python:3.12-alpine",
		client:    client,
	}, nil
}

// getBaseContainer returns a configured base container with Poetry installed
func (m *Poetry) getBaseContainer(source *dagger.Directory) *dagger.Container {
	return m.client.Container().
		From(m.baseImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})
}

// Install installs project dependencies using Poetry.
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// Returns:
// - *dagger.Directory: The directory with installed dependencies
// - error: Any error that occurred during installation
func (m *Poetry) Install(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	container := m.getBaseContainer(source).
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
func (m *Poetry) Build(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	container := m.getBaseContainer(source).
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
func (m *Poetry) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	container := m.getBaseContainer(source).
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
func (m *Poetry) Lock(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	container := m.getBaseContainer(source).
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
func (m *Poetry) Update(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	container := m.getBaseContainer(source).
		WithExec([]string{"poetry", "update"})

	return container.Directory("/src"), nil
}