# Registry Module for Dagger

A Dagger module that provides integration with Docker Registry, a stateless, highly scalable server-side application for storing and distributing container images. This module enables you to run and manage a Docker Registry instance within your Dagger pipelines.

## Features

- Docker Registry deployment and configuration
- Version selection from official images
- Custom container support
- Port configuration
- Data persistence with cache volumes
- Service exposure
- Environment variable management

## Usage

### Basic Setup

```go
// Initialize Registry with default settings
registry, err := dag.Registry().New(ctx)
if err != nil {
    return err
}

// Get the Registry service
service := registry.Service()
```

### Custom Configuration

```go
// Initialize Registry with custom settings
registry, err := dag.Registry().New(
    ctx,
    "2.8",                    // Version
    nil,                      // Custom container (optional)
    5000,                     // Port
    dag.CacheVolume("data"), // Data volume (optional)
)
if err != nil {
    return err
}
```

### Production Setup

```go
func ProductionSetup(ctx context.Context) error {
    // Use custom container with production settings
    container := dag.Container().
        From("registry:2.8").
        WithEnvVariable("REGISTRY_LOGLEVEL", "info").
        WithEnvVariable("REGISTRY_STORAGE_DELETE_ENABLED", "true")

    // Initialize Registry with production settings
    registry, err := dag.Registry().New(
        ctx,
        "",
        container,
        5000,
        persistentVolume,
    )
    if err != nil {
        return err
    }

    return nil
}
```

## Configuration Options

### Version

- Specifies the version (tag) of the official Docker Registry image to use
- Default: "2.8"
- Optional: Set to empty string when using a custom container

### Custom Container

- Allows using a custom container instead of the official image
- Takes precedence over version setting
- Optional: Default is nil

### Port

- Port number to expose the registry on
- Default: 5000
- Optional: Set to 0 to use default

### Data Volume

- Cache volume for persisting registry data between runs
- Mounted at /var/lib/registry
- Optional: Default is nil (no persistence)

## Dependencies

The module requires:

- Dagger SDK
- Internet access to pull Registry images
- Optional: Cache volume for data persistence

## Testing

The module includes tests that verify:

- Basic Registry initialization
- Custom configuration handling
- Service operations
- Data persistence

To run the tests:

```bash
dagger do test
```

## License

This module is part of the Dagger project and is subject to its licensing terms.

Copyright Camptocamp SA
SPDX-License-Identifier: AGPL-3.0-or-later
