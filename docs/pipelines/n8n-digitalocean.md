# N8N DigitalOcean Pipeline

This module provides a specialized deployment pipeline for n8n to DigitalOcean App Platform, with integrated Caddy server for HTTPS and advanced configuration options.

## Features

- Automated deployment to DigitalOcean App Platform
- Integrated Caddy server for HTTPS
- Persistent volume management
- Environment variable configuration
- Basic authentication support
- SQLite database integration
- Health check configuration
- Registry authentication
- Deployment status monitoring

## Installation

```bash
dagger mod use github.com/felipepimentel/daggerverse/pipelines/n8n-digitalocean@latest
```

## Usage

### Basic Example

```go
// Initialize the module
n8n := dag.N8nDigitalocean().
    WithSource(dag.Host().Directory(".")).
    WithDOConfig(&DOConfig{
        Token:        dag.SetSecret("do_token", "your-token"),
        Region:       "nyc",
        AppName:      "my-n8n",
        InstanceSize: "basic-xxs",
    })

// Deploy to DigitalOcean
container, err := n8n.Deploy(ctx)
```

### Configuration Options

The module supports the following configuration:

```go
type N8nDigitalocean struct {
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

type DOConfig struct {
    Token        *dagger.Secret
    Region       string
    AppName      string
    InstanceSize string
}

type CaddyConfig struct {
    Domain string
}
```

## Deployment Architecture

The module deploys two services:

1. **n8n Service**:
   - Node.js-based n8n instance
   - Persistent volume for data storage
   - Health check endpoint
   - Configurable instance size

2. **Caddy Service**:
   - Automatic HTTPS
   - Reverse proxy to n8n
   - Basic-xxs instance size
   - Persistent volume for certificates

## Complete Deployment Example

```go
n8n := dag.N8nDigitalocean().
    WithSource(dag.Host().Directory(".")).
    WithRegistry("registry.digitalocean.com/your-registry").
    WithTag("latest").
    WithDOConfig(&DOConfig{
        Token:        dag.SetSecret("do_token", "your-token"),
        Region:       "nyc",
        AppName:      "my-n8n",
        InstanceSize: "basic-xxs",
    }).
    WithEnvVars([]EnvVar{
        {Name: "N8N_HOST", Value: "n8n.example.com"},
        {Name: "N8N_PROTOCOL", Value: "https"},
        {Name: "N8N_PORT", Value: "5678"},
        {Name: "N8N_BASIC_AUTH_ACTIVE", Value: "true"},
        {Name: "N8N_BASIC_AUTH_USER", Value: "admin"},
        {Name: "N8N_BASIC_AUTH_PASSWORD", Value: "your-password"},
        {Name: "N8N_ENCRYPTION_KEY", Value: "your-encryption-key"},
        {Name: "NODE_ENV", Value: "production"},
        {Name: "DB_TYPE", Value: "sqlite"},
        {Name: "DB_SQLITE_PATH", Value: "/home/node/.n8n/database.sqlite"},
    })

container, err := n8n.Deploy(ctx)
```

## GitHub Actions Integration

Create a workflow file `.github/workflows/n8n-do.yml`:

```yaml
name: N8N DigitalOcean Deployment
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install Dagger CLI
        uses: dagger/dagger-for-github@v5
        with:
          version: "0.15.2"
      
      - name: Deploy to DigitalOcean
        env:
          DO_TOKEN: ${{ secrets.DO_TOKEN }}
          N8N_DOMAIN: ${{ secrets.N8N_DOMAIN }}
          N8N_BASIC_AUTH_PASSWORD: ${{ secrets.N8N_BASIC_AUTH_PASSWORD }}
          N8N_ENCRYPTION_KEY: ${{ secrets.N8N_ENCRYPTION_KEY }}
        run: |
          dagger call --progress=plain \
            ci \
            --source . \
            --region "nyc" \
            --app-name "my-n8n" \
            --token "$DO_TOKEN" \
            --domain "$N8N_DOMAIN" \
            --basic-auth-password "$N8N_BASIC_AUTH_PASSWORD" \
            --encryption-key "$N8N_ENCRYPTION_KEY"
```

## Monitoring Deployment Status

```go
status, err := n8n.GetStatus(ctx, "your-app-id")
```

## Best Practices

1. **Security**:
   - Always use HTTPS (enabled by default with Caddy)
   - Enable basic authentication
   - Use strong encryption keys
   - Keep secrets in secure environment variables

2. **Performance**:
   - Choose appropriate instance sizes
   - Monitor resource usage
   - Use persistent volumes for data

3. **Maintenance**:
   - Regularly backup the SQLite database
   - Keep n8n version updated
   - Monitor health checks

## Common Issues

1. **Deployment Failures**:
   - Verify DigitalOcean token permissions
   - Check registry authentication
   - Validate resource quotas

2. **Configuration Issues**:
   - Ensure all required environment variables are set
   - Check domain DNS configuration
   - Verify SSL certificate provisioning

3. **Runtime Problems**:
   - Monitor application logs
   - Check health check status
   - Verify volume persistence

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](../CONTRIBUTING.md) for details on how to submit pull requests.

## License

This module is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details. 