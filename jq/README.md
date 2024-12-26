# JQ Module for Dagger

A Dagger module that provides integration with `jq`, a lightweight and portable command-line JSON processor. This module allows you to install and use `jq` in your Dagger pipelines across different platforms and container environments.

## Features

- Cross-platform support (Linux, macOS, Windows)
- Multiple architecture support (amd64, i386, armhf, arm64, riscv64)
- Binary verification with SHA256 checksums
- Container integration
- Red Hat Universal Base Image support
- Customizable installation prefix
- Platform-specific binary handling

## Usage

### Basic Setup

```typescript
import { jq } from "@felipepimentel/daggerverse/jq";

// Initialize the JQ module with a specific version
const client = jq("1.7");
```

### Get JQ Binary

```typescript
// Get jq binary for a specific platform
const binary = await client.binary(dag.platform("linux/amd64"));
```

### Create Root Filesystem Overlay

```typescript
// Create an overlay with jq installed
const overlay = await client.overlay(
  dag.platform("linux/amd64"),
  "/usr/local" // Optional: custom prefix
);
```

### Install in Container

```typescript
// Install jq in an existing container
const container = await client.installation(baseContainer);
```

### Get JQ Container

```typescript
// Get a container with jq as entrypoint
const container = await client.container(baseContainer);
```

## Red Hat Container Integration

### Standard UBI Container

```typescript
// Get a Red Hat Universal Base Image container with jq
const container = await client.redhatContainer(dag.platform("linux/amd64"));
```

### Minimal UBI Container

```typescript
// Get a Red Hat Minimal Universal Base Image container with jq
const container = await client.redhatMinimalContainer(
  dag.platform("linux/amd64")
);
```

### Micro UBI Container

```typescript
// Get a Red Hat Micro Universal Base Image container with jq
const container = await client.redhatMicroContainer(
  dag.platform("linux/amd64")
);
```

## Examples

### Process JSON in Pipeline

```typescript
import { jq } from "@felipepimentel/daggerverse/jq";

export async function processJSON() {
  // Initialize JQ
  const client = jq("1.7");

  // Get a container with jq
  const container = await client.redhatContainer();

  // Process JSON
  const result = await container
    .withMountedFile("/work/data.json", jsonFile)
    .withExec(["-r", ".field", "/work/data.json"]);
}
```

### Custom Installation

```typescript
import { jq } from "@felipepimentel/daggerverse/jq";

export async function customInstall() {
  // Initialize JQ
  const client = jq("1.7");

  // Create overlay with custom prefix
  const overlay = await client.overlay(
    dag.platform("linux/amd64"),
    "/opt/tools"
  );

  // Use overlay in container
  const container = baseContainer.withDirectory("/", overlay);
}
```

### Multi-Platform Build

```typescript
import { jq } from "@felipepimentel/daggerverse/jq";

export async function multiPlatformBuild() {
  // Initialize JQ
  const client = jq("1.7");

  const platforms = [
    "linux/amd64",
    "linux/arm64",
    "darwin/amd64",
    "windows/amd64",
  ];

  for (const platform of platforms) {
    const binary = await client.binary(dag.platform(platform));
    // Use binary...
  }
}
```

## Configuration

### Constructor Options

The module accepts:

- `version`: Version of jq to install (e.g., "1.7")

### Binary Options

The `binary` method accepts:

- `platform`: Target platform in format "os/arch" (optional, defaults to host platform)

### Overlay Options

The `overlay` method accepts:

- `platform`: Target platform (optional)
- `prefix`: Installation prefix (optional, defaults to "/usr/local")

### Container Options

The `container` method accepts:

- `container`: Base container to install jq into

## Dependencies

The module requires:

- Dagger SDK
- Internet access to download jq binaries
- Red Hat module (for Red Hat container integration)

## License

See [LICENSE](../LICENSE) file in the root directory.
