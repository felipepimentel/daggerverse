// Package main provides functionality for managing Python projects with Poetry.
package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/libraries/poetry/internal/dagger"
)

// Poetry handles Python project management with Poetry.
type Poetry struct {
	// Base image for Poetry operations
	// +private
	BaseImage string
}

// New creates a new instance of Poetry with the provided configuration.
func New(
	// Base Python image to use
	// +optional
	// +default="python:3.12-alpine"
	baseImage string,
) *Poetry {
	if baseImage == "" {
		baseImage = "python:3.12-alpine"
	}

	return &Poetry{
		BaseImage: baseImage,
	}
}

// getBaseContainer returns a configured base container with Poetry installed
func (m *Poetry) getBaseContainer(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From(m.BaseImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})
}

// Install installs project dependencies using Poetry.
func (m *Poetry) Install(source *dagger.Directory) *dagger.Directory {
	container := m.getBaseContainer(source).
		WithExec([]string{"poetry", "config", "virtualenvs.create", "false"}).
		WithExec([]string{"poetry", "install", "--no-interaction"})

	return container.Directory("/src")
}

// Build builds the Python package using Poetry.
func (m *Poetry) Build(source *dagger.Directory) *dagger.Directory {
	container := m.getBaseContainer(source).
		WithExec([]string{"poetry", "config", "virtualenvs.create", "false"}).
		WithExec([]string{"poetry", "install", "--no-interaction"}).
		WithExec([]string{"poetry", "build"})

	return container.Directory("/src/dist")
}

// BuildWithVersion builds the Python package using Poetry with a specific version.
func (m *Poetry) BuildWithVersion(source *dagger.Directory, version string) *dagger.Directory {
	container := m.getBaseContainer(source).
		WithExec([]string{"poetry", "config", "virtualenvs.create", "false"}).
		WithExec([]string{"poetry", "version", version}).
		WithExec([]string{"poetry", "install", "--no-interaction"}).
		WithExec([]string{"poetry", "build"})

	return container.Directory("/src/dist")
}

// Test runs tests using Poetry.
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
func (m *Poetry) Lock(source *dagger.Directory) *dagger.Directory {
	container := m.getBaseContainer(source).
		WithExec([]string{"poetry", "lock", "--no-update"})

	return container.Directory("/src")
}

// Update updates dependencies to their latest versions.
func (m *Poetry) Update(source *dagger.Directory) *dagger.Directory {
	container := m.getBaseContainer(source).
		WithExec([]string{"poetry", "update"})

	return container.Directory("/src")
}