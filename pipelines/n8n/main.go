// Package main provides a Dagger module for n8n CI/CD operations
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/felipepimentel/daggerverse/pipelines/n8n/internal/dagger"
)

// Error messages for common failures
const (
	errSetupRegistry   = "failed to setup registry"
	errPushImages      = "failed to push images"
	errCreateDroplet   = "failed to create droplet"
	errConfigureDNS    = "failed to configure DNS"
	errSetupServer     = "failed to setup server"
	errDeployServices  = "failed to deploy services"
	errSetupMonitoring = "failed to setup monitoring"
	errSetupBackups    = "failed to setup backups"
)

// Log messages for progress tracking
const (
	logStartDeploy     = "Starting n8n deployment..."
	logSetupRegistry   = "Setting up container registry..."
	logPushImages      = "Pushing required images..."
	logCreateDroplet   = "Creating DigitalOcean droplet..."
	logConfigureDNS    = "Configuring DNS..."
	logSetupServer     = "Setting up server..."
	logDeployServices  = "Deploying services..."
	logSetupMonitoring = "Setting up monitoring..."
	logSetupBackups    = "Setting up backups..."
	logSuccess         = "Deployment completed successfully!"
)

// N8N represents the n8n deployment pipeline configuration
type N8N struct {
	// Base configuration
	Domain     string // e.g., "example.com"
	Subdomain  string // e.g., "n8n"
	Region     string // e.g., "nyc1"
	Size       string // e.g., "s-2vcpu-4gb"
	SSHKeyName string // Name of the SSH key in DigitalOcean

	// Registry configuration
	RegistryName string // e.g., "n8n-registry"

	// n8n specific configuration
	N8nVersion string // e.g., "0.234.0"
	N8nPort    int    // e.g., 5678

	// Database configuration
	PostgresVersion string // e.g., "15-alpine"
	PostgresUser    string
	PostgresPass    string
	PostgresDB      string

	// Backup configuration
	BackupEnabled    bool
	BackupCron       string // e.g., "0 0 * * *"
	BackupRetention  int    // Number of days to retain backups

	// Monitoring configuration
	MonitoringEnabled bool
	CAdvisorPort     int // e.g., 8080

	// SSL/TLS configuration
	SSLEmail string // Email for Let's Encrypt
}

// New creates a new instance of the N8N pipeline with default values
func New(
	// Domain for the n8n instance
	// +optional
	domain string,
	// Subdomain for the n8n instance
	// +optional
	subdomain string,
	// SSH key name in DigitalOcean
	// +optional
	sshKeyName string,
	// Region for deployment
	// +optional
	// +default="nyc1"
	region string,
	// Droplet size
	// +optional
	// +default="s-2vcpu-4gb"
	size string,
) *N8N {
	if region == "" {
		region = "nyc1"
	}
	if size == "" {
		size = "s-2vcpu-4gb"
	}

	return &N8N{
		Domain:           domain,
		Subdomain:       subdomain,
		Region:          region,
		Size:            size,
		SSHKeyName:      sshKeyName,
		RegistryName:    "n8n-registry",
		N8nVersion:      "0.234.0",
		N8nPort:         5678,
		PostgresVersion: "15-alpine",
		PostgresUser:    "n8n",
		PostgresPass:    "n8n",
		PostgresDB:      "n8n",
		BackupEnabled:   true,
		BackupCron:      "0 0 * * *",
		BackupRetention: 7,
		MonitoringEnabled: true,
		CAdvisorPort:    8080,
	}
}

// Deploy orchestrates the complete n8n deployment pipeline
func (n *N8N) Deploy(ctx context.Context, doToken, sshKey *dagger.Secret) error {
	fmt.Println(logStartDeploy)

	// Step 1: Set up container registry
	fmt.Println(logSetupRegistry)
	registry := dag.Container().From("registry:2")
	if err := n.setupRegistry(ctx, registry); err != nil {
		return fmt.Errorf("%s: %w", errSetupRegistry, err)
	}

	// Step 2: Push required images
	fmt.Println(logPushImages)
	if err := n.pushImages(ctx, registry); err != nil {
		return fmt.Errorf("%s: %w", errPushImages, err)
	}

	// Step 3: Create droplet
	fmt.Println(logCreateDroplet)
	dropletIP, err := n.createDroplet(ctx, doToken)
	if err != nil {
		return fmt.Errorf("%s: %w", errCreateDroplet, err)
	}

	// Step 4: Configure DNS
	fmt.Println(logConfigureDNS)
	if err := n.configureDNS(ctx, doToken, dropletIP); err != nil {
		return fmt.Errorf("%s: %w", errConfigureDNS, err)
	}

	// Step 5: Wait for DNS propagation
	time.Sleep(30 * time.Second)

	// Step 6: Setup server
	fmt.Println(logSetupServer)
	if err := n.setupServer(ctx, sshKey, dropletIP); err != nil {
		return fmt.Errorf("%s: %w", errSetupServer, err)
	}

	// Step 7: Deploy services
	fmt.Println(logDeployServices)
	if err := n.deployServices(ctx, sshKey, dropletIP); err != nil {
		return fmt.Errorf("%s: %w", errDeployServices, err)
	}

	// Step 8: Configure monitoring if enabled
	if n.MonitoringEnabled {
		fmt.Println(logSetupMonitoring)
		if err := n.setupMonitoring(ctx, sshKey, dropletIP); err != nil {
			return fmt.Errorf("%s: %w", errSetupMonitoring, err)
		}
	}

	// Step 9: Configure backups if enabled
	if n.BackupEnabled {
		fmt.Println(logSetupBackups)
		if err := n.setupBackups(ctx, sshKey, dropletIP); err != nil {
			return fmt.Errorf("%s: %w", errSetupBackups, err)
		}
	}

	fmt.Println(logSuccess)
	return nil
}

// GetURL returns the URL of the n8n instance
func (n *N8N) GetURL() string {
	return fmt.Sprintf("https://%s.%s", n.Subdomain, n.Domain)
}

// GetStatus returns the status of the n8n deployment
func (n *N8N) GetStatus(ctx context.Context, doToken *dagger.Secret) (string, error) {
	container := dag.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", doToken).
		WithExec([]string{
			"compute", "droplet", "get",
			fmt.Sprintf("%s-%s", n.Subdomain, n.Domain),
			"--format", "json",
		})

	return container.Stdout(ctx)
}

// Destroy removes the n8n deployment
func (n *N8N) Destroy(ctx context.Context, doToken *dagger.Secret) error {
	container := dag.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", doToken).
		WithExec([]string{
			"compute", "droplet", "delete",
			fmt.Sprintf("%s-%s", n.Subdomain, n.Domain),
			"--force",
		})

	_, err := container.Sync(ctx)
	return err
}

// setupRegistry configures the container registry
func (n *N8N) setupRegistry(ctx context.Context, registry *dagger.Container) error {
	// Create registry container with authentication
	container := registry.
		WithExposedPort(5000).
		WithExec([]string{
			"sh", "-c",
			"mkdir -p /auth && htpasswd -Bbn admin admin > /auth/htpasswd",
		})

	// Run the registry in the background
	_, err := container.WithExec([]string{"registry", "serve", "/etc/docker/registry/config.yml"}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to start registry: %w", err)
	}

	return nil
}

// pushImages builds and pushes required images
func (n *N8N) pushImages(ctx context.Context, registry *dagger.Container) error {
	// Base images needed for n8n deployment
	images := []struct {
		name string
		tag  string
	}{
		{name: "n8nio/n8n", tag: n.N8nVersion},
		{name: "postgres", tag: n.PostgresVersion},
		{name: "caddy", tag: "2.7.6-alpine"},
	}

	// Add cAdvisor if monitoring is enabled
	if n.MonitoringEnabled {
		images = append(images, struct{ name, tag string }{
			name: "gcr.io/cadvisor/cadvisor",
			tag:  "v0.47.2",
		})
	}

	// Pull and push each image
	for _, img := range images {
		// Pull the image
		container := dag.Container().From(fmt.Sprintf("%s:%s", img.name, img.tag))

		// Add standard labels
		container = container.
			WithLabel("org.opencontainers.image.created", time.Now().Format(time.RFC3339)).
			WithLabel("org.opencontainers.image.source", "https://github.com/felipepimentel/daggerverse").
			WithLabel("org.opencontainers.image.version", img.tag)

		// Add image-specific labels
		switch img.name {
		case "n8nio/n8n":
			container = container.
				WithLabel("org.opencontainers.image.title", "n8n Workflow Automation").
				WithLabel("org.opencontainers.image.description", "n8n is an extendable workflow automation tool")
		case "postgres":
			container = container.
				WithLabel("org.opencontainers.image.title", "PostgreSQL Database").
				WithLabel("org.opencontainers.image.description", "PostgreSQL database for n8n")
		case "caddy":
			container = container.
				WithLabel("org.opencontainers.image.title", "Caddy Web Server").
				WithLabel("org.opencontainers.image.description", "Caddy web server for SSL/TLS and reverse proxy")
		case "gcr.io/cadvisor/cadvisor":
			container = container.
				WithLabel("org.opencontainers.image.title", "cAdvisor").
				WithLabel("org.opencontainers.image.description", "Container monitoring and performance analysis")
		}

		// Push the image to our registry using the registry name from config
		targetRef := fmt.Sprintf("%s/%s:%s", n.RegistryName, img.name, img.tag)
		_, err := container.Publish(ctx, targetRef)
		if err != nil {
			return fmt.Errorf("failed to push image %s: %w", targetRef, err)
		}
	}

	return nil
}

// createDroplet creates a new DigitalOcean droplet
func (n *N8N) createDroplet(ctx context.Context, doToken *dagger.Secret) (string, error) {
	container := dag.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", doToken).
		WithExec([]string{
			"compute", "droplet", "create",
			fmt.Sprintf("%s-%s", n.Subdomain, n.Domain),
			"--region", n.Region,
			"--size", n.Size,
			"--image", "docker-20-04",
			"--ssh-keys", n.SSHKeyName,
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", err
	}

	return output, nil
}

// configureDNS sets up DNS records
func (n *N8N) configureDNS(ctx context.Context, doToken *dagger.Secret, dropletIP string) error {
	container := dag.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", doToken).
		WithExec([]string{
			"compute", "domain", "records", "create",
			n.Domain,
			"--record-type", "A",
			"--record-name", n.Subdomain,
			"--record-data", dropletIP,
		})

	_, err := container.Sync(ctx)
	return err
}

// setupServer prepares the server for deployment
func (n *N8N) setupServer(ctx context.Context, sshKey *dagger.Secret, dropletIP string) error {
	container := dag.Container().
		From("alpine:latest").
		WithSecretVariable("SSH_KEY", sshKey).
		WithExec([]string{"sh", "-c", "mkdir -p /root/.ssh && echo \"$SSH_KEY\" > /root/.ssh/id_rsa"}).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
		WithExec([]string{"ssh", "-o", "StrictHostKeyChecking=no", fmt.Sprintf("root@%s", dropletIP), "mkdir", "-p", "/opt/n8n"})

	_, err := container.Sync(ctx)
	return err
}

// deployServices deploys n8n and related services
func (n *N8N) deployServices(ctx context.Context, sshKey *dagger.Secret, dropletIP string) error {
	// Implementation here
	return nil
}

// setupMonitoring configures monitoring services
func (n *N8N) setupMonitoring(ctx context.Context, sshKey *dagger.Secret, dropletIP string) error {
	// Implementation here
	return nil
}

// setupBackups configures backup services
func (n *N8N) setupBackups(ctx context.Context, sshKey *dagger.Secret, dropletIP string) error {
	// Implementation here
	return nil
}
