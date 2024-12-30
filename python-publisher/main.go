// Package main provides functionality for publishing Python packages to PyPI.
// It handles authentication and package publishing using Poetry in a containerized environment.
package main

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
)

// PythonPublisher handles the publishing process for Python packages to PyPI.
// It uses Poetry for package publishing and requires a PyPI authentication token.
type PythonPublisher struct {
	// PackagePath specifies the path to the Python package within the source directory.
	// This path should contain the pyproject.toml file.
	PackagePath string
}

// Publish uploads the Python package to PyPI.
// The process includes:
// 1. Locating the pyproject.toml file
// 2. Setting up a Python container with Poetry
// 3. Configuring PyPI authentication
// 4. Building and publishing the package
//
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// - token: The PyPI authentication token
//
// Returns:
// - error: Any error that occurred during the publishing process
func (m *PythonPublisher) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	// Find pyproject.toml location
	packagePath, err := m.findPyProjectToml(source)
	if err != nil {
		return fmt.Errorf("error finding pyproject.toml: %v", err)
	}

	client, err := dagger.Connect(ctx)
	if err != nil {
		return fmt.Errorf("error connecting to dagger: %v", err)
	}
	defer client.Close()

	// Setup Python container with Poetry
	container := client.Container().
		From("python:3.11-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src/" + packagePath).
		WithExec([]string{"pip", "install", "poetry"})

	// Configure PyPI token
	tokenValue, err := token.Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("error getting token value: %v", err)
	}

	container = container.WithExec([]string{
		"poetry", "config", "pypi-token.pypi", tokenValue,
	})

	// Publish to PyPI
	_, err = container.WithExec([]string{"poetry", "publish", "--build"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error publishing package: %v", err)
	}

	return nil
}

// findPyProjectToml searches for the pyproject.toml file in the specified package path.
// If the file is not found in the package path, it returns an error.
//
// Parameters:
// - source: The source directory to search in
//
// Returns:
// - string: The path where pyproject.toml was found
// - error: An error if the file was not found
func (m *PythonPublisher) findPyProjectToml(source *dagger.Directory) (string, error) {
	// Try package path first
	if m.PackagePath != "" {
		if _, err := source.File(filepath.Join(m.PackagePath, "pyproject.toml")).Contents(context.Background()); err == nil {
			return m.PackagePath, nil
		}
	}

	// Return error if not found
	return "", fmt.Errorf("pyproject.toml not found")
} 