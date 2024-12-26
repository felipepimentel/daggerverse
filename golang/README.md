# Golang Module for Dagger

A Dagger module that provides integration with Go in Red Hat Universal Base Images (UBI). This module allows you to install and configure Go in your Dagger pipelines, with support for both standard and minimal UBI containers.

## Features

- Go installation in Red Hat Universal Base Images
- Automatic cache configuration for Go modules and build
- Git integration for package management
- Support for both standard and minimal UBI containers
- Cross-platform support through platform specification
- Persistent cache management for improved performance

## Usage

### Basic Setup

```typescript
import { golang } from "@felipepimentel/daggerverse/golang";

// Initialize the Golang module
const client = golang();

// Get a container with Go installed
const container = await client.redhatContainer();
```

### Container Integration

The module provides two main container types:

#### Standard Red Hat UBI

```typescript
// Get a standard Red Hat UBI container with Go
const container = await client.redhatContainer({
  platform: "linux/amd64", // Optional platform specification
});

// Configure an existing container
const configured = await client.redhatInstallation(existingContainer);
```

#### Minimal Red Hat UBI

```typescript
// Get a minimal Red Hat UBI container with Go
const container = await client.redhatMinimalContainer({
  platform: "linux/amd64", // Optional platform specification
});

// Configure an existing minimal container
const configured = await client.redhatMinimalInstallation(existingContainer);
```

## Configuration

### Cache Configuration

The module automatically configures Go caches:

```typescript
// Default cache configuration
const configured = await client.configuration(container);
```

This sets up:

- `GOPATH` at `/var/cache/go`
- `GOCACHE` at `/var/cache/go/build`
- Persistent cache volume named "golang"

### Installation Options

Both standard and minimal installations include:

- Go compiler and tools
- Git for package management
- Cache configuration
- Environment variable setup

## Platform Support

The module supports various platforms through the platform parameter in container creation:

- Linux (amd64, arm64)
- Other platforms supported by Red Hat UBI

## Examples

### Create a Development Container

```typescript
import { golang } from "@felipepimentel/daggerverse/golang";

export async function createDevContainer() {
  const client = golang();

  // Get a container with Go installed
  const container = await client.redhatContainer({
    platform: "linux/amd64",
  });

  return container;
}
```

### Configure Existing Container

```typescript
import { golang } from "@felipepimentel/daggerverse/golang";

export async function configureContainer(base: Container) {
  const client = golang();

  // Install Go in the base container
  const container = await client.redhatInstallation(base);

  return container;
}
```

### Use Minimal Container

```typescript
import { golang } from "@felipepimentel/daggerverse/golang";

export async function createMinimalContainer() {
  const client = golang();

  // Get a minimal container with Go
  const container = await client.redhatMinimalContainer({
    platform: "linux/amd64",
  });

  return container;
}
```

## Dependencies

The module requires:

- Dagger SDK
- Red Hat module (automatically included as a dependency)
- Internet access to install packages

## Cache Management

The module uses Dagger's cache volume system to persist:

- Go module cache
- Go build cache
- Package downloads

Cache location: `/var/cache/go`

## License

See [LICENSE](../LICENSE) file in the root directory.
