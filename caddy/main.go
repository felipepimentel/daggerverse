// Caddy
//
// Get a container or service running Caddy HTTP server to serve static content without TLS.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"dagger/caddy/internal/dagger"
)

const (
	// Caddy container image registry
	ImageRegistry string = "docker.io"
	// Caddy container image repository
	ImageRepository string = "caddy"
	// Caddy container image tag
	ImageTag string = "2.8.4"
	// Caddy container image digest
	ImageDigest string = "sha256:69f9a2cd92221b45258a5728b3f08c9b03ba03ed27e0ac791b4343400c3e7385"
)

// Caddy
type Caddy struct {
	// +private
	Directory *dagger.Directory
}

// Caddy constructor
func New(
	// Directory containing static content to serve
	directory *dagger.Directory,
) *Caddy {
	caddy := &Caddy{
		Directory: directory,
	}

	return caddy
}

// Get a Caddy container ready to serve the static content
//
// Static content is mounted under `/usr/share/caddy` and container exposes port 8080.
func (caddy *Caddy) Container(
	// Platform to get container for
	// +optional
	platform dagger.Platform,
	// Mount the directory containing the static content instead of copying it.
	// +optional
	mountDirectory bool,
) *dagger.Container {
	caddyfile := dag.CurrentModule().Source().File("Caddyfile")

	container := dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(ImageRegistry+"/"+ImageRepository+":"+ImageTag+"@"+ImageDigest).
		WithExec([]string{"chown", "65535:65535", "/config/caddy"}).
		WithExec([]string{"chown", "65535:65535", "/data/caddy"}).
		WithEntrypoint([]string{"caddy"}).
		WithDefaultArgs([]string{"run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"}).
		WithFile("/etc/caddy/Caddyfile", caddyfile)

	if mountDirectory {
		container = container.
			WithMountedDirectory("/usr/share/caddy", caddy.Directory)
	} else {
		container = container.
			WithDirectory("/usr/share/caddy", caddy.Directory)
	}

	container = container.
		WithUser("65535").
		WithExposedPort(8080)

	return container
}

// Get a Caddy service serving the static content
//
// See `container()` for details.
func (caddy *Caddy) Server() *dagger.Service {
	return caddy.Container("", true).AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true})
}
