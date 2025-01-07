// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/python-pipeline/internal/dagger"
)

// PythonPipeline orchestrates Python project workflows using Poetry and PyPI.
type PythonPipeline struct {
	client *dagger.Client
}

// New creates a new instance of PythonPipeline.
func New(ctx context.Context) (*PythonPipeline, error) {
	client := dagger.Connect()
	if client == nil {
		return nil, fmt.Errorf("failed to initialize Dagger client")
	}
	return &PythonPipeline{
		client: client,
	}, nil
}

// ContainerConfig holds configuration for the base container.
type ContainerConfig struct {
	pythonVersion string
	gitEmail      string
	gitName       string
}

// DefaultContainerConfig returns default container configuration.
func DefaultContainerConfig() ContainerConfig {
	return ContainerConfig{
		pythonVersion: "3.12-slim",
		gitEmail:      "github-actions[bot]@users.noreply.github.com",
		gitName:       "github-actions[bot]",
	}
}

// setupContainer configures the base container with required dependencies.
func (p *PythonPipeline) setupContainer(ctx context.Context, source *dagger.Directory, config ContainerConfig) (*dagger.Container, error) {
	if p.client == nil {
		return nil, fmt.Errorf("Dagger client is not initialized")
	}

	container := p.client.Container().
		From(fmt.Sprintf("python:%s", config.pythonVersion)).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "curl", "ca-certificates"}).
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"}).
		WithExec([]string{"git", "config", "--global", "user.email", config.gitEmail}).
		WithExec([]string{"git", "config", "--global", "user.name", config.gitName})

	return container, nil
}

// getVersion retrieves the next version using the versioner module.
func (p *PythonPipeline) getVersion(ctx context.Context, source *dagger.Directory) (string, error) {
	versionerModule := p.client.Versioner()
	version, err := versionerModule.BumpVersion(ctx, source, true)
	if err != nil {
		return "", fmt.Errorf("error running versioner module: %w", err)
	}

	if version == "" {
		return "", fmt.Errorf("invalid version returned from versioner module")
	}

	fmt.Printf("Using version: %s\n", version)
	return version, nil
}

// QualityCheck represents a code quality check to be performed.
type QualityCheck struct {
	name    string
	command []string
}

// DefaultQualityChecks returns the standard set of quality checks.
func DefaultQualityChecks() []QualityCheck {
	return []QualityCheck{
		{"tests", []string{"poetry", "run", "pytest"}},
		{"black", []string{"poetry", "run", "black", ".", "--check"}},
		{"ruff", []string{"poetry", "run", "ruff", "check", "."}},
	}
}

// runQualityChecks executes tests and code quality checks.
func (p *PythonPipeline) runQualityChecks(ctx context.Context, container *dagger.Container, checks []QualityCheck) error {
	for _, check := range checks {
		if _, err := container.WithExec(check.command).Stdout(ctx); err != nil {
			return fmt.Errorf("error running %s: %w", check.name, err)
		}
	}
	return nil
}

// publishToPyPI handles the PyPI publishing process.
func (p *PythonPipeline) publishToPyPI(ctx context.Context, container *dagger.Container, token *dagger.Secret) error {
	tokenValue, err := token.Plaintext(ctx)
	if err != nil || tokenValue == "" {
		return fmt.Errorf("invalid PYPI_TOKEN: %w", err)
	}

	fmt.Printf("Successfully read PYPI_TOKEN: %s\n", tokenValue[:5])

	_, err = container.
		WithSecretVariable("PYPI_TOKEN", token).
		WithExec([]string{"poetry", "build"}).
		WithExec([]string{"poetry", "publish", "--no-interaction"}).
		Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to publish to PyPI: %w", err)
	}

	return nil
}

// CICD runs the complete CI/CD pipeline for a Python project.
func (p *PythonPipeline) CICD(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	config := DefaultContainerConfig()

	container, err := p.setupContainer(ctx, source, config)
	if err != nil {
		return fmt.Errorf("failed to setup container: %w", err)
	}

	version, err := p.getVersion(ctx, source)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	container = container.
		WithEnvVariable("VERSION", version).
		WithExec([]string{"poetry", "install", "--no-interaction"})

	if err := p.runQualityChecks(ctx, container, DefaultQualityChecks()); err != nil {
		return fmt.Errorf("quality checks failed: %w", err)
	}

	if token != nil {
		if err := p.publishToPyPI(ctx, container, token); err != nil {
			return fmt.Errorf("failed to publish to PyPI: %w", err)
		}
	}

	return nil
}
