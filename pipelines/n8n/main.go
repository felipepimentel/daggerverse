// Package main provides a Dagger module for n8n CI/CD operations
package main

import (
	"context"
	"fmt"

	"dagger/n8n/internal/dagger"
)

// Provider defines the interface for deployment providers
type Provider interface {
	DaggerObject
	Deploy(ctx context.Context, container *dagger.Container, registry string, tag string) error
	GetStatus(ctx context.Context) (*dagger.Container, error)
}

// N8N provides methods for building and deploying n8n
type N8N struct {
	DaggerObject
	// Source directory containing n8n configuration
	Source *dagger.Directory
	// Environment variables for n8n
	EnvVars []EnvVar
	// Port to expose n8n on
	Port int
	// Registry to publish to
	Registry string
	// Image tag
	Tag string
	// Registry auth token
	RegistryAuth *dagger.Secret
	// Deployment provider
	Provider Provider
	// DigitalOcean configuration (for backward compatibility)
	DOConfig *DOConfig
}

// EnvVar represents an environment variable
type EnvVar struct {
	DaggerObject
	Name  string
	Value string
}

// DOConfig represents DigitalOcean-specific configuration
type DOConfig struct {
	DaggerObject
	Token        *dagger.Secret
	Region       string
	AppName      string
	InstanceSize string
}

// DigitalOceanProvider implements the Provider interface for DigitalOcean
type DigitalOceanProvider struct {
	DaggerObject
	Token        *dagger.Secret
	Region       string
	AppName      string
	InstanceSize string
	Domain       string
}

// Build creates a container with n8n installed and configured
func (n *N8N) Build(ctx context.Context) (*dagger.Container, error) {
	if n.Source == nil {
		return nil, fmt.Errorf("source directory is required")
	}

	// Create n8n volume for persistence
	n8nVolume := dag.CacheVolume("n8n_data")

	container := dag.Container().
		From("node:18-alpine").
		WithMountedDirectory("/app", n.Source).
		WithMountedCache("/data", n8nVolume).
		WithWorkdir("/app").
		WithExec([]string{"npm", "install", "-g", "n8n"})

	// Set environment variables
	for _, env := range n.EnvVars {
		container = container.WithEnvVariable(env.Name, env.Value)
	}

	// Set default port if not specified
	if n.Port == 0 {
		n.Port = 5678
	}

	return container.
		WithExposedPort(n.Port).
		WithEntrypoint([]string{"n8n", "start"}), nil
}

// Test runs n8n tests
func (n *N8N) Test(ctx context.Context) error {
	if n.Source == nil {
		return fmt.Errorf("source directory is required")
	}

	container := dag.Container().
		From("node:18-alpine").
		WithMountedDirectory("/app", n.Source).
		WithWorkdir("/app")

	// Install dependencies
	_, err := container.WithExec([]string{"npm", "install"}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	// Run tests if package.json exists
	_, err = container.WithExec([]string{"test", "-f", "package.json"}).Sync(ctx)
	if err == nil {
		_, err = container.WithExec([]string{"npm", "test"}).Sync(ctx)
		if err != nil {
			return fmt.Errorf("tests failed: %w", err)
		}
	}

	return nil
}

// Publish creates and publishes a Docker image for n8n
func (n *N8N) Publish(ctx context.Context) (*dagger.Container, error) {
	if n.Source == nil {
		return nil, fmt.Errorf("source directory is required")
	}
	if n.Registry == "" {
		return nil, fmt.Errorf("registry is required")
	}
	if n.Tag == "" {
		return nil, fmt.Errorf("tag is required")
	}

	// Build the container
	container, err := n.Build(ctx)
	if err != nil {
		return nil, err
	}

	// Add registry auth if provided
	if n.RegistryAuth != nil {
		container = container.WithRegistryAuth(n.Registry, "", n.RegistryAuth)
	}

	// Publish the container
	publishedRef, err := container.Publish(ctx, fmt.Sprintf("%s:%s", n.Registry, n.Tag))
	if err != nil {
		return nil, fmt.Errorf("failed to publish container: %w", err)
	}

	return dag.Container().From(publishedRef), nil
}

// Deploy deploys n8n using the configured provider
func (n *N8N) Deploy(ctx context.Context) (*dagger.Container, error) {
	if n.Provider == nil {
		return nil, fmt.Errorf("deployment provider is required")
	}

	// Build and publish the container
	container, err := n.Publish(ctx)
	if err != nil {
		return nil, fmt.Errorf("deployment failed: %w", err)
	}

	// Deploy using the provider
	if err := n.Provider.Deploy(ctx, container, n.Registry, n.Tag); err != nil {
		return nil, fmt.Errorf("deployment failed: %w", err)
	}

	return container, nil
}

// Deploy implements the Provider interface for DigitalOcean
func (p *DigitalOceanProvider) Deploy(ctx context.Context, container *dagger.Container, registry string, tag string) error {
	if p.Token == nil {
		return fmt.Errorf("DigitalOcean token is required")
	}

	// Get token value
	token, err := p.Token.Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Create app spec with Caddy configuration if domain is provided
	services := []map[string]interface{}{
		{
			"name":              "n8n",
			"instance_size_slug": p.InstanceSize,
			"instance_count":    1,
			"image": map[string]string{
				"registry_type": "DOCR",
				"repository":    registry,
				"tag":          tag,
			},
			"health_check": map[string]interface{}{
				"http_path":             "/healthz",
				"initial_delay_seconds": 30,
			},
			"volumes": []map[string]string{
				{
					"name":       "n8n-data",
					"mount_path": "/data",
				},
			},
		},
	}

	// Add Caddy service if domain is provided
	if p.Domain != "" {
		services = append(services, map[string]interface{}{
			"name":              "caddy",
			"instance_size_slug": "basic-xxs",
			"instance_count":    1,
			"image": map[string]string{
				"registry_type": "DOCKER_HUB",
				"repository":    "caddy",
				"tag":          "2.7-alpine",
			},
			"volumes": []map[string]string{
				{
					"name":       "caddy-data",
					"mount_path": "/data",
				},
			},
			"env": []map[string]string{
				{
					"key":   "DOMAIN",
					"value": p.Domain,
				},
			},
		})
	}

	// Create app spec
	appSpec := map[string]interface{}{
		"name":     p.AppName,
		"region":   p.Region,
		"services": services,
	}

	// Create a new container with doctl
	_, err = dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{"apps", "create", "--spec", fmt.Sprintf("%v", appSpec)}).
		Sync(ctx)

	return err
}

// GetStatus implements the Provider interface for DigitalOcean
func (p *DigitalOceanProvider) GetStatus(ctx context.Context) (*dagger.Container, error) {
	if p.Token == nil {
		return nil, fmt.Errorf("DigitalOcean token is required")
	}

	// Get token value
	token, err := p.Token.Plaintext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Get app status
	return dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{"apps", "list", "--format", "json"}), nil
}

// CI runs the CI pipeline for n8n
func (n *N8N) CI(ctx context.Context) error {
	if n.Source == nil && n.Registry == "" && n.Tag == "" {
		return fmt.Errorf("source, registry and tag are required")
	}
	return n.Test(ctx)
}

// CD runs the CD pipeline for n8n
func (n *N8N) CD(ctx context.Context) (*dagger.Container, error) {
	if n.Source == nil && n.Registry == "" && n.Tag == "" {
		return nil, fmt.Errorf("source, registry and tag are required")
	}
	return n.Deploy(ctx)
}

// WithSource sets the source directory
func (n *N8N) WithSource(source *dagger.Directory) *N8N {
	n.Source = source
	return n
}

// WithRegistry sets the registry
func (n *N8N) WithRegistry(registry string) *N8N {
	n.Registry = registry
	return n
}

// WithTag sets the tag
func (n *N8N) WithTag(tag string) *N8N {
	n.Tag = tag
	return n
}

// WithRegistryAuth sets the registry auth token
func (n *N8N) WithRegistryAuth(auth *dagger.Secret) *N8N {
	n.RegistryAuth = auth
	return n
}

// WithProvider sets the deployment provider
func (n *N8N) WithProvider(provider Provider) *N8N {
	n.Provider = provider
	return n
}
