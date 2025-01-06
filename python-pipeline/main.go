// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"

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

	// Call the versioner module to get the next version
	versionerModule := dag.Versioner()
	version, err := versionerModule.BumpVersion(ctx, source, true)
	if err != nil {
		return fmt.Errorf("error running versioner module: %w", err)
	}

	if version == "" {
		return fmt.Errorf("invalid version returned from versioner module")
	}

	fmt.Printf("Using version: %s\n", version)

	// Pass version as an environment variable for subsequent steps
	container = container.WithEnvVariable("VERSION", version)

	// Install dependencies
	container = container.WithExec([]string{"poetry", "install", "--no-interaction"})

	// Run tests
	_, err = container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running tests: %v", err)
	}

	// Run black check
	_, err = container.WithExec([]string{"poetry", "run", "black", ".", "--check"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running black check: %v", err)
	}

	// Run ruff check
	_, err = container.WithExec([]string{"poetry", "run", "ruff", "check", "."}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running ruff check: %v", err)
	}

	// If token is provided, publish to PyPI
	if token != nil {
		tokenValue, err := token.Plaintext(ctx)
		if err != nil {
			fmt.Println("Error reading PYPI_TOKEN from environment:", err)
			return fmt.Errorf("failed to read PYPI_TOKEN value: %w", err)
		}

		if tokenValue == "" {
			fmt.Println("PYPI_TOKEN is set but empty")
			return fmt.Errorf("PYPI_TOKEN is empty")
		}

		fmt.Printf("Successfully read PYPI_TOKEN: %s\n", tokenValue[:5]) // Exibe apenas os primeiros 5 caracteres

		container = container.WithSecretVariable("PYPI_TOKEN", token)

		// Build the package before publishing
		_, err = container.WithExec([]string{"poetry", "build"}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("failed to build the package: %w", err)
		}

		// Publish the package to PyPI
		_, err = container.WithExec([]string{"poetry", "publish", "--no-interaction"}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("failed to publish to PyPI: %w", err)
		}
	}

	return nil
}
