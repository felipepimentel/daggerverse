---
layout: default
title: Wolfi Module
parent: Libraries
nav_order: 20
---

# Wolfi Module

The Wolfi module provides integration with [Wolfi](https://wolfi.dev/), a community Linux OS designed for container and cloud-native workloads. This module allows you to build and manage Wolfi-based containers in your Dagger pipelines.

## Features

- Container building
- Package management
- Security hardening
- Base image creation
- Multi-stage builds
- Layer optimization
- Dependency management
- Image signing

## Installation

To use the Wolfi module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/wolfi"
)
```

## Usage Examples

### Basic Container Creation

```go
func (m *MyModule) Example(ctx context.Context) (*Container, error) {
    wolfi := dag.Wolfi().New()
    
    // Create container
    return wolfi.Container(
        ctx,
        []string{
            "python-3.11",
            "git",
            "curl",
        },
    )
}
```

### Custom Package Installation

```go
func (m *MyModule) InstallPackages(ctx context.Context) (*Container, error) {
    wolfi := dag.Wolfi().New()
    
    // Install specific packages
    return wolfi.WithPackages(
        ctx,
        []string{
            "nodejs-20",
            "npm",
            "build-base",
        },
        map[string]string{
            "repository": "wolfi-edge",
            "keyring": "/usr/share/apk/keys/wolfi-signing.rsa.pub",
        },
    )
}
```

### Multi-stage Build

```go
func (m *MyModule) MultistageBuild(ctx context.Context) (*Container, error) {
    wolfi := dag.Wolfi().New()
    
    // Create multi-stage build
    builder := wolfi.Container(ctx, []string{"go", "build-base"})
    
    return wolfi.Container(
        ctx,
        []string{"ssl_client", "ca-certificates"},
    ).WithFile(
        "/app",
        builder.WithDirectory(".", dag.Directory("./src")).
            WithWorkdir("/src").
            WithExec([]string{"go", "build", "-o", "/app"}),
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Container Build
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build Container
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/wolfi
          args: |
            do -p '
              wolfi := Wolfi().New()
              wolfi.Container(
                ctx,
                []string{
                  "python-3.11",
                  "git",
                  "curl",
                },
              )
            '
```

## API Reference

### Wolfi

Main module struct that provides access to Wolfi functionality.

#### Constructor

- `New() *Wolfi`
  - Creates a new Wolfi instance
  - Default version: "latest"
  - Default platform: "linux/amd64"

#### Methods

- `Container(ctx context.Context, packages []string) (*Container, error)`
  - Creates a Wolfi container
  - Parameters:
    - `packages`: List of packages to install
  
- `WithPackages(ctx context.Context, packages []string, config map[string]string) (*Container, error)`
  - Installs additional packages
  - Parameters:
    - `packages`: List of packages
    - `config`: Package configuration
  
- `WithSecurity(ctx context.Context, config map[string]string) (*Container, error)`
  - Applies security configurations
  - Parameters:
    - `config`: Security settings

## Best Practices

1. **Container Building**
   - Minimize layers
   - Use multi-stage builds
   - Clean package cache

2. **Security**
   - Keep base minimal
   - Update regularly
   - Verify signatures

3. **Dependencies**
   - Pin versions
   - Use official repos
   - Document requirements

4. **Optimization**
   - Reduce image size
   - Layer caching
   - Remove unnecessary files

## Troubleshooting

Common issues and solutions:

1. **Package Issues**
   ```
   Error: package not found
   Solution: Check package name and repository
   ```

2. **Build Problems**
   ```
   Error: build failed
   Solution: Verify build dependencies
   ```

3. **Security Errors**
   ```
   Error: signature verification failed
   Solution: Check keyring configuration
   ```

## Configuration Example

```yaml
# wolfi.yaml
repository:
  name: wolfi-base
  url: https://packages.wolfi.dev/os
  keyring: /usr/share/apk/keys/wolfi-signing.rsa.pub

packages:
  base:
    - ca-certificates
    - ssl_client
    - tzdata
  build:
    - build-base
    - git
    - make
  runtime:
    - python-3.11
    - nodejs-20

security:
  verify_signatures: true
  no_cache: false
  update_index: true
```

## Advanced Usage

### Custom Repository Configuration

```go
func (m *MyModule) CustomRepo(ctx context.Context) (*Container, error) {
    wolfi := dag.Wolfi().New()
    
    // Use custom repository
    return wolfi.WithRepository(
        ctx,
        "custom-repo",
        "https://custom.repository.dev/wolfi",
        dag.File("./repo-key.rsa.pub"),
        map[string]string{
            "priority": "100",
            "verify": "true",
        },
    )
}
```

### Security Hardening

```go
func (m *MyModule) SecureContainer(ctx context.Context) (*Container, error) {
    wolfi := dag.Wolfi().New()
    
    // Apply security configurations
    return wolfi.WithSecurity(
        ctx,
        map[string]string{
            "no_root": "true",
            "readonly_root": "true",
            "seccomp_profile": "default.json",
            "capabilities": "none",
        },
    ).Container(ctx, []string{"app-runtime"})
}
``` 