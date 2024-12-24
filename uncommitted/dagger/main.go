// This module checks for uncommitted git changes

package main

import "dagger/uncommitted/internal/dagger"

type Uncommitted struct {
	// +private
	Ctr *dagger.Container
}

func New(
	// Python image to use.
	// +optional
	// renovate image: datasource=docker depName=python versioning=docker
	// +default="python:3.13.1-alpine"
	Image string,
) *Uncommitted {
	return &Uncommitted{
		Ctr: dag.Container().From(Image).
			WithExec([]string{"apk", "add", "git"}).
			WithExec([]string{"pip", "install", "setuptools", "check-uncommitted-git-changes"}),
	}
}

// CheckUncommitted runs check_uncommitted_git_changes
//
// Example usage: dagger call check-uncommitted --source /path/to/your/repo
func (m *Uncommitted) CheckUncommitted(source *dagger.Directory) *dagger.Container {
	return m.Ctr.
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"check_uncommitted_git_changes"})
}
