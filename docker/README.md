# Docker Dagger Module

This Dagger module provides integration with Docker Engine, allowing you to run Docker operations within your Dagger pipelines. It includes support for running an ephemeral Docker Engine, executing Docker CLI commands, and managing images.

## Features

- Spawn ephemeral Docker Engine instances
- Execute Docker CLI commands
- Pull and push Docker images
- Import Dagger containers into Docker Engine
- Manage Docker images and containers
- Persistent engine state with cache volumes
- Customizable Docker Engine and CLI versions

## Usage

### Basic Usage

```typescript
import { docker } from "@felipepimentel/daggerverse/docker";

// Initialize Docker module
const client = docker();

// Create a Docker CLI instance with default engine
const cli = client.cli();

// Pull an image
const image = await cli.pull("alpine", "latest");
```

### Advanced Usage

#### Custom Engine Configuration

```typescript
// Create an engine with custom version and persistence
const engine = client.engine({
  version: "24.0", // version
  persist: true, // persist state
  namespace: "my-namespace", // namespace for persistence
});

// Create CLI with custom engine
const cli = client.cli({
  version: "24.0", // CLI version
  engine, // custom engine
});
```

#### Working with Images

```typescript
// Pull an image
const image = await cli.pull("nginx", "latest");

// Push an image
const ref = await cli.push("my-registry/nginx", "custom-tag");

// Import a Dagger container
const image = await cli.import(myDaggerContainer);

// Run a container
const output = await cli.run("nginx", "latest", ["--rm", "-p", "8080:80"]);
```

## API Reference

### Docker Module

#### `docker()`

Creates a new instance of the Docker module.

#### `engine(opts?: EngineOpts): Service`

Creates a new Docker Engine instance.

Parameters:

- `version` (optional): Docker Engine version (default: "24.0")
- `persist` (optional): Whether to persist engine state (default: true)
- `namespace` (optional): Namespace for persistence

#### `cli(opts?: CLIOpts): CLI`

Creates a new Docker CLI instance.

Parameters:

- `version` (optional): Docker CLI version (default: "24.0")
- `engine` (optional): Custom Docker Engine service

### CLI Methods

#### `pull(repository: string, tag: string): Promise<Image>`

Pulls a Docker image.

#### `push(repository: string, tag: string): Promise<string>`

Pushes a Docker image.

#### `import(container: Container): Promise<Image>`

Imports a Dagger container into Docker Engine.

#### `run(name: string, tag: string, args: string[]): Promise<string>`

Runs a Docker container.

#### `image(repository: string, tag: string, localID?: string): Promise<Image>`

Looks up an image in the local Docker Engine cache.

## Examples

### Complete Pipeline Example

```typescript
import { docker } from "@felipepimentel/daggerverse/docker";

export async function main() {
  // Initialize Docker
  const client = docker();

  // Create CLI with persistent engine
  const cli = client.cli({
    version: "24.0",
    engine: client.engine({
      version: "24.0",
      persist: true,
      namespace: "my-project",
    }),
  });

  // Pull base image
  const baseImage = await cli.pull("alpine", "latest");

  // Create and import a custom container
  const customContainer = dag
    .container()
    .from("alpine:latest")
    .withExec(["apk", "add", "nginx"]);

  const customImage = await cli.import(customContainer);

  // Push the custom image
  const ref = await cli.push("my-registry/custom-nginx", "v1.0");

  console.log("Pushed image:", ref);
}
```

### Working with Multiple Engines

```typescript
// Create two engines with different versions
const engine1 = client.engine({
  version: "24.0",
  persist: true,
  namespace: "prod",
});

const engine2 = client.engine({
  version: "23.0",
  persist: true,
  namespace: "dev",
});

// Create CLIs for each engine
const cli1 = client.cli({
  version: "24.0",
  engine: engine1,
});

const cli2 = client.cli({
  version: "23.0",
  engine: engine2,
});

// Use different engines for different operations
const prodImage = await cli1.pull("nginx", "latest");
const devImage = await cli2.pull("nginx", "alpine");
```

## License

See [LICENSE](../LICENSE) file in the root directory.
