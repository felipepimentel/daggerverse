// main.go
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

// Run executes the release pipeline for all modules
func (m *Release) Run(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	// Detect modules
	modules, err := m.detectModules(ctx, source)
	if err != nil {
		return fmt.Errorf("error detecting modules: %v", err)
	}

	if len(modules) == 0 {
		return fmt.Errorf("no modules found")
	}

	// Setup base container for Git operations
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "git", "openssh"}).
		WithEnvVariable("GIT_AUTHOR_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_AUTHOR_EMAIL", "github-actions[bot]@users.noreply.github.com").
		WithEnvVariable("GIT_COMMITTER_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_COMMITTER_EMAIL", "github-actions[bot]@users.noreply.github.com").
		WithSecretVariable("token", token)

	// Fetch tags and reset to main branch
	container = container.WithExec([]string{
		"sh", "-c", `
		git fetch --tags --force
		git checkout main
		git reset --hard origin/main
		`,
	})

	// Get the last commit message
	commitMsg, err := container.WithExec([]string{
		"git", "log", "-1", "--pretty=%B",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error getting commit message: %v", err)
	}

	// Process each module
	for _, module := range modules {
		moduleContainer := container.WithWorkdir(filepath.Join("/src", module))

		// Get current version
		currentVersion, err := m.getCurrentVersion(ctx, moduleContainer, module)
		if err != nil {
			return fmt.Errorf("error getting current version for %s: %v", module, err)
		}

		// Determine version bump type
		commitType := m.getCommitType(commitMsg)
		newVersion, err := m.bumpVersion(currentVersion, commitType)
		if err != nil {
			return fmt.Errorf("error bumping version for %s: %v", module, err)
		}

		// Create and push tag
		tagName := fmt.Sprintf("%s/v%s", module, newVersion)
		if err := m.createAndPushTag(ctx, moduleContainer, tagName, commitMsg); err != nil {
			return fmt.Errorf("error handling tag for %s: %v", module, err)
		}

		// Publish to Daggerverse
		publishContainer := dag.Container().
			From("alpine:latest").
			WithDirectory("/src", source).
			WithWorkdir(filepath.Join("/src", module)).
			WithExec([]string{"apk", "add", "--no-cache", "dagger"}).
			WithExec([]string{"dagger", "publish"})

		if _, err := publishContainer.Sync(ctx); err != nil {
			return fmt.Errorf("error publishing module %s: %v", module, err)
		}
	}

	return nil
}

// detectModules finds all Dagger modules in the repository
func (m *Release) detectModules(ctx context.Context, source *dagger.Directory) ([]string, error) {
	container := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"find", ".", "-name", "dagger.json", "-exec", "dirname", "{}", ";"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("error finding modules: %v", err)
	}

	var modules []string
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		module := strings.TrimPrefix(line, "./")
		if module != "" {
			modules = append(modules, module)
		}
	}

	return modules, nil
}

// createAndPushTag creates and pushes a Git tag with the commit message
func (m *Release) createAndPushTag(ctx context.Context, container *dagger.Container, tagName, commitMsg string) error {
	// Configure Git to use HTTPS instead of SSH for pushing
	container = container.WithExec([]string{
		"git", "config", "--global", "url.https://github.com/.insteadOf", "git@github.com:",
	})

	// Set up credentials for HTTPS using token directly in the URL
	container = container.WithExec([]string{
		"git", "remote", "set-url", "origin",
		"https://${token}@github.com/felipepimentel/daggerverse.git",
	})

	// Make sure we're on the main branch and up to date
	container = container.WithExec([]string{
		"git", "checkout", "main",
	})

	// Get current commit hash
	hash, err := container.WithExec([]string{
		"git", "rev-parse", "HEAD",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error getting current commit hash: %v", err)
	}
	hash = strings.TrimSpace(hash)

	// Try to delete existing tag locally (ignore errors)
	container.WithExec([]string{
		"git", "tag", "-d", tagName,
	})

	// Create the tag
	if _, err := container.WithExec([]string{
		"git", "tag", "-a", tagName, "-m", commitMsg, hash,
	}).Stdout(ctx); err != nil {
		return fmt.Errorf("error creating tag: %v", err)
	}

	// Pack refs to ensure tag is visible
	if _, err := container.WithExec([]string{
		"git", "pack-refs", "--all",
	}).Stdout(ctx); err != nil {
		return fmt.Errorf("error packing refs: %v", err)
	}

	// Verify the tag exists
	output, err := container.WithExec([]string{
		"git", "tag", "-l", tagName,
	}).Stdout(ctx)
	if err != nil || !strings.Contains(output, tagName) {
		return fmt.Errorf("tag %s was not created successfully", tagName)
	}

	// Push the tag with force to update any existing remote tag
	if _, err := container.WithExec([]string{
		"git", "push", "-f", "origin", tagName,
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
