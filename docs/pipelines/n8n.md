# N8N Pipeline

This module provides functionality for building, testing, and deploying n8n instances. It supports multiple deployment providers through a flexible provider interface.

## Features

- Build n8n containers with proper configuration
- Run tests on n8n instances
- Deploy to multiple cloud providers (currently supports DigitalOcean)
- Configurable environment variables
- Volume persistence for data
- Optional reverse proxy (Caddy) configuration

## Usage

### Basic Example

```go
// Create a new n8n instance
n8n := dag.N8N().
    WithSource(source).
    WithRegistry("registry.example.com/n8n").
    WithTag("latest")

// Run tests
if err := n8n.Test(ctx); err != nil {
    return err
}

// Build and verify the container
container, err := n8n.Build(ctx)
if err != nil {
    return err
}
```

### Deployment with DigitalOcean

```go
// Create a DigitalOcean provider
provider := &DigitalOceanProvider{
    Token:        token,
    Region:       "nyc1",
    AppName:      "my-n8n",
    InstanceSize: "basic-xxs",
    Domain:       "n8n.example.com",  // Optional, adds Caddy reverse proxy if specified
}

// Configure n8n with the provider
n8n := dag.N8N().
    WithSource(source).
    WithRegistry("registry.digitalocean.com/myregistry/n8n").
    WithTag("latest").
    WithProvider(provider)

// Deploy
container, err := n8n.Deploy(ctx)
```

## Environment Variables

The module supports all standard n8n environment variables. Common ones include:

- `N8N_HOST`: The host where n8n will be accessible
- `N8N_PROTOCOL`: Protocol (http/https)
- `N8N_PORT`: Port to expose n8n on (default: 5678)
- `N8N_BASIC_AUTH_ACTIVE`: Enable basic auth
- `N8N_BASIC_AUTH_USER`: Basic auth username
- `N8N_BASIC_AUTH_PASSWORD`: Basic auth password
- `N8N_ENCRYPTION_KEY`: Key for encrypting credentials

## Provider Interface

The module uses a provider interface that allows implementing different deployment targets:

```go
type Provider interface {
    Deploy(ctx context.Context, container *dagger.Container, registry string, tag string) error
    GetStatus(ctx context.Context) (*dagger.Container, error)
}
```

Currently supported providers:
- DigitalOcean: Deploys n8n to DigitalOcean Apps platform with optional Caddy reverse proxy

To implement a new provider, create a struct that implements the Provider interface.

## Methods

### Build

Creates a container with n8n installed and configured:

```go
container, err := n8n.Build(ctx)
```

### Test

Runs tests if a package.json exists:

```go
err := n8n.Test(ctx)
```

### Deploy

Deploys n8n using the configured provider:

```go
container, err := n8n.Deploy(ctx)
```

### CI/CD

The module provides CI/CD pipeline methods:

```go
// Run CI pipeline (tests)
err := n8n.CI(ctx)

// Run CD pipeline (deploy)
container, err := n8n.CD(ctx)
```

## Configuration Methods

The module uses a builder pattern for configuration:

- `WithSource(source)`: Sets the source directory
- `WithRegistry(registry)`: Sets the container registry
- `WithTag(tag)`: Sets the container tag
- `WithRegistryAuth(auth)`: Sets registry authentication
- `WithProvider(provider)`: Sets the deployment provider

## Data Persistence

The module automatically configures a persistent volume for n8n data at `/data` in the container. 