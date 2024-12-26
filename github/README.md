# GitHub CLI Module

This Dagger module provides integration with the [GitHub CLI](https://cli.github.com/), allowing you to interact with GitHub from your Dagger pipelines.

## Usage

```typescript
import { github } from "@felipepimentel/daggerverse/github";

// Create a new instance with specific version
const gh = github("2.40.1");

// Get a container with GitHub CLI installed
const container = await gh.redhatContainer();
```

## Functions

### `new(version: string)`

Creates a new instance of the GitHub CLI module.

- `version`: The version of GitHub CLI to use (without the 'v' prefix)

### `binary(platform?: Platform)`

Gets the GitHub CLI executable binary for the specified platform.

- `platform`: Optional. Target platform in the format "os/arch"

### `overlay(platform?: Platform, prefix?: string)`

Gets a root filesystem overlay with GitHub CLI.

- `platform`: Optional. Target platform
- `prefix`: Optional. Filesystem prefix under which to install GitHub CLI (defaults to "/usr/local")

### `installation(container: Container)`

Installs GitHub CLI in a container.

- `container`: The container in which to install GitHub CLI

### `container(container: Container)`

Gets a container with GitHub CLI from a base container.

- `container`: Base container to use

### `redhatContainer(platform?: Platform)`

Gets a Red Hat Universal Base Image container with GitHub CLI.

- `platform`: Optional. Target platform

### `redhatMinimalContainer(platform?: Platform)`

Gets a Red Hat Minimal Universal Base Image container with GitHub CLI.

- `platform`: Optional. Target platform

### `redhatMicroContainer(platform?: Platform)`

Gets a Red Hat Micro Universal Base Image container with GitHub CLI.

Note: Features requiring Git will not work in the micro container.

- `platform`: Optional. Target platform

## Example

```typescript
import { github } from "@felipepimentel/daggerverse/github";

export default async function example() {
  const gh = github("2.40.1");

  // Get a container with GitHub CLI
  const container = await gh.redhatContainer();

  // Use GitHub CLI
  await container.withExec(["auth", "status"]);
}
```

## License

See [LICENSE](../LICENSE) file in the root directory.
