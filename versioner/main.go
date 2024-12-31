// Package main provides functionality for semantic versioning of projects.
// It uses semantic-release for version management.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger/versioner/internal/dagger"
)

// Versioner handles semantic versioning for projects.
// It uses semantic-release to analyze commit messages and determine version bumps.
type Versioner struct{}

// New creates a new instance of Versioner.
func New(ctx context.Context) (*Versioner, error) {
	return &Versioner{}, nil
}

// BumpVersion increments the project version using semantic-release.
// The process includes:
// 1. Setting up a container with semantic-release
// 2. Running semantic-release to determine the next version
// 3. Updating the package version
//
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the project
//
// Returns:
// - string: The new version number
// - error: Any error that occurred during the versioning process
func (m *Versioner) BumpVersion(ctx context.Context, source *dagger.Directory) (string, error) {
	client := dagger.Connect()

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
		"name": "@daggerverse/versioner",
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