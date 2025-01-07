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
func (p *PythonPipeline) setupContainer(ctx context.Context, client *dagger.Client, source *dagger.Directory, config ContainerConfig) (*dagger.Container, error) {
	fmt.Println("Setting up container...")

	container := client.Container().
		From(fmt.Sprintf("python:%s", config.pythonVersion)).
		WithDirectory("/project", source).
		WithWorkdir("/project").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "curl", "ca-certificates"}).
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"}).
		WithExec([]string{"git", "config", "--global", "user.email", config.gitEmail}).
		WithExec([]string{"git", "config", "--global", "user.name", config.gitName}).
		WithExec([]string{"poetry", "install", "--with", "dev", "--no-interaction"})

	fmt.Println("Container setup completed!")
	return container, nil
}

// getVersion retrieves the next version using the versioner module.
func (p *PythonPipeline) getVersion(ctx context.Context, client *dagger.Client, source *dagger.Directory) (string, error) {
	fmt.Println("Getting version from versioner module...")

	versionerModule := client.Versioner()
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

// runQualityChecks executes quality checks such as tests and linters.
func (p *PythonPipeline) runQualityChecks(ctx context.Context, client *dagger.Client, container *dagger.Container) error {
	qualityChecks := []struct {
		name    string
		command []string
	}{
		{"pytest", []string{"poetry", "run", "pytest"}},
		{"black", []string{"poetry", "run", "black", ".", "--check"}},
		{"ruff", []string{"poetry", "run", "ruff", "check", "."}},
	}

	for _, check := range qualityChecks {
		fmt.Printf("Running %s...\n", check.name)
		_, err := container.WithExec(check.command).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error running %s: %w", check.name, err)
		}
		fmt.Printf("%s passed successfully!\n", check.name)
	}
	return nil
}

// publishToPyPI publishes the package to PyPI if a token is provided.
func (p *PythonPipeline) publishToPyPI(ctx context.Context, client *dagger.Client, container *dagger.Container, token *dagger.Secret) error {
	fmt.Println("Publishing to PyPI...")

	container = container.WithSecretVariable("PYPI_TOKEN", token)

	// Build the package
	_, err := container.WithExec([]string{"poetry", "build"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to build the package: %w", err)
	}

	// Check dist directory content
	output, err := container.WithExec([]string{"ls", "-la", "/project/dist"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to list dist directory contents: %w", err)
	}

	fmt.Printf("Dist directory contents:\n%s\n", output)

	// Publish the package
	_, err = container.WithExec([]string{"poetry", "publish", "--no-interaction"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish to PyPI: %w", err)
	}

	fmt.Println("Package published successfully to PyPI!")
	return nil
}

// CICD runs the complete CI/CD pipeline for a Python project.
func (p *PythonPipeline) CICD(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	fmt.Println("Starting CI/CD pipeline...")

	// Initialize the Dagger client
	client := dagger.Connect()
	if client == nil {
		return fmt.Errorf("failed to initialize Dagger client")
	}

	// Debug: Log the source directory structure
	dirContents, err := client.Container().WithDirectory("/project", source).WithExec([]string{"ls", "-la", "/project"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to list source directory contents: %w", err)
	}
	fmt.Printf("Source directory contents:\n%s\n", dirContents)

	// Load default container configuration
	config := DefaultContainerConfig()

	// Setup the container
	container, err := p.setupContainer(ctx, client, source, config)
	if err != nil {
		return fmt.Errorf("failed to setup container: %w", err)
	}

	// Get the next version
	version, err := p.getVersion(ctx, client, source)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	// Set the version in the container environment
	container = container.WithEnvVariable("VERSION", version)

	// Install dependencies
	_, err = container.WithExec([]string{"poetry", "install", "--no-interaction"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	// Run quality checks
	err = p.runQualityChecks(ctx, client, container)
	if err != nil {
		return fmt.Errorf("quality checks failed: %w", err)
	}

	// Publish to PyPI if token is provided
	if token != nil {
		err = p.publishToPyPI(ctx, client, container, token)
		if err != nil {
			return fmt.Errorf("failed to publish to PyPI: %w", err)
		}
	}

	fmt.Println("CI/CD pipeline completed successfully!")
	return nil
}
