---
layout: default
title: Registry Config Module
parent: Essentials
nav_order: 12
---

# Registry Config Module

The Registry Config module provides functionality for managing container registry configurations in your Dagger pipelines. It helps you handle authentication, credentials, and registry settings for various container registries.

## Features

- Registry authentication
- Credential management
- Multiple registry support
- Secure secret handling
- Configuration generation
- Docker compatibility
- Token management
- URL validation
- Error handling
- Configuration persistence

## Installation

To use the Registry Config module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/registry-config"
)
```

## Usage Examples

### Basic Registry Configuration

```go
func (m *MyModule) Example(ctx context.Context) (*File, error) {
    registryConfig := dag.RegistryConfig()
    
    // Add registry credentials
    return registryConfig.AddAuth(
        "registry.example.com",
        "username",
        dag.SetSecret("REGISTRY_TOKEN", "your-token"),
    )
}
```

### Multiple Registries

```go
func (m *MyModule) MultiRegistry(ctx context.Context) (*File, error) {
    registryConfig := dag.RegistryConfig()
    
    // Add first registry
    config, err := registryConfig.AddAuth(
        "registry1.example.com",
        "user1",
        dag.SetSecret("REGISTRY1_TOKEN", "token1"),
    )
    if err != nil {
        return nil, err
    }
    
    // Add second registry
    return registryConfig.AddAuth(
        "registry2.example.com",
        "user2",
        dag.SetSecret("REGISTRY2_TOKEN", "token2"),
    )
}
```

### Use in Container

```go
func (m *MyModule) UseConfig(ctx context.Context) error {
    registryConfig := dag.RegistryConfig()
    
    // Create config
    config, err := registryConfig.AddAuth(
        "registry.example.com",
        "username",
        dag.SetSecret("REGISTRY_TOKEN", "token"),
    )
    if err != nil {
        return err
    }
    
    // Use in container
    return dag.Container().
        From("alpine:latest").
        WithMountedFile("/root/.docker/config.json", config).
        WithExec([]string{
            "docker",
            "pull",
            "registry.example.com/image:tag",
        }).
        Sync(ctx)
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Registry Auth
on: [push]

jobs:
  auth:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Configure Registry
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/registry-config
          args: |
            do -p '
              registryConfig := RegistryConfig()
              registryConfig.AddAuth(
                "registry.example.com",
                "username",
                SetSecret("REGISTRY_TOKEN", "${{ secrets.REGISTRY_TOKEN }}"),
              )
            '
```

## API Reference

### RegistryConfig

Main module struct that provides access to registry configuration functionality.

#### Methods

- `AddAuth(registry string, username string, secret *Secret) (*File, error)`
  - Adds registry authentication
  - Parameters:
    - `registry`: Registry URL
    - `username`: Registry username
    - `secret`: Authentication secret
  - Returns configuration file

## Best Practices

1. **Credential Management**
   - Use secrets
   - Rotate tokens
   - Limit access

2. **Configuration**
   - Validate URLs
   - Check credentials
   - Document settings

3. **Security**
   - Secure storage
   - Encrypt secrets
   - Audit access

4. **Integration**
   - Test connections
   - Verify permissions
   - Monitor usage

## Troubleshooting

Common issues and solutions:

1. **Authentication Issues**
   ```
   Error: authentication failed
   Solution: Check credentials
   ```

2. **Registry Problems**
   ```
   Error: registry not found
   Solution: Verify registry URL
   ```

3. **Configuration Errors**
   ```
   Error: invalid config format
   Solution: Check JSON structure
   ```

## Configuration Example

```json
{
  "auths": {
    "registry.example.com": {
      "auth": "base64_encoded_credentials",
      "email": "user@example.com"
    }
  },
  "credHelpers": {
    "gcr.io": "gcloud",
    "*.azurecr.io": "acr-helper"
  }
}
```

## Advanced Usage

### Custom Authentication Helper

```go
func (m *MyModule) CustomAuth(ctx context.Context) error {
    registryConfig := dag.RegistryConfig()
    
    // Create base config
    config, err := registryConfig.AddAuth(
        "registry.example.com",
        "username",
        dag.SetSecret("REGISTRY_TOKEN", "token"),
    )
    if err != nil {
        return err
    }
    
    // Use custom helper
    return dag.Container().
        From("alpine:latest").
        WithMountedFile("/root/.docker/config.json", config).
        WithEnvVariable("DOCKER_CONFIG", "/root/.docker").
        WithExec([]string{
            "sh", "-c",
            `
            # Install custom helper
            apk add --no-cache jq
            
            # Modify config
            jq '.credHelpers."registry.example.com"="custom-helper"' \
                /root/.docker/config.json > /tmp/config.json
            
            # Use new config
            mv /tmp/config.json /root/.docker/config.json
            `,
        }).
        Sync(ctx)
}
```

### Registry Management

```go
func (m *MyModule) ManageRegistries(ctx context.Context) error {
    registryConfig := dag.RegistryConfig()
    
    // Define registries
    registries := []struct {
        url      string
        username string
        secret   string
    }{
        {"registry1.example.com", "user1", "REGISTRY1_TOKEN"},
        {"registry2.example.com", "user2", "REGISTRY2_TOKEN"},
        {"registry3.example.com", "user3", "REGISTRY3_TOKEN"},
    }
    
    // Configure all registries
    var config *File
    var err error
    
    for _, reg := range registries {
        config, err = registryConfig.AddAuth(
            reg.url,
            reg.username,
            dag.SetSecret(reg.secret, "token"),
        )
        if err != nil {
            return err
        }
    }
    
    // Verify configuration
    return dag.Container().
        From("alpine:latest").
        WithMountedFile("/config.json", config).
        WithExec([]string{
            "sh", "-c",
            "jq '.auths | keys[]' /config.json",
        }).
        Sync(ctx)
} 