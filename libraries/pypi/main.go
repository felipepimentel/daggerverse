// Package main provides functionality for publishing Python packages to PyPI.
package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/libraries/pypi/internal/dagger"
)

// Pypi handles publishing Python packages to PyPI.
type Pypi struct {
	// Base image for PyPI operations
	// +private
	BaseImage string
}

// New creates a new instance of Pypi.
func New(
	// Base Python image to use
	// +optional
	// +default="python:3.12-slim"
	baseImage string,
) *Pypi {
	if baseImage == "" {
		baseImage = "python:3.12-slim"
	}

	return &Pypi{
		BaseImage: baseImage,
	}
}

// Publish publishes a Python package to PyPI.
// It expects the source directory to contain the built package (dist directory).
func (m *Pypi) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	container := dag.Container().
		From(m.BaseImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"}).
		WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token)

	// Publish package
	_, err := container.WithExec([]string{"poetry", "publish"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error publishing package: %v", err)
	}

	return nil
}