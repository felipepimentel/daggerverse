// jq
//
// Get jq command-line JSON processor.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"dagger/jq/internal/dagger"
	"fmt"
	"strings"
)

const (
	// Name of jq executable binary
	BinaryName string = "jq"
)

// jq
type Jq struct {
	// +private
	Version string
}

// jq constructor
func New(
	// jq version to get
	version string,
) *Jq {
	jq := &Jq{
		Version: version,
	}

	return jq
}

// Get jq executable binary
func (jq *Jq) Binary(
	ctx context.Context,
	// Platform to get jq for
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
		"darwin":  "macos",
		"windows": "windows",
	}[platformElements[0]]

	arch := map[string]string{
		"amd64":   "amd64",
		"386":     "i386",
		"arm":     "armhf",
		"arm64":   "arm64",
		"riscv64": "riscv64",
	}[platformElements[1]]

	downloadURL := "https://github.com/jqlang/jq/releases/download/jq-" + jq.Version

	binaryName := fmt.Sprintf("jq-%s-%s", os, arch)

	if os == "windows" {
		binaryName += ".exe"
	}

	checksumsName := "sha256sum.txt"

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

// Get a root filesystem overlay with jq
func (jq *Jq) Overlay(
	ctx context.Context,
	// Platform to get jq for
	// +optional
	platform dagger.Platform,
	// Filesystem prefix under which to install jq
	// +optional
	prefix string,
) (*dagger.Directory, error) {
	if prefix == "" {
		prefix = "/usr/local"
	}

	binary, err := jq.Binary(ctx, platform)

	if err != nil {
		return nil, fmt.Errorf("failed to get jq binary: %s", err)
	}

	overlay := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithDirectory("bin", dag.Directory().
				WithFile(BinaryName, binary),
			),
		)

	return overlay, nil
}

// Install jq in a container
func (jq *Jq) Installation(
	ctx context.Context,
	// Container in which to install jq
	container *dagger.Container,
) (*dagger.Container, error) {
	platform, err := container.Platform(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get container platform: %s", err)
	}

	overlay, err := jq.Overlay(ctx, platform, "")

	if err != nil {
		return nil, fmt.Errorf("failed to get jq overlay: %s", err)
	}

	container = container.
		WithDirectory("/", overlay)

	return container, nil
}

// Get a jq container from a base container
func (jq *Jq) Container(
	ctx context.Context,
	// Base container
	container *dagger.Container,
) (*dagger.Container, error) {
	container, err := jq.Installation(ctx, container)

	if err != nil {
		return nil, fmt.Errorf("failed to install jq: %s", err)
	}

	container = container.
		WithEntrypoint([]string{BinaryName}).
		WithoutDefaultArgs()

	return container, nil
}

// Get a Red Hat Universal Base Image container with jq
func (jq *Jq) RedhatContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Container(dagger.RedhatContainerOpts{Platform: platform})

	return jq.Container(ctx, container)
}

// Get a Red Hat Minimal Universal Base Image container with jq
func (jq *Jq) RedhatMinimalContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Minimal().Container(dagger.RedhatMinimalContainerOpts{Platform: platform})

	return jq.Container(ctx, container)
}

// Get a Red Hat Micro Universal Base Image container with jq
func (jq *Jq) RedhatMicroContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Micro().Container(dagger.RedhatMicroContainerOpts{Platform: platform})

	return jq.Container(ctx, container)
}
