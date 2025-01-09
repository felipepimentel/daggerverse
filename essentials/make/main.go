package main

import (
	"github.com/felipepimentel/daggerverse/essentials/make/internal/dagger"
)

// A Dagger module to use make
type Make struct{}

// Execute the command 'make' in a directory, and return the modified directory
func (m *Make) Make(
	dir *dagger.Directory,
	args []string,
	// +optional
	// +default="Makefile"
	makefile string,
) *dagger.Directory {
	return dag.
		Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "make"}).
		WithEntrypoint([]string{"make"}).
		WithMountedDirectory("/src", dir).
		WithWorkdir("/src").
		WithExec(append([]string{"-f", makefile}, args...)).
		Directory(".")
}
