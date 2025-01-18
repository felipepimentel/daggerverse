package main

import (
	"context"
	"slices"

	"github.com/felipepimentel/daggerverse/essentials/inline-node/internal/dagger"
)

type InlineNode struct {
	Ctr *dagger.Container
}

// New creates a new Node.js environment
func New(
	// Node.js version to use
	// +optional
	// +default="20"
	version string,
) *InlineNode {
	if version == "" {
		version = "20"
	}
	return &InlineNode{
		Ctr: dag.Container().From("node:" + version + "-alpine"),
	}
}

func (m *InlineNode) WithPackage(name string) (*InlineNode, error) {
	m.Ctr = m.Ctr.WithExec([]string{"npm", "install", "--no-package-lock", name})
	return m, nil
}

func (m *InlineNode) WithPackages(packages []string) (*InlineNode, error) {
	// sort the requirements to optimize caching
	slices.Sort(packages)

	for _, name := range packages {
		m.WithPackage(name)
	}

	return m, nil
}

func (m *InlineNode) Code(ctx context.Context, code string) (*dagger.Container, error) {
	return m.Ctr.
		WithNewFile("script.js", code).
		WithExec([]string{"node", "script.js"}).
		Sync(ctx)
}
