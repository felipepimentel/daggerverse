// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"

	"dagger/python-runtime/internal/dagger"
)

// Default configurations
const (
	// DefaultPythonVersion is the default Python version to use
	DefaultPythonVersion = "3.12-slim"
	// DefaultGitEmail is the default Git email for commits
	DefaultGitEmail = "github-actions[bot]@users.noreply.github.com"
	// DefaultGitName is the default Git username for commits
	DefaultGitName = "github-actions[bot]"
)

// PythonPipeline provides a complete pipeline for Python projects.
type PythonPipeline struct {
	// Python version to use
	pythonVersion string
	// Git configuration
	gitEmail string
	gitName  string
}

// Option configures a PythonPipeline
type Option func(*PythonPipeline)

// WithPythonVersion sets the Python version
func WithPythonVersion(version string) Option {
	return func(p *PythonPipeline) {
		p.pythonVersion = version
	}
}

// WithGitConfig sets the Git configuration
func WithGitConfig(email, name string) Option {
	return func(p *PythonPipeline) {
		p.gitEmail = email
		p.gitName = name
	}
}

// New creates a new PythonPipeline instance with the provided options.
// If no options are provided, default values will be used.
func New(opts ...Option) *PythonPipeline {
	p := &PythonPipeline{
		pythonVersion: DefaultPythonVersion,
		gitEmail:      DefaultGitEmail,
		gitName:       DefaultGitName,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Build creates a development environment with all dependencies installed.
func (m *PythonPipeline) Build(ctx context.Context, source *dagger.Directory) (*dagger.Container, error) {
	// Use Poetry module directly
	poetry := dag.Poetry().WithPythonVersion(m.pythonVersion)
	
	// Install dependencies using Poetry module
	container, err := poetry.Install(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("failed to install dependencies: %w", err)
	}

	return container, nil
}

// Test runs the test suite.
func (m *PythonPipeline) Test(ctx context.Context, source *dagger.Directory) error {
	// Use Poetry module directly
	poetry := dag.Poetry().WithPythonVersion(m.pythonVersion)

	// Run tests using Poetry module
	output, err := poetry.Test(ctx, source)
	if err != nil {
		return fmt.Errorf("test execution failed: %w", err)
	}

	fmt.Println(output)
	return nil
}

// Lint runs the linter.
func (m *PythonPipeline) Lint(ctx context.Context, source *dagger.Directory) error {
	// Use Ruff module directly
	ruff := dag.Ruff()
	lintRun := ruff.Lint(source)

	// Print summary
	summary, err := lintRun.Summary(ctx)
	if err != nil {
		return fmt.Errorf("failed to get lint summary: %w", err)
	}
	fmt.Println(summary)

	// Assert no errors
	return lintRun.Assert(ctx)
}

// Publish builds and publishes the package to PyPI.
func (m *PythonPipeline) Publish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
	// Run tests first
	if err := m.Test(ctx, source); err != nil {
		return "", fmt.Errorf("tests failed: %w", err)
	}

	// Run linting
	if err := m.Lint(ctx, source); err != nil {
		return "", fmt.Errorf("linting failed: %w", err)
	}

	// Get version using Versioner module
	version, err := dag.Versioner().BumpVersion(ctx, source, true)
	if err != nil {
		return "", fmt.Errorf("failed to bump version: %w", err)
	}
	fmt.Println("Using version:", version)

	// Configure git using Git module
	gitRepo := dag.Git(".")
	gitRepo = gitRepo.
		WithConfig("user.email", m.gitEmail).
		WithConfig("user.name", m.gitName)

	// Build using Poetry module
	poetry := dag.Poetry().WithPythonVersion(m.pythonVersion)
	buildDir, err := poetry.Build(ctx, source)
	if err != nil {
		return "", fmt.Errorf("failed to build package: %w", err)
	}

	// Publish using PyPI module
	pypi := dag.Pypi()
	address, err := pypi.Publish(ctx, buildDir, token)
	if err != nil {
		return "", fmt.Errorf("failed to publish package: %w", err)
	}

	fmt.Println("Package published successfully to:", address)
	return address, nil
} 