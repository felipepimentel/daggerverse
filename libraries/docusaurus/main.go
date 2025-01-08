package main

import (
	"fmt"

	"github.com/felipepimentel/daggerverse/libraries/docusaurus/internal/dagger"
)

func New(
	// The source directory of your docusaurus site, this can be a local directory or a remote git repo
	src *dagger.Directory,
	// Optional working directory if you need to execute docusaurus commands outside of your root
	// +optional
	// +default="/src"
	dir string,
	// Optional flag to disable cache
	// +optional
	// +default=false
	disableCache bool,
	// Optional cache volume name; this is useful if you work with multiple projects
	// or have node_dependencies that are rapidly changing to avoid issues with
	// npm having failures.
	// +optional
	// +default="node-docusaurus-docs"
	cacheVolumeName string,
	// Optional flag to use yarn instead of npm
	// +optional
	// +default=false
	yarn bool,
) *Docusaurus {
	return &Docusaurus{
		Src:             src,
		Dir:             dir,
		DisableCache:    disableCache,
		CacheVolumeName: cacheVolumeName,
		Yarn:            yarn,
	}
}

// Docusaurus
type Docusaurus struct {
	Src             *dagger.Directory
	Dir             string
	DisableCache    bool
	CacheVolumeName string
	Yarn            bool
}

// Return base container for running docusaurus with docs mounted and docusaurus
// dependencies installed.
func (m *Docusaurus) Base() *dagger.Container {
	ctr := dag.Container().
		From("node:lts-alpine").
		WithoutEntrypoint().
		WithMountedDirectory("/src", m.Src).
		WithWorkdir(m.Dir)

	if !m.DisableCache {
		ctr = ctr.
			WithMountedCache(
				fmt.Sprintf("%s/node_modules", m.Dir),
				dag.CacheVolume(m.CacheVolumeName),
			).
			WithMountedCache(
				fmt.Sprintf("%s/build", m.Dir),
				dag.CacheVolume(m.CacheVolumeName+"-build"),
			).
			WithMountedCache(
				"/root/.npm",
				dag.CacheVolume("node-docusaurus-root"),
			).
			WithMountedCache(
				"/root/.yarn",
				dag.CacheVolume("node-docusaurus-root-yarn"),
			)
	}

	return ctr.
		WithExposedPort(3000).
		WithExec([]string{m.packageManager(), "install"})
}

// Build production docs
func (m *Docusaurus) Build() *dagger.Directory {
	return m.Base().
		WithExec([]string{m.packageManager(), "run", "build"}).
		// copying build to a temp directory because
		// cache volumes cannot be exported. This is totally
		// worth vs the time it takes to build on a cold cache
		WithMountedDirectory("/tmp/build", dag.Directory()).
		WithExec([]string{"cp", "-r", "build/.", "/tmp/build"}).
		Directory("/tmp/build")
}

// Serve production docs locally as a service
func (m *Docusaurus) Serve() *dagger.Service {
	return m.Base().
		WithExec([]string{m.packageManager(), "run", "build"}).
		WithExec([]string{m.packageManager(), "run", "serve", "--build"}).
		AsService()
}

// Build and serve development docs as a service
func (m *Docusaurus) ServeDev() *dagger.Service {
	return m.Base().
		WithExec([]string{m.packageManager(), "start", "--", "--host", "0.0.0.0"}).
		AsService()
}

func (m *Docusaurus) packageManager() string {
	if m.Yarn {
		return "yarn"
	}
	return "npm"
}
