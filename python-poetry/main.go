// Package main provides a Dagger module for Python Poetry projects
package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// Poetry is the main type for the module
type Poetry struct {
	// Base container for Python operations
	Base *dagger.Container
}

// New creates a new instance of the Poetry module
func New() *Poetry {
	return &Poetry{}
}

// InstallDeps installs Poetry and project dependencies
func (p *Poetry) InstallDeps(ctx context.Context, src *dagger.Directory, pythonVersion string) *dagger.Container {
	return p.Base.From(fmt.Sprintf("python:%s-slim", pythonVersion)).
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "poetry"}).
		WithExec([]string{"poetry", "install"})
}

// Lint runs Ruff linter on the project
func (p *Poetry) Lint(ctx context.Context, src *dagger.Directory, pythonVersion string) *dagger.Container {
	return p.InstallDeps(ctx, src, pythonVersion).
		WithExec([]string{"poetry", "run", "ruff", "check", "."})
}

// Build builds the Python project using Poetry
func (p *Poetry) Build(ctx context.Context, src *dagger.Directory, pythonVersion string) *dagger.Container {
	return p.InstallDeps(ctx, src, pythonVersion).
		WithExec([]string{"poetry", "build"})
}

// Publish publishes the package to the specified repository
func (p *Poetry) Publish(ctx context.Context, src *dagger.Directory, pythonVersion, repository string) *dagger.Container {
	return p.Build(ctx, src, pythonVersion).
		WithExec([]string{"poetry", "config", "repositories.custom", repository}).
		WithExec([]string{"poetry", "publish", "--repository", "custom", "--build"})
} 