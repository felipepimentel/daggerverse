package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Versioner struct{}

func New() *Versioner {
	return &Versioner{}
}

func (m *Versioner) BumpVersion(ctx context.Context, source string) (string, error) {
	// Get the latest tag
	cmd := exec.Command("git", "tag", "-l", "v*")
	cmd.Dir = source
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting latest tag: %w", err)
	}

	// Parse version
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	var newTag string
	if len(tags) == 0 || tags[0] == "" {
		newTag = "v0.1.0"
	} else {
		latestTag := tags[len(tags)-1]
		version := strings.TrimPrefix(latestTag, "v")
		var major, minor, patch int
		_, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
		if err != nil {
			return "", fmt.Errorf("error parsing version: %w", err)
		}
		patch++
		newTag = fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	}

	// Create new tag
	cmd = exec.Command("git", "tag", "-a", newTag, "-m", fmt.Sprintf("Release %s", newTag))
	cmd.Dir = source
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error creating tag: %w", err)
	}

	return newTag, nil
}

func (m *Versioner) GetCurrentVersion(ctx context.Context, source string) (string, error) {
	cmd := exec.Command("git", "tag", "-l", "v*")
	cmd.Dir = source
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting current version: %w", err)
	}

	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tags) == 0 || tags[0] == "" {
		return "", nil
	}
	return tags[len(tags)-1], nil
} 