// Package main provides functionality for publishing Python packages to PyPI.
// It uses Poetry for package publishing.
package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/pypi/internal/dagger"
)

// Pypi handles publishing Python packages to PyPI.
type Pypi struct{}

// New creates a new instance of Pypi.
func New(ctx context.Context) (*Pypi, error) {
	return &Pypi{}, nil
}

// Publish publishes a Python package to PyPI using Poetry.
// The process includes:
// 1. Setting up a Python container with Poetry
// 2. Building the package
// 3. Publishing to PyPI
//
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// - token: The PyPI authentication token
//
// Returns:
// - error: Any error that occurred during the publishing process
func (m *Pypi) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	client := dagger.Connect()

	// Setup Python container with Poetry
	container := client.Container().
		From("python:3.12-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Configure PyPI credentials
	container = container.WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token)

	// Build and publish package
	_, publishErr := container.WithExec([]string{
		"poetry", "build",
	}).WithExec([]string{
		"poetry", "publish",
	}).Stdout(ctx)

	if publishErr != nil {
		return fmt.Errorf("error publishing package: %v", publishErr)
	}

	return nil
}