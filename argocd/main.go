// Argo CD
//
// Get Argo CD command-line interface.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"dagger/argocd/internal/dagger"
	"fmt"
	"strings"
)

const (
	// Name of Argo CD executable binary
	BinaryName string = "argocd"
)

// Argo CD
type Argocd struct {
	// +private
	Version string
}

// Argo CD constructor
func New(
	// Argo CD version to get
	version string,
) *Argocd {
	argocd := &Argocd{
		Version: version,
	}

	return argocd
}

// Get Argo CD executable binary
func (argocd *Argocd) Binary(
	ctx context.Context,
	// Platform to get Argo CD for
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

	os := platformElements[0]
	arch := platformElements[1]

	downloadURL := "https://github.com/argoproj/argo-cd/releases/download/v" + argocd.Version

	binaryName := fmt.Sprintf("argocd-%s-%s", os, arch)

	if os == "windows" {
		binaryName += ".exe"
	}

	checksumsName := "cli_checksums.txt"

	binary := dag.HTTP(downloadURL + "/" + binaryName)
	checksums := dag.HTTP(downloadURL + "/" + checksumsName)

	container := dag.Redhat().Container().
		WithMountedFile(binaryName, binary).
		WithMountedFile(checksumsName, checksums).
		WithExec([]string{"sh", "-c", "grep -w " + binaryName + " " + checksumsName + " | sha256sum -c"}).
		WithExec([]string{"chmod", "a+x", binaryName})

	binary = container.File(binaryName)

	return binary, nil
}

// Get a root filesystem overlay with Argo CD
func (argocd *Argocd) Overlay(
	ctx context.Context,
	// Platform to get Argo CD for
	// +optional
	platform dagger.Platform,
	// Filesystem prefix under which to install Argo CD
	// +optional
	prefix string,
) (*dagger.Directory, error) {
	if prefix == "" {
		prefix = "/usr/local"
	}

	binary, err := argocd.Binary(ctx, platform)

	if err != nil {
		return nil, fmt.Errorf("failed to get Argo CD binary: %s", err)
	}

	overlay := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithDirectory("bin", dag.Directory().
				WithFile(BinaryName, binary),
			),
		)

	return overlay, nil
}

// Install Argo CD in a container
func (argocd *Argocd) Installation(
	ctx context.Context,
	// Container in which to install Argo CD
	container *dagger.Container,
) (*dagger.Container, error) {
	platform, err := container.Platform(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get container platform: %s", err)
	}

	overlay, err := argocd.Overlay(ctx, platform, "")

	if err != nil {
		return nil, fmt.Errorf("failed to get Argo CD overlay: %s", err)
	}

	container = container.
		WithDirectory("/", overlay)

	return container, nil
}

// Get a Argo CD container from a base container
func (argocd *Argocd) Container(
	ctx context.Context,
	// Base container
	container *dagger.Container,
) (*dagger.Container, error) {
	container, err := argocd.Installation(ctx, container)

	if err != nil {
		return nil, fmt.Errorf("failed to install Argo CD: %s", err)
	}

	container = container.
		WithEntrypoint([]string{BinaryName}).
		WithoutDefaultArgs()

	return container, nil
}

// Get a Red Hat Universal Base Image container with Argo CD
func (argocd *Argocd) RedhatContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Container(dagger.RedhatContainerOpts{Platform: platform})

	return argocd.Container(ctx, container)
}

// Get a Red Hat Minimal Universal Base Image container with Argo CD
func (argocd *Argocd) RedhatMinimalContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Minimal().Container(dagger.RedhatMinimalContainerOpts{Platform: platform})

	return argocd.Container(ctx, container)
}

// Get a Red Hat Micro Universal Base Image container with Argo CD
func (argocd *Argocd) RedhatMicroContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Micro().Container(dagger.RedhatMicroContainerOpts{Platform: platform})

	return argocd.Container(ctx, container)
}
