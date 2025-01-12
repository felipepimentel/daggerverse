// Package main provides a Dagger module for n8n CI/CD operations
package main

import (
	"context"
	"fmt"

	"dagger/n8n/internal/dagger"
)

// N8N provides methods for building and deploying n8n
type N8N struct {
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
	// DigitalOcean configuration
	DOConfig *DOConfig
}

// EnvVar represents an environment variable
type EnvVar struct {
	Name  string
	Value string
}

// DOConfig represents DigitalOcean-specific configuration
type DOConfig struct {
	Token        *dagger.Secret
	Region       string
	AppName      string
	InstanceSize string
}

// Build creates a container with n8n installed and configured
func (n *N8N) Build(ctx context.Context) (*dagger.Container, error) {
	if n.Source == nil {
		return nil, fmt.Errorf("source directory is required")
	}

	container := dag.Container().
		From("node:18-alpine").
		WithMountedDirectory("/app", n.Source).
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

// CI runs the CI pipeline for n8n
func (n *N8N) CI(ctx context.Context) error {
	// Run tests
	if err := n.Test(ctx); err != nil {
		return fmt.Errorf("CI failed: %w", err)
	}

	// Build container to verify it works
	_, err := n.Build(ctx)
	if err != nil {
		return fmt.Errorf("CI failed: %w", err)
	}

	return nil
}

// CD runs the CD pipeline for n8n
func (n *N8N) CD(ctx context.Context) (*dagger.Container, error) {
	// Run CI first
	if err := n.CI(ctx); err != nil {
		return nil, fmt.Errorf("CD failed: %w", err)
	}

	// Build and publish the container
	container, err := n.Publish(ctx)
	if err != nil {
		return nil, fmt.Errorf("CD failed: %w", err)
	}

	// If DigitalOcean config is provided, deploy to DO
	if n.DOConfig != nil {
		if n.DOConfig.Token == nil {
			return nil, fmt.Errorf("DigitalOcean token is required for deployment")
		}

		// Get token value
		token, err := n.DOConfig.Token.Plaintext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get token: %w", err)
		}

		// Create a new container with doctl
		doctlContainer := dag.Container().
			From("digitalocean/doctl:latest").
			WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token)

		// Create app spec
		appSpec := fmt.Sprintf(`{
			"name": "%s",
			"region": "%s",
			"services": [{
				"name": "%s",
				"instance_size_slug": "%s",
				"instance_count": 1,
				"image": {
					"registry_type": "DOCR",
					"repository": "%s",
					"tag": "%s"
				},
				"health_check": {
					"http_path": "/healthz"
				}
			}]
		}`, n.DOConfig.AppName, n.DOConfig.Region, n.DOConfig.AppName, n.DOConfig.InstanceSize, n.Registry, n.Tag)

		// Create app
		_, err = doctlContainer.
			WithExec([]string{"apps", "create", "--spec", appSpec}).
			Sync(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create app: %w", err)
		}

		return doctlContainer, nil
	}

	return container, nil
}
