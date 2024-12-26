# Red Hat Module for Dagger

A Dagger module that provides integration with Red Hat Universal Base Images (UBI). This module allows you to create and customize containers based on Red Hat UBI, with support for standard, minimal, and micro variants.

## Features

- Support for multiple UBI variants:
  - Standard UBI (ubi9)
  - Minimal UBI (ubi9-minimal)
  - Micro UBI (ubi9-micro)
- Package management with DNF/MicroDNF
- Module management for Red Hat packages
- CA certificates management
- Cross-platform support
- Optimized package installation with minimal dependencies

## Usage

### Basic Setup

```go
// Initialize the Red Hat module
redhat := dag.Redhat().New()

// Get a standard UBI container
container := redhat.Container(
    "linux/amd64", // Optional platform specification
)
```

### Container Variants

#### Standard UBI

```go
// Get a standard UBI container
container := redhat.Container()

// Install packages
container = redhat.Packages([]string{
    "httpd",
    "openssl",
}).Installed(container)
```

#### Minimal UBI

```go
// Get a minimal UBI container
minimal := redhat.Minimal()
container := minimal.Container()

// Install packages
container = minimal.Packages([]string{
    "curl",
    "tar",
}).Installed(container)
```

#### Micro UBI

```go
// Get a micro UBI container
micro := redhat.Micro()
container := micro.Container()
```

## Package Management

### Standard UBI

```go
// Install packages
container = redhat.Packages([]string{
    "git",
    "make",
    "gcc",
}).Installed(container)

// Remove packages
container = redhat.Packages([]string{
    "git",
}).Removed(container)
```

### Minimal UBI

```go
minimal := redhat.Minimal()

// Install packages
container = minimal.Packages([]string{
    "python3",
    "pip",
}).Installed(container)

// Remove packages
container = minimal.Packages([]string{
    "python3",
}).Removed(container)
```

## Module Management

### Standard UBI

```go
// Enable a module
container = redhat.Module("nodejs:18").Enabled(container)

// Disable a module
container = redhat.Module("nodejs:18").Disabled(container)
```

### Minimal UBI

```go
minimal := redhat.Minimal()

// Enable a module
container = minimal.Module("nodejs:18").Enabled(container)

// Disable a module
container = minimal.Module("nodejs:18").Disabled(container)
```

## CA Certificates

```go
// Get CA certificates from UBI
caCertificates := redhat.CaCertificates()
```

## Container Configuration

All container variants support:

- Platform specification
- Working directory configuration (default: /home)
- Package installation optimization flags

## Image Details

### Standard UBI

- Registry: registry.access.redhat.com
- Repository: ubi9
- Tag: 9.5-1732804088
- Digest: sha256:b632d0cc6263372a90e9097dcac0a369e456b144a66026b9eac029a22f0f6e07

### Minimal UBI

- Registry: registry.access.redhat.com
- Repository: ubi9-minimal
- Tag: 9.5-1733767867
- Digest: sha256:f598528219a1be07cf520fbe82a2d2434dc9841e1f0a878382c8a13bf42cb486

### Micro UBI

- Registry: registry.access.redhat.com
- Repository: ubi9-micro
- Tag: 9.5-1733767087
- Digest: sha256:3313e52bb1aad4017a0c35f9f2ae35cf8526eeeb83f6ecbec449ba9c5cb9cb07

## Examples

### Create a Development Container

```go
func CreateDevContainer(ctx context.Context) (*Container, error) {
    redhat := dag.Redhat().New()

    // Get a container with development tools
    container := redhat.Container().
        With(redhat.Packages([]string{
            "git",
            "make",
            "gcc",
            "gcc-c++",
        }).Installed)

    return container, nil
}
```

### Create a Minimal Web Server

```go
func CreateWebServer(ctx context.Context) (*Container, error) {
    redhat := dag.Redhat().New()
    minimal := redhat.Minimal()

    // Get a minimal container with httpd
    container := minimal.Container().
        With(minimal.Packages([]string{
            "httpd",
        }).Installed)

    return container, nil
}
```

## Package Installation Options

The module uses optimized package installation flags:

- `--nodocs`: Exclude documentation
- `--setopt install_weak_deps=0`: Skip weak dependencies
- `--assumeyes`: Automatic yes to prompts

## Dependencies

The module requires:

- Dagger SDK
- Internet access to Red Hat container registry
- DNF/MicroDNF package manager (included in base images)

## License

This module is part of the Dagger project and is subject to its licensing terms.

Copyright Camptocamp SA
SPDX-License-Identifier: AGPL-3.0-or-later
