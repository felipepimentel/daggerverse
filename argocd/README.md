# Argo CD Module for Dagger

A Dagger module that provides integration with Argo CD CLI (`argocd`), allowing you to interact with Argo CD in your Dagger pipelines. This module enables you to install and use the Argo CD command-line interface across different platforms and container environments.

## Features

- Cross-platform support (Linux, macOS, Windows)
- Multiple architecture support
- Binary verification with SHA256 checksums
- Container integration
- Red Hat Universal Base Image support
- Customizable installation prefix
- Platform-specific binary handling

## Usage

### Basic Setup

```typescript
import { argocd } from "@felipepimentel/daggerverse/argocd";

// Initialize the Argo CD module with a specific version
const client = argocd("2.10.0");
```

### Get Argo CD Binary

```typescript
// Get argocd binary for a specific platform
const binary = await client.binary({
  platform: "linux/amd64",
});
```

### Create Root Filesystem Overlay

```typescript
// Create an overlay with argocd installed
const overlay = await client.overlay({
  platform: "linux/amd64",
  prefix: "/usr/local", // Optional: custom prefix
});
```

### Install in Container

```typescript
// Install argocd in an existing container
const container = await client.installation(baseContainer);
```

### Get Argo CD Container

```typescript
// Get a container with argocd as entrypoint
const container = await client.container(baseContainer);
```

## Red Hat Container Integration

### Standard UBI Container

```typescript
// Get a Red Hat Universal Base Image container with argocd
const container = await client.redhatContainer({
  platform: "linux/amd64",
});
```

### Minimal UBI Container

```typescript
// Get a Red Hat Minimal Universal Base Image container with argocd
const container = await client.redhatMinimalContainer({
  platform: "linux/amd64",
});
```

### Micro UBI Container

```typescript
// Get a Red Hat Micro Universal Base Image container with argocd
const container = await client.redhatMicroContainer({
  platform: "linux/amd64",
});
```

## Examples

### Manage Applications

```typescript
import { argocd } from "@felipepimentel/daggerverse/argocd";

export async function manageApplication() {
  // Initialize Argo CD
  const client = argocd("2.10.0");

  // Get a container with argocd
  const container = await client.redhatContainer();

  // Login and manage applications
  await container
    .withSecretVariable("ARGOCD_AUTH_TOKEN", authToken)
    .withExec([
      "app",
      "create",
      "myapp",
      "--repo",
      "https://github.com/org/repo",
      "--path",
      "manifests",
      "--dest-server",
      "https://kubernetes.default.svc",
      "--dest-namespace",
      "default",
    ]);
}
```

### Custom Installation

```typescript
import { argocd } from "@felipepimentel/daggerverse/argocd";

export async function customInstall() {
  // Initialize Argo CD
  const client = argocd("2.10.0");

  // Create overlay with custom prefix
  const overlay = await client.overlay({
    platform: "linux/amd64",
    prefix: "/opt/tools",
  });

  // Use overlay in container
  const container = baseContainer.withDirectory("/", overlay);
}
```

### Multi-Platform Build

```typescript
import { argocd } from "@felipepimentel/daggerverse/argocd";

export async function multiPlatformBuild() {
  // Initialize Argo CD
  const client = argocd("2.10.0");

  const platforms = [
    "linux/amd64",
    "linux/arm64",
    "darwin/amd64",
    "windows/amd64",
  ];

  for (const platform of platforms) {
    const binary = await client.binary({ platform });
    // Use binary...
  }
}
```

## Configuration

### Constructor Options

The module accepts:

- `version`: Version of Argo CD CLI to install (e.g., "2.10.0")

### Binary Options

The `binary` method accepts:

- `platform`: Target platform in format "os/arch" (optional, defaults to host platform)

### Overlay Options

The `overlay` method accepts:

- `platform`: Target platform (optional)
- `prefix`: Installation prefix (optional, defaults to "/usr/local")

### Container Options

The `container` method accepts:

- `container`: Base container to install Argo CD CLI into

## Dependencies

The module requires:

- Dagger SDK
- Internet access to download Argo CD binaries
- Red Hat module (for Red Hat container integration)

## License

See [LICENSE](../LICENSE) file in the root directory.
