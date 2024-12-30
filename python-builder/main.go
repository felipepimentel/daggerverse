// Package main provides functionality for building and testing Python projects using Poetry.
// It handles dependency management, test execution, and package building in a containerized environment.
package main

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
)

// PythonBuilder handles the build and test process for Python projects.
// It uses Poetry for dependency management and pytest for test execution.
type PythonBuilder struct {
	// PackagePath specifies the path to the Python package within the source directory.
	// This path should contain the pyproject.toml file.
	PackagePath string
}

// Build compiles and tests the Python project.
// The process includes:
// 1. Locating the pyproject.toml file
// 2. Setting up a Python container with Poetry
// 3. Installing project dependencies
// 4. Running tests with pytest
// 5. Building the package
//
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
//
// Returns:
// - error: Any error that occurred during the build process
func (m *PythonBuilder) Build(ctx context.Context, source *dagger.Directory) error {
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

	// Install dependencies
	container = container.WithExec([]string{"poetry", "install", "--no-interaction"})

	// Run tests
	_, err = container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running tests: %v", err)
	}

	// Build package
	_, err = container.WithExec([]string{"poetry", "build"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error building package: %v", err)
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
func (m *PythonBuilder) findPyProjectToml(source *dagger.Directory) (string, error) {
	// Try package path first
	if m.PackagePath != "" {
		if _, err := source.File(filepath.Join(m.PackagePath, "pyproject.toml")).Contents(context.Background()); err == nil {
			return m.PackagePath, nil
		}
	}

	// Return error if not found
	return "", fmt.Errorf("pyproject.toml not found")
} 