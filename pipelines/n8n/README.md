# n8n Pipeline Module

This module provides a reusable CI/CD pipeline for deploying n8n to DigitalOcean. It automates the entire deployment process, including infrastructure provisioning, DNS configuration, and service deployment.

## Features

- Automated deployment of n8n to DigitalOcean
- PostgreSQL database setup and configuration
- Caddy reverse proxy with automatic SSL/TLS
- DNS configuration
- Container monitoring with cAdvisor
- Complete cleanup functionality

## Prerequisites

1. DigitalOcean account with API token
2. Domain managed by DigitalOcean DNS

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
    n8n := n8n.New().
        WithRegion("nyc1").
        WithSize("s-2vcpu-2gb").
        WithImage("ubuntu-20-04-x64")

    // Deploy n8n
    url, err := n8n.Deploy(ctx, os.Getenv("DIGITALOCEAN_TOKEN"))
    if err != nil {
        panic(err)
    }

    fmt.Printf("n8n is now available at: %s\n", url)
}
```

## Configuration Methods

- `WithRegion(region string) *N8N`: Set the DigitalOcean region (default: "nyc1")
- `WithSize(size string) *N8N`: Set the droplet size (default: "s-2vcpu-2gb")
- `WithImage(image string) *N8N`: Set the droplet image (default: "ubuntu-20-04-x64")

## Deployment Process

1. **Cleanup**: Remove any existing resources with the same name
2. **SSH Keys**: Generate and register SSH keys with DigitalOcean
3. **Infrastructure**: Create a droplet with the specified configuration
4. **DNS**: Configure DNS records for the n8n subdomain
5. **Docker**: Install Docker and Docker Compose on the droplet
6. **Configuration**: Create necessary configuration files
7. **Services**: Start n8n and Caddy services using Docker Compose

## Configuration Files

### docker-compose.yml
- n8n service configuration
- Caddy reverse proxy setup
- Volume and network definitions

### Caddyfile
- Automatic HTTPS configuration
- WebSocket support
- Security headers
- Logging configuration

### .env
- n8n host configuration
- Basic authentication settings
- Encryption key generation

## Security Features

1. **SSL/TLS**: Automatic HTTPS with Let's Encrypt
2. **Basic Auth**: Enabled by default with configurable credentials
3. **Security Headers**:
   - HSTS
   - XSS Protection
   - Frame Options
   - Content Type Options
   - Referrer Policy

## Dependencies

This module uses the following reusable modules:
- `digitalocean`: For managing DigitalOcean resources
- `docker`: For Docker and Docker Compose operations
- `ssh`: For SSH key management and remote execution

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes using semantic commit messages
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 