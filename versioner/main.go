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
	container = container.
		WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"}).
		WithExec([]string{"git", "config", "--global", "safe.directory", "*"}).
		WithExec([]string{"git", "config", "--global", "init.defaultBranch", "main"}).
		WithExec([]string{"git", "config", "--global", "core.sshCommand", "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"}).
		WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "/src"}).
		WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "/work"}).
		WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "."}).
		WithExec([]string{"git", "config", "--global", "core.fileMode", "false"}).
		WithExec([]string{"git", "config", "--global", "core.autocrlf", "false"}).
		WithExec([]string{"git", "config", "--global", "core.longpaths", "true"}).
		WithExec([]string{"git", "config", "--global", "http.postBuffer", "524288000"}).
		WithExec([]string{"git", "config", "--global", "http.sslVerify", "false"}).
		WithExec([]string{"git", "config", "--global", "http.followRedirects", "true"}).
		WithExec([]string{"git", "config", "--global", "pack.windowMemory", "100m"}).
		WithExec([]string{"git", "config", "--global", "pack.packSizeLimit", "100m"}).
		WithExec([]string{"git", "config", "--global", "pack.threads", "1"}).
		WithExec([]string{"git", "config", "--global", "pack.deltaCacheSize", "100m"}).
		WithExec([]string{"git", "config", "--global", "core.compression", "0"}).
		WithExec([]string{"git", "config", "--global", "core.bigFileThreshold", "50m"}).
		WithExec([]string{"git", "config", "--global", "core.preloadIndex", "true"}).
		WithExec([]string{"git", "config", "--global", "core.fscache", "true"}).
		WithExec([]string{"git", "config", "--global", "gc.auto", "0"}).
		WithExec([]string{"git", "config", "--global", "gc.autoDetach", "false"}).
		WithExec([]string{"git", "config", "--global", "gc.pruneExpire", "now"}).
		WithExec([]string{"git", "config", "--global", "fetch.parallel", "1"}).
		WithExec([]string{"git", "config", "--global", "http.lowSpeedLimit", "1000"}).
		WithExec([]string{"git", "config", "--global", "http.lowSpeedTime", "60"}).
		WithExec([]string{"git", "config", "--global", "http.maxRequests", "1"}).
		WithExec([]string{"git", "config", "--global", "http.minSessions", "1"}).
		WithExec([]string{"git", "config", "--global", "protocol.version", "2"}).
		WithExec([]string{"git", "config", "--global", "transfer.fsckObjects", "false"}).
		WithExec([]string{"git", "config", "--global", "advice.detachedHead", "false"}).
		WithExec([]string{"git", "config", "--global", "advice.pushUpdateRejected", "false"}).
		WithExec([]string{"git", "config", "--global", "http.retryCount", "3"}).
		WithExec([]string{"git", "config", "--global", "http.retryDelay", "2"}).
		WithExec([]string{"git", "config", "--global", "http.maxRequestBuffer", "100M"}).
		WithExec([]string{"git", "config", "--global", "http.version", "HTTP/1.1"})

	// Check if git is already initialized
	gitStatus, err := container.WithExec([]string{"sh", "-c", "[ -d .git ] && echo 'true' || echo 'false'"}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error checking git status: %w", err)
	}

	if strings.TrimSpace(gitStatus) == "false" {
		container = container.
			WithExec([]string{"git", "init"}).
			WithExec([]string{"git", "add", "."}).
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
		WithExec([]string{"apk", "add", "--no-cache", "git", "openssh"})

	// Configure git
	container = container.
		WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"}).
		WithExec([]string{"git", "config", "--global", "safe.directory", "*"}).
		WithExec([]string{"git", "config", "--global", "init.defaultBranch", "main"}).
		WithExec([]string{"git", "config", "--global", "core.sshCommand", "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"}).
		WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "/src"}).
		WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "/work"}).
		WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "."}).
		WithExec([]string{"git", "config", "--global", "core.fileMode", "false"}).
		WithExec([]string{"git", "config", "--global", "core.autocrlf", "false"}).
		WithExec([]string{"git", "config", "--global", "core.longpaths", "true"}).
		WithExec([]string{"git", "config", "--global", "http.postBuffer", "524288000"}).
		WithExec([]string{"git", "config", "--global", "http.sslVerify", "false"}).
		WithExec([]string{"git", "config", "--global", "http.followRedirects", "true"}).
		WithExec([]string{"git", "config", "--global", "pack.windowMemory", "100m"}).
		WithExec([]string{"git", "config", "--global", "pack.packSizeLimit", "100m"}).
		WithExec([]string{"git", "config", "--global", "pack.threads", "1"}).
		WithExec([]string{"git", "config", "--global", "pack.deltaCacheSize", "100m"}).
		WithExec([]string{"git", "config", "--global", "core.compression", "0"}).
		WithExec([]string{"git", "config", "--global", "core.bigFileThreshold", "50m"}).
		WithExec([]string{"git", "config", "--global", "core.preloadIndex", "true"}).
		WithExec([]string{"git", "config", "--global", "core.fscache", "true"}).
		WithExec([]string{"git", "config", "--global", "gc.auto", "0"}).
		WithExec([]string{"git", "config", "--global", "gc.autoDetach", "false"}).
		WithExec([]string{"git", "config", "--global", "gc.pruneExpire", "now"}).
		WithExec([]string{"git", "config", "--global", "fetch.parallel", "1"}).
		WithExec([]string{"git", "config", "--global", "http.lowSpeedLimit", "1000"}).
		WithExec([]string{"git", "config", "--global", "http.lowSpeedTime", "60"}).
		WithExec([]string{"git", "config", "--global", "http.maxRequests", "1"}).
		WithExec([]string{"git", "config", "--global", "http.minSessions", "1"}).
		WithExec([]string{"git", "config", "--global", "protocol.version", "2"}).
		WithExec([]string{"git", "config", "--global", "transfer.fsckObjects", "false"}).
		WithExec([]string{"git", "config", "--global", "advice.detachedHead", "false"}).
		WithExec([]string{"git", "config", "--global", "advice.pushUpdateRejected", "false"}).
		WithExec([]string{"git", "config", "--global", "http.retryCount", "3"}).
		WithExec([]string{"git", "config", "--global", "http.retryDelay", "2"}).
		WithExec([]string{"git", "config", "--global", "http.maxRequestBuffer", "100M"}).
		WithExec([]string{"git", "config", "--global", "http.version", "HTTP/1.1"})

	// Check if git is already initialized
	gitStatus, err := container.WithExec([]string{"sh", "-c", "[ -d .git ] && echo 'true' || echo 'false'"}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("error checking git status: %w", err)
	}

	if strings.TrimSpace(gitStatus) == "false" {
		container = container.
			WithExec([]string{"git", "init"}).
			WithExec([]string{"git", "add", "."}).
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