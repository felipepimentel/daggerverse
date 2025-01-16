// Package main provides a Dagger module for Docker operations
package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/libraries/docker/internal/dagger"
)

// Docker provides methods for Docker operations
type Docker struct {
	client   *dagger.Client
	registry *RegistryConfig
}

// Label represents a Docker label
type Label struct {
	Key   string
	Value string
}

// BuildArg represents a Docker build argument
type BuildArg struct {
	Key   string
	Value string
}

// DriverOpt represents a driver-specific option
type DriverOpt struct {
	Key   string
	Value string
}

// ImageConfig represents configuration for Docker image operations
type ImageConfig struct {
	Source      string       // Source image name
	Target      string       // Target image name
	Tag         string       // Image tag
	Labels      []Label      // Image labels
	BuildArgs   []BuildArg   // Build arguments
	Dockerfile  string       // Path to Dockerfile
	Context     *dagger.Directory // Build context
	PullPolicy  string       // Pull policy (always, never, if-not-present)
	Registry    string       // Registry URL
}

// RegistryConfig represents configuration for Docker registry operations
type RegistryConfig struct {
	URL      string         // Registry URL
	Username string         // Registry username
	Password *dagger.Secret // Registry password
}

// VolumeConfig represents configuration for Docker volume operations
type VolumeConfig struct {
	Name       string       // Volume name
	Driver     string       // Volume driver
	Labels     []Label      // Volume labels
	DriverOpts []DriverOpt  // Driver-specific options
}

// NetworkConfig represents configuration for Docker network operations
type NetworkConfig struct {
	Name       string       // Network name
	Driver     string       // Network driver
	Labels     []Label      // Network labels
	DriverOpts []DriverOpt  // Driver-specific options
	Internal   bool         // Internal network
	IPv6       bool         // Enable IPv6
}

// New creates a new Docker client
func New() *Docker {
	return &Docker{
		client: dag,
	}
}

// WithRegistry configures registry authentication
func (d *Docker) WithRegistry(config RegistryConfig) *Docker {
	d.registry = &config
	return d
}

// PullImage pulls a Docker image
func (d *Docker) PullImage(ctx context.Context, image, tag string, config *ImageConfig) (*dagger.Container, error) {
	container := d.client.Container()
	
	if config != nil {
		if config.Registry != "" {
			image = fmt.Sprintf("%s/%s", config.Registry, image)
		}
		
		if config.Labels != nil {
			for _, label := range config.Labels {
				container = container.WithLabel(label.Key, label.Value)
			}
		}
	}
	
	return container.From(fmt.Sprintf("%s:%s", image, tag)), nil
}

// PushImage pushes a Docker image to a registry
func (d *Docker) PushImage(ctx context.Context, config ImageConfig) error {
	container := d.client.Container().From(config.Source)
	
	if config.Labels != nil {
		for _, label := range config.Labels {
			container = container.WithLabel(label.Key, label.Value)
		}
	}

	if d.registry != nil {
		container = container.WithRegistryAuth(
			d.registry.URL,
			d.registry.Username,
			d.registry.Password,
		)
	}

	_, err := container.Publish(ctx, config.Target)
	if err != nil {
		return fmt.Errorf("failed to push image %s: %w", config.Target, err)
	}

	return nil
}

// BuildImage builds a Docker image from a context
func (d *Docker) BuildImage(ctx context.Context, config ImageConfig) (*dagger.Container, error) {
	if config.Context == nil {
		return nil, fmt.Errorf("build context is required")
	}

	container := d.client.Container()
	
	if config.Dockerfile != "" {
		container = container.Build(config.Context, dagger.ContainerBuildOpts{
			Dockerfile: config.Dockerfile,
		})
	} else {
		container = container.Build(config.Context)
	}
	
	if config.BuildArgs != nil {
		for _, buildArg := range config.BuildArgs {
			container = container.WithEnvVariable(buildArg.Key, buildArg.Value)
		}
	}

	if config.Labels != nil {
		for _, label := range config.Labels {
			container = container.WithLabel(label.Key, label.Value)
		}
	}

	return container, nil
}

// TagImage tags a Docker image
func (d *Docker) TagImage(ctx context.Context, source, target string) error {
	container := d.client.Container().From(source)
	
	if d.registry != nil {
		container = container.WithRegistryAuth(
			d.registry.URL,
			d.registry.Username,
			d.registry.Password,
		)
	}

	_, err := container.Publish(ctx, target)
	if err != nil {
		return fmt.Errorf("failed to tag image %s as %s: %w", source, target, err)
	}

	return nil
}

// InspectImage returns information about a Docker image
func (d *Docker) InspectImage(ctx context.Context, image string) (*dagger.Container, error) {
	container := d.client.Container().From(image)
	
	if d.registry != nil {
		container = container.WithRegistryAuth(
			d.registry.URL,
			d.registry.Username,
			d.registry.Password,
		)
	}

	return container, nil
}

// VerifyImageExists checks if an image exists in the registry
func (d *Docker) VerifyImageExists(ctx context.Context, image string) (bool, error) {
	container := d.client.Container().From(image)
	
	if d.registry != nil {
		container = container.WithRegistryAuth(
			d.registry.URL,
			d.registry.Username,
			d.registry.Password,
		)
	}

	_, err := container.ID(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to verify image %s: %w", image, err)
	}

	return true, nil
}

// CreateVolume creates a new Docker volume
func (d *Docker) CreateVolume(ctx context.Context, config VolumeConfig) (*dagger.CacheVolume, error) {
	volume := d.client.CacheVolume(config.Name)
	return volume, nil
}

// RemoveVolume removes a Docker volume
func (d *Docker) RemoveVolume(ctx context.Context, name string) error {
	// Note: Currently, Dagger doesn't provide direct volume removal.
	// This is a placeholder for future implementation.
	return fmt.Errorf("volume removal not supported yet")
}

// CreateNetwork creates a new Docker network
func (d *Docker) CreateNetwork(ctx context.Context, config NetworkConfig) error {
	// Note: Currently, Dagger doesn't provide direct network management.
	// This is a placeholder for future implementation.
	return fmt.Errorf("network creation not supported yet")
}

// RemoveNetwork removes a Docker network
func (d *Docker) RemoveNetwork(ctx context.Context, name string) error {
	// Note: Currently, Dagger doesn't provide direct network management.
	// This is a placeholder for future implementation.
	return fmt.Errorf("network removal not supported yet")
}

// ListNetworks lists all Docker networks
func (d *Docker) ListNetworks(ctx context.Context) error {
	// Note: Currently, Dagger doesn't provide direct network management.
	// This is a placeholder for future implementation.
	return fmt.Errorf("network listing not supported yet")
}

// ListVolumes lists all Docker volumes
func (d *Docker) ListVolumes(ctx context.Context) ([]*dagger.CacheVolume, error) {
	// Note: Currently, Dagger doesn't provide volume listing.
	// This is a placeholder for future implementation.
	return nil, fmt.Errorf("volume listing not supported yet")
}
