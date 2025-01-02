package main

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

type Versioner struct {
	dag *dagger.Client
}

func New(c *dagger.Client) *Versioner {
	return &Versioner{dag: c}
}

func (m *Versioner) BumpVersion(ctx context.Context, source string) (string, error) {
	hostDir := m.dag.Host().Directory(source)
	container := m.dag.Container().
		From("alpine:latest").
		WithDirectory("/src", hostDir).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithExec([]string{"git", "config", "--global", "user.email", "dagger@example.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "Dagger"})

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
	}).Sync(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating tag: %w", err)
	}

	return newTag, nil
}

func (m *Versioner) GetCurrentVersion(ctx context.Context, source string) (string, error) {
	hostDir := m.dag.Host().Directory(source)
	container := m.dag.Container().
		From("alpine:latest").
		WithDirectory("/src", hostDir).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git"})

	output, err := container.WithExec([]string{
		"sh", "-c",
		"git tag -l 'v*' | sort -V | tail -n 1",
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting current version: %w", err)
	}

	return strings.TrimSpace(output), nil
} 