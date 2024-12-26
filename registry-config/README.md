# Registry Config Module for Dagger

A Dagger module that safely creates and manages OCI registry configuration files for tools like Helm and Oras. This module addresses security concerns around storing registry credentials in the filesystem by providing a secure way to handle authentication through Dagger's secret management.

## Features

- Secure registry authentication configuration
- Support for multiple registries
- Compatible with tools like Helm and Oras
- Safe credential handling through Dagger secrets
- Prevents credentials from leaking into layer cache
- Read-only configuration mounting

## Usage

### Basic Configuration

```go
// Initialize the Registry Config module
config := dag.RegistryConfig()

// Add registry authentication
config = config.WithRegistryAuth(
    "ghcr.io",           // Registry address
    "username",          // Registry username
    dag.SetSecret("GITHUB_TOKEN", token), // Registry password/token as a secret
)
```

### Multiple Registries

```go
// Configure multiple registries
config := dag.RegistryConfig().
    WithRegistryAuth(
        "ghcr.io",
        "username1",
        dag.SetSecret("GITHUB_TOKEN", githubToken),
    ).
    WithRegistryAuth(
        "docker.io",
        "username2",
        dag.SetSecret("DOCKER_TOKEN", dockerToken),
    )
```

### Mounting Configuration

```go
// Mount the config in a container
container := dag.Container().From("alpine")
container = container.WithMountedSecret(
    "/config.json",
    config.Secret(),
)
```

## Examples

### Using with Helm

```go
func HelmWithRegistry(ctx context.Context) error {
    // Create registry config
    config := dag.RegistryConfig().
        WithRegistryAuth(
            "ghcr.io",
            "username",
            dag.SetSecret("REGISTRY_TOKEN", token),
        )

    // Use with Helm
    container := dag.Container().From("alpine/helm")

    // Mount the config
    container = container.WithMountedSecret(
        "/config.json",
        config.Secret(),
    )

    // Run Helm commands
    result := container.WithExec([]string{
        "pull",
        "oci://ghcr.io/org/chart",
    })

    return nil
}
```

### Using with Oras

```go
func OrasWithRegistry(ctx context.Context) error {
    // Create registry config
    config := dag.RegistryConfig().
        WithRegistryAuth(
            "docker.io",
            "username",
            dag.SetSecret("DOCKER_TOKEN", token),
        )

    // Use with Oras
    container := dag.Container().From("oras/oras")

    // Mount the config
    container = container.WithMountedSecret(
        "/config.json",
        config.Secret(),
    )

    // Run Oras commands
    result := container.WithExec([]string{
        "pull",
        "docker.io/org/artifact:tag",
    })

    return nil
}
```

## API Reference

### RegistryConfig

The main module interface:

```go
type RegistryConfig struct{}

// Add authentication for a registry
func (m *RegistryConfig) WithRegistryAuth(
    address string,
    username string,
    secret *Secret,
) *RegistryConfig

// Get the configuration as a secret
func (m *RegistryConfig) Secret() *Secret

// Mount the configuration as a secret in a container
func (m *RegistryConfig) MountSecret(
    container *Container,
    path string,
) *Container
```

## Important Notes

1. **Security**: This module is designed to handle registry credentials securely through Dagger's secret management system, preventing credentials from being exposed in the filesystem or layer cache.

2. **Read-only Configuration**: The configuration file is mounted read-only, which means tools' built-in authentication mechanisms (like `helm registry login`) may not work as they often try to modify the config file.

3. **Compatibility**: While designed primarily for Helm and Oras, this module should work with any tool that uses the standard OCI registry configuration format.

## Testing

The module includes a test suite that can be run using:

```bash
dagger do test
```

The test suite verifies:

- Basic configuration generation
- Multiple registry support
- Secret mounting
- Configuration format compliance

## Dependencies

The module requires:

- Go 1.22 or later
- Dagger SDK
- Tools that use standard OCI registry configuration format

## Implementation Details

- Uses JSON format for registry configuration
- Implements secure secret handling
- Supports multiple registry configurations
- Provides read-only configuration mounting
- Prevents credential leakage into layer cache

## License

This module is part of the Dagger project and is subject to its licensing terms.

Copyright Camptocamp SA
SPDX-License-Identifier: AGPL-3.0-or-later
