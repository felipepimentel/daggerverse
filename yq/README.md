# YQ Module for Dagger

A Dagger module that provides integration with `yq`, a lightweight and portable command-line YAML processor. This module allows you to read, update, and manipulate YAML files in your Dagger pipelines.

## Features

- YAML file manipulation and querying
- Support for in-place file editing
- Directory mounting for batch operations
- Container-based execution
- Configurable image and version selection
- Shell access for debugging

## Usage

### Basic Setup

```typescript
import { yq } from "@felipepimentel/daggerverse/yq";

// Initialize the YQ module with default settings
const client = yq({
  image: "", // Optional: custom image
  version: "", // Optional: specific version
  source: sourceDirectory, // Directory containing YAML files
});
```

### Reading YAML Values

```typescript
// Get a value from a YAML file
const value = await client.get(
  ".metadata.name", // YQ expression
  "config.yaml" // YAML file path
);
```

### Modifying YAML Files

```typescript
// Update a value in a YAML file
const withUpdate = client.set(
  ".spec.replicas = 3", // YQ expression
  "deployment.yaml" // YAML file path
);

// Get the modified directory
const result = withUpdate.state();
```

### Working with Directories

```typescript
// Override the source directory
const withDir = client.withDirectory(newSourceDir);

// Get the current directory state
const currentDir = withDir.state();
```

## Configuration

### Constructor Options

The module accepts:

- `image`: Container image to use (default: "mikefarah/yq")
- `version`: YQ version to use (default: "4.35.2")
- `source`: Directory containing YAML files

### Default Settings

- Working directory: `/opt/`
- Default entrypoint: `yq`
- File ownership: `yq` user

## Container Access

### Get Container

```typescript
// Access the underlying container
const container = client.container();
```

### Debug Shell

```typescript
// Open an interactive shell
const shell = client.shell();
```

## Examples

### Update Kubernetes Configuration

```typescript
import { yq } from "@felipepimentel/daggerverse/yq";

export async function updateKubeConfig() {
  // Initialize YQ with Kubernetes manifests
  const client = yq({
    source: dag.directory().withFiles({
      "deployment.yaml": `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 1
`,
    }),
  });

  // Update replicas
  const withUpdate = client.set(".spec.replicas = 3", "deployment.yaml");

  // Get the result
  const result = withUpdate.state();
}
```

### Read Multiple Values

```typescript
import { yq } from "@felipepimentel/daggerverse/yq";

export async function readConfig() {
  const client = yq({
    source: configDir,
  });

  // Get multiple values
  const name = await client.get(".metadata.name", "config.yaml");
  const version = await client.get(".spec.version", "config.yaml");

  return `${name}:${version}`;
}
```

### Batch Processing

```typescript
import { yq } from "@felipepimentel/daggerverse/yq";

export async function processDirectory() {
  const client = yq({
    source: yamlDir,
  });

  // Update all deployments
  const withUpdate = client.set(
    '.metadata.labels.environment = "production"',
    "*.yaml"
  );
}
```

## Dependencies

The module requires:

- Dagger SDK
- Internet access to pull the YQ container image
- YAML files to process

## License

See [LICENSE](../LICENSE) file in the root directory.
