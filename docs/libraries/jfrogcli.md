---
layout: default
title: JFrog CLI Module
parent: Libraries
nav_order: 6
---

# JFrog CLI Module

The JFrog CLI module provides integration with [JFrog CLI](https://jfrog.com/cli/), a compact and smart client that provides a simple interface to automate access to JFrog products. This module allows you to interact with Artifactory, Xray, and other JFrog services in your Dagger pipelines.

## Features

- Artifactory integration
- Package management
- Build info tracking
- Security scanning
- Repository management
- Multi-platform support
- Authentication handling
- Custom configuration

## Installation

To use the JFrog CLI module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/jfrogcli"
)
```

## Usage Examples

### Basic Authentication Setup

```go
func (m *MyModule) Example(ctx context.Context) error {
    jfrog := dag.JFrogCLI().New()
    
    // Configure JFrog CLI
    return jfrog.Configure(
        ctx,
        "my-server",
        "https://artifactory.example.com",
        dag.SetSecret("JFROG_USER", "user"),
        dag.SetSecret("JFROG_PASSWORD", "password"),
    )
}
```

### Upload Artifacts

```go
func (m *MyModule) UploadArtifacts(ctx context.Context) error {
    jfrog := dag.JFrogCLI().New()
    
    // Upload artifacts to Artifactory
    return jfrog.Upload(
        ctx,
        "./dist/*.jar",
        "libs-release-local",
        "--build-name=my-build",
        "--build-number=1",
    )
}
```

### Download Artifacts

```go
func (m *MyModule) DownloadArtifacts(ctx context.Context) error {
    jfrog := dag.JFrogCLI().New()
    
    // Download artifacts from Artifactory
    return jfrog.Download(
        ctx,
        "libs-release-local/org/example/app/*.jar",
        "./deps/",
        "--build-name=my-build",
        "--build-number=1",
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: JFrog Operations
on: [push]

jobs:
  jfrog:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Upload to Artifactory
        uses: dagger/dagger-action@v1
        env:
          JFROG_USER: ${{ secrets.JFROG_USER }}
          JFROG_PASSWORD: ${{ secrets.JFROG_PASSWORD }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/jfrogcli
          args: |
            do -p '
              jfrog := JFrogCLI().New()
              jfrog.Configure(
                ctx,
                "my-server",
                "https://artifactory.example.com",
                dag.SetSecret("JFROG_USER", JFROG_USER),
                dag.SetSecret("JFROG_PASSWORD", JFROG_PASSWORD),
              ).Upload(
                ctx,
                "./dist/*.jar",
                "libs-release-local",
                "--build-name=my-build",
                "--build-number=1",
              )
            '
```

## API Reference

### JFrogCLI

Main module struct that provides access to JFrog CLI functionality.

#### Constructor

- `New() *JFrogCLI`
  - Creates a new JFrog CLI instance
  - Default version: "latest"
  - Default platform: "linux/amd64"

#### Methods

- `Configure(ctx context.Context, serverID string, url string, user *Secret, password *Secret) *JFrogCLI`
  - Configures JFrog CLI server connection
  - Parameters:
    - `serverID`: Server identifier
    - `url`: Artifactory URL
    - `user`: Username secret
    - `password`: Password secret
  
- `Upload(ctx context.Context, pattern string, target string, args ...string) error`
  - Uploads artifacts to Artifactory
  - Parameters:
    - `pattern`: File pattern to upload
    - `target`: Target repository path
    - `args`: Additional upload arguments
  
- `Download(ctx context.Context, pattern string, target string, args ...string) error`
  - Downloads artifacts from Artifactory
  - Parameters:
    - `pattern`: File pattern to download
    - `target`: Target local path
    - `args`: Additional download arguments

## Best Practices

1. **Authentication**
   - Use secrets for credentials
   - Rotate credentials regularly
   - Use access tokens when possible

2. **Build Info**
   - Track build information
   - Use consistent naming
   - Include relevant metadata

3. **Repository Management**
   - Follow naming conventions
   - Clean up old artifacts
   - Use appropriate repository types

4. **Security**
   - Enable Xray scanning
   - Follow security best practices
   - Monitor security issues

## Troubleshooting

Common issues and solutions:

1. **Authentication Issues**
   ```
   Error: unauthorized access
   Solution: Verify credentials and permissions
   ```

2. **Upload Failures**
   ```
   Error: failed to upload artifacts
   Solution: Check file patterns and repository permissions
   ```

3. **Download Issues**
   ```
   Error: artifact not found
   Solution: Verify artifact path and repository configuration
   ```

## Configuration Example

```yaml
# jfrog-cli.conf.yaml
version: 1
artifactory:
  serverID: my-server
  url: https://artifactory.example.com
  user: ${JFROG_USER}
  password: ${JFROG_PASSWORD}
  defaultRepo: libs-release-local
```

## Advanced Usage

### Build Integration

```go
func (m *MyModule) BuildIntegration(ctx context.Context) error {
    jfrog := dag.JFrogCLI().New()
    
    // Start build tracking
    return jfrog.
        BuildAdd(ctx, "my-build", "1").
        Upload(
            ctx,
            "./dist/*.jar",
            "libs-release-local",
            "--build-name=my-build",
            "--build-number=1",
        ).
        BuildPublish(ctx, "my-build", "1")
}
```

### Security Scanning

```go
func (m *MyModule) SecurityScan(ctx context.Context) error {
    jfrog := dag.JFrogCLI().New()
    
    // Scan artifacts with Xray
    return jfrog.XrayScan(
        ctx,
        "my-build",
        "1",
        "--fail=high",
    )
}
``` 