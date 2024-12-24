// Auto detect runtime information for a specific project

package main

import (
	"context"
	"main/internal/dagger"
)

type Autodetection struct{}

// Expose node auto dection runtime information
func (a *Autodetection) Node(
	ctx context.Context,
	// The path to the project to analyze
	src *dagger.Directory,
	// Define patterns to exclude from the analysis
	// +optional
	patternExclusions []string,
) (*NodeAnalyzer, error) {
	return newNodeAnalyzer(ctx, src, patternExclusions)
}

// Expose OCI dection runtime information
func (a *Autodetection) Oci(
	ctx context.Context,
	// The path to the project to analyze
	src *dagger.Directory,
	// Define patterns to exclude from the analysis
	// +optional
	patternExclusions []string,
) (*OciAnalyzer, error) {
	return newOciAnalyzer(ctx, src, patternExclusions)
}
