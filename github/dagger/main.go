// GitHub
//
// Get GitHub command-line interface.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"dagger/github/internal/dagger"
	"fmt"
	"strings"
)

const (
	// Name of GitHub executable binary
	BinaryName string = "gh"
)

// GitHub
type Github struct {
	// +private
	Version string
}

// GitHub constructor
func New(
	// GitHub version to get
	version string,
) *Github {
	github := &Github{
		Version: version,
	}

	return github
}

// Get GitHub executable binary
func (github *Github) Binary(
	ctx context.Context,
	// Platform to get GitHub for
	// +optional
	platform dagger.Platform,
) (*dagger.File, error) {
	if platform == "" {
		defaultPlatform, err := dag.DefaultPlatform(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get platform: %s", err)
		}

		platform = defaultPlatform
	}

	platformElements := strings.Split(string(platform), "/")

	os := map[string]string{
		"linux":   "linux",
		"darwin":  "macOS",
		"windows": "windows",
	}[platformElements[0]]

	arch := platformElements[1]

	downloadURL := "https://github.com/cli/cli/releases/download/v" + github.Version

	archiveBaseName := fmt.Sprintf("gh_%s_%s_%s", github.Version, os, arch)

	archiveName := archiveBaseName + func() string {
		if os == "linux" {
			return ".tar.gz"
		} else {
			return ".zip"
		}
	}()

	checksumsName := fmt.Sprintf("gh_%s_checksums.txt", github.Version)

	archive := dag.HTTP(downloadURL + "/" + archiveName)
	checksums := dag.HTTP(downloadURL + "/" + checksumsName)

	container := dag.Redhat().Container().
		WithMountedFile(archiveName, archive).
		WithMountedFile(checksumsName, checksums).
		WithExec([]string{"sh", "-c", "grep -w " + archiveName + " " + checksumsName + " | sha256sum -c"})

	if os == "linux" {
		container = container.
			WithExec([]string{"tar", "--extract", "--file", archiveName})
	} else {
		container = container.
			With(dag.Redhat().Packages([]string{
				"unzip",
			}).Installed).
			WithExec([]string{"unzip", archiveName})
	}

	binary := container.Directory(archiveBaseName).Directory("bin").File(BinaryName)

	return binary, nil
}

// Get a root filesystem overlay with GitHub
func (github *Github) Overlay(
	ctx context.Context,
	// Platform to get GitHub for
	// +optional
	platform dagger.Platform,
	// Filesystem prefix under which to install GitHub
	// +optional
	prefix string,
) (*dagger.Directory, error) {
	if prefix == "" {
		prefix = "/usr/local"
	}

	binary, err := github.Binary(ctx, platform)

	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub binary: %s", err)
	}

	overlay := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithDirectory("bin", dag.Directory().
				WithFile(BinaryName, binary),
			),
		)

	return overlay, nil
}

// Install GitHub in a container
func (github *Github) Installation(
	ctx context.Context,
	// Container in which to install GitHub
	container *dagger.Container,
) (*dagger.Container, error) {
	platform, err := container.Platform(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get container platform: %s", err)
	}

	overlay, err := github.Overlay(ctx, platform, "")

	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub overlay: %s", err)
	}

	container = container.
		WithDirectory("/", overlay)

	return container, nil
}

// Get a GitHub container from a base container
func (github *Github) Container(
	ctx context.Context,
	// Base container
	container *dagger.Container,
) (*dagger.Container, error) {
	container, err := github.Installation(ctx, container)

	if err != nil {
		return nil, fmt.Errorf("failed to install GitHub: %s", err)
	}

	container = container.
		WithEntrypoint([]string{BinaryName}).
		WithoutDefaultArgs()

	return container, nil
}

// Get a Red Hat Universal Base Image container with GitHub
func (github *Github) RedhatContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Container(dagger.RedhatContainerOpts{Platform: platform}).
		With(dag.Redhat().Packages([]string{
			"git",
		}).Installed)

	return github.Container(ctx, container)
}

// Get a Red Hat Minimal Universal Base Image container with GitHub
func (github *Github) RedhatMinimalContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Minimal().Container(dagger.RedhatMinimalContainerOpts{Platform: platform}).
		With(dag.Redhat().Minimal().Packages([]string{
			"git",
		}).Installed)

	return github.Container(ctx, container)
}

// Get a Red Hat Micro Universal Base Image container with GitHub
//
// Features requiring Git will not work.
func (github *Github) RedhatMicroContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Micro().Container(dagger.RedhatMicroContainerOpts{Platform: platform})

	return github.Container(ctx, container)
}
