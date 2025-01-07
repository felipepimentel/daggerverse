// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"
	"strings"

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

	// Create base container with Python and mount source code
	container := client.Container().
		From(fmt.Sprintf("python:%s", config.pythonVersion)).
		WithDirectory("/src", source).
		WithWorkdir("/src")

	// Install system dependencies and tools
	container = container.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "curl", "ca-certificates"}).
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Configure git
	container = container.
		WithExec([]string{"git", "config", "--global", "user.email", config.gitEmail}).
		WithExec([]string{"git", "config", "--global", "user.name", config.gitName})

	// Install poetry dependencies from the source project
	container = container.
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
	fmt.Println("Building and publishing to PyPI...")

	// Get current version from environment
	version, err := container.EnvVariable(ctx, "VERSION")
	if err != nil {
		return fmt.Errorf("failed to get VERSION environment variable: %w", err)
	}
	if version == "" {
		return fmt.Errorf("VERSION environment variable not set")
	}

	container = container.WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token)

	// Update version in pyproject.toml before building
	_, err = container.WithExec([]string{
		"poetry", "version", version,
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to update version in pyproject.toml: %w", err)
	}

	// Verify the version was updated correctly
	output, err := container.WithExec([]string{
		"poetry", "version", "--short",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to verify version: %w", err)
	}
	currentVersion := strings.TrimSpace(output)
	if currentVersion != version {
		return fmt.Errorf("version mismatch: expected %s, got %s", version, currentVersion)
	}

	// Clean dist directory before building
	_, err = container.WithExec([]string{
		"rm", "-rf", "dist",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to clean dist directory: %w", err)
	}

	// Build and publish in sequence to ensure version consistency
	commands := []struct {
		name    string
		command []string
	}{
		{"configure token", []string{"poetry", "config", "pypi-token.pypi", "$POETRY_PYPI_TOKEN_PYPI"}},
		{"build package", []string{"poetry", "build"}},
		{"verify dist", []string{"ls", "-la", "dist"}},
		{"publish package", []string{"poetry", "publish", "--username", "__token__", "--no-interaction"}},
	}

	for _, cmd := range commands {
		fmt.Printf("Executing: %s\n", cmd.name)
		_, err = container.WithExec(cmd.command).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("failed to %s: %w", cmd.name, err)
		}
	}

	fmt.Printf("Package version %s built and published successfully to PyPI!\n", version)
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

	// Load default container configuration
	config := DefaultContainerConfig()

	// Setup the container
	container, err := p.setupContainer(ctx, client, source, config)
	if err != nil {
		return fmt.Errorf("failed to setup container: %w", err)
	}

	// Debug: Log both root and source directory structures
	fmt.Println("Root directory contents:")
	rootContents, err := container.WithExec([]string{"ls", "-la", "/"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to list root directory contents: %w", err)
	}
	fmt.Printf("%s\n", rootContents)

	fmt.Println("\nSource directory contents (/src):")
	srcContents, err := container.WithExec([]string{"ls", "-la", "/src"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to list source directory contents: %w", err)
	}
	fmt.Printf("%s\n", srcContents)

	// Get the next version
	version, err := p.getVersion(ctx, client, source)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	fmt.Printf("Using version: %s\n", version)

	// Set the version in the container environment
	container = container.WithEnvVariable("VERSION", version)

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
