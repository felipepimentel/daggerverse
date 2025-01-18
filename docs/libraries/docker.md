---
layout: default
title: Docker Module
parent: Libraries
nav_order: 1
---

# Docker Module

This module provides a comprehensive integration with Docker Engine, allowing you to manage containers, images, and Docker operations directly through Dagger.

## Features

- Ephemeral Docker Engine management
- Docker CLI integration
- Image pulling and pushing
- Container import/export
- State persistence
- Registry authentication
- Image management
- Container execution

## Installation

```bash
dagger mod use github.com/felipepimentel/daggerverse/libraries/docker@latest
```

## Usage

### Basic Example

```go
// Initialize the module
docker := dag.Docker()

// Create a Docker engine
engine := docker.Engine("24.0", true, "my-namespace")

// Get a CLI instance
cli := docker.CLI("24.0", engine)
```

## Docker Engine

### Starting an Engine

```go
// Start an ephemeral engine
engine := docker.Engine("24.0", false, "")

// Start a persistent engine
engine := docker.Engine("24.0", true, "my-namespace")
```

### Configuration Options

```go
func (e *Docker) Engine(
    // Docker Engine version
    version string,    // default: "24.0"
    // Persist the state
    persist bool,      // default: true
    // Namespace for state
    namespace string,
) *dagger.Service
```

## Docker CLI

### Basic CLI Operations

```go
// Get a CLI instance
cli := docker.CLI("24.0", nil)  // nil for ephemeral engine

// Pull an image
image, err := cli.Pull(ctx, "alpine", "latest")

// Push an image
ref, err := cli.Push(ctx, "my-registry/alpine", "latest")
```

### Chaining Operations

```go
// Pull and chain operations
cli, err := cli.WithPull(ctx, "alpine", "latest")
if err != nil {
    return err
}

// Push and chain operations
cli, err = cli.WithPush(ctx, "my-registry/alpine", "latest")
```

### Container Management

```go
// Import a container
container := dag.Container().From("alpine:latest")
image, err := cli.Import(ctx, container)

// Run a container
output, err := cli.Run(ctx, "alpine", "latest", []string{"echo", "hello"})
```

## Image Management

### Image Operations

```go
// Look up an image
image, err := cli.Image(ctx, "alpine", "latest", "")

// List images
images, err := cli.Images(ctx, "alpine", "latest", "")

// Duplicate an image
newImage, err := image.Duplicate(ctx, "my-registry/alpine", "v2")

// Export an image
container := image.Export()

// Push an image
ref, err := image.Push(ctx)
```

## Best Practices

1. **Engine Management**:
   - Use persistent engines for long-running operations
   - Use namespaces to isolate different workloads
   - Clean up unused engines

2. **Image Handling**:
   - Tag images appropriately
   - Use specific versions instead of 'latest'
   - Clean up unused images

3. **Performance**:
   - Reuse CLI instances
   - Use image caching
   - Optimize container layers

## Common Issues

1. **Engine Connection**:
   - Verify engine version compatibility
   - Check network connectivity
   - Validate service bindings

2. **Image Operations**:
   - Ensure registry authentication
   - Check image name format
   - Verify pull/push permissions

3. **Container Operations**:
   - Monitor resource usage
   - Check container logs
   - Validate mount points

## Examples

### Complete Workflow

```go
// Initialize Docker
docker := dag.Docker()

// Create persistent engine
engine := docker.Engine("24.0", true, "production")

// Get CLI
cli := docker.CLI("24.0", engine)

// Pull base image
baseImage, err := cli.Pull(ctx, "alpine", "latest")
if err != nil {
    return err
}

// Create custom image
container := dag.Container().
    From("alpine:latest").
    WithExec([]string{"apk", "add", "python3"})

// Import custom image
customImage, err := cli.Import(ctx, container)
if err != nil {
    return err
}

// Tag and push
newImage, err := customImage.Duplicate(ctx, "my-registry/python-alpine", "v1")
if err != nil {
    return err
}

// Push to registry
ref, err := newImage.Push(ctx)
```

### Registry Authentication

```go
// Create authenticated CLI
cli := docker.CLI("24.0", nil).
    WithEnvVariable("DOCKER_USERNAME", "user").
    WithEnvVariable("DOCKER_PASSWORD", dag.SetSecret("docker_password", "password"))

// Push to private registry
ref, err := cli.Push(ctx, "private-registry/image", "latest")
```

### Custom Engine Configuration

```go
// Create engine with custom configuration
engine := docker.Engine("24.0", true, "custom").
    WithEnvVariable("DOCKER_TLS_VERIFY", "1").
    WithMountedDirectory("/certs", dag.Host().Directory("./certs"))

// Use custom engine
cli := docker.CLI("24.0", engine)
```

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](../CONTRIBUTING.md) for details on how to submit pull requests.

## License

This module is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details. 