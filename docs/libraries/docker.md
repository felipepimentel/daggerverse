---
layout: default
title: Docker Module
parent: Libraries
nav_order: 1
---

# Docker Module

The Docker module provides a seamless integration with Docker Engine, allowing you to manage containers, images, and Docker operations within your Dagger pipelines.

## Features

- Spawn ephemeral Docker Engine instances
- Execute Docker CLI commands
- Manage Docker images (pull, push, import)
- Run containers
- State persistence options

## Installation

To use the Docker module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/docker"
)
```

## Usage Examples

### Basic Docker Engine Setup

```go
func (m *MyModule) Example(ctx context.Context) (*Container, error) {
    docker := dag.Docker()
    
    // Start a Docker engine with default settings
    engine := docker.Engine()
    
    // Create a CLI instance connected to the engine
    cli := docker.CLI()
    
    return cli.Container(), nil
}
```

### Custom Engine Configuration

```go
func (m *MyModule) CustomEngine(ctx context.Context) (*Container, error) {
    docker := dag.Docker()
    
    // Start a Docker engine with custom version and persistence
    engine := docker.Engine(
        dagger.EngineOpts{
            Version: "24.0",
            Persist: true,
            Namespace: "my-project",
        },
    )
    
    return docker.CLI(dagger.CLIOpts{
        Engine: engine,
    }).Container(), nil
}
```

### Image Operations

```go
func (m *MyModule) ImageOps(ctx context.Context) error {
    docker := dag.Docker()
    cli := docker.CLI()
    
    // Pull an image
    image, err := cli.Pull(ctx, "alpine", "latest")
    if err != nil {
        return err
    }
    
    // Push an image
    _, err = cli.Push(ctx, "my-registry/alpine", "custom-tag")
    if err != nil {
        return err
    }
    
    return nil
}
```

### Running Containers

```go
func (m *MyModule) RunContainer(ctx context.Context) (string, error) {
    docker := dag.Docker()
    cli := docker.CLI()
    
    // Run a container with custom arguments
    output, err := cli.Run(ctx, "alpine", "latest", []string{
        "echo",
        "Hello from Docker!",
    })
    
    return output, err
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Docker Operations
on: [push]

jobs:
  docker-ops:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Docker Operations with Dagger
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/docker
          args: |
            do -p '
              docker := Docker()
              cli := docker.CLI()
              cli.Pull(ctx, "alpine", "latest")
            '
```

## API Reference

### Docker

Main module struct that provides access to Docker functionality.

#### Methods

- `Engine(version string, persist bool, namespace string) *Service`
  - Spawns an ephemeral Docker Engine
  - Parameters:
    - `version`: Docker Engine version (default: "24.0")
    - `persist`: Whether to persist engine state (default: true)
    - `namespace`: Namespace for persistence

- `CLI(version string, engine *Service) *CLI`
  - Creates a Docker CLI instance
  - Parameters:
    - `version`: Docker CLI version (default: "24.0")
    - `engine`: Optional custom engine instance

### CLI

Docker CLI wrapper providing command execution capabilities.

#### Methods

- `Pull(repository string, tag string) (*Image, error)`
  - Pulls an image from a registry
  
- `Push(repository string, tag string) (string, error)`
  - Pushes an image to a registry
  
- `Import(container *Container) (*Image, error)`
  - Imports a container as an image
  
- `Run(name string, tag string, args []string) (string, error)`
  - Runs a container with specified arguments

### Image

Represents a Docker image in the engine.

#### Methods

- `Export() *Container`
  - Exports the image as a container
  
- `Duplicate(repository string, tag string) (*Image, error)`
  - Creates a copy of the image with new tags
  
- `Push(ctx context.Context) (string, error)`
  - Pushes the image to a registry

## Best Practices

1. **Engine Persistence**
   - Use persistence for long-running pipelines
   - Use namespaces to isolate different projects

2. **Resource Management**
   - Clean up unused images and containers
   - Use appropriate tags for versioning

3. **Security**
   - Avoid running containers with root privileges
   - Use specific versions instead of 'latest' tag

## Troubleshooting

Common issues and solutions:

1. **Connection Issues**
   ```
   Error: Cannot connect to the Docker daemon
   Solution: Ensure the Docker engine is running and accessible
   ```

2. **Permission Issues**
   ```
   Error: Permission denied
   Solution: Check if the necessary capabilities are granted
   ```

3. **Image Pull Failures**
   ```
   Error: Image pull failed
   Solution: Verify registry credentials and image name/tag
   ``` 