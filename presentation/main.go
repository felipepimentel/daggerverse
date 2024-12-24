// Camptocamp branded presentation
//
// Build a presentation directory, container or service using reveal.js with Camptocamp theme.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"dagger/presentation/internal/dagger"
	"encoding/json"
	"fmt"
)

// Presentation
type Presentation struct{}

// Get a directory containing a newly initialized presentation
func (*Presentation) Init() *dagger.Directory {
	return dag.CurrentModule().Source().Directory("template")
}

// Presentation builder
type PresentationBuilder struct {
	// Get a container ready to build the presentation
	*dagger.Container
}

// Presentation builder constructor
func (*Presentation) Builder(
	ctx context.Context,
	// Directory containing presentation to build
	directory *dagger.Directory,
	// npm configuration file (used to pass GitHub registry credentials)
	npmrc *dagger.Secret,
) (*PresentationBuilder, error) {
	const packageJsonFilename string = "package.json"

	packageJsonString, err := directory.File(packageJsonFilename).Contents(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to read %q file: %w", packageJsonFilename, err)
	}

	var configuration struct{}

	err = json.Unmarshal([]byte(packageJsonString), &configuration)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %q file: %w", packageJsonFilename, err)
	}

	builder := &PresentationBuilder{}

	kroki := dag.Kroki()

	builder.Container = dag.Redhat().Minimal().Container().
		With(dag.Nodejs(dagger.NodejsOpts{
			Npmrc: npmrc,
		}).RedhatMinimalInstallation).
		WithServiceBinding("kroki", kroki.Server()).
		WithMountedDirectory(".", directory).
		WithExec([]string{"npm", "clean-install"}).
		WithEntrypoint([]string{"npm", "run", "all"}).
		WithoutDefaultArgs()

	return builder, nil
}

// Presentation build result
type PresentationBuildResult struct {
	// Get a directory containing the presentation build result
	*dagger.Directory
}

// Build the presentation
func (builder *PresentationBuilder) Build() *PresentationBuildResult {
	build := &PresentationBuildResult{
		Directory: builder.WithExec(nil, dagger.ContainerWithExecOpts{UseEntrypoint: true}).Directory("dist"),
	}

	return build
}

// Get a container ready to serve the presentation
//
// Container exposes port 8080.
func (build *PresentationBuildResult) Container(
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) *dagger.Container {
	return dag.Caddy(build.Directory).Container(dagger.CaddyContainerOpts{Platform: platform})
}

// Get a service serving the presentation
//
// See `container()` for details.
func (build *PresentationBuildResult) Server() *dagger.Service {
	return dag.Caddy(build.Directory).Server()
}
