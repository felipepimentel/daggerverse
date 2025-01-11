// Package main provides a Dagger module for DigitalOcean operations
package main

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/felipepimentel/daggerverse/libraries/digitalocean/internal/dagger"
	"golang.org/x/oauth2"
)

// Digitalocean provides methods for interacting with DigitalOcean services
type Digitalocean struct {
	Token string // DigitalOcean API token
}

// AppConfig represents configuration for a DigitalOcean app
type AppConfig struct {
	Name             string
	Region           string
	InstanceSize     string
	InstanceCount    int64
	Container        Container
	EnvVars         map[string]string
	HealthCheckPath  string
	HTTPPort        int
}

// N8NAppConfig represents specific configuration for n8n deployment
type N8NAppConfig struct {
	AppConfig
	DatabaseURL string
	WebhookURL  string
	EncKey      string
}

// Container represents a Dagger container
type Container interface {
	WithRegistryAuth(registry, username string, secret Secret) Container
	Publish(ctx context.Context, ref string) (string, error)
}

// Secret represents a Dagger secret
type Secret interface {
	Plaintext(ctx context.Context) (string, error)
}

// StringSecret is a simple string-based secret
type StringSecret struct {
	value string
}

// Plaintext returns the secret value
func (s *StringSecret) Plaintext(ctx context.Context) (string, error) {
	return s.value, nil
}

// getClient creates a new DigitalOcean API client
func (d *Digitalocean) getClient() (*godo.Client, error) {
	if d.Token == "" {
		return nil, fmt.Errorf("DigitalOcean API token is required")
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: d.Token})
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	return godo.NewClient(oauthClient), nil
}

// CreateDroplet creates a new DigitalOcean droplet
func (d *Digitalocean) CreateDroplet(
	ctx context.Context,
	name string,
	region string,
	size string,
	image string,
	sshKeys []string,
) (*dagger.Container, error) {
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}

	createRequest := &godo.DropletCreateRequest{
		Name:    name,
		Region:  region,
		Size:    size,
		Image:   godo.DropletCreateImage{Slug: image},
		SSHKeys: make([]godo.DropletCreateSSHKey, len(sshKeys)),
	}

	for i, key := range sshKeys {
		createRequest.SSHKeys[i] = godo.DropletCreateSSHKey{Fingerprint: key}
	}

	droplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	// Return a container with DigitalOcean CLI configured
	return dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", d.Token).
		WithExec([]string{"compute", "droplet", "get", fmt.Sprintf("%d", droplet.ID), "--format", "json"}), nil
}

// ListDroplets returns a list of all droplets in the account
func (d *Digitalocean) ListDroplets(ctx context.Context) (*dagger.Container, error) {
	return dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", d.Token).
		WithExec([]string{"compute", "droplet", "list", "--format", "json"}), nil
}

// DeleteDroplet deletes a DigitalOcean droplet by ID
func (d *Digitalocean) DeleteDroplet(ctx context.Context, dropletID int) error {
	client, err := d.getClient()
	if err != nil {
		return err
	}

	_, err = client.Droplets.Delete(ctx, dropletID)
	return err
}

// GetDropletIP returns the public IP of a droplet
func (d *Digitalocean) GetDropletIP(ctx context.Context, dropletID int) (string, error) {
	client, err := d.getClient()
	if err != nil {
		return "", err
	}

	droplet, _, err := client.Droplets.Get(ctx, dropletID)
	if err != nil {
		return "", err
	}

	for _, network := range droplet.Networks.V4 {
		if network.Type == "public" {
			return network.IPAddress, nil
		}
	}

	return "", fmt.Errorf("no public IP found for droplet %d", dropletID)
}

// DeployApp deploys a container as a DigitalOcean app with more configuration options
func (d *Digitalocean) DeployApp(ctx context.Context, config AppConfig) (*dagger.Container, error) {
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}

	if config.InstanceSize == "" {
		config.InstanceSize = "basic-xxs"
	}
	if config.InstanceCount == 0 {
		config.InstanceCount = 1
	}
	if config.HTTPPort == 0 {
		config.HTTPPort = 80
	}

	// Convert environment variables to EnvVar format
	envVars := make([]*godo.AppVariableDefinition, 0, len(config.EnvVars))
	for key, value := range config.EnvVars {
		envVars = append(envVars, &godo.AppVariableDefinition{
			Key:   key,
			Value: value,
			Type:  godo.AppVariableType_General,
			Scope: godo.AppVariableScope_RunTime,
		})
	}

	// Create app spec
	spec := &godo.AppSpec{
		Name:   config.Name,
		Region: config.Region,
		Services: []*godo.AppServiceSpec{
			{
				Name:             config.Name,
				InstanceSizeSlug: config.InstanceSize,
				InstanceCount:    config.InstanceCount,
				Image: &godo.ImageSourceSpec{
					Registry:   "registry.digitalocean.com",
					Repository: config.Name,
					Tag:       "latest",
				},
				HealthCheck: &godo.AppServiceSpecHealthCheck{
					Path: config.HealthCheckPath,
				},
				Envs: envVars,
			},
		},
	}

	// Create the app
	app, _, err := client.Apps.Create(ctx, &godo.AppCreateRequest{Spec: spec})
	if err != nil {
		return nil, fmt.Errorf("failed to create app: %w", err)
	}

	// Push container to DO registry
	registryContainer := config.Container.WithRegistryAuth(
		"registry.digitalocean.com",
		"",
		&StringSecret{value: d.Token},
	)

	_, err = registryContainer.Publish(ctx, fmt.Sprintf("registry.digitalocean.com/%s:latest", config.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to push container: %w", err)
	}

	// Return container with app info
	return dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", d.Token).
		WithExec([]string{"apps", "get", app.ID, "--format", "json"}), nil
}

// WaitForAppDeployment waits for a DigitalOcean app deployment to complete
func (d *Digitalocean) WaitForAppDeployment(ctx context.Context, appID string) error {
	client, err := d.getClient()
	if err != nil {
		return err
	}

	for {
		app, _, err := client.Apps.Get(ctx, appID)
		if err != nil {
			return err
		}

		if app.ActiveDeployment != nil && app.ActiveDeployment.Phase == "ACTIVE" {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Continue polling
		}
	}
}

// GetAppURL returns the URL of a deployed app
func (d *Digitalocean) GetAppURL(ctx context.Context, appID string) (string, error) {
	client, err := d.getClient()
	if err != nil {
		return "", err
	}

	app, _, err := client.Apps.Get(ctx, appID)
	if err != nil {
		return "", err
	}

	if app.DefaultIngress != "" {
		return app.DefaultIngress, nil
	}

	return "", fmt.Errorf("no URL found for app %s", appID)
}

// DeployN8N deploys n8n as a DigitalOcean app with specific configurations
func (d *Digitalocean) DeployN8N(ctx context.Context, config N8NAppConfig) (*dagger.Container, error) {
	if config.InstanceSize == "" {
		config.InstanceSize = "professional-xs" // Recommended size for n8n
	}
	if config.Region == "" {
		config.Region = "nyc1" // Default region
	}
	if config.HTTPPort == 0 {
		config.HTTPPort = 5678 // Default n8n port
	}

	// Add n8n-specific environment variables
	if config.EnvVars == nil {
		config.EnvVars = make(map[string]string)
	}
	
	// Set required n8n environment variables
	config.EnvVars["N8N_PORT"] = fmt.Sprintf("%d", config.HTTPPort)
	config.EnvVars["N8N_PROTOCOL"] = "https"
	config.EnvVars["NODE_ENV"] = "production"
	
	if config.DatabaseURL != "" {
		config.EnvVars["DB_TYPE"] = "postgresdb"
		config.EnvVars["DB_POSTGRESDB_DATABASE"] = config.DatabaseURL
	}
	
	if config.WebhookURL != "" {
		config.EnvVars["WEBHOOK_URL"] = config.WebhookURL
	}
	
	if config.EncKey != "" {
		config.EnvVars["N8N_ENCRYPTION_KEY"] = config.EncKey
	}

	// Set health check path for n8n
	config.HealthCheckPath = "/healthz"

	// Deploy using the base app deployment method
	return d.DeployApp(ctx, config.AppConfig)
}

// GetN8NAppStatus returns the status and URL of a deployed n8n app
func (d *Digitalocean) GetN8NAppStatus(ctx context.Context, appID string) (string, string, error) {
	client, err := d.getClient()
	if err != nil {
		return "", "", err
	}

	app, _, err := client.Apps.Get(ctx, appID)
	if err != nil {
		return "", "", err
	}

	var url string
	if app.LiveDomain != "" {
		url = fmt.Sprintf("https://%s", app.LiveDomain)
	}

	return string(app.ActiveDeployment.Phase), url, nil
}
