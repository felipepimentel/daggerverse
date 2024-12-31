package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/release/internal/dagger"
)

// Release handles the CI/CD pipeline for all modules
type Release struct{}

// New creates a new instance of Release
func New() *Release {
	return &Release{}
}

// Run executes the release pipeline for all modules
func (m *Release) Run(ctx context.Context, source *dagger.Directory) error {
	// Setup Git container
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git"})

	// Get the last commit message
	commitMsg, err := container.WithExec([]string{
		"git", "log", "-1", "--pretty=%B",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error getting commit message: %v", err)
	}

	// Determine commit type
	commitType := m.getCommitType(commitMsg)

	// List of modules to version
	modules := []string{
		"python-poetry",
		"python-pypi",
		"python-pipeline",
	}

	// Process each module
	for _, module := range modules {
		// Get current version
		version, err := m.getCurrentVersion(ctx, container, module)
		if err != nil {
			return fmt.Errorf("error getting current version for %s: %v", module, err)
		}

		// Bump version
		newVersion, err := m.bumpVersion(version, commitType)
		if err != nil {
			return fmt.Errorf("error bumping version for %s: %v", module, err)
		}

		// Create tag
		tagName := fmt.Sprintf("%s/v%s", module, newVersion)
		_, err = container.WithExec([]string{
			"git", "tag", "-a", tagName,
			"-m", fmt.Sprintf("Release %s", tagName),
		}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error creating tag for %s: %v", module, err)
		}

		// Push tag
		_, err = container.WithExec([]string{
			"git", "push", "origin", tagName,
		}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error pushing tag for %s: %v", module, err)
		}
	}

	return nil
}

// getCommitType determines the type of commit from the message
func (m *Release) getCommitType(msg string) string {
	if strings.Contains(msg, "BREAKING CHANGE") {
		return "BREAKING CHANGE"
	}
	if strings.HasPrefix(msg, "feat") {
		return "feat"
	}
	if strings.HasPrefix(msg, "fix") || strings.HasPrefix(msg, "perf") {
		return "fix"
	}
	return "patch"
}

// getCurrentVersion gets the current version of a module
func (m *Release) getCurrentVersion(ctx context.Context, container *dagger.Container, module string) (string, error) {
	output, err := container.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("git tag -l '%s/v*' | sort -V | tail -n 1", module),
	}).Stdout(ctx)
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(output)
	if version == "" {
		return "0.1.0", nil
	}

	version = strings.TrimPrefix(version, fmt.Sprintf("%s/v", module))
	return version, nil
}

// bumpVersion increments the version based on commit type
func (m *Release) bumpVersion(version, commitType string) (string, error) {
	var major, minor, patch int
	_, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		return "", fmt.Errorf("error parsing version: %v", err)
	}

	switch commitType {
	case "BREAKING CHANGE":
		major++
		minor = 0
		patch = 0
	case "feat":
		minor++
		patch = 0
	case "fix", "perf":
		patch++
	default:
		patch++
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
} 