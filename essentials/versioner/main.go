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
		WithExec([]string{"apk", "add", "--no-cache", "git", "openssh"})

	// Configure git
	container = container.WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"})

	// Ensure repository is initialized
	gitStatus, err := container.WithExec([]string{"sh", "-c", "[ -d .git ] && echo 'true' || echo 'false'"}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error checking git status: %w", err)
	}

	if strings.TrimSpace(gitStatus) == "false" {
		container = container.
			WithExec([]string{"git", "init"}).
			WithExec([]string{"git", "remote", "add", "origin", "https://<username>:<token>@github.com/<user>/<repo>.git"}).
			WithExec([]string{"git", "fetch", "origin"}).
			WithExec([]string{"git", "checkout", "-b", "main"}).
			WithExec([]string{"git", "pull", "--rebase", "origin", "main"}).
			WithExec([]string{"git", "add", "."}).
			WithExec([]string{"git", "commit", "-m", "Initial commit"})
	} else {
		// Sync with remote branch
		container = container.WithExec([]string{"git", "fetch", "origin"}).
			WithExec([]string{"git", "pull", "--rebase", "origin", "main"})
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
		major, minor, patch = 0, 1, 0
	} else {
		version := strings.TrimPrefix(strings.TrimSpace(output), "v")
		_, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
		if err != nil {
			return "", fmt.Errorf("error parsing version: %w", err)
		}
		patch++
	}

	// Create new tag
	newTag := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	_, err = container.WithExec([]string{
		"git", "tag", "-a", newTag, "-m", fmt.Sprintf("Release %s", newTag),
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating tag: %w", err)
	}

	// Push branch and ensure sync with remote
	_, err = container.WithExec([]string{"git", "push", "--set-upstream", "origin", "main"}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error pushing branch to remote: %w", err)
	}

	// Push all tags to remote
	_, err = container.WithExec([]string{"git", "push", "origin", "--tags"}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error pushing tags to remote: %w", err)
	}

	if outputVersion {
		return strings.TrimPrefix(newTag, "v"), nil
	}

	return "", nil
}
