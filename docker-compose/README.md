# Docker Compose Module for Dagger

A native Dagger reimplementation of Docker Compose that allows you to use Docker Compose configurations in your Dagger pipelines. This module provides a subset of Docker Compose functionality with full compatibility for supported features.

## Features

Currently supported Docker Compose features:

- Fully compatible configuration parser
- Service image specification (`services.X.image`)
- Service build configuration (`services.X.build`)
- Port mapping (`services.X.ports`)
- Environment variables (`services.X.environment`)
- Custom entrypoint (`services.X.entrypoint`)
- Custom command (`services.X.command`)

## Usage

### Basic Setup

```typescript
import { dockerCompose } from "@felipepimentel/daggerverse/docker-compose";

// Initialize the Docker Compose module
const compose = dockerCompose();

// Load a project from a directory containing docker-compose.yml
const project = compose.project(sourceDir);
```

### Working with Services

```typescript
// Get a specific service by name
const service = project.service("api");

// Get all services in the project
const services = await project.services();
```

### Service Configuration

```typescript
// Get the raw configuration for a service
const config = await service.config();
```

### Container Operations

```typescript
// Get the base container for a service (without compose-specific modifications)
const baseContainer = await service.baseContainer();

// Get the fully configured container for a service
const container = await service.container();

// Run the service directly on the Dagger Engine
const svc = await service.up();
```

## Example Project

The module includes an example project that demonstrates various Docker Compose features:

```typescript
// Load the example project
const example = compose.example();

// Access the example project's configuration
const config = await example.config();
```

Example `docker-compose.yml`:

```yaml
version: "3.1"

services:
  api:
    image: golang:1.21
    ports:
      - 8020:8020
    command: ["go", "run", "/api/main.go"]
    environment:
      LOCAL: "true"
      DAGGER_API_SERVER_URL: http://localhost:8080
    volumes:
      - .:/api
    working_dir: /api

  db:
    image: postgres:13.8
    command: ["postgres", "-c", "log_statement=all"]
    ports:
      - 5432:5432
    environment:
      POSTGRES_PASSWORD: dagger
      POSTGRES_USER: dagger
```

## Complete Example

```typescript
import { dockerCompose } from "@felipepimentel/daggerverse/docker-compose";

export async function runComposeProject() {
  // Initialize Docker Compose
  const compose = dockerCompose();

  // Load project from source directory
  const project = compose.project(sourceDir);

  // Get API service
  const apiService = project.service("api");

  // Get container with all compose configurations applied
  const container = await apiService.container();

  // Run the service
  const service = await apiService.up();
}
```

## Configuration

### Project Options

The `project` function accepts:

- `source`: Directory containing the Docker Compose configuration (optional)

### Service Methods

- `config`: Returns the raw YAML configuration for a service
- `baseContainer`: Returns the base container without compose-specific modifications
- `container`: Returns the fully configured container with all compose settings applied
- `up`: Runs the service directly on the Dagger Engine

## Dependencies

The module requires:

- Dagger SDK
- Docker Compose specification parser
- YAML parser

## Limitations

Current limitations include:

- Not all Docker Compose features are supported
- Volume mounts are not yet implemented
- Networks are not yet supported
- Dependencies and healthchecks are not yet implemented
- Some Docker Compose-specific features may work differently in the Dagger environment

## License

See [LICENSE](../LICENSE) file in the root directory.
