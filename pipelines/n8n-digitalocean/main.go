// Package main provides a Dagger module for deploying n8n to DigitalOcean
package main

import (
	"context"
	"fmt"

	"dagger/n8n-digitalocean/internal/dagger"
)

var client *dagger.Client

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

// New creates a new instance of the n8n-digitalocean module
func New(c *dagger.Client) *N8NDigitalOcean {
	client = c
	return &N8NDigitalOcean{}
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

// CaddyConfig represents the Caddy server configuration
type CaddyConfig struct {
	Domain string
}

// generateCaddyfile creates a Caddyfile configuration
func (n *N8NDigitalOcean) generateCaddyfile(domain string) string {
	return fmt.Sprintf(`%s {
    reverse_proxy n8n:5678 {
        flush_interval -1
    }
}`, domain)
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

	// Create n8n volume
	n8nVolume := client.CacheVolume("n8n_data")

	// Build n8n container
	container := client.Container().
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

	// Create app spec with Caddy configuration
	appSpec := fmt.Sprintf(`{
		"name": "%s",
		"region": "%s",
		"services": [
			{
				"name": "n8n",
				"instance_size_slug": "%s",
				"instance_count": 1,
				"image": {
					"registry_type": "DOCR",
					"repository": "%s",
					"tag": "%s"
				},
				"health_check": {
					"http_path": "/healthz",
					"initial_delay_seconds": 30
				},
				"volumes": [
					{
						"name": "n8n-data",
						"mount_path": "/data"
					}
				]
			},
			{
				"name": "caddy",
				"instance_size_slug": "basic-xxs",
				"instance_count": 1,
				"image": {
					"registry_type": "DOCKER_HUB",
					"repository": "caddy",
					"tag": "2.7-alpine"
				},
				"volumes": [
					{
						"name": "caddy-data",
						"mount_path": "/data"
					}
				],
				"env": [
					{
						"key": "DOMAIN",
						"value": "%s"
					}
				]
			}
		]
	}`, n.DOConfig.AppName, n.DOConfig.Region, n.DOConfig.InstanceSize, n.Registry, n.Tag, n.EnvVars[0].Value)

	// Get token value
	token, err := n.DOConfig.Token.Plaintext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Create a new container with doctl
	return client.Container().
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
	return client.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{"apps", "get", appID, "--format", "json"}), nil
}

// CI runs the CI pipeline for n8n
func (n *N8NDigitalOcean) CI(ctx context.Context, source *dagger.Directory, region string, appName string, token *dagger.Secret, domain string, basicAuthPassword string, encryptionKey string, sshKey *dagger.Secret, sshKeyFingerprint string, sshKeyID string) (*dagger.Container, error) {
	n.Source = source
	n.Registry = "registry.digitalocean.com/pimentel/n8n"
	n.Tag = "latest"
	n.RegistryAuth = token
	n.DOConfig = &DOConfig{
		Token:        token,
		Region:       region,
		AppName:      appName,
		InstanceSize: "basic-xxs",
	}

	// Set environment variables
	n.EnvVars = []EnvVar{
		{Name: "N8N_HOST", Value: domain},
		{Name: "N8N_PROTOCOL", Value: "https"},
		{Name: "N8N_PORT", Value: "5678"},
		{Name: "N8N_BASIC_AUTH_ACTIVE", Value: "true"},
		{Name: "N8N_BASIC_AUTH_USER", Value: "admin"},
		{Name: "N8N_BASIC_AUTH_PASSWORD", Value: basicAuthPassword},
		{Name: "N8N_ENCRYPTION_KEY", Value: encryptionKey},
		{Name: "NODE_ENV", Value: "production"},
		{Name: "DB_TYPE", Value: "sqlite"},
		{Name: "DB_SQLITE_PATH", Value: "/home/node/.n8n/database.sqlite"},
	}

	return n.Deploy(ctx)
}
