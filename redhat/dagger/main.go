// Red Hat
//
// Get and customize containers based on Red Hat Universal Base Images.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"dagger/redhat/internal/dagger"
	"strings"
)

const (
	// Red Hat Universal Base Image container registry
	ImageRegistry string = "registry.access.redhat.com"

	// Red Hat Universal Base Image container repository
	ImageRepository string = "ubi9"
	// Red Hat Universal Base Image container tag
	ImageTag string = "9.5-1732804088"
	// Red Hat Universal Base Image container digest
	ImageDigest string = "sha256:b632d0cc6263372a90e9097dcac0a369e456b144a66026b9eac029a22f0f6e07"

	// Red Hat Minimal Universal Base Image container repository
	MinimalImageRepository string = "ubi9-minimal"
	// Red Hat Minimal Universal Base Image container tag
	MinimalImageTag string = "9.5-1733767867"
	// Red Hat Minimal Universal Base Image container digest
	MinimalImageDigest string = "sha256:f598528219a1be07cf520fbe82a2d2434dc9841e1f0a878382c8a13bf42cb486"

	// Red Hat Micro Universal Base Image container repository
	MicroImageRepository string = "ubi9-micro"
	// Red Hat Micro Universal Base Image container tag
	MicroImageTag string = "9.5-1733767087"
	// Red Hat Micro Universal Base Image container digest
	MicroImageDigest string = "sha256:3313e52bb1aad4017a0c35f9f2ae35cf8526eeeb83f6ecbec449ba9c5cb9cb07"
)

// Red Hat Universal Base Image
type Redhat struct{}

// Red Hat Universal Base Image constructor
func New() *Redhat {
	return &Redhat{}
}

// Get a Red Hat Universal Base Image container
func (*Redhat) Container(
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) *dagger.Container {
	container := dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(ImageRegistry + "/" + ImageRepository + ":" + ImageTag + "@" + ImageDigest).
		WithWorkdir("/home")

	return container
}

// Red Hat Universal Base Image module
type RedhatModule struct {
	// +private
	Name string
}

// Red Hat Universal Base Image module constructor
func (*Redhat) Module(
	// Module name
	name string,
) *RedhatModule {
	module := &RedhatModule{
		Name: name,
	}

	return module
}

// Enable a module in a Red Hat Universal Base Image container
func (module *RedhatModule) Enabled(
	// Container in which to enable the module
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "dnf module enable --assumeyes " + module.Name + " && dnf clean all"})
}

// Disable a module in a Red Hat Universal Base Image container
func (module *RedhatModule) Disabled(
	// Container in which to disable the module
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "dnf module disable --assumeyes " + module.Name + " && dnf clean all"})
}

// Red Hat Universal Base Image packages
type RedhatPackages struct {
	// +private
	Names []string
}

// Red Hat Universal Base Image packages constructor
func (*Redhat) Packages(
	// Packages name
	names []string,
) *RedhatPackages {
	packages := &RedhatPackages{
		Names: names,
	}

	return packages
}

// Install packages in a Red Hat Universal Base Image container
func (packages *RedhatPackages) Installed(
	// Container in which to install the packages
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "dnf install --nodocs --setopt install_weak_deps=0 --assumeyes " + strings.Join(packages.Names, " ") + " && dnf clean all"})
}

// Remove packages in a Red Hat Universal Base Image container
func (packages *RedhatPackages) Removed(
	// Container in which to remove the packages
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "dnf remove --assumeyes " + strings.Join(packages.Names, " ") + " && dnf clean all"})
}

// Get Red Hat Universal Base Image CA certificates
func (redhat *Redhat) CaCertificates() *dagger.Directory {
	const installroot string = "/tmp/rootfs"

	caCertificates := redhat.Container("").
		WithExec([]string{"sh", "-c", "mkdir " + installroot + " && dnf --installroot " + installroot + " install --nodocs --setopt install_weak_deps=0 --assumeyes ca-certificates && dnf --installroot " + installroot + " clean all"}).
		Directory(installroot + "/etc/pki/ca-trust")

	return caCertificates
}

// Red Hat Minimal Universal Base Image
type RedhatMinimal struct{}

// Red Hat Minimal Universal Base Image constructor
func (*Redhat) Minimal() *RedhatMinimal {
	return &RedhatMinimal{}
}

// Get a Red Hat Minimal Universal Base Image container
func (*RedhatMinimal) Container(
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) *dagger.Container {
	container := dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(ImageRegistry + "/" + MinimalImageRepository + ":" + MinimalImageTag + "@" + MinimalImageDigest).
		WithWorkdir("/home")

	return container
}

// Red Hat Minimal Universal Base Image module
type RedhatMinimalModule struct {
	// +private
	Name string
}

// Red Hat Minimal Universal Base Image module constructor
func (*RedhatMinimal) Module(
	// Module name
	name string,
) *RedhatMinimalModule {
	module := &RedhatMinimalModule{
		Name: name,
	}

	return module
}

// Enable a module in a Red Hat Minimal Universal Base Image container
func (module *RedhatMinimalModule) Enabled(
	// Container in which to enable the module
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "microdnf module enable --assumeyes " + module.Name + " && microdnf clean all"})
}

// Disable a module in a Red Hat Minimal Universal Base Image container
func (module *RedhatMinimalModule) Disabled(
	// Container in which to disable the module
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "microdnf module disable --assumeyes " + module.Name + " && microdnf clean all"})
}

// Red Hat Minimal Universal Base Image packages
type RedhatMinimalPackages struct {
	// +private
	Names []string
}

// Red Hat Minimal Universal Base Image packages constructor
func (*RedhatMinimal) Packages(
	// Packages name
	names []string,
) *RedhatMinimalPackages {
	packages := &RedhatMinimalPackages{
		Names: names,
	}

	return packages
}

// Install packages in a Red Hat Minimal Universal Base Image container
func (packages *RedhatMinimalPackages) Installed(
	// Container in which to install the packages
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "microdnf install --nodocs --setopt install_weak_deps=0 --assumeyes " + strings.Join(packages.Names, " ") + " && microdnf clean all"})
}

// Remove packages in a Red Hat Minimal Universal Base Image container
func (packages *RedhatMinimalPackages) Removed(
	// Container in which to remove the packages
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "microdnf remove --assumeyes " + strings.Join(packages.Names, " ") + " && microdnf clean all"})
}

// Red Hat Micro Universal Base Image
type RedhatMicro struct{}

// Red Hat Micro Universal Base Image constructor
func (*Redhat) Micro() *RedhatMicro {
	return &RedhatMicro{}
}

// Get a Red Hat Micro Universal Base Image container
func (*RedhatMicro) Container(
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) *dagger.Container {
	container := dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(ImageRegistry + "/" + MicroImageRepository + ":" + MicroImageTag + "@" + MicroImageDigest).
		WithWorkdir("/home")

	return container
}
