// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"

	"dagger/python-pipeline/internal/dagger"
)

// PythonPipeline orchestrates Python project workflows using Poetry and PyPI.
type PythonPipeline struct{}

// New creates a new instance of PythonPipeline.
func New(ctx context.Context) (*PythonPipeline, error) {
	return &PythonPipeline{}, nil
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
	client := dagger.Connect()

	// Get Poetry module
	poetry := dag.PythonPoetry()
	
	// Get PyPI module
	pypi := dag.PythonPypi()

	// Install dependencies
	installed, err := poetry.With(source).Install(ctx)
	if err != nil {
		return fmt.Errorf("error installing dependencies: %v", err)
	}

	// Run tests
	testOutput, err := poetry.With(installed).Test(ctx)
	if err != nil {
		return fmt.Errorf("error running tests: %v", err)
	}
	fmt.Println("Test output:", testOutput)

	// Build package
	built, err := poetry.With(installed).Build(ctx)
	if err != nil {
		return fmt.Errorf("error building package: %v", err)
	}

	// Publish to PyPI
	err = pypi.With(built).WithSecret("PYPI_TOKEN", token).Publish(ctx)
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
	client := dagger.Connect()

	// Get Poetry module
	poetry := client.Container().Import("python-poetry")

	// Update dependencies
	updated, err := poetry.With(source).Update(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating dependencies: %v", err)
	}

	// Update lock file
	locked, err := poetry.With(updated).Lock(ctx)
	if err != nil {
		return nil, fmt.Errorf("error updating lock file: %v", err)
	}

	return locked, nil
} 