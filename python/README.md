# Python Module for Dagger

A Dagger module that provides integration with Python, allowing you to run Python applications and manage Python environments in your Dagger pipelines.

## Features

- Python environment setup and configuration
- Pip package management with cache support
- Source code mounting and management
- Custom base container support
- Cross-platform compatibility
- Workspace directory management

## Usage

### Basic Setup

```typescript
import { python } from "@felipepimentel/daggerverse/python";

// Initialize Python with default settings
const client = python();
```

### Custom Version

```typescript
// Initialize Python with a specific version
const client = python("3.12");
```

### Custom Container

```typescript
// Initialize Python with a custom container
const container = dag
  .container()
  .from("python:3.12-slim")
  .withEnvVariable("PYTHONUNBUFFERED", "1");

const client = python("", container);
```

### With Pip Cache

```typescript
// Create a cache volume for pip
const pipCache = dag.cacheVolume("pip-cache");

// Initialize Python with pip cache
const client = python().withPipCache(pipCache);
```

### With Source Code

```typescript
// Mount source code directory
const client = python().withSource(sourceDir);
```

## Configuration

### Constructor Options

The module accepts:

- `version`: Python version to use (default: "latest")
- `container`: Custom container configuration (optional)

### Additional Methods

- `withPipCache`: Configure pip cache with options:

  - `cache`: Cache volume for pip packages
  - `source`: Optional directory to use as cache root
  - `sharing`: Cache sharing mode

- `withSource`: Mount source code with options:
  - `source`: Directory containing Python source code
  - Default working directory: `/work`

## Examples

### Development Setup

```typescript
import { python } from "@felipepimentel/daggerverse/python";

export async function devSetup() {
  // Create cache volume for pip
  const pipCache = dag.cacheVolume("pip-cache");

  // Initialize Python with development configuration
  const client = python("3.12").withPipCache(pipCache);
}
```

### Project Setup

```typescript
import { python } from "@felipepimentel/daggerverse/python";

export async function projectSetup() {
  // Create source directory
  const source = dag
    .directory()
    .withNewFile(
      "requirements.txt",
      `
      flask==3.0.0
      requests==2.31.0
    `
    )
    .withNewFile(
      "app.py",
      `
      from flask import Flask
      app = Flask(__name__)

      @app.route('/')
      def hello():
          return 'Hello, World!'
    `
    );

  // Initialize Python with project configuration
  const client = python()
    .withSource(source)
    .withPipCache(dag.cacheVolume("pip-cache"));
}
```

### Production Setup

```typescript
import { python } from "@felipepimentel/daggerverse/python";

export async function productionSetup() {
  // Create custom container with production settings
  const container = dag
    .container()
    .from("python:3.12-slim")
    .withEnvVariable("PYTHONUNBUFFERED", "1")
    .withEnvVariable("PYTHONOPTIMIZE", "2");

  // Initialize Python with production configuration
  const client = python("", container)
    .withSource(productionCode)
    .withPipCache(productionCache);
}
```

## Dependencies

The module requires:

- Dagger SDK
- Internet access to pull Python images
- Optional: pip cache volume for package management

## Testing

The module includes a test suite that verifies:

- Python version management
- Pip cache functionality
- Source code mounting
- Container configuration
- Environment variables

To run the tests:

```bash
dagger do test
```

## License

See [LICENSE](../LICENSE) file in the root directory.
