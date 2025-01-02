package main

import (
	"context"
	"fmt"
	"strings"

	"versioner/internal/dagger"
)

type Versioner struct {
	dag *dagger.Client
}

func New(c *dagger.Client) *Versioner {
	return &Versioner{dag: c}
}

// BumpVersion increments the version of a module based on the commit type
func (m *Versioner) BumpVersion(ctx context.Context, source *dagger.Directory, module, commitType string) (string, error) {
	container := m.dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git"})

	// Get the latest tag for the module
	output, err := container.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("git tag -l '%s/v*' | sort -V | tail -n 1", module),
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting latest tag: %v", err)
	}

	var major, minor, patch int
	if output == "" {
		// No existing tag, start with v0.1.0
		major, minor, patch = 0, 1, 0
	} else {
		// Parse existing version
		version := strings.TrimPrefix(strings.TrimSpace(output), fmt.Sprintf("%s/v", module))
		_, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
		if err != nil {
			return "", fmt.Errorf("error parsing version: %v", err)
		}

		// Increment version based on commit type
		switch commitType {
		case "feat":
			minor++
			patch = 0
		case "fix", "perf":
			patch++
		case "BREAKING CHANGE":
			major++
			minor = 0
			patch = 0
		default:
			patch++
		}
	}

	// Format new version tag
	newTag := fmt.Sprintf("%s/v%d.%d.%d", module, major, minor, patch)

	// Create and push new tag
	_, err = container.WithExec([]string{
		"git", "tag", "-a", newTag, "-m", fmt.Sprintf("Release %s", newTag),
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating tag: %v", err)
	}

	return newTag, nil
}

// GetCurrentVersion gets the current version of a module
func (m *Versioner) GetCurrentVersion(ctx context.Context, source *dagger.Directory, module string) (string, error) {
	container := m.dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git"})

	output, err := container.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("git tag -l '%s/v*' | sort -V | tail -n 1", module),
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting current version: %v", err)
	}

	return strings.TrimSpace(output), nil
} 