package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"dagger/release/internal/dagger"
)

// Release handles the CI/CD pipeline for all modules
type Release struct{}

// New creates a new instance of Release
func New() *Release {
	return &Release{}
}

// DetectModules finds all Dagger modules in the repository
func (m *Release) DetectModules(ctx context.Context, source *dagger.Directory) ([]string, error) {
	// Use alpine container to find modules
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"find", ".", "-name", "dagger.json", "-exec", "dirname", "{}", ";"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("error finding modules: %v", err)
	}

	// Process the output to get module paths
	var modules []string
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		// Remove "./" prefix if present
		module := strings.TrimPrefix(line, "./")
		if module != "" {
			modules = append(modules, module)
		}
	}

	return modules, nil
}

// ReleaseModule handles the release process for a single module
func (m *Release) ReleaseModule(ctx context.Context, source *dagger.Directory, modulePath string, token *dagger.Secret) error {
	// Setup Git container
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git", "openssh"}).
		WithEnvVariable("GIT_AUTHOR_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_AUTHOR_EMAIL", "github-actions[bot]@users.noreply.github.com").
		WithEnvVariable("GIT_COMMITTER_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_COMMITTER_EMAIL", "github-actions[bot]@users.noreply.github.com")

	// Configure Git with token
	container = container.WithExec([]string{
		"git", "config", "--global", "url.https://oauth2:token@github.com/.insteadOf", "https://github.com/",
	}).WithSecretVariable("token", token)

	// Get the last commit message
	commitMsg, err := container.WithExec([]string{
		"git", "log", "-1", "--pretty=%B",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error getting commit message: %v", err)
	}

	// Get current version
	version, err := m.getCurrentVersion(ctx, container, modulePath)
	if err != nil {
		return fmt.Errorf("error getting current version: %v", err)
	}

	// Determine version bump type
	commitType := m.getCommitType(commitMsg)
	newVersion, err := m.bumpVersion(version, commitType)
	if err != nil {
		return fmt.Errorf("error bumping version: %v", err)
	}

	// Create and push tag
	tagName := fmt.Sprintf("%s/v%s", modulePath, newVersion)
	if err := m.createAndPushTag(ctx, container, tagName, commitMsg); err != nil {
		return fmt.Errorf("error handling tag: %v", err)
	}

	// Publish to Daggerverse
	publishContainer := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir(filepath.Join("/src", modulePath)).
		WithExec([]string{"apk", "add", "--no-cache", "dagger"}).
		WithExec([]string{"dagger", "publish"})

	if _, err := publishContainer.Sync(ctx); err != nil {
		return fmt.Errorf("error publishing module: %v", err)
	}

	return nil
}

// createAndPushTag creates and pushes a Git tag with the commit message
func (m *Release) createAndPushTag(ctx context.Context, container *dagger.Container, tagName, commitMsg string) error {
	// Create tag with commit message
	if _, err := container.WithExec([]string{
		"git", "tag", "-a", tagName,
		"-m", commitMsg,
	}).Stdout(ctx); err != nil {
		return fmt.Errorf("error creating tag: %v", err)
	}

	// Push tag
	if _, err := container.WithExec([]string{
		"git", "push", "origin", tagName,
	}).Stdout(ctx); err != nil {
		return fmt.Errorf("error pushing tag: %v", err)
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