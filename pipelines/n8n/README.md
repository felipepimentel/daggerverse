# n8n Pipeline Module

This module provides a reusable CI/CD pipeline for deploying n8n to DigitalOcean. It automates the entire deployment process, including infrastructure provisioning, DNS configuration, and service deployment.

## Features

- Automated deployment of n8n to DigitalOcean
- PostgreSQL database setup and configuration
- Caddy reverse proxy with automatic SSL/TLS
- DNS configuration
- Backup management with retention policies
- Container monitoring with cAdvisor
- CI validation checks
- Complete cleanup functionality

## Prerequisites

1. DigitalOcean account with API token
2. SSH key added to DigitalOcean
3. Domain managed by DigitalOcean DNS

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/felipepimentel/daggerverse/pipelines/n8n"
)

func main() {
    ctx := context.Background()

    // Create a new n8n deployment
    n8n := n8n.New(
        "example.com",      // Domain
        "n8n",             // Subdomain
        "ssh-key-name",    // SSH key name in DigitalOcean
    ).WithSSLEmail("admin@example.com")

    // Run CI checks
    if err := n8n.CI(ctx); err != nil {
        fmt.Fprintf(os.Stderr, "CI validation failed: %v\n", err)
        os.Exit(1)
    }

    // Deploy n8n
    if err := n8n.Deploy(ctx); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to deploy n8n: %v\n", err)
        os.Exit(1)
    }

    // Get the n8n URL
    fmt.Printf("n8n is available at: %s\n", n8n.GetURL())
}
```

### Advanced Configuration

```go
n8n := n8n.New("example.com", "n8n", "ssh-key-name").
    WithRegion("nyc1").
    WithSize("s-2vcpu-4gb").
    WithN8NVersion("0.234.0").
    WithN8NPort(5678).
    WithPostgresConfig(
        "15-alpine",    // Version
        "n8n",         // User
        "password123", // Password
        "n8n",        // Database
    ).
    WithBackupConfig(
        true,         // Enabled
        "0 0 * * *", // Cron schedule
    ).
    WithRegistryConfig("n8n-registry").
    WithSSLEmail("admin@example.com")
```

## Configuration Options

### Base Configuration

- `Domain`: Your domain name (e.g., "example.com")
- `Subdomain`: Subdomain for n8n (e.g., "n8n")
- `Region`: DigitalOcean region (default: "nyc1")
- `Size`: Droplet size (default: "s-2vcpu-4gb")
- `SSHKeyName`: Name of your SSH key in DigitalOcean

### Registry Configuration

- `RegistryName`: Name of the DigitalOcean container registry (default: "n8n-registry")

### n8n Configuration

- `N8nVersion`: n8n version to deploy (default: "0.234.0")
- `N8nPort`: Port for n8n (default: 5678)

### Database Configuration

- `PostgresVersion`: PostgreSQL version (default: "15-alpine")
- `PostgresUser`: Database user (default: "n8n")
- `PostgresPass`: Database password (default: "n8n")
- `PostgresDB`: Database name (default: "n8n")

### Backup Configuration

- `BackupEnabled`: Enable automated backups (default: true)
- `BackupCron`: Backup schedule in cron format (default: "0 0 * * *")
- `BackupRetention`: Number of days to retain backups (default: 7)

### Monitoring Configuration

- `MonitoringEnabled`: Enable cAdvisor monitoring (default: true)
- `CAdvisorPort`: Port for cAdvisor metrics (default: 8080)

### SSL/TLS Configuration

- `SSLEmail`: Email for Let's Encrypt certificates

## Methods

### Configuration Methods

- `WithRegion(region string) *N8N`
- `WithSize(size string) *N8N`
- `WithN8NVersion(version string) *N8N`
- `WithN8NPort(port int) *N8N`
- `WithPostgresConfig(version, user, pass, db string) *N8N`
- `WithBackupConfig(enabled bool, cron string) *N8N`
- `WithRegistryConfig(name string) *N8N`
- `WithSSLEmail(email string) *N8N`

### Deployment Methods

- `CI(ctx context.Context) error`: Run validation checks
- `Deploy(ctx context.Context) error`: Deploy n8n
- `Destroy(ctx context.Context) error`: Remove specific deployment
- `Cleanup(ctx context.Context) error`: Remove all resources
- `GetStatus(ctx context.Context) (*dropletInfo, error)`: Get deployment status
- `GetURL() string`: Get the n8n URL

## Required Secrets

The module requires the following secrets to be set:

1. `do_token`: DigitalOcean API token
2. `ssh_key`: SSH private key for server access
3. `ssh_key_fingerprint`: SSH key fingerprint registered in DigitalOcean
4. `ssh_key_id`: SSH key ID registered in DigitalOcean

Optional secrets:
1. `n8n_basic_auth_password`: Password for n8n basic auth
2. `n8n_encryption_key`: Encryption key for n8n

## Example Usage with Secrets

```go
// Set secrets in your environment
os.Setenv("DAGGER_DO_TOKEN", "your-digitalocean-token")
os.Setenv("DAGGER_SSH_KEY", "your-ssh-private-key")
os.Setenv("DAGGER_SSH_KEY_FINGERPRINT", "your-ssh-key-fingerprint")
os.Setenv("DAGGER_SSH_KEY_ID", "your-ssh-key-id")

// Create and deploy n8n
n8n := n8n.New("example.com", "n8n", "ssh-key-name").
    WithSSLEmail("admin@example.com")

// Run CI checks
if err := n8n.CI(ctx); err != nil {
    log.Fatal(err)
}

// Deploy
if err := n8n.Deploy(ctx); err != nil {
    log.Fatal(err)
}
```

## Deployment Steps

1. **CI Validation**
   - Validate all configurations
   - Test image pulls
   - Verify secret availability

2. **Registry Setup**
   - Create DigitalOcean container registry
   - Configure registry authentication
   - Push required images

3. **Infrastructure**
   - Create droplet with monitoring enabled
   - Configure DNS records
   - Set up SSH access

4. **Service Deployment**
   - Deploy PostgreSQL database
   - Configure n8n with database connection
   - Set up Caddy reverse proxy
   - Enable SSL/TLS with Let's Encrypt

5. **Monitoring & Backups**
   - Deploy cAdvisor for container monitoring
   - Configure automated database backups
   - Set up backup retention policies

6. **Verification**
   - Check service health
   - Verify HTTPS access
   - Validate monitoring metrics
   - Test backup functionality

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat(n8n): add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 