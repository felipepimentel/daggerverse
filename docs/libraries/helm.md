---
layout: default
title: Helm Module
parent: Libraries
nav_order: 5
---

# Helm Module

The Helm module provides integration with [Helm](https://helm.sh/), the package manager for Kubernetes. This module allows you to manage Helm charts, repositories, and deployments in your Dagger pipelines.

## Features

- Helm chart management
- Repository handling
- Chart installation and upgrades
- Package dependencies
- Version control
- Custom values support
- Multi-platform compatibility

## Installation

To use the Helm module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/helm"
)
```

## Usage Examples

### Basic Helm Chart Installation

```go
func (m *MyModule) Example(ctx context.Context) error {
    helm := dag.Helm().New()
    
    // Add repository and install chart
    return helm.
        AddRepo(ctx, "bitnami", "https://charts.bitnami.com/bitnami").
        Install(
            ctx,
            "my-release",
            "bitnami/nginx",
            "1.0.0",
            dag.File("./values.yaml"),
        )
}
```

### Chart Upgrade

```go
func (m *MyModule) UpgradeChart(ctx context.Context) error {
    helm := dag.Helm().New()
    
    // Upgrade existing release
    return helm.Upgrade(
        ctx,
        "my-release",
        "bitnami/nginx",
        "1.1.0",
        dag.File("./values.yaml"),
    )
}
```

### Repository Management

```go
func (m *MyModule) ManageRepos(ctx context.Context) error {
    helm := dag.Helm().New()
    
    // Add and update repositories
    return helm.
        AddRepo(ctx, "stable", "https://charts.helm.sh/stable").
        AddRepo(ctx, "jetstack", "https://charts.jetstack.io").
        UpdateRepos(ctx)
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Helm Operations
on: [push]

jobs:
  helm:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install Helm Chart
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/helm
          args: |
            do -p '
              helm := Helm().New()
              helm.AddRepo(
                ctx,
                "bitnami",
                "https://charts.bitnami.com/bitnami",
              ).Install(
                ctx,
                "my-release",
                "bitnami/nginx",
                "1.0.0",
                dag.File("./values.yaml"),
              )
            '
```

## API Reference

### Helm

Main module struct that provides access to Helm functionality.

#### Constructor

- `New() *Helm`
  - Creates a new Helm instance
  - Default version: "latest"
  - Default platform: "linux/amd64"

#### Methods

- `AddRepo(ctx context.Context, name string, url string) *Helm`
  - Adds a Helm repository
  - Parameters:
    - `name`: Repository name
    - `url`: Repository URL
  
- `UpdateRepos(ctx context.Context) error`
  - Updates all configured repositories
  
- `Install(ctx context.Context, release string, chart string, version string, values *File) error`
  - Installs a Helm chart
  - Parameters:
    - `release`: Release name
    - `chart`: Chart reference
    - `version`: Chart version
    - `values`: Values file
  
- `Upgrade(ctx context.Context, release string, chart string, version string, values *File) error`
  - Upgrades a Helm release
  - Parameters:
    - `release`: Release name
    - `chart`: Chart reference
    - `version`: Chart version
    - `values`: Values file

## Best Practices

1. **Repository Management**
   - Use official and trusted repositories
   - Keep repositories updated
   - Document repository sources

2. **Version Control**
   - Pin chart versions
   - Test upgrades in staging
   - Track version changes

3. **Values Management**
   - Use version-controlled values files
   - Document values changes
   - Validate values before deployment

4. **Security**
   - Follow Helm security best practices
   - Use signed charts when possible
   - Regularly update dependencies

## Troubleshooting

Common issues and solutions:

1. **Repository Issues**
   ```
   Error: failed to fetch repository
   Solution: Check repository URL and connectivity
   ```

2. **Chart Installation Failures**
   ```
   Error: chart not found
   Solution: Verify chart name and repository configuration
   ```

3. **Version Conflicts**
   ```
   Error: incompatible versions
   Solution: Check version compatibility and dependencies
   ```

## Values File Example

```yaml
# values.yaml
replicaCount: 3
image:
  repository: nginx
  tag: 1.21.0
  pullPolicy: IfNotPresent
service:
  type: LoadBalancer
  port: 80
resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

## Advanced Usage

### Custom Chart Development

```go
func (m *MyModule) DevelopChart(ctx context.Context) error {
    helm := dag.Helm().New()
    
    // Package and install local chart
    return helm.
        Package(ctx, "./my-chart", "1.0.0").
        Install(
            ctx,
            "my-release",
            "./my-chart-1.0.0.tgz",
            "",
            dag.File("./values.yaml"),
        )
}
```

### Dependency Management

```go
func (m *MyModule) ManageDeps(ctx context.Context) error {
    helm := dag.Helm().New()
    
    // Update chart dependencies
    return helm.
        AddRepo(ctx, "deps", "https://charts.deps.io").
        UpdateDependencies(ctx, "./my-chart")
}
``` 