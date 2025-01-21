package n8n

import (
	"context"
	"fmt"
	"strings"
	"encoding/base64"

	"github.com/felipepimentel/daggerverse/pipelines/n8n/internal/dagger"
)

var _ = dagger.Container{}

type N8n struct {
	// Base Configuration
	domain     string
	subdomain  string
	region     string
	size       string
	sshKeyName string

	// n8n Configuration
	n8nVersion string
	n8nPort    int

	// Database Configuration
	postgresVersion string
	postgresUser    string
	postgresPass    string
	postgresDB      string

	// Backup Configuration
	backupEnabled    bool
	backupCron       string
	backupRetention  int

	// Monitoring Configuration
	monitoringEnabled bool
	cAdvisorPort     int

	// SSL Configuration
	sslEmail string
}

func New(domain, subdomain, sshKeyName string) *N8n {
	return &N8n{
		domain:           domain,
		subdomain:        subdomain,
		sshKeyName:      sshKeyName,
		region:          "nyc1",
		size:            "s-2vcpu-4gb",
		n8nVersion:      "0.234.0",
		n8nPort:         5678,
		postgresVersion: "15-alpine",
		postgresUser:    "n8n",
		postgresPass:    "n8n",
		postgresDB:      "n8n",
		backupEnabled:   true,
		backupCron:      "0 0 * * *",
		backupRetention: 7,
		monitoringEnabled: true,
		cAdvisorPort:     8080,
	}
}

// Configuration methods
func (n *N8n) WithRegion(region string) *N8n {
	n.region = region
	return n
}

func (n *N8n) WithSize(size string) *N8n {
	n.size = size
	return n
}

func (n *N8n) WithN8NVersion(version string) *N8n {
	n.n8nVersion = version
	return n
}

func (n *N8n) WithN8NPort(port int) *N8n {
	n.n8nPort = port
	return n
}

func (n *N8n) WithPostgresConfig(version, user, pass, db string) *N8n {
	n.postgresVersion = version
	n.postgresUser = user
	n.postgresPass = pass
	n.postgresDB = db
	return n
}

func (n *N8n) WithBackupConfig(enabled bool, cron string, retention int) *N8n {
	n.backupEnabled = enabled
	n.backupCron = cron
	n.backupRetention = retention
	return n
}

func (n *N8n) WithMonitoringConfig(enabled bool, cAdvisorPort int) *N8n {
	n.monitoringEnabled = enabled
	n.cAdvisorPort = cAdvisorPort
	return n
}

func (n *N8n) WithSSLEmail(email string) *N8n {
	n.sslEmail = email
	return n
}

// Helper methods
func (n *N8n) GetURL() string {
	return fmt.Sprintf("https://%s.%s", n.subdomain, n.domain)
}

// getExistingDroplet checks if a droplet with the given name already exists
func (n *N8n) getExistingDroplet(ctx context.Context, do *dagger.Container, dropletName string) (string, error) {
	dropletOutput, err := do.WithExec([]string{
		"/app/doctl", "compute", "droplet", "list",
		"--format", "ID,Name",
		"--no-header",
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list droplets: %w", err)
	}

	for _, line := range strings.Split(dropletOutput, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.Join(fields[1:], " ") == dropletName {
			return fields[0], nil
		}
	}

	return "", nil
}

func (n *N8n) Deploy(ctx context.Context, token string, sshPrivateKey string) error {
	// Set up DigitalOcean client
	doToken := dag.SetSecret("do_token", token)
	do := dag.Container().From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", doToken).
		WithDefaultArgs([]string{"/app/doctl"})

	// Get SSH key ID
	sshKeyIDOutput, err := do.WithExec([]string{
		"/app/doctl", "compute", "ssh-key", "list",
		"--format", "ID,Name",
		"--no-header",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get SSH key ID: %w", err)
	}

	// Find the SSH key ID by name
	var sshKeyID string
	for _, line := range strings.Split(sshKeyIDOutput, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.Join(fields[1:], " ") == n.sshKeyName {
			sshKeyID = fields[0]
			break
		}
	}
	if sshKeyID == "" {
		return fmt.Errorf("SSH key with name %q not found", n.sshKeyName)
	}
	fmt.Printf("Debug: Using SSH Key ID: %s\n", sshKeyID)

	// Check for existing droplet
	dropletName := fmt.Sprintf("%s-n8n", n.subdomain)
	existingDropletID, err := n.getExistingDroplet(ctx, do, dropletName)
	if err != nil {
		return fmt.Errorf("failed to check for existing droplet: %w", err)
	}

	var dropletID string
	if existingDropletID != "" {
		fmt.Printf("Found existing droplet with ID: %s\n", existingDropletID)
		dropletID = existingDropletID
	} else {
		// Create new droplet
		dropletOutput, err := do.WithExec([]string{
			"/app/doctl", "compute", "droplet", "create",
			"--region", n.region,
			"--size", n.size,
			"--image", "ubuntu-20-04-x64",
			"--ssh-keys", sshKeyID,
			"--format", "ID",
			"--no-header",
			dropletName,
		}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("failed to create droplet: %w", err)
		}
		dropletID = strings.TrimSpace(dropletOutput)
		fmt.Printf("Created new droplet with ID: %s\n", dropletID)
	}

	// Wait for droplet to be ready
	_, err = do.WithExec([]string{
		"/app/doctl", "compute", "droplet", "wait",
		dropletID,
	}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for droplet: %w", err)
	}

	// Additional wait to ensure SSH is ready
	_, err = do.WithExec([]string{"sleep", "30"}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to sleep: %w", err)
	}

	// Get droplet status
	statusOutput, err := do.WithExec([]string{
		"/app/doctl", "compute", "droplet", "get",
		dropletID,
		"--format", "Status",
		"--no-header",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get droplet status: %w", err)
	}
	status := strings.TrimSpace(statusOutput)
	fmt.Printf("Debug: Droplet status is '%s'\n", status)

	if status != "active" {
		return fmt.Errorf("droplet is not active, current status: %s", status)
	}

	// Get droplet IP
	ipOutput, err := do.WithExec([]string{
		"/app/doctl", "compute", "droplet", "get",
		dropletID,
		"--format", "PublicIPv4",
		"--no-header",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get droplet IP: %w", err)
	}
	fmt.Printf("Debug: Raw IP output is '%s'\n", ipOutput)
	ipAddress := strings.TrimSpace(ipOutput)
	fmt.Printf("Debug: Trimmed IP address is '%s'\n", ipAddress)

	if ipAddress == "" {
		return fmt.Errorf("failed to get droplet IP: IP address is empty")
	}

	// Check if domain exists
	_, err = do.WithExec([]string{
		"/app/doctl", "compute", "domain", "get",
		n.domain,
	}).Sync(ctx)
	if err != nil {
		fmt.Printf("Domain does not exist, creating it...\n")
		_, err = do.WithExec([]string{
			"/app/doctl", "compute", "domain", "create",
			n.domain,
		}).Sync(ctx)
		if err != nil {
			return fmt.Errorf("failed to create domain: %w", err)
		}
	}

	// Create DNS record
	_, err = do.WithExec([]string{
		"/app/doctl", "compute", "domain", "records", "create",
		n.domain,
		"--record-type", "A",
		"--record-name", n.subdomain,
		"--record-data", ipAddress,
	}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create DNS record: %w", err)
	}

	// Prepare docker-compose.yml content
	composeContent := fmt.Sprintf(`version: "3.8"

services:
  n8n:
    image: n8nio/n8n:%s
    restart: always
    ports:
      - "%d:%d"
    environment:
      - N8N_HOST=%s.%s
      - N8N_PROTOCOL=https
      - N8N_PORT=%d
      - N8N_ENCRYPTION_KEY=your-secret-key
      - WEBHOOK_TUNNEL_URL=https://%s.%s/
      - DB_TYPE=postgresdb
      - DB_POSTGRESDB_HOST=postgres
      - DB_POSTGRESDB_DATABASE=%s
      - DB_POSTGRESDB_USER=%s
      - DB_POSTGRESDB_PASSWORD=%s
    volumes:
      - n8n_data:/home/node/.n8n
    networks:
      - n8n-network
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
      - postgres_data:/var/lib/postgresql/data
    networks:
      - n8n-network

  caddy:
    image: caddy:2
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data
      - caddy_config:/config
    networks:
      - n8n-network%s

volumes:
  n8n_data:
  postgres_data:
  caddy_data:
  caddy_config:%s

networks:
  n8n-network:
    driver: bridge`,
		n.n8nVersion,
		n.n8nPort, n.n8nPort,
		n.subdomain, n.domain,
		n.n8nPort,
		n.subdomain, n.domain,
		n.postgresDB,
		n.postgresUser,
		n.postgresPass,
		n.postgresVersion,
		n.postgresUser,
		n.postgresPass,
		n.postgresDB,
		n.getMonitoringService(),
		n.getMonitoringVolumes(),
	)

	// Prepare Caddyfile content
	caddyContent := fmt.Sprintf(`%s.%s {
    reverse_proxy n8n:%d
}`, n.subdomain, n.domain, n.n8nPort)

	fmt.Printf("Debug: docker-compose.yml content:\n%s\n", composeContent)
	fmt.Printf("Debug: Caddyfile content:\n%s\n", caddyContent)

	// Set up SSH client
	sshContainer := dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "openssh-client"}).
		WithExec([]string{"mkdir", "-p", "/root/.ssh"})

	// Decode SSH private key if it's base64 encoded
	decodedKey := sshPrivateKey
	if !strings.HasPrefix(sshPrivateKey, "-----BEGIN") {
		decoded, err := base64.StdEncoding.DecodeString(sshPrivateKey)
		if err != nil {
			return fmt.Errorf("failed to decode SSH private key: %w", err)
		}
		decodedKey = string(decoded)
	}

	// Ensure key has proper line endings
	decodedKey = strings.ReplaceAll(decodedKey, "\r\n", "\n")
	if !strings.HasSuffix(decodedKey, "\n") {
		decodedKey += "\n"
	}

	// Write SSH key file
	sshKeyFile := dag.Directory().
		WithNewFile("/root/.ssh/id_rsa", decodedKey, dagger.DirectoryWithNewFileOpts{
			Permissions: 0600,
		})

	sshContainer = sshContainer.
		WithMountedDirectory("/root/.ssh", sshKeyFile).
		WithExec([]string{"sh", "-c", fmt.Sprintf("ssh-keyscan -H %s >> /root/.ssh/known_hosts 2>/dev/null", ipAddress)})

	// Create directories and write files
	sshCmd := fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i /root/.ssh/id_rsa root@%s 'mkdir -p /root/n8n && echo %q > /root/n8n/docker-compose.yml && echo %q > /root/n8n/Caddyfile'",
		ipAddress, composeContent, caddyContent)
	_, err = sshContainer.WithExec([]string{"sh", "-c", sshCmd}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create directories and files: %w", err)
	}

	// Install Docker and Docker Compose
	sshCmd = fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i /root/.ssh/id_rsa root@%s 'while pgrep apt-get > /dev/null; do sleep 5; done && curl -fsSL https://get.docker.com | sh && curl -L https://github.com/docker/compose/releases/download/v2.5.0/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose'",
		ipAddress)
	_, err = sshContainer.WithExec([]string{"sh", "-c", sshCmd}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to install Docker and Docker Compose: %w", err)
	}

	// Set up backup if enabled
	if n.backupEnabled {
		backupScript := n.getBackupScript()
		sshCmd = fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i /root/.ssh/id_rsa root@%s 'echo %q > /root/n8n/backup.sh && chmod +x /root/n8n/backup.sh && (crontab -l 2>/dev/null; echo \"%s /root/n8n/backup.sh\") | crontab -'",
			ipAddress, backupScript, n.backupCron)
		_, err = sshContainer.WithExec([]string{"sh", "-c", sshCmd}).Sync(ctx)
		if err != nil {
			return fmt.Errorf("failed to set up backup: %w", err)
		}
	}

	// Start services
	sshCmd = fmt.Sprintf("ssh -o StrictHostKeyChecking=no -i /root/.ssh/id_rsa root@%s 'cd /root/n8n && docker-compose up -d'",
		ipAddress)
	_, err = sshContainer.WithExec([]string{"sh", "-c", sshCmd}).Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	fmt.Printf("Successfully deployed n8n to %s\n", n.GetURL())
	return nil
}

func (n *N8n) getMonitoringService() string {
	if !n.monitoringEnabled {
		return ""
	}

	return fmt.Sprintf(`

  cadvisor:
    image: gcr.io/cadvisor/cadvisor:latest
    restart: always
    privileged: true
    ports:
      - "%d:8080"
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:ro
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro
      - /dev/disk/:/dev/disk:ro
    networks:
      - n8n-network`, n.cAdvisorPort)
}

func (n *N8n) getMonitoringVolumes() string {
	if !n.monitoringEnabled {
		return ""
	}

	return `
  cadvisor_data:`
}

func (n *N8n) getBackupScript() string {
	return fmt.Sprintf(`#!/bin/bash
set -e

# Configuration
BACKUP_DIR="/root/n8n/backups"
RETENTION_DAYS=%d

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Backup timestamp
TIMESTAMP=$(date +%%Y%%m%%d_%%H%%M%%S)

# Backup PostgreSQL database
docker exec n8n-postgres-%s pg_dump -U %s %s > "$BACKUP_DIR/n8n_db_$TIMESTAMP.sql"

# Backup n8n data
tar -czf "$BACKUP_DIR/n8n_data_$TIMESTAMP.tar.gz" -C /root/n8n n8n_data

# Clean up old backups
find "$BACKUP_DIR" -type f -mtime +$RETENTION_DAYS -delete
`, n.backupRetention, n.postgresDB, n.postgresUser, n.postgresDB)
}

// CI performs validation checks before deployment
func (n *N8n) CI(ctx context.Context) error {
	// Validate base configuration
	if n.domain == "" {
		return fmt.Errorf("domain is required")
	}
	if n.subdomain == "" {
		return fmt.Errorf("subdomain is required")
	}
	if n.sshKeyName == "" {
		return fmt.Errorf("SSH key name is required")
	}

	// Validate n8n configuration
	if n.n8nVersion == "" {
		return fmt.Errorf("n8n version is required")
	}
	if n.n8nPort <= 0 {
		return fmt.Errorf("invalid n8n port: %d", n.n8nPort)
	}

	// Validate database configuration
	if n.postgresVersion == "" {
		return fmt.Errorf("PostgreSQL version is required")
	}
	if n.postgresUser == "" {
		return fmt.Errorf("PostgreSQL user is required")
	}
	if n.postgresPass == "" {
		return fmt.Errorf("PostgreSQL password is required")
	}
	if n.postgresDB == "" {
		return fmt.Errorf("PostgreSQL database name is required")
	}

	// Validate backup configuration
	if n.backupEnabled {
		if n.backupCron == "" {
			return fmt.Errorf("backup cron schedule is required when backup is enabled")
		}
		if n.backupRetention <= 0 {
			return fmt.Errorf("invalid backup retention days: %d", n.backupRetention)
		}
	}

	// Test pulling required images
	images := []string{
		fmt.Sprintf("n8nio/n8n:%s", n.n8nVersion),
		fmt.Sprintf("postgres:%s", n.postgresVersion),
		"caddy:2",
	}
	if n.monitoringEnabled {
		images = append(images, "gcr.io/cadvisor/cadvisor:latest")
	}

	for _, image := range images {
		container := dag.Container().From(image)
		_, err := container.ID(ctx)
		if err != nil {
			return fmt.Errorf("failed to pull image %s: %w", image, err)
		}
	}

	return nil
}

// Cleanup removes the n8n deployment and associated resources
func (n *N8n) Cleanup(ctx context.Context, token string) error {
	// Set up DigitalOcean client
	doToken := dag.SetSecret("do_token", token)
	do := dag.Container().From("digitalocean/doctl:latest").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", doToken).
		WithDefaultArgs([]string{"/app/doctl"})

	// Get droplet ID
	dropletOutput, err := do.WithExec([]string{
		"/app/doctl", "compute", "droplet", "list",
		"--format", "ID",
		"--no-header",
		"--tag-name", fmt.Sprintf("%s-n8n", n.subdomain),
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to list droplets: %w", err)
	}

	dropletID := strings.TrimSpace(dropletOutput)
	if dropletID != "" {
		// Delete droplet
		_, err = do.WithExec([]string{
			"/app/doctl", "compute", "droplet", "delete",
			"--force",
			dropletID,
		}).Sync(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete droplet: %w", err)
		}
	}

	// Delete DNS record
	records, err := do.WithExec([]string{
		"/app/doctl", "compute", "domain", "records", "list",
		n.domain,
		"--format", "ID,Name",
		"--no-header",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to list DNS records: %w", err)
	}

	for _, record := range strings.Split(records, "\n") {
		fields := strings.Fields(record)
		if len(fields) >= 2 && fields[1] == n.subdomain {
			_, err = do.WithExec([]string{
				"/app/doctl", "compute", "domain", "records", "delete",
				n.domain,
				"--force",
				fields[0],
			}).Sync(ctx)
			if err != nil {
				return fmt.Errorf("failed to delete DNS record: %w", err)
			}
			break
		}
	}

	return nil
} 