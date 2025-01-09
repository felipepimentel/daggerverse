---
layout: default
title: Melange Module
parent: Essentials
nav_order: 11
---

# Melange Module

The Melange module provides integration with Melange, a build system for creating APK packages. This module allows you to build and manage APK packages in your Dagger pipelines.

## Features

- APK package building
- YAML configuration
- Pipeline integration
- Build caching
- Version management
- Dependency resolution
- Repository configuration
- Multi-architecture support
- Build signing
- Output management

## Installation

To use the Melange module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/melange"
)
```

## Usage Examples

### Basic Package Building

```go
func (m *MyModule) Example(ctx context.Context) (*Directory, error) {
    melange := dag.Melange().New(nil, false)
    
    // Build package from YAML
    return melange.Build(
        ctx,
        dag.Directory(".").File("package.yaml"),
        "",           // arch (optional)
        nil,          // keyring (optional)
        nil,          // repository append (optional)
        "",           // version (optional)
        nil,          // work dir (optional)
    )
}
```

### Custom Architecture Build

```go
func (m *MyModule) ArchSpecific(ctx context.Context) (*Directory, error) {
    melange := dag.Melange().New(nil, false)
    
    // Build for specific architecture
    return melange.Build(
        ctx,
        dag.Directory(".").File("package.yaml"),
        "aarch64",    // target architecture
        nil,          // keyring
        nil,          // repository append
        "1.0.0",      // version
        nil,          // work dir
    )
}
```

### With Custom Repository

```go
func (m *MyModule) CustomRepo(ctx context.Context) (*Directory, error) {
    melange := dag.Melange().New(nil, false)
    
    // Add custom repository
    repos := []string{
        "https://dl-cdn.alpinelinux.org/alpine/edge/main",
        "https://custom.repo/packages",
    }
    
    // Build with custom repos
    return melange.Build(
        ctx,
        dag.Directory(".").File("package.yaml"),
        "",           // arch
        nil,          // keyring
        repos,        // repository append
        "",           // version
        nil,          // work dir
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Package Build
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build Package
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/melange
          args: |
            do -p '
              melange := Melange().New(nil, false)
              melange.Build(
                ctx,
                Directory(".").File("package.yaml"),
                "",
                nil,
                nil,
                "",
                nil,
              )
            '
```

## API Reference

### Melange

Main module struct that provides access to Melange functionality.

#### Constructor

- `New(container *Container, withoutCache bool) *Melange`
  - Creates a new Melange instance
  - Parameters:
    - `container`: Custom base container (optional)
    - `withoutCache`: Disable default cache volume (optional)

#### Methods

- `Build(ctx context.Context, config *File, arch string, keyring *Directory, repositoryAppend []string, version string, workDir *Directory) (*Directory, error)`
  - Builds APK package
  - Parameters:
    - `config`: YAML configuration file
    - `arch`: Target architecture (optional)
    - `keyring`: GPG keyring directory (optional)
    - `repositoryAppend`: Additional repositories (optional)
    - `version`: Package version (optional)
    - `workDir`: Working directory (optional)

## Best Practices

1. **Package Configuration**
   - Use version control
   - Document dependencies
   - Validate YAML

2. **Build Management**
   - Cache builds
   - Sign packages
   - Test outputs

3. **Repository Setup**
   - Verify sources
   - Pin versions
   - Document URLs

4. **Architecture Support**
   - Test all targets
   - Validate builds
   - Document limits

## Troubleshooting

Common issues and solutions:

1. **Build Errors**
   ```
   Error: build failed
   Solution: Check dependencies
   ```

2. **Repository Issues**
   ```
   Error: repository not found
   Solution: Verify repository URL
   ```

3. **Version Problems**
   ```
   Error: invalid version
   Solution: Use valid version string
   ```

## Configuration Example

```yaml
# package.yaml
package:
  name: example-pkg
  version: 1.0.0
  description: Example package
  target-architecture:
    - all
  copyright:
    - paths:
      - "*"
      license: MIT

environment:
  contents:
    repositories:
      - https://dl-cdn.alpinelinux.org/alpine/edge/main
    packages:
      - alpine-baselayout
      - busybox

pipeline:
  - uses: git-checkout
    with:
      repository: https://github.com/example/repo
      tag: v1.0.0
      
  - uses: autoconf/make
    with:
      opts:
        - --prefix=/usr
```

## Advanced Usage

### Multi-Architecture Build

```go
func (m *MyModule) MultiArchBuild(ctx context.Context) error {
    melange := dag.Melange().New(nil, false)
    
    // Define architectures
    arches := []string{"amd64", "arm64", "ppc64le"}
    
    // Build for each architecture
    for _, arch := range arches {
        buildDir, err := melange.Build(
            ctx,
            dag.Directory(".").File("package.yaml"),
            arch,
            nil,
            nil,
            "",
            nil,
        )
        if err != nil {
            return err
        }
        
        // Process build output
        err = dag.Container().
            From("alpine:latest").
            WithMountedDirectory("/build", buildDir).
            WithExec([]string{
                "sh", "-c",
                fmt.Sprintf("cp -r /build/* /out/%s/", arch),
            }).
            Sync(ctx)
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### Custom Build Pipeline

```go
func (m *MyModule) CustomBuild(ctx context.Context) error {
    melange := dag.Melange().New(nil, false)
    
    // Add custom keyring
    keyring := dag.Directory("./keys")
    
    // Add custom repositories
    repos := []string{
        "https://custom.repo/main",
        "https://custom.repo/testing",
    }
    
    // Build with custom configuration
    buildDir, err := melange.Build(
        ctx,
        dag.Directory(".").File("package.yaml"),
        "amd64",
        keyring,
        repos,
        "2.0.0",
        dag.Directory("./work"),
    )
    if err != nil {
        return err
    }
    
    // Process and verify build
    return dag.Container().
        From("alpine:latest").
        WithMountedDirectory("/build", buildDir).
        WithExec([]string{
            "sh", "-c",
            `
            # Verify package
            apk verify /build/*.apk
            
            # Extract package info
            tar -tzf /build/*.apk > /build/contents.txt
            
            # Copy to output
            mkdir -p /out/verified
            cp /build/*.apk /out/verified/
            `,
        }).
        Sync(ctx)
} 