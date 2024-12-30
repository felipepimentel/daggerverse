// Package main provides functionality for semantic versioning of Python projects.
// It uses semantic-release to automatically determine the next version number
// based on commit messages and updates the project's version accordingly.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
)

// PythonVersioner handles semantic versioning for Python projects.
// It uses semantic-release to analyze commit messages and determine version bumps.
type PythonVersioner struct{}

// BumpVersion increments the project version using semantic-release.
// The process includes:
// 1. Setting up a Node.js container with required tools
// 2. Installing semantic-release and plugins
// 3. Configuring Git for commit operations
// 4. Creating a package.json with semantic-release configuration
// 5. Running semantic-release to determine and apply version bump
//
// Required environment variables:
// - GITHUB_TOKEN: GitHub token for semantic-release operations
//
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
//
// Returns:
// - string: The new version number
// - error: Any error that occurred during the versioning process
func (m *PythonVersioner) BumpVersion(ctx context.Context, source *dagger.Directory) (string, error) {
	client, err := dagger.Connect(ctx)
	if err != nil {
		return "", fmt.Errorf("error connecting to dagger: %v", err)
	}
	defer client.Close()

	// Setup Node.js container with required tools
	container := client.Container().
		From("node:lts-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithEnvVariable("GIT_AUTHOR_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_AUTHOR_EMAIL", "github-actions[bot]@users.noreply.github.com").
		WithEnvVariable("GIT_COMMITTER_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_COMMITTER_EMAIL", "github-actions[bot]@users.noreply.github.com")

	// Install required packages
	container = container.
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "openssh-client"})

	// Install semantic-release and plugins
	container = container.WithExec([]string{
		"npm", "install", "-g",
		"semantic-release",
		"@semantic-release/commit-analyzer",
		"@semantic-release/release-notes-generator",
		"@semantic-release/changelog",
		"@semantic-release/git",
		"@semantic-release/github",
	})

	// Configure Git
	container = container.
		WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"})

	// Create package.json with semantic-release configuration
	packageJSON := `{
		"name": "@daggerverse/python",
		"version": "0.0.0-development",
		"private": true,
		"release": {
			"branches": ["main"],
			"plugins": [
				"@semantic-release/commit-analyzer",
				"@semantic-release/release-notes-generator",
				"@semantic-release/changelog",
				["@semantic-release/git", {
					"assets": ["CHANGELOG.md"],
					"message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}"
				}],
				["@semantic-release/github", {
					"assets": []
				}]
			]
		}
	}`

	container = container.
		WithExec([]string{"bash", "-c", fmt.Sprintf("echo '%s' > /src/package.json", packageJSON)})

	// Run semantic-release
	output, err := container.
		WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN")).
		WithEnvVariable("GH_TOKEN", os.Getenv("GITHUB_TOKEN")).
		WithExec([]string{
			"npx", "semantic-release",
			"--branches", "main",
			"--ci", "false",
			"--debug",
		}).Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("error running semantic-release: %v", err)
	}

	// Extract version from output
	version := strings.TrimSpace(output)
	if version == "" {
		return "", fmt.Errorf("no version found in semantic-release output")
	}

	return version, nil
} 