package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipepimentel/daggerverse/essentials/versioner/internal/dagger"
)

// Versioner implements version management for repositories
type Versioner struct{}

// New creates a new Versioner instance
func New() *Versioner {
	return &Versioner{}
}

// BumpVersion creates a new version tag based on the latest tag and commit type
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

	// Get all tags sorted by version
	output, err := container.WithExec([]string{
		"sh", "-c",
		"git tag -l 'v*' | sort -V",
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting tags: %w", err)
	}

	// Get the latest commit message and hash
	commitMsg, err := container.WithExec([]string{
		"sh", "-c",
		"git log -1 --pretty=%B%n%H",
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting commit message: %w", err)
	}

	// Split commit message and hash
	parts := strings.Split(strings.TrimSpace(commitMsg), "\n")
	commitHash := parts[len(parts)-1]
	commitMsg = strings.Join(parts[:len(parts)-1], "\n")

	// Find the highest version among all tags
	var major, minor, patch int
	tags := strings.Split(strings.TrimSpace(output), "\n")
	if len(tags) == 0 || (len(tags) == 1 && tags[0] == "") {
		major, minor, patch = 0, 1, 0
	} else {
		for _, tag := range tags {
			var m, n, p int
			version := strings.TrimPrefix(strings.TrimSpace(tag), "v")
			_, err := fmt.Sscanf(version, "%d.%d.%d", &m, &n, &p)
			if err == nil {
				// Update highest version found
				if m > major || (m == major && n > minor) || (m == major && n == minor && p > patch) {
					major, minor, patch = m, n, p
				}
			}
		}

		// Check if the latest tag points to the current commit
		for _, tag := range tags {
			tagHash, err := container.WithExec([]string{
				"git", "rev-parse", tag,
			}).Stdout(ctx)
			if err == nil && strings.TrimSpace(tagHash) == commitHash {
				// This commit already has a tag, increment patch version
				patch++
				break
			}
		}

		// Determine version bump based on commit message
		commitMsg = strings.ToLower(strings.TrimSpace(commitMsg))
		if strings.Contains(commitMsg, "breaking change") || strings.Contains(commitMsg, "!:") {
			major++
			minor = 0
			patch = 0
		} else if strings.HasPrefix(commitMsg, "feat:") || strings.HasPrefix(commitMsg, "feat(") {
			minor++
			patch = 0
		} else {
			// For any other commit type (including non-semantic commits), increment patch
			patch++
		}
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
