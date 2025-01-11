# DigitalOcean Module

This module provides a Dagger interface for interacting with DigitalOcean services, with special focus on n8n deployments.

## Features

- 🚀 Container deployment to App Platform
- 🔄 n8n-specific deployment configurations
- 🖥️ Droplet management
- 📊 Status monitoring
- 🔐 Secure token handling

## Installation

```go
import (
    "github.com/felipepimentel/daggerverse/libraries/digitalocean"
)
```

## Usage

### Basic Example

```go
func main() {
    ctx := context.Background()
    
    do := digitalocean.New().WithToken(os.Getenv("DIGITALOCEAN_TOKEN"))
    
    // Deploy n8n
    config := digitalocean.N8NAppConfig{
        AppConfig: digitalocean.AppConfig{
            Name:   "my-n8n",
            Region: "nyc",
            EnvVars: []digitalocean.EnvVar{
                {Key: "N8N_BASIC_AUTH_ACTIVE", Value: "true"},
                {Key: "N8N_BASIC_AUTH_USER", Value: "admin"},
            },
        },
        WebhookURL: "https://my-n8n.example.com",
    }
    
    container, err := do.DeployN8N(ctx, config)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Advanced Example

```go
func deployN8N(ctx context.Context, token, appName string) error {
    do := digitalocean.New().WithToken(token)
    
    config := digitalocean.N8NAppConfig{
        AppConfig: digitalocean.AppConfig{
            Name:          appName,
            Region:       "fra1",
            InstanceSize: "basic-xs",
            InstanceCount: 2,
            EnvVars: []digitalocean.EnvVar{
                {Key: "N8N_BASIC_AUTH_ACTIVE", Value: "true"},
                {Key: "N8N_BASIC_AUTH_USER", Value: "admin"},
                {Key: "N8N_PROTOCOL", Value: "https"},
            },
            HealthCheckPath: "/healthz",
            HTTPPort:       5678,
        },
        WebhookURL:   "https://n8n.example.com",
        EncKey:       "your-encryption-key",
        DatabaseURL:  "postgresql://user:pass@host:5432/n8n",
    }
    
    container, err := do.DeployN8N(ctx, config)
    if err != nil {
        return fmt.Errorf("deploy failed: %w", err)
    }
    
    // Monitor deployment status
    status, err := do.GetN8NAppStatus(ctx, appID)
    if err != nil {
        return fmt.Errorf("status check failed: %w", err)
    }
    
    fmt.Printf("Deployment Status: %s\n", status.Status)
    fmt.Printf("Application URL: %s\n", status.URL)
    
    return nil
}
```

## API Reference

### Types

#### `Digitalocean`

Main type for interacting with DigitalOcean services.

```go
type Digitalocean struct {
    Token string
}
```

#### `EnvVar`

Environment variable representation.

```go
type EnvVar struct {
    Key   string
    Value string
}
```

#### `AppConfig`

Configuration for DigitalOcean app deployments.

```go
type AppConfig struct {
    Name             string
    Region           string
    InstanceSize     string
    InstanceCount    int64
    Container        Container
    EnvVars         []EnvVar
    HealthCheckPath  string
    HTTPPort        int
}
```

#### `N8NAppConfig`

n8n-specific configuration extending `AppConfig`.

```go
type N8NAppConfig struct {
    AppConfig
    DatabaseURL string
    WebhookURL  string
    EncKey      string
}
```

#### `N8NAppStatus`

Status information for n8n deployments.

```go
type N8NAppStatus struct {
    Status string
    URL    string
}
```

### Methods

#### `DeployN8N`

```go
func (d *Digitalocean) DeployN8N(ctx context.Context, config N8NAppConfig) (*Container, error)
```

Deploys an n8n instance to DigitalOcean App Platform.

#### `GetN8NAppStatus`

```go
func (d *Digitalocean) GetN8NAppStatus(ctx context.Context, appID string) (*N8NAppStatus, error)
```

Returns the status and URL of a deployed n8n application.

#### `WaitForAppDeployment`

```go
func (d *Digitalocean) WaitForAppDeployment(ctx context.Context, appID string) error
```

Waits for an app deployment to complete.

## Environment Variables

| Name | Description | Required |
|------|-------------|----------|
| `DIGITALOCEAN_TOKEN` | DigitalOcean API token | Yes |

## Best Practices

1. **Security**
   - Never hardcode the DigitalOcean token
   - Use environment variables for sensitive data
   - Rotate tokens regularly

2. **Deployment**
   - Use appropriate instance sizes
   - Configure health checks
   - Set proper environment variables
   - Use HTTPS for production

3. **Monitoring**
   - Check deployment status regularly
   - Monitor resource usage
   - Set up alerts

## Error Handling

The module provides detailed error messages. Common errors:

```go
// Authentication error
if err != nil && strings.Contains(err.Error(), "authentication failed") {
    // Handle invalid token
}

// Resource error
if err != nil && strings.Contains(err.Error(), "insufficient resources") {
    // Handle resource limits
}
```

## Related Modules

- [Docker Module](./docker.md) - For container management
- [n8n Pipeline](../pipelines/n8n.md) - For n8n-specific pipelines
- [n8n DigitalOcean Pipeline](../pipelines/n8n-digitalocean.md) - For complete n8n deployment pipelines

## Support

For issues and feature requests:
- Open an issue in the [repository](https://github.com/felipepimentel/daggerverse)
- Include error messages and logs
- Provide steps to reproduce 