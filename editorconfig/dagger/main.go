// A generated module for Editorconfig functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.
package main

import (
	"dagger/editorconfig/internal/dagger"
)

const (
	defaultImageRepository = "mstruebing/editorconfig-checker:latest"
)

// Editorconfig is a module for checking editorconfig files.
type Editorconfig struct {
	Image string
	Ctr   *dagger.Container
}

// New creates a new instance of the Editorconfig struct
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Editorconfig {
	return &Editorconfig{
		Image: image,
	}
}

// Container returns the underlying Dagger container
func (m *Editorconfig) Container() *dagger.Container {
	if m.Ctr != nil {
		return m.Ctr
	}

	image := m.Image
	if image == "" {
		image = defaultImageRepository
	}

	m.Ctr = dag.Container().From(image)
	return m.Ctr
}

// Check runs the editorconfig-checker command.
func (m *Editorconfig) Check(
	// Source directory
	source *dagger.Directory,
	// excludeDirectoryPattern is an optional argument that specifies a pattern to exclude directories.
	// +default=".git"
	excludeDirectoryPattern string,
) *dagger.Container {

	return m.Container().
		WithMountedDirectory("/tmp", source.WithoutDirectory(excludeDirectoryPattern)).
		WithWorkdir("/tmp").
		WithExec([]string{"ec"})
}
