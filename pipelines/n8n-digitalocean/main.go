// Package main provides a Dagger module for deploying n8n to DigitalOcean
package main

import (
	"context"
	"fmt"

	"dagger/n8n-digitalocean/internal/dagger"
)

// N8NDigitalOcean represents the n8n deployment pipeline for DigitalOcean
type N8NDigitalOcean struct {
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

// Deploy builds and deploys n8n to DigitalOcean
func (n *N8NDigitalOcean) Deploy(ctx context.Context) (*dagger.Container, error) {
	if n.Source == nil {
		return nil, fmt.Errorf("source directory is required")
	}
	if n.DOConfig == nil {
		return nil, fmt.Errorf("digitalocean configuration is required")
	}
	if n.DOConfig.Token == nil {
		return nil, fmt.Errorf("digitalocean token is required")
	}

	// Build n8n container
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

	container = container.
		WithExposedPort(n.Port).
		WithEntrypoint([]string{"n8n", "start"})

	// Add registry auth if provided
	if n.RegistryAuth != nil {
		container = container.WithRegistryAuth(n.Registry, "", n.RegistryAuth)
	}

	// Publish the container
	_, err := container.Publish(ctx, fmt.Sprintf("%s:%s", n.Registry, n.Tag))
	if err != nil {
		return nil, fmt.Errorf("failed to publish container: %w", err)
	}

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

	// Get token value
	token, err := n.DOConfig.Token.Plaintext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Create a new container with doctl
	return dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{"apps", "create", "--spec", appSpec}), nil
}

// GetStatus returns the deployment status and URL of the n8n app
func (n *N8NDigitalOcean) GetStatus(ctx context.Context, appID string) (*dagger.Container, error) {
	if n.DOConfig == nil || n.DOConfig.Token == nil {
		return nil, fmt.Errorf("digitalocean configuration is required")
	}

	// Get token value
	token, err := n.DOConfig.Token.Plaintext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Get app status
	return dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{"apps", "get", appID, "--format", "json"}), nil
}
