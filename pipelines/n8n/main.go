package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/felipepimentel/daggerverse/pipelines/n8n/internal/dagger"
)

// N8N represents a module for deploying N8N to DigitalOcean
type N8N struct {
	doToken   *dagger.Secret
	Domain    string
	Subdomain string
	Region    string
	Size      string
	Image     string
}



// New creates a new N8N module with default values
func New() *N8N {
	return &N8N{
		Domain:    "pepper88.com",
		Subdomain: "n8n",
		Region:    "nyc1",
		Size:      "s-2vcpu-2gb",
		Image:     "ubuntu-20-04-x64",
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

// WithImage sets the droplet image
func (n *N8N) WithImage(image string) *N8N {
	n.Image = image
	return n
}

// Deploy deploys n8n to DigitalOcean
func (n *N8N) Deploy(ctx context.Context, doToken string) (string, error) {
	// If the token is an environment variable reference, get the actual value
	if strings.HasPrefix(doToken, "env:") {
		envVar := strings.TrimPrefix(doToken, "env:")
		envValue := os.Getenv(envVar)
		if envValue == "" {
			return "", fmt.Errorf("environment variable %s is not set", envVar)
		}
		doToken = envValue
	}

	// Store the token as a secret for use in other methods
	n.doToken = dag.SetSecret("do-token", doToken)

	fmt.Println("🚀 Starting n8n deployment...")

	// Clean up old resources
	fmt.Println("🧹 Cleaning up old resources...")
	if err := n.cleanupOldResources(ctx); err != nil {
		return "", fmt.Errorf("failed to clean up old resources: %w", err)
	}

	// Generate SSH keys
	fmt.Println("🔑 Generating SSH keys...")
	sshKeys, err := n.Ssh().GenerateKeys()
	if err != nil {
		return "", fmt.Errorf("failed to generate SSH keys: %w", err)
	}

	// Register SSH key with DigitalOcean
	fmt.Println("📝 Registering SSH key with DigitalOcean...")
	keyID, err := n.DigitalOcean().RegisterSSHKey(ctx, sshKeys.name, sshKeys.publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to register SSH key: %w", err)
	}
	sshKeys.keyID = keyID
	fmt.Printf("✅ SSH key registered with ID: %s\n", keyID)

	// Create droplet using the DigitalOcean module
	fmt.Println("🚀 Creating droplet...")
	dag.DigitalOcean().CreateDroplet(ctx, "n8n", n.Region, n.Size, n.Image, keyID, n.getUserData())

	dag.DigitalOcean().CreateDNSRecord(ctx, n.Domain, n.Subdomain, n.Region, n.Size, n.Image, keyID, n.getUserData())

	// Wait for SSH to be available
	fmt.Println("⏳ Waiting for SSH to be available...")
	ssh := dag.DigitalOcean().WaitForDroplet(ctx, "n8n", "active", 180*time.Second)
}

func (n *N8N) cleanupOldResources(ctx context.Context) error {
	fmt.Println("🧹 Cleaning up old resources...")
	do := dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithWorkdir("/app").
		WithExec([]string{"ln", "-s", "/app/doctl", "/usr/local/bin/doctl"})

	// List and delete old droplets
	output, err := do.WithExec([]string{
		"doctl", "compute", "droplet", "list",
		"--format", "ID,Name",
		"--no-header",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to list droplets: %w", err)
	}

	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		dropletID := parts[0]
		dropletName := parts[1]

		if dropletName == "n8n" {
			fmt.Printf("🗑️ Deleting old droplet: %s (%s)\n", dropletName, dropletID)
			_, err := do.WithExec([]string{
				"doctl", "compute", "droplet", "delete",
				dropletID,
				"--force",
			}).Sync(ctx)
			if err != nil {
				return fmt.Errorf("failed to delete droplet %s: %w", dropletID, err)
			}
		}
	}

	return nil
}



func (n *N8N) registerSSHKey(ctx context.Context, name, publicKey string) (string, error) {
	fmt.Println("📝 Registering SSH key with DigitalOcean...")
	do := dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", n.doToken).
		WithWorkdir("/app").
		WithExec([]string{"ln", "-s", "/app/doctl", "/usr/local/bin/doctl"})

	output, err := do.WithExec([]string{
		"doctl", "compute", "ssh-key", "create",
		name,
		"--public-key", publicKey,
		"--format", "ID",
		"--no-header",
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to register SSH key: %w", err)
	}

	return strings.TrimSpace(output), nil
}

func (n *N8N) getUserData() string {
	return `#!/bin/bash
set -euxo pipefail

# Wait for cloud-init to complete
cloud-init status --wait

# Configure system
echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
echo 'debconf debconf/priority select critical' | debconf-set-selections

# Install required packages
apt-get update
DEBIAN_FRONTEND=noninteractive apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    openssh-server \
    ufw

# Configure SSH
mkdir -p /root/.ssh
chmod 700 /root/.ssh
sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin prohibit-password/' /etc/ssh/sshd_config
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
systemctl enable ssh
systemctl restart ssh

# Configure firewall
ufw allow OpenSSH
ufw allow 5678/tcp  # n8n
ufw --force enable

# Install Docker
curl -fsSL https://get.docker.com | sh
systemctl enable docker
systemctl start docker
DEBIAN_FRONTEND=noninteractive apt-get install -y docker-compose-plugin

# Create n8n directory
mkdir -p /opt/n8n
chown root:root /opt/n8n
chmod 755 /opt/n8n

# Signal completion
touch /root/.cloud-init-complete
sync`
}

func (n *N8N) createConfigFiles(ctx context.Context, dropletIP string, doToken *dagger.Secret) error {
	fmt.Println("📝 Creating configuration files...")

	ssh := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "openssh-client"}).
		WithNewFile("/root/.ssh/id_ed25519", n.sshKeys.privateKey+"\n", dagger.ContainerWithNewFileOpts{
			Permissions: 0600,
		}).
		WithExec([]string{"chmod", "700", "/root/.ssh"}).
		WithExec([]string{"sh", "-c", "echo 'Host *\n\tStrictHostKeyChecking no\n\tUserKnownHostsFile /dev/null' > /root/.ssh/config"}).
		WithExec([]string{"chmod", "600", "/root/.ssh/config"})

	// Create directory structure
	_, err := ssh.WithExec([]string{
		"ssh",
		"-i", "/root/.ssh/id_ed25519",
		fmt.Sprintf("root@%s", dropletIP),
		"mkdir -p /opt/n8n && chmod 755 /opt/n8n",
	}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Create configuration files
	files := map[string]string{
		"docker-compose.yml": n.getDockerComposeContent(),
		".env":              n.getEnvContent(),
		"Caddyfile":         n.getCaddyfileContent(),
	}

	for filename, content := range files {
		fmt.Printf("📝 Creating %s...\n", filename)
		mode := "644"
		if filename == ".env" {
			mode = "600"
		}

		tempFile := fmt.Sprintf("/tmp/%s", filename)
		ssh = ssh.WithNewFile(tempFile, content)

		_, err := ssh.WithExec([]string{
			"scp",
			"-i", "/root/.ssh/id_ed25519",
			tempFile,
			fmt.Sprintf("root@%s:/opt/n8n/%s", dropletIP, filename),
		}).Sync(ctx)
		if err != nil {
			return fmt.Errorf("failed to copy %s: %w", filename, err)
		}

		_, err = ssh.WithExec([]string{
			"ssh",
			"-i", "/root/.ssh/id_ed25519",
			fmt.Sprintf("root@%s", dropletIP),
			fmt.Sprintf("chmod %s /opt/n8n/%s", mode, filename),
		}).Sync(ctx)
		if err != nil {
			return fmt.Errorf("failed to set permissions for %s: %w", filename, err)
		}
	}

	fmt.Println("✅ Configuration files created successfully")
	return nil
}

func (n *N8N) getDockerComposeContent() string {
	return `version: '3.8'

services:
  n8n:
    image: n8nio/n8n:latest
    restart: always
    ports:
      - "127.0.0.1:5678:5678"
    env_file:
      - .env
    volumes:
      - n8n_data:/home/node/.n8n
    networks:
      - n8n-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:5678/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 30s

  caddy:
    image: caddy:2.7.6
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy_data:/data
      - caddy_config:/config
    networks:
      - n8n-network
    depends_on:
      - n8n
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:80"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 30s

volumes:
  n8n_data:
  caddy_data:
  caddy_config:

networks:
  n8n-network:
    driver: bridge`
}

func (n *N8N) getEnvContent() string {
	return fmt.Sprintf(`# N8N Configuration
N8N_HOST=%s.%s
N8N_PORT=5678
N8N_PROTOCOL=https
NODE_ENV=production
WEBHOOK_URL=https://%s.%s/
GENERIC_TIMEZONE=America/Sao_Paulo

# Security Settings
N8N_BASIC_AUTH_ACTIVE=true
N8N_BASIC_AUTH_USER=admin
N8N_BASIC_AUTH_PASSWORD=admin123
N8N_ENCRYPTION_KEY=%s`, n.Subdomain, n.Domain, n.Subdomain, n.Domain, generateRandomString(32))
}

func (n *N8N) getCaddyfileContent() string {
	return fmt.Sprintf(`%s.%s {
    # Use HTTPS with automatic certificate management
    tls internal

    # Enable Gzip compression
    encode gzip

    # Reverse proxy to n8n with improved settings
    reverse_proxy n8n:5678 {
        # Enable WebSocket support
        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-For {remote_host}
        header_up X-Forwarded-Proto {scheme}
        header_up X-Forwarded-Host {host}

        # Timeouts
        flush_interval -1
        transport http {
            keepalive 30s
            keepalive_idle_conns 10
        }
    }

    # Security headers
    header {
        # Enable HSTS
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        # Disable FLoC tracking
        Permissions-Policy "interest-cohort=()"
        # XSS protection
        X-XSS-Protection "1; mode=block"
        # Prevent clickjacking
        X-Frame-Options "SAMEORIGIN"
        # Prevent MIME type sniffing
        X-Content-Type-Options "nosniff"
        # Referrer policy
        Referrer-Policy "strict-origin-when-cross-origin"
    }

    # Basic logging
    log {
        output file /data/access.log {
            roll_size 10MB
            roll_keep 10
        }
    }
}`, n.Subdomain, n.Domain)
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // Ensure unique values
	}
	return string(result)
}