package main

import (
	"context"
	"slices"

	"github.com/felipepimentel/daggerverse/essentials/inline-python/internal/dagger"
)

type InlinePython struct {
	Ctr *dagger.Container
}

// New creates a new Python environment
func New(
	// Python version to use
	// +optional
	// +default="3"
	version string,
) *InlinePython {
	if version == "" {
		version = "3"
	}
	return &InlinePython{
		Ctr: dag.Container().From("python:" + version + "-alpine"),
	}
}

// fixme: add usage
func (m *InlinePython) WithPackage(name string) (*InlinePython, error) {
	m.Ctr = m.Ctr.WithExec([]string{"pip", "install", name})
	return m, nil
}

// fixme: add usage
func (m *InlinePython) WithPackages(packages []string) (*InlinePython, error) {
	// sort the requirements to optimize caching
	slices.Sort(packages)

	for _, name := range packages {
		m.WithPackage(name)
	}

	return m, nil
}

func (m *InlinePython) Code(ctx context.Context, code string) (*dagger.Container, error) {
	return m.Ctr.WithExec([]string{"python", "-c", code}).Sync(ctx)
}
