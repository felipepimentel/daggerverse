// Package main provides a Dagger module for n8n CI/CD operations
package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/felipepimentel/daggerverse/pipelines/n8n/internal/dagger"
)

// N8N provides methods for building and deploying n8n
type N8N struct {
	// Source directory containing n8n configuration
	Source *dagger.Directory
	// Environment variables for n8n
	EnvVars []EnvVar
	// Domain for n8n instance
	Domain string
	// Registry to publish to
	Registry string
	// Image tag
	Tag string
	// Registry auth token
	RegistryAuth *dagger.Secret
	// SSH key for server access
	SSHKey *dagger.Secret
	// SSH key fingerprint
	SSHKeyFingerprint string
	// Server configuration
	Server *ServerConfig
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Name         string
	Region       string
	Size         string
	Image        string
	SSHKeyID     string
	Monitoring   bool
	IPv6Enabled  bool
	Tags         []string
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string
	Value string
}

// NewN8N creates a new N8N instance with default configuration
func NewN8N() *N8N {
	return &N8N{
		EnvVars: make([]EnvVar, 0),
		Server: &ServerConfig{
			Size:        "s-1vcpu-1gb",
			Image:       "docker-20-04",
			Monitoring:  true,
			IPv6Enabled: true,
			Tags:        []string{"n8n"},
		},
	}
}

// WithSource sets the source directory
func (n *N8N) WithSource(source *dagger.Directory) *N8N {
	n.Source = source
	return n
}

// WithEnvVar adds an environment variable
func (n *N8N) WithEnvVar(name, value string) *N8N {
	if n.EnvVars == nil {
		n.EnvVars = make([]EnvVar, 0)
	}
	n.EnvVars = append(n.EnvVars, EnvVar{Key: name, Value: value})
	return n
}

// WithEnvVars adds multiple environment variables
func (n *N8N) WithEnvVars(envVars []EnvVar) *N8N {
	if n.EnvVars == nil {
		n.EnvVars = make([]EnvVar, 0, len(envVars))
	}
	n.EnvVars = append(n.EnvVars, envVars...)
	return n
}

// WithDefaultEnvVars sets default environment variables for n8n
func (n *N8N) WithDefaultEnvVars() *N8N {
	defaultEnvs := []EnvVar{
		{Key: "N8N_PORT", Value: "5678"},
		{Key: "N8N_PROTOCOL", Value: "https"},
		{Key: "NODE_ENV", Value: "production"},
		{Key: "N8N_BASIC_AUTH_ACTIVE", Value: "true"},
		{Key: "N8N_DIAGNOSTICS_ENABLED", Value: "true"},
		{Key: "N8N_DIAGNOSTICS_CONFIG", Value: "true"},
		{Key: "N8N_USER_MANAGEMENT_DISABLED", Value: "false"},
	}

	if n.EnvVars == nil {
		n.EnvVars = make([]EnvVar, 0, len(defaultEnvs))
	}
	n.EnvVars = append(n.EnvVars, defaultEnvs...)
	return n
}

// WithDomain sets the domain
func (n *N8N) WithDomain(domain string) *N8N {
	n.Domain = domain
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

// WithSSHKey sets the SSH key
func (n *N8N) WithSSHKey(key *dagger.Secret) *N8N {
	n.SSHKey = key
	return n
}

// WithSSHKeyFingerprint sets the SSH key fingerprint
func (n *N8N) WithSSHKeyFingerprint(fingerprint string) *N8N {
	n.SSHKeyFingerprint = fingerprint
	return n
}

// WithServerConfig sets the server configuration
func (n *N8N) WithServerConfig(config *ServerConfig) *N8N {
	n.Server = config
	return n
}

// Setup prepares the environment for n8n deployment
func (n *N8N) Setup(ctx context.Context) error {
	// Get registry auth token
	token, err := n.RegistryAuth.Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registry auth token: %w", err)
	}

	// Create registry
	_, err = dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{"registry", "create", "n8n-registry"}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	// Pull and tag images
	n8nImage := fmt.Sprintf("n8nio/n8n:%s", n.Tag)
	registryN8N := fmt.Sprintf("%s/n8n:%s", n.Registry, n.Tag)
	registryCaddy := fmt.Sprintf("%s/caddy:2.7.6", n.Registry)

	// Pull n8n
	_, err = dag.Container().From(n8nImage).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to pull n8n image: %w", err)
	}

	// Pull Caddy
	_, err = dag.Container().From("caddy:2.7.6").Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to pull caddy image: %w", err)
	}

	// Tag and push images
	container := dag.Container().
		From("docker:cli").
		WithRegistryAuth(n.Registry, "", n.RegistryAuth)

	// Tag and push n8n
	_, err = container.
		WithExec([]string{"tag", n8nImage, registryN8N}).
		WithExec([]string{"push", registryN8N}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to tag and push n8n image: %w", err)
	}

	// Tag and push Caddy
	_, err = container.
		WithExec([]string{"tag", "caddy:2.7.6", registryCaddy}).
		WithExec([]string{"push", registryCaddy}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to tag and push caddy image: %w", err)
	}

	return nil
}

// Deploy deploys n8n to DigitalOcean
func (n *N8N) Deploy(ctx context.Context) error {
	// Get registry auth token
	token, err := n.RegistryAuth.Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registry auth token: %w", err)
	}

	// Create droplet
	_, err = dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{
			"compute", "droplet", "create",
			n.Server.Name,
			"--size", n.Server.Size,
			"--image", n.Server.Image,
			"--region", n.Server.Region,
			"--ssh-keys", n.SSHKeyFingerprint,
			"--enable-monitoring",
			"--enable-ipv6",
			"--tag-name", "n8n",
		}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create droplet: %w", err)
	}

	// Get droplet IP
	output, err := dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{
			"compute", "droplet", "get",
			n.Server.Name,
			"--format", "PublicIPv4",
			"--no-header",
		}).
		Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get droplet IP: %w", err)
	}

	// Create DNS record
	_, err = dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{
			"compute", "domain", "records", "create",
			filepath.Dir(n.Domain),
			"--record-type", "A",
			"--record-name", filepath.Base(n.Domain),
			"--record-data", output,
		}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create DNS record: %w", err)
	}

	// Create configuration files
	err = n.createConfig(ctx, output)
	if err != nil {
		return fmt.Errorf("failed to create configuration: %w", err)
	}

	// Deploy services
	err = n.deployServices(ctx, output)
	if err != nil {
		return fmt.Errorf("failed to deploy services: %w", err)
	}

	return nil
}

// createConfig creates the necessary configuration files on the server
func (n *N8N) createConfig(ctx context.Context, serverIP string) error {
	// Create directories
	_, err := dag.Container().
		From("alpine").
		WithExec([]string{"mkdir", "-p", "/opt/n8n/caddy_config"}).
		WithExec([]string{"mkdir", "-p", "/opt/n8n/data"}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Create docker-compose.yml
	composeConfig := fmt.Sprintf(`version: '3.8'

services:
  n8n:
    image: %s/n8n:%s
    restart: always
    ports:
      - "5678:5678"
    environment:
      - N8N_HOST=%s
      - N8N_PORT=5678
      - N8N_PROTOCOL=https
%s
    volumes:
      - /opt/n8n/data:/home/node/.n8n
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:5678/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3

  caddy:
    image: %s/caddy:2.7.6
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /opt/n8n/caddy_config:/etc/caddy
      - caddy_data:/data
      - caddy_config:/config
    depends_on:
      - n8n
    healthcheck:
      test: ["CMD", "caddy", "version"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  caddy_data:
  caddy_config:
`, n.Registry, n.Tag, n.Domain, n.formatEnvVars(), n.Registry)

	// Create Caddyfile
	var caddyEmail string
	for _, env := range n.EnvVars {
		if env.Key == "CADDY_ACME_EMAIL" {
			caddyEmail = env.Value
			break
		}
	}

	caddyConfig := fmt.Sprintf(`{
    email %s
}

%s {
    reverse_proxy n8n:5678 {
        flush_interval -1
        health_timeout 5s
    }
    log {
        output file /var/log/caddy/access.log
        format json
    }
    tls {
        protocols tls1.2 tls1.3
    }
}
`, caddyEmail, n.Domain)

	// Get SSH key plaintext
	sshKey, err := n.SSHKey.Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get SSH key: %w", err)
	}

	// Write configuration files to server
	container := dag.Container().From("alpine")
	container = container.WithNewFile("/opt/n8n/docker-compose.yml", composeConfig)
	container = container.WithNewFile("/opt/n8n/caddy_config/Caddyfile", caddyConfig)

	_, err = container.
		WithExec([]string{"scp", "-i", sshKey,
			"/opt/n8n/docker-compose.yml",
			"/opt/n8n/caddy_config/Caddyfile",
			fmt.Sprintf("root@%s:/opt/n8n/", serverIP),
		}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to write configuration files: %w", err)
	}

	return nil
}

// deployServices deploys the services on the server
func (n *N8N) deployServices(ctx context.Context, serverIP string) error {
	// Get SSH key plaintext
	sshKey, err := n.SSHKey.Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get SSH key: %w", err)
	}

	// Deploy services using docker compose
	_, err = dag.Container().
		From("alpine").
		WithExec([]string{"ssh", "-i", sshKey,
			fmt.Sprintf("root@%s", serverIP),
			"cd /opt/n8n && docker compose pull && docker compose up -d",
		}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to deploy services: %w", err)
	}

	return nil
}

// formatEnvVars formats environment variables for docker-compose
func (n *N8N) formatEnvVars() string {
	var envVars string
	for _, env := range n.EnvVars {
		envVars += fmt.Sprintf("      - %s=%s\n", env.Key, env.Value)
	}
	return envVars
}

// GetStatus returns the status of the n8n deployment
func (n *N8N) GetStatus(ctx context.Context) error {
	// Get registry auth token
	token, err := n.RegistryAuth.Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registry auth token: %w", err)
	}

	_, err = dag.Container().
		From("digitalocean/doctl:latest").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", token).
		WithExec([]string{
			"compute", "droplet", "get",
			n.Server.Name,
			"--format", "ID,Name,PublicIPv4,Status",
		}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	return nil
}
