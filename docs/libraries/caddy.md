---
layout: default
title: Caddy Module
parent: Libraries
nav_order: 2
---

# Caddy Module

The Caddy module provides integration with Caddy, a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go. This module allows you to easily set up and configure Caddy as a reverse proxy in your Dagger pipelines.

## Features

- Easy Caddy server setup
- Reverse proxy configuration
- Multiple upstream services support
- Automatic Caddyfile generation
- Service binding integration
- Port exposure management
- Container and service modes

## Installation

To use the Caddy module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/caddy"
)
```

## Usage Examples

### Basic Caddy Setup

```go
func (m *MyModule) Example(ctx context.Context) (*Service, error) {
    caddy := dag.Caddy().New()
    
    // Add an upstream service
    caddy = caddy.WithService(
        ctx,
        upstreamService,  // *dagger.Service
        "api",           // upstream name
        8080,           // upstream port
    )
    
    return caddy.Serve(ctx), nil
}
```

### Multiple Upstream Services

```go
func (m *MyModule) MultipleServices(ctx context.Context) (*Service, error) {
    caddy := dag.Caddy().New()
    
    // Add multiple upstream services
    caddy = caddy.
        WithService(ctx, apiService, "api", 8080).
        WithService(ctx, webService, "web", 3000).
        WithService(ctx, dbService, "db", 5432)
    
    return caddy.Serve(ctx), nil
}
```

### Custom Container Configuration

```go
func (m *MyModule) CustomConfig(ctx context.Context) (*Container, error) {
    caddy := dag.Caddy().New()
    
    // Add services
    caddy = caddy.WithService(ctx, apiService, "api", 8080)
    
    // Get container for custom configuration
    container := caddy.Container(ctx)
    
    return container, nil
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Caddy Proxy
on: [push]

jobs:
  proxy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup Caddy Proxy
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/caddy
          args: |
            do -p '
              caddy := Caddy().New()
              caddy.WithService(
                ctx,
                apiService,
                "api",
                8080,
              ).Serve(ctx)
            '
```

## API Reference

### Caddy

Main module struct that provides access to Caddy functionality.

#### Constructor

- `New() *Caddy`
  - Creates a new Caddy instance
  - Returns an empty Caddy configuration

#### Methods

- `WithService(ctx context.Context, upstreamService *Service, upstreamName string, upstreamPort int32) *Caddy`
  - Adds an upstream service to the Caddy configuration
  - Parameters:
    - `upstreamService`: The service to proxy to
    - `upstreamName`: Name for the upstream service
    - `upstreamPort`: Port the upstream service listens on
  
- `GetCaddyFile(ctx context.Context) string`
  - Generates a Caddyfile based on the configured services
  
- `Container(ctx context.Context) *Container`
  - Returns a container with Caddy configured
  
- `Serve(ctx context.Context) *Service`
  - Returns a service running Caddy

### ServiceConfig

Configuration struct for upstream services.

#### Fields

- `UpstreamName`: Name of the upstream service
- `UpstreamPort`: Port the upstream service listens on
- `UpstreamSvc`: Reference to the upstream service

## Best Practices

1. **Service Configuration**
   - Use meaningful upstream names
   - Document port mappings
   - Keep services organized

2. **Resource Management**
   - Monitor proxy performance
   - Handle service dependencies
   - Clean up unused services

3. **Security**
   - Use HTTPS when possible
   - Configure appropriate headers
   - Follow security best practices

4. **Networking**
   - Use appropriate port mappings
   - Handle service discovery
   - Monitor connection health

## Troubleshooting

Common issues and solutions:

1. **Connection Issues**
   ```
   Error: dial tcp: lookup upstream: no such host
   Solution: Verify upstream service name and binding
   ```

2. **Port Conflicts**
   ```
   Error: bind: address already in use
   Solution: Check for port conflicts and use unique ports
   ```

3. **Configuration Problems**
   ```
   Error: parsing caddyfile: no addresses
   Solution: Verify Caddyfile syntax and service configuration
   ```

## Generated Caddyfile Example

For a configuration with multiple services:

```caddyfile
:8080 {
    reverse_proxy api:8080
}

:3000 {
    reverse_proxy web:3000
}

:5432 {
    reverse_proxy db:5432
}
```

## Advanced Usage

### Custom Caddy Container

```go
container := caddy.Container(ctx).
    WithEnvVariable("CADDY_OPTION", "value").
    WithMountedFile("/etc/caddy/custom.conf", configFile)
```

### Service with Health Checks

```go
func (m *MyModule) WithHealthCheck(ctx context.Context) (*Service, error) {
    caddy := dag.Caddy().New()
    
    // Add service with health check
    healthCheck := dag.Container().
        From("healthcheck:latest").
        WithExposedPort(8081)
    
    caddy = caddy.
        WithService(ctx, mainService, "api", 8080).
        WithService(ctx, healthCheck.AsService(), "health", 8081)
    
    return caddy.Serve(ctx), nil
}
``` 