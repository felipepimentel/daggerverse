// A generated module for Ruff functions
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
	"dagger/ruff/internal/dagger"
)

type Ruff struct{}

// CheckWithConfig runs the ruff check command with a configuration file.
func (m *Ruff) CheckWithConfig(
	// source is an optional argument that specifies a directory.
	source *dagger.Directory,
	// file is an optional argument that specifies ruff configuration file.
	file *dagger.File,
) *dagger.Container {
	return base().
		WithMountedDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithFile("/.config/.ruff.toml", file).
		WithExec([]string{"/ruff", "check", "--config", "/.config/.ruff.toml"})
}

// Check runs the ruff check command.
func (m *Ruff) Check(
	// source is an optional argument that specifies a directory.
	source *dagger.Directory,
) *dagger.Container {
	return base().
		WithMountedDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithExec([]string{"/ruff", "check"})
}

// base returns the ruff base container
func base() *dagger.Container {
	return dag.Container().
		From("ghcr.io/astral-sh/ruff:latest")
}
