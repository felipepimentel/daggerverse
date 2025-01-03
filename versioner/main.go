package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipepimentel/daggerverse/versioner/internal/dagger"
)

// Versioner implements version management for repositories
type Versioner struct{}

// New creates a new Versioner instance
func New() *Versioner {
	return &Versioner{}
}

// BumpVersion creates a new version tag based on the latest tag
func (m *Versioner) BumpVersion(ctx context.Context, source *dagger.Directory, outputVersion bool) (string, error) {
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git"})

	// Check if git is already initialized
	gitStatus, err := container.WithExec([]string{"sh", "-c", "[ -d .git ] && echo 'true' || echo 'false'"}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error checking git status: %w", err)
	}

	if strings.TrimSpace(gitStatus) == "false" {
		container = container.
			WithExec([]string{"git", "init"}).
			WithExec([]string{"git", "add", "."}).
			WithExec([]string{"git", "config", "--global", "user.email", "dagger@example.com"}).
			WithExec([]string{"git", "config", "--global", "user.name", "Dagger"}).
			WithExec([]string{"git", "commit", "-m", "Initial commit"})
	}

	// Get the latest tag
	output, err := container.WithExec([]string{
		"sh", "-c",
		"git tag -l 'v*' | sort -V | tail -n 1",
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting latest tag: %w", err)
	}

	var major, minor, patch int
	if output == "" {
		// No existing tag, start with v0.1.0
		major, minor, patch = 0, 1, 0
	} else {
		// Parse existing version
		version := strings.TrimPrefix(strings.TrimSpace(output), "v")
		_, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
		if err != nil {
			return "", fmt.Errorf("error parsing version: %w", err)
		}

		// Increment patch version
		patch++
	}

	// Format new version tag
	newTag := fmt.Sprintf("v%d.%d.%d", major, minor, patch)

	// Create new tag
	_, err = container.WithExec([]string{
		"git", "tag", "-a", newTag, "-m", fmt.Sprintf("Release %s", newTag),
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating tag: %w", err)
	}

	if outputVersion {
		return strings.TrimPrefix(newTag, "v"), nil
	}

	return "", nil
}

// GetCurrentVersion returns the latest version tag in the repository
func (m *Versioner) GetCurrentVersion(ctx context.Context, source *dagger.Directory) (string, error) {
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git"})

	// Check if git is already initialized
	gitStatus, err := container.WithExec([]string{"sh", "-c", "[ -d .git ] && echo 'true' || echo 'false'"}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error checking git status: %w", err)
	}

	if strings.TrimSpace(gitStatus) == "false" {
		container = container.
			WithExec([]string{"git", "init"}).
			WithExec([]string{"git", "add", "."}).
			WithExec([]string{"git", "config", "--global", "user.email", "dagger@example.com"}).
			WithExec([]string{"git", "config", "--global", "user.name", "Dagger"}).
			WithExec([]string{"git", "commit", "-m", "Initial commit"})
	}

	output, err := container.WithExec([]string{
		"sh", "-c",
		"git tag -l 'v*' | sort -V | tail -n 1",
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting current version: %w", err)
	}

	return strings.TrimSpace(output), nil
} 