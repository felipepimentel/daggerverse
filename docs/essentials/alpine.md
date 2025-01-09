---
layout: default
title: Alpine Module
parent: Essentials
nav_order: 1
---

# Alpine Module

The Alpine module provides integration with Alpine Linux, allowing you to create and manage Alpine-based containers in your Dagger pipelines. This module supports both Alpine Linux and Wolfi distributions.

## Features

- Alpine Linux container creation
- Package management
- Multi-architecture support
- Branch selection
- Base image configuration
- Package dependency resolution
- Environment setup
- Wolfi distribution support

## Installation

To use the Alpine module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/alpine"
)
```

## Usage Examples

### Basic Container Creation

```go
func (m *MyModule) Example(ctx context.Context) (*Container, error) {
    alpine, err := dag.Alpine().New(
        "",           // default architecture
        "edge",       // branch
        []string{     // packages
            "python3",
            "git",
            "curl",
        },
        "",           // default distro (Alpine)
    )
    if err != nil {
        return nil, err
    }
    
    return alpine.Container(ctx)
}
```

### Custom Architecture and Branch

```go
func (m *MyModule) CustomArch(ctx context.Context) (*Container, error) {
    alpine, err := dag.Alpine().New(
        "arm64",      // architecture
        "3.18",       // branch
        []string{     // packages
            "nodejs",
            "npm",
        },
        "",           // default distro (Alpine)
    )
    if err != nil {
        return nil, err
    }
    
    return alpine.Container(ctx)
}
```

### Using Wolfi Distribution

```go
func (m *MyModule) WolfiContainer(ctx context.Context) (*Container, error) {
    alpine, err := dag.Alpine().New(
        "",           // default architecture
        "edge",       // branch
        []string{     // packages
            "python3",
            "go",
        },
        "DISTRO_WOLFI", // use Wolfi distribution
    )
    if err != nil {
        return nil, err
    }
    
    return alpine.Container(ctx)
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
      - name: Build Alpine Container
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/alpine
          args: |
            do -p '
              alpine, err := Alpine().New(
                "",
                "edge",
                []string{"python3", "git"},
                "",
              )
              if err != nil {
                return nil, err
              }
              alpine.Container(ctx)
            '
```

## API Reference

### Alpine

Main module struct that provides access to Alpine Linux functionality.

#### Constructor

- `New(arch string, branch string, packages []string, distro Distro) (Alpine, error)`
  - Creates a new Alpine instance
  - Parameters:
    - `arch`: Hardware architecture (optional, defaults to runtime.GOARCH)
    - `branch`: Alpine branch (optional, defaults to "edge")
    - `packages`: List of packages to install (optional)
    - `distro`: Distribution type (optional, defaults to DISTRO_ALPINE)

#### Methods

- `Container(ctx context.Context) (*Container, error)`
  - Creates an Alpine container
  - Returns a configured container with specified packages and settings

#### Constants

- `DistroAlpine`: Use Alpine Linux distribution
- `DistroWolfi`: Use Wolfi distribution

## Best Practices

1. **Package Management**
   - Install only necessary packages
   - Use specific package versions when needed
   - Clean package cache after installation

2. **Architecture Support**
   - Test on target architectures
   - Use cross-compilation when needed
   - Verify package availability

3. **Branch Selection**
   - Use stable branches for production
   - Test with edge for latest features
   - Pin versions for reproducibility

4. **Container Optimization**
   - Minimize layer count
   - Remove unnecessary files
   - Use multi-stage builds

## Troubleshooting

Common issues and solutions:

1. **Package Issues**
   ```
   Error: package not found
   Solution: Verify package name and repository
   ```

2. **Architecture Problems**
   ```
   Error: unsupported architecture
   Solution: Check architecture compatibility
   ```

3. **Branch Errors**
   ```
   Error: invalid branch
   Solution: Use valid branch format (e.g., "3.18" or "edge")
   ```

## Configuration Example

```yaml
# alpine.yaml
architecture: amd64
branch: "3.18"
packages:
  - alpine-baselayout
  - alpine-keys
  - apk-tools
  - busybox
  - python3
  - git
distro: DISTRO_ALPINE
```

## Advanced Usage

### Custom Package Configuration

```go
func (m *MyModule) CustomPackages(ctx context.Context) (*Container, error) {
    alpine, err := dag.Alpine().New(
        "",
        "edge",
        []string{
            "build-base",
            "gcc",
            "musl-dev",
            "python3-dev",
        },
        "",
    )
    if err != nil {
        return nil, err
    }
    
    return alpine.Container(ctx)
}
``` 