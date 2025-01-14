package main

import (
	"context"
	"dagger/mkdocs/internal/dagger"
	"fmt"
)

// MkDocs represents a module for building and deploying MkDocs documentation
type MkDocs struct {
	// Python version to use
	PythonVersion string
}

// Default configuration values
const (
	defaultPythonVersion = "3.11"
	defaultMkdocsTheme  = "material"
)

type MkDocsConfig struct {
	// Source directory containing mkdocs.yml and docs/
	Source *dagger.Directory
	// Custom requirements file (optional)
	RequirementsFile *dagger.File
	// Output directory name (default: "site")
	OutputDir string
	// Base URL for the documentation
	BaseURL string
	// Whether to use strict mode
	Strict bool
	// Whether to minify HTML
	Minify bool
	// Whether to include git revision date
	GitRevisionDate bool
}

// Container returns a base Python container with MkDocs dependencies
func (m *MkDocs) Container() *dagger.Container {
	pythonVersion := m.PythonVersion
	if pythonVersion == "" {
		pythonVersion = defaultPythonVersion
	}

	return dag.Container().
		From(fmt.Sprintf("python:%s-slim", pythonVersion)).
		WithExec([]string{"pip", "install", "--no-cache-dir",
			"mkdocs-material",
			"mkdocs-minify-plugin",
			"mkdocs-git-revision-date-localized-plugin",
			"pillow",
			"cairosvg",
		})
}

// Build builds the MkDocs documentation
func (m *MkDocs) Build(ctx context.Context, config *MkDocsConfig) (*dagger.Directory, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}

	if config.Source == nil {
		return nil, fmt.Errorf("source directory is required")
	}

	container := m.Container()

	// Mount source directory
	container = container.WithMountedDirectory("/src", config.Source)
	container = container.WithWorkdir("/src")

	// Install custom requirements if provided
	if config.RequirementsFile != nil {
		container = container.
			WithMountedFile("/src/requirements.txt", config.RequirementsFile).
			WithExec([]string{"pip", "install", "--no-cache-dir", "-r", "requirements.txt"})
	}

	// Build command
	buildCmd := []string{"mkdocs", "build"}

	if config.Strict {
		buildCmd = append(buildCmd, "--strict")
	}

	if config.BaseURL != "" {
		buildCmd = append(buildCmd, "--site-url", config.BaseURL)
	}

	// Execute build
	container = container.WithExec(buildCmd)

	// Return the built site directory
	outputDir := "site"
	if config.OutputDir != "" {
		outputDir = config.OutputDir
	}

	return container.Directory(outputDir), nil
}

// Serve starts a development server (useful for local development)
func (m *MkDocs) Serve(config *MkDocsConfig) *dagger.Container {
	container := m.Container()

	if config.Source != nil {
		container = container.WithMountedDirectory("/src", config.Source)
		container = container.WithWorkdir("/src")
	}

	return container.WithExec([]string{"mkdocs", "serve", "--dev-addr", "0.0.0.0:8000"})
}

// Deploy builds the documentation and returns a container ready to deploy
func (m *MkDocs) Deploy(ctx context.Context, config *MkDocsConfig) (*dagger.Directory, error) {
	return m.Build(ctx, config)
}

// ValidateConfig validates the MkDocs configuration
func (m *MkDocs) ValidateConfig(ctx context.Context, config *MkDocsConfig) (bool, error) {
	if config == nil || config.Source == nil {
		return false, fmt.Errorf("invalid configuration: source directory is required")
	}

	container := m.Container().
		WithMountedDirectory("/src", config.Source).
		WithWorkdir("/src")

	_, err := container.WithExec([]string{"mkdocs", "build", "--strict", "--dry-run"}).
		Stdout(ctx)

	return err == nil, err
}
