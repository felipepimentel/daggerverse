// Envoy proxy module for Dagger.

// The Envoy proxy Dagger module enables you to do the following:
// - run an instance of Envoy proxy using the provided configuration
// - validate the Envoy configuration file
package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/libraries/envoy/internal/dagger"
)

type Envoy struct {
	Version  string
	Platform dagger.Platform
}

func New() *Envoy {
	return &Envoy{
		Version:  "v1.30-latest",
		Platform: "linux/arm64",
	}
}

// EnvoyProxyService creates a new service that runs the Envoy proxy with the given configuration.
// Example usage:
//  1. Starts the Envoy proxy with the given config and exposes port 10000 to the host.
//     dagger call envoy-proxy-service --config ./examples/httpbin-sample.yaml --port 10000 up
func (m *Envoy) EnvoyProxyService(
	ctx context.Context,
	// +optional
	// +default="v1.30-latest"
	version string,
	// +optional
	// +default="linux/arm64"
	platform dagger.Platform,
	// +required
	config *dagger.File,
	port []int,
) (*dagger.Service, error) {
	opts := dagger.ContainerOpts{
		Platform: platform,
	}

	configContents, err := config.Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read config contents: %w", err)
	}

	container := dag.Container(opts).
		From("envoyproxy/envoy:"+m.Version).
		WithNewFile("/etc/envoy/envoy.yaml", configContents)

	for _, p := range port {
		container = container.WithExposedPort(p)
	}
	return container.AsService(), nil
}

// ValidateConfig validates the given Envoy configuration.
func (m *Envoy) ValidateConfig(
	ctx context.Context,
	// +optional
	// +default="v1.30-latest"
	version string,
	// +optional
	// +default="linux/arm64"
	platform dagger.Platform,
	// +required
	config *dagger.File) (string, error) {

	opts := dagger.ContainerOpts{
		Platform: platform,
	}

	configContents, err := config.Contents(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to read config contents: %w", err)
	}

	// Run Envoy container with --mode validate and pipe the response to stdout
	return dag.Container(opts).
		From("envoyproxy/envoy:"+m.Version).
		WithNewFile("/etc/envoy/envoy.yaml", configContents).
		WithExec([]string{"envoy", "--mode", "validate", "-c", "/etc/envoy/envoy.yaml"}).
		Stdout(ctx)
}
