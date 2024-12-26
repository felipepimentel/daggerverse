# Kustomize Module for Dagger

A Dagger module that provides native Kubernetes configuration management using Kustomize. This module allows you to build, edit, and manage Kubernetes manifests using Kustomize's powerful customization features.

## Features

- Build Kustomize targets from directories or URLs
- Edit Kustomization files programmatically
- Support for all Kustomize configuration options
- Custom container and version selection
- Clean path handling for subdirectories
- Secure and isolated builds

## Usage

### Basic Build

```typescript
import { kustomize } from "@felipepimentel/daggerverse/kustomize";

// Initialize the Kustomize module
const client = kustomize();

// Build a kustomization target
const result = await client.build({
  source, // Source directory containing kustomization.yaml
  dir: "path/to/overlay", // Optional subdirectory
});
```

### Custom Version

```typescript
// Use a specific version of Kustomize
const client = kustomize({
  version: "v5.0.1",
});
```

### Custom Container

```typescript
// Use a custom container
const container = dag
  .container()
  .from("custom-image:latest")
  .withEnvVariable("CUSTOM_VAR", "value");

const client = kustomize({
  container: container,
});
```

## Examples

### Building a Basic Kustomization

```typescript
import { kustomize } from "@felipepimentel/daggerverse/kustomize";

export async function buildKustomization() {
  // Initialize module
  const client = kustomize();

  // Get source directory
  const source = dag.host().directory("./k8s");

  // Build kustomization
  const result = await client.build({
    source,
    dir: "overlays/production",
  });

  // Get the output
  const output = await result.contents();
}
```

### Editing Kustomization Files

```typescript
import { kustomize } from "@felipepimentel/daggerverse/kustomize";

export async function editKustomization() {
  // Initialize module
  const client = kustomize();

  // Get source directory
  const source = dag.host().directory("./k8s/base");

  // Edit kustomization
  const edited = await client
    .edit(source)
    .set()
    .annotation("environment", "production")
    .set()
    .image("nginx:1.16")
    .set()
    .namespace("prod")
    .directory();
}
```

### Complex Build Configuration

```typescript
import { kustomize } from "@felipepimentel/daggerverse/kustomize";

export async function complexBuild() {
  // Initialize with custom version
  const client = kustomize({
    version: "v5.0.1",
  });

  // Get source directory
  const source = dag.host().directory("./k8s");

  // Build with multiple layers
  const base = await client.build({ source, dir: "base" });
  const dev = await client.build({ source, dir: "overlays/dev" });
  const prod = await client.build({ source, dir: "overlays/prod" });
}
```

## API Reference

### Constructor Options

The module accepts:

- `version`: Kustomize version to use (optional)
- `container`: Custom container to use (optional)

### Methods

#### `build(options: BuildOptions)`

Build a kustomization target.

Parameters:

- `source`: Source directory containing kustomization.yaml
- `dir`: Optional subdirectory path

Returns a file with the built manifests.

#### `edit(source: Directory, dir?: string)`

Edit a kustomization file.

Parameters:

- `source`: Source directory containing kustomization.yaml
- `dir`: Optional subdirectory path

Returns an Edit interface for modifying the kustomization.

### Edit Interface

Methods for editing kustomization files:

- `directory()`: Get the modified directory
- `set()`: Set values in the kustomization

### Set Interface

Methods for setting values in kustomization files:

- `annotation(key: string, value: string)`: Set annotations
- `image(image: string)`: Set image
- `namespace(namespace: string)`: Set namespace

## Testing

The module includes a test suite that can be run using:

```bash
dagger do test
```

The test suite verifies:

- Basic kustomization building
- Kustomization file editing
- Path handling
- Configuration options

## Dependencies

The module requires:

- Dagger SDK
- Kustomize image (pulled automatically)
- Internet access for initial container pulling

## Implementation Details

- Uses official Kustomize image from registry.k8s.io
- Supports Kustomize v5.0.1 by default
- Implements clean path handling for subdirectories
- Provides both build and edit functionality
- Supports custom container configurations

## License

See [LICENSE](../LICENSE) file in the root directory.
