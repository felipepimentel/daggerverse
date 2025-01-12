// Package main provides a Dagger module for deploying n8n to DigitalOcean
package main

import (
	"context"
	"fmt"

	do "github.com/felipepimentel/daggerverse/libraries/digitalocean"
	docker "github.com/felipepimentel/daggerverse/libraries/docker"
	sshmanager "github.com/felipepimentel/daggerverse/libraries/ssh-manager"
	n8n "github.com/felipepimentel/daggerverse/pipelines/n8n"

	"dagger.io/dagger"
)

// N8NDigitalOcean represents the n8n deployment pipeline for DigitalOcean
type N8NDigitalOcean struct {
	// N8N pipeline configuration
	N8N *n8n.N8N
	// DigitalOcean configuration
	DigitalOcean *do.Digitalocean
	// Docker configuration for building and pushing images
	Docker *docker.Docker
	// SSH Manager for key management
	SSHManager *sshmanager.SSHManager
	// Application configuration
	AppName      string
	Region       string
	InstanceSize string
	DatabaseURL  string
	WebhookURL   string
	EncKey       string
}

// Deploy builds and deploys n8n to DigitalOcean
func (n *N8NDigitalOcean) Deploy(ctx context.Context) (*dagger.Container, error) {
	if n.N8N == nil {
		return nil, fmt.Errorf("n8n configuration is required")
	}
	if n.DigitalOcean == nil {
		return nil, fmt.Errorf("digitalocean configuration is required")
	}
	if n.Docker == nil {
		return nil, fmt.Errorf("docker configuration is required")
	}
	if n.SSHManager == nil {
		return nil, fmt.Errorf("ssh manager is required")
	}

	// Generate ephemeral SSH key
	key, err := n.SSHManager.GenerateKey(ctx, fmt.Sprintf("n8n-%s", n.AppName))
	if err != nil {
		return nil, fmt.Errorf("failed to generate SSH key: %w", err)
	}

	// Clean up the key after deployment
	defer func() {
		if err := n.SSHManager.DeleteKey(ctx, key.Fingerprint); err != nil {
			fmt.Printf("Warning: failed to delete SSH key: %v\n", err)
		}
	}()

	// Build and publish n8n container
	container, err := n.N8N.CD(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build and publish n8n: %w", err)
	}

	// Configure n8n app deployment
	appConfig := do.N8NAppConfig{
		AppConfig: do.AppConfig{
			Name:         n.AppName,
			Region:       n.Region,
			InstanceSize: n.InstanceSize,
			Container:    container,
		},
		DatabaseURL: n.DatabaseURL,
		WebhookURL:  n.WebhookURL,
		EncKey:      n.EncKey,
	}

	// Deploy to DigitalOcean
	deployedApp, err := n.DigitalOcean.DeployN8N(ctx, appConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy n8n to DigitalOcean: %w", err)
	}

	return deployedApp, nil
}

// GetStatus returns the deployment status and URL of the n8n app
func (n *N8NDigitalOcean) GetStatus(ctx context.Context, appID string) (string, string, error) {
	if n.DigitalOcean == nil {
		return "", "", fmt.Errorf("digitalocean configuration is required")
	}

	return n.DigitalOcean.GetN8NAppStatus(ctx, appID)
}
