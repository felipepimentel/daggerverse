// Package main provides a Dagger module for n8n CI/CD operations
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/felipepimentel/daggerverse/pipelines/n8n/internal/dagger"
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
	BackupEnabled bool
	BackupCron    string // e.g., "0 0 * * *"
	BackupRetention int  // Number of days to retain backups

	// Monitoring configuration
	MonitoringEnabled bool
	CAdvisorPort     int // e.g., 8080

	// SSL/TLS configuration
	SSLEmail string // Email for Let's Encrypt

	// Clients and secrets
	client   *dagger.Client
	doToken  *dagger.Secret
	sshKey   *dagger.Secret
}

// dropletInfo represents information about a DigitalOcean droplet
type dropletInfo struct {
	ID      string
	Name    string
	IPv4    string
	Status  string
}

// New creates a new instance of the N8N pipeline with default values
func New(domain, subdomain, sshKeyName string) *N8N {
	return &N8N{
		Domain:          domain,
		Subdomain:       subdomain,
		Region:         "nyc1",
		Size:           "s-2vcpu-4gb",
		SSHKeyName:     sshKeyName,
		RegistryName:   "n8n-registry",
		N8nVersion:     "0.234.0",
		N8nPort:        5678,
		PostgresVersion: "15-alpine",
		PostgresUser:    "n8n",
		PostgresPass:    "n8n",
		PostgresDB:     "n8n",
		BackupEnabled:  true,
		BackupCron:     "0 0 * * *",
		BackupRetention: 7,
		MonitoringEnabled: true,
		CAdvisorPort:     8080,
		client:         dag,
		doToken:        dag.SetSecret("do_token", ""),
		sshKey:         dag.SetSecret("ssh_key", ""),
	}
}

// WithRegion sets the DigitalOcean region
func (n *N8N) WithRegion(region string) *N8N {
	n.Region = region
	return n
}

// WithSize sets the droplet size
func (n *N8N) WithSize(size string) *N8N {
	n.Size = size
	return n
}

// WithN8NVersion sets the n8n version
func (n *N8N) WithN8NVersion(version string) *N8N {
	n.N8nVersion = version
	return n
}

// WithN8NPort sets the n8n port
func (n *N8N) WithN8NPort(port int) *N8N {
	n.N8nPort = port
	return n
}

// WithPostgresConfig sets the PostgreSQL configuration
func (n *N8N) WithPostgresConfig(version, user, pass, db string) *N8N {
	n.PostgresVersion = version
	n.PostgresUser = user
	n.PostgresPass = pass
	n.PostgresDB = db
	return n
}

// WithBackupConfig sets the backup configuration
func (n *N8N) WithBackupConfig(enabled bool, cron string) *N8N {
	n.BackupEnabled = enabled
	n.BackupCron = cron
	return n
}

// WithSSLEmail sets the SSL/TLS email for Let's Encrypt
func (n *N8N) WithSSLEmail(email string) *N8N {
	n.SSLEmail = email
	return n
}

// WithRegistryConfig sets the registry configuration
func (n *N8N) WithRegistryConfig(name string) *N8N {
	n.RegistryName = name
	return n
}

// GetStatus returns the status of the n8n deployment
func (n *N8N) GetStatus(ctx context.Context) (*dropletInfo, error) {
	container := n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{
			"compute", "droplet", "get",
			fmt.Sprintf("%s-%s", n.Subdomain, n.Domain),
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get droplet status: %w", err)
	}

	return &dropletInfo{
		Name:   fmt.Sprintf("%s-%s", n.Subdomain, n.Domain),
		Status: output,
	}, nil
}

// GetURL returns the URL of the n8n instance
func (n *N8N) GetURL() string {
	return fmt.Sprintf("https://%s.%s", n.Subdomain, n.Domain)
}

// Destroy removes the n8n deployment
func (n *N8N) Destroy(ctx context.Context) error {
	container := n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{
			"compute", "droplet", "delete",
			fmt.Sprintf("%s-%s", n.Subdomain, n.Domain),
			"--force",
		})

	_, err := container.Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to destroy droplet: %w", err)
	}

	// Remove DNS record
	container = n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{
			"compute", "domain", "records", "delete",
			n.Domain,
			"--record-name", n.Subdomain,
			"--force",
		})

	_, err = container.Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove DNS record: %w", err)
	}

	return nil
}

// Deploy orchestrates the complete n8n deployment pipeline
func (n *N8N) Deploy(ctx context.Context) error {
	// Step 1: Set up container registry
	if err := n.setupRegistry(ctx); err != nil {
		return fmt.Errorf("failed to setup registry: %w", err)
	}

	// Step 2: Push required images
	if err := n.pushImages(ctx); err != nil {
		return fmt.Errorf("failed to push images: %w", err)
	}

	// Step 3: Create droplet
	dropletIP, err := n.createDroplet(ctx)
	if err != nil {
		return fmt.Errorf("failed to create droplet: %w", err)
	}

	// Step 4: Configure DNS
	if err := n.configureDNS(ctx, dropletIP); err != nil {
		return fmt.Errorf("failed to configure DNS: %w", err)
	}

	// Step 5: Wait for DNS propagation
	time.Sleep(30 * time.Second)

	// Step 6: Create directories and copy configurations
	if err := n.setupServerDirectories(ctx, dropletIP); err != nil {
		return fmt.Errorf("failed to setup server directories: %w", err)
	}

	// Step 7: Deploy services
	if err := n.deployServices(ctx, dropletIP); err != nil {
		return fmt.Errorf("failed to deploy services: %w", err)
	}

	// Step 8: Configure monitoring if enabled
	if n.MonitoringEnabled {
		if err := n.setupMonitoring(ctx, dropletIP); err != nil {
			return fmt.Errorf("failed to setup monitoring: %w", err)
		}
	}

	// Step 9: Configure backups if enabled
	if n.BackupEnabled {
		if err := n.setupBackups(ctx, dropletIP); err != nil {
			return fmt.Errorf("failed to setup backups: %w", err)
		}
	}

	// Step 10: Verify deployment
	if err := n.verifyDeployment(ctx, dropletIP); err != nil {
		return fmt.Errorf("deployment verification failed: %w", err)
	}

	return nil
}

// setupServerDirectories creates necessary directories on the server
func (n *N8N) setupServerDirectories(ctx context.Context, serverIP string) error {
	dirs := []string{
		"/opt/n8n/data",
		"/opt/n8n/postgres",
		"/opt/n8n/caddy_data",
		"/opt/n8n/caddy_config",
		"/opt/n8n/backups",
	}

	for _, dir := range dirs {
		container := n.client.Container().
			From("alpine").
			WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
			WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
			WithExec([]string{
				"ssh", "-o", "StrictHostKeyChecking=no",
				fmt.Sprintf("root@%s", serverIP),
				fmt.Sprintf("mkdir -p %s", dir),
			})

		if _, err := container.Sync(ctx); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// setupMonitoring configures cAdvisor for container monitoring
func (n *N8N) setupMonitoring(ctx context.Context, serverIP string) error {
	container := n.client.Container().
		From("alpine").
		WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
		WithNewFile("/opt/n8n/docker-compose.monitoring.yml", fmt.Sprintf(`
version: '3.8'
services:
  cadvisor:
    image: gcr.io/cadvisor/cadvisor:latest
    container_name: cadvisor
    ports:
      - "%d:8080"
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:ro
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro
      - /dev/disk/:/dev/disk:ro
    restart: unless-stopped
`, n.CAdvisorPort)).
		WithExec([]string{
			"scp", "-o", "StrictHostKeyChecking=no",
			"/opt/n8n/docker-compose.monitoring.yml",
			fmt.Sprintf("root@%s:/opt/n8n/", serverIP),
		})

	if _, err := container.Sync(ctx); err != nil {
		return fmt.Errorf("failed to setup monitoring: %w", err)
	}

	// Start monitoring services
	container = n.client.Container().
		From("alpine").
		WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
		WithExec([]string{
			"ssh", "-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("root@%s", serverIP),
			"cd /opt/n8n && docker compose -f docker-compose.monitoring.yml up -d",
		})

	if _, err := container.Sync(ctx); err != nil {
		return fmt.Errorf("failed to start monitoring services: %w", err)
	}

	return nil
}

// setupBackups configures automated backups
func (n *N8N) setupBackups(ctx context.Context, serverIP string) error {
	// Create backup script
	backupScript := fmt.Sprintf(`#!/bin/bash
BACKUP_DIR="/opt/n8n/backups"
RETENTION_DAYS=%d

# Create backup
TIMESTAMP=$(date +%%Y%%m%%d_%%H%%M%%S)
docker compose -f /opt/n8n/docker-compose.yml exec -T postgres pg_dumpall -U n8n > $BACKUP_DIR/n8n_backup_$TIMESTAMP.sql

# Cleanup old backups
find $BACKUP_DIR -name "n8n_backup_*.sql" -type f -mtime +$RETENTION_DAYS -delete
`, n.BackupRetention)

	container := n.client.Container().
		From("alpine").
		WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
		WithNewFile("/opt/n8n/backup.sh", backupScript).
		WithExec([]string{
			"scp", "-o", "StrictHostKeyChecking=no",
			"/opt/n8n/backup.sh",
			fmt.Sprintf("root@%s:/opt/n8n/", serverIP),
		})

	if _, err := container.Sync(ctx); err != nil {
		return fmt.Errorf("failed to copy backup script: %w", err)
	}

	// Make script executable and setup cron
	container = n.client.Container().
		From("alpine").
		WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
		WithExec([]string{
			"ssh", "-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("root@%s", serverIP),
			fmt.Sprintf(`chmod +x /opt/n8n/backup.sh && 
(crontab -l 2>/dev/null; echo "%s /opt/n8n/backup.sh") | crontab -`, n.BackupCron),
		})

	if _, err := container.Sync(ctx); err != nil {
		return fmt.Errorf("failed to setup backup cron: %w", err)
	}

	return nil
}

// verifyDeployment checks if all services are running correctly
func (n *N8N) verifyDeployment(ctx context.Context, serverIP string) error {
	// Check container status
	container := n.client.Container().
		From("alpine").
		WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
		WithExec([]string{
			"ssh", "-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("root@%s", serverIP),
			"cd /opt/n8n && docker compose ps --format json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container status: %w", err)
	}

	if output == "" {
		return fmt.Errorf("no containers found running")
	}

	// Verify HTTPS access
	container = n.client.Container().
		From("alpine").
		WithExec([]string{
			"wget", "--spider", "--no-check-certificate",
			fmt.Sprintf("https://%s.%s", n.Subdomain, n.Domain),
		})

	if _, err := container.Sync(ctx); err != nil {
		return fmt.Errorf("HTTPS verification failed: %w", err)
	}

	return nil
}

// createDroplet provisions a new DigitalOcean droplet
func (n *N8N) createDroplet(ctx context.Context) (string, error) {
	container := n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{
			"compute", "droplet", "create",
			fmt.Sprintf("%s-%s", n.Subdomain, n.Domain),
			"--size", n.Size,
			"--image", "docker-20-04",
			"--region", n.Region,
			"--ssh-keys", n.SSHKeyName,
			"--enable-monitoring",
			"--enable-ipv6",
			"--tag-name", "n8n",
			"--wait",
			"--format", "json",
		})

	_, err := container.Sync(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create droplet: %w", err)
	}

	// Get droplet IP
	ipContainer := n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{
			"compute", "droplet", "get",
			fmt.Sprintf("%s-%s", n.Subdomain, n.Domain),
			"--format", "PublicIPv4",
			"--no-header",
		})

	ip, err := ipContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get droplet IP: %w", err)
	}

	return ip, nil
}

// configureDNS sets up DNS records for the n8n instance
func (n *N8N) configureDNS(ctx context.Context, ipAddress string) error {
	container := n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{
			"compute", "domain", "records", "create",
			n.Domain,
			"--record-type", "A",
			"--record-name", n.Subdomain,
			"--record-data", ipAddress,
			"--record-ttl", "3600",
		})

	_, err := container.Sync(ctx)
	return err
}

// deployServices sets up and starts all required services
func (n *N8N) deployServices(ctx context.Context, serverIP string) error {
	// Create necessary directories
	setupCmds := []string{
		"mkdir -p /opt/n8n/data",
		"mkdir -p /opt/n8n/postgres",
		"mkdir -p /opt/n8n/caddy_data",
		"mkdir -p /opt/n8n/caddy_config",
	}

	for _, cmd := range setupCmds {
		container := n.client.Container().
			From("alpine").
			WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
			WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
			WithExec([]string{"ssh", "-o", "StrictHostKeyChecking=no", fmt.Sprintf("root@%s", serverIP), cmd})

		if _, err := container.Sync(ctx); err != nil {
			return fmt.Errorf("failed to execute setup command %s: %w", cmd, err)
		}
	}

	// Generate and upload docker-compose.yml
	composeConfig := n.generateDockerCompose()
	container := n.client.Container().
		From("alpine").
		WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
		WithNewFile("/opt/n8n/docker-compose.yml", composeConfig).
		WithExec([]string{
			"scp", "-o", "StrictHostKeyChecking=no",
			"/opt/n8n/docker-compose.yml",
			fmt.Sprintf("root@%s:/opt/n8n/", serverIP),
		})

	if _, err := container.Sync(ctx); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}

	// Generate and upload Caddyfile
	caddyConfig := n.generateCaddyfile()
	container = n.client.Container().
		From("alpine").
		WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
		WithNewFile("/opt/n8n/Caddyfile", caddyConfig).
		WithExec([]string{
			"scp", "-o", "StrictHostKeyChecking=no",
			"/opt/n8n/Caddyfile",
			fmt.Sprintf("root@%s:/opt/n8n/", serverIP),
		})

	if _, err := container.Sync(ctx); err != nil {
		return fmt.Errorf("failed to write Caddyfile: %w", err)
	}

	// Start services
	startCmds := []string{
		"cd /opt/n8n",
		"docker compose pull",
		"docker compose up -d",
	}

	for _, cmd := range startCmds {
		container := n.client.Container().
			From("alpine").
			WithMountedSecret("/root/.ssh/id_rsa", n.sshKey).
			WithExec([]string{"chmod", "600", "/root/.ssh/id_rsa"}).
			WithExec([]string{"ssh", "-o", "StrictHostKeyChecking=no", fmt.Sprintf("root@%s", serverIP), cmd})

		if _, err := container.Sync(ctx); err != nil {
			return fmt.Errorf("failed to execute start command %s: %w", cmd, err)
		}
	}

	return nil
}

// generateDockerCompose creates the docker-compose.yml configuration
func (n *N8N) generateDockerCompose() string {
	return fmt.Sprintf(`version: '3.8'

services:
  n8n:
    image: n8nio/n8n:%s
    restart: always
    ports:
      - "%d:5678"
    environment:
      - DB_TYPE=postgresdb
      - DB_POSTGRESDB_HOST=postgres
      - DB_POSTGRESDB_DATABASE=%s
      - DB_POSTGRESDB_USER=%s
      - DB_POSTGRESDB_PASSWORD=%s
      - N8N_HOST=${n.Subdomain}.${n.Domain}
      - N8N_PROTOCOL=https
      - NODE_ENV=production
    volumes:
      - /opt/n8n/data:/home/node/.n8n
    depends_on:
      - postgres

  postgres:
    image: postgres:%s
    restart: always
    environment:
      - POSTGRES_USER=%s
      - POSTGRES_PASSWORD=%s
      - POSTGRES_DB=%s
    volumes:
      - /opt/n8n/postgres:/var/lib/postgresql/data

  caddy:
    image: caddy:2.7.6
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /opt/n8n/Caddyfile:/etc/caddy/Caddyfile
      - /opt/n8n/caddy_data:/data
      - /opt/n8n/caddy_config:/config
    environment:
      - EMAIL=%s
`, n.N8nVersion, n.N8nPort, n.PostgresDB, n.PostgresUser, n.PostgresPass,
		n.PostgresVersion, n.PostgresUser, n.PostgresPass, n.PostgresDB,
		n.SSLEmail)
}

// generateCaddyfile creates the Caddy reverse proxy configuration
func (n *N8N) generateCaddyfile() string {
	return fmt.Sprintf(`%s.%s {
    reverse_proxy n8n:%d
    tls %s
}`, n.Subdomain, n.Domain, n.N8nPort, n.SSLEmail)
}

// setupRegistry creates and configures the DigitalOcean container registry
func (n *N8N) setupRegistry(ctx context.Context) error {
	// Create registry
	container := n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{
			"registry", "create", n.RegistryName,
			"--subscription-tier", "basic",
		})

	_, err := container.Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	// Login to registry
	container = n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{"registry", "login"})

	_, err = container.Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to login to registry: %w", err)
	}

	return nil
}

// pushImages pushes n8n and Caddy images to the DigitalOcean registry
func (n *N8N) pushImages(ctx context.Context) error {
	images := []struct {
		source string
		target string
	}{
		{
			source: fmt.Sprintf("n8nio/n8n:%s", n.N8nVersion),
			target: fmt.Sprintf("registry.digitalocean.com/%s/n8n:%s", n.RegistryName, n.N8nVersion),
		},
		{
			source: "caddy:2.7.6",
			target: fmt.Sprintf("registry.digitalocean.com/%s/caddy:2.7.6", n.RegistryName),
		},
	}

	for _, img := range images {
		// Pull and push image using Docker client
		container := n.client.Container().
			From("docker:cli").
			WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
			WithExec([]string{"pull", img.source}).
			WithExec([]string{"tag", img.source, img.target}).
			WithExec([]string{"push", img.target})

		_, err := container.Sync(ctx)
		if err != nil {
			return fmt.Errorf("failed to process image %s: %w", img.source, err)
		}
	}

	return nil
}

// CI runs validation checks before deployment
func (n *N8N) CI(ctx context.Context) error {
	// Validate domain configuration
	if n.Domain == "" || n.Subdomain == "" {
		return fmt.Errorf("domain and subdomain must be set")
	}

	// Validate registry configuration
	if n.RegistryName == "" {
		return fmt.Errorf("registry name must be set")
	}

	// Validate SSH configuration
	if n.SSHKeyName == "" {
		return fmt.Errorf("SSH key name must be set")
	}

	// Validate n8n configuration
	if n.N8nVersion == "" {
		return fmt.Errorf("n8n version must be set")
	}

	// Validate PostgreSQL configuration
	if n.PostgresUser == "" || n.PostgresPass == "" || n.PostgresDB == "" {
		return fmt.Errorf("PostgreSQL configuration must be set")
	}

	// Validate SSL configuration
	if n.SSLEmail == "" {
		return fmt.Errorf("SSL email must be set for Let's Encrypt")
	}

	// Test image pulls
	images := []string{
		fmt.Sprintf("n8nio/n8n:%s", n.N8nVersion),
		"caddy:2.7.6",
		"postgres:" + n.PostgresVersion,
	}

	for _, img := range images {
		container := n.client.Container().From(img)
		if _, err := container.Sync(ctx); err != nil {
			return fmt.Errorf("failed to pull image %s: %w", img, err)
		}
	}

	return nil
}

// Cleanup removes all resources created by the deployment
func (n *N8N) Cleanup(ctx context.Context) error {
	// Step 1: Get droplet info
	info, err := n.GetStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to get droplet status: %w", err)
	}

	// Step 2: Remove DNS records
	container := n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{
			"compute", "domain", "records", "delete",
			n.Domain,
			"--record-name", n.Subdomain,
			"--force",
		})

	if _, err := container.Sync(ctx); err != nil {
		return fmt.Errorf("failed to remove DNS records: %w", err)
	}

	// Step 3: Delete droplet
	if info != nil {
		container = n.client.Container().
			From("digitalocean/doctl:latest").
			WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
			WithExec([]string{
				"compute", "droplet", "delete",
				fmt.Sprintf("%s-%s", n.Subdomain, n.Domain),
				"--force",
			})

		if _, err := container.Sync(ctx); err != nil {
			return fmt.Errorf("failed to delete droplet: %w", err)
		}
	}

	// Step 4: Delete registry
	container = n.client.Container().
		From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithExec([]string{
			"registry", "delete", n.RegistryName,
			"--force",
		})

	if _, err := container.Sync(ctx); err != nil {
		return fmt.Errorf("failed to delete registry: %w", err)
	}

	return nil
}
