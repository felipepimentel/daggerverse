// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"

	"dagger/python-runtime/internal/dagger"
)

// PythonPipeline provides a complete pipeline for Python projects.
type PythonPipeline struct {
	// Git configuration
	gitEmail string
	gitName  string
	// +private
	client *dagger.Client
}

// New creates a new PythonPipeline instance.
func New(
	// Git user email
	gitEmail string,
	// Git user name
	gitName string,
) *PythonPipeline {
	client := dagger.Connect()
	return &PythonPipeline{
		gitEmail: gitEmail,
		gitName:  gitName,
		client:   client,
	}
}

// Build creates a development environment with all dependencies installed.
func (m *PythonPipeline) Build(ctx context.Context, source *dagger.Directory) *dagger.Container {
	poetry := m.client.Poetry()

	// Install dependencies using Poetry module
	dir, err := poetry.Install(ctx, source)
	if err != nil {
		return nil
	}

	return dir.AsContainer()
}

// Test runs the test suite.
func (m *PythonPipeline) Test(ctx context.Context, source *dagger.Directory) error {
	poetry := m.client.Poetry()

	// Run tests using Poetry module
	output, err := poetry.Test(ctx, source)
	if err != nil {
		return err
	}

	fmt.Println(output)
	return nil
}

// Lint runs the linter.
func (m *PythonPipeline) Lint(ctx context.Context, source *dagger.Directory) error {
	// Use Ruff module for linting
	ruff := m.client.Ruff()
	lintRun := ruff.Lint(source)

	// Print summary
	summary, err := lintRun.Summary(ctx)
	if err != nil {
		return err
	}
	fmt.Println(summary)

	// Assert no errors
	return lintRun.Assert(ctx)
}

// Publish builds and publishes the package to PyPI.
func (m *PythonPipeline) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
	// Get version using Versioner module
	version, err := m.client.Versioner().BumpVersion(ctx, source, true)
	if err != nil {
		return "", err
	}
	fmt.Println("Using version:", version)

	// Configure git using Git module
	gitRepo := m.client.Git(".")
	gitRepo = gitRepo.
		WithConfig("user.email", m.gitEmail).
		WithConfig("user.name", m.gitName)

	// Build using Poetry module
	poetry := m.client.Poetry()
	buildDir, err := poetry.Build(ctx, source)
	if err != nil {
		return "", err
	}

	// Publish using PyPI module
	pypi := m.client.Pypi()
	address, err := pypi.Publish(ctx, buildDir, token)
	if err != nil {
		return "", err
	}

	fmt.Println("Package published successfully to:", address)
	return address, nil
} 