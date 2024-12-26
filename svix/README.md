# Svix Module for Dagger

A Dagger module that provides integration with [Svix](https://www.svix.com/), the enterprise-ready webhook service. This module allows you to run Svix for development and testing purposes in your Dagger pipelines.

## Features

- Svix server deployment and configuration
- PostgreSQL database integration
- JWT authentication support
- Custom configuration file support
- Container-based execution
- Service endpoint management
- Automatic secret generation

## Usage

### Basic Setup

```go
// Initialize the Svix module with default settings
svix, err := dag.Svix().New(ctx)
if err != nil {
    return err
}

// Get the Svix service
service := svix.Service()
```

### With PostgreSQL Integration

```go
// Initialize Svix with PostgreSQL
svix, err := dag.Svix().New(ctx,
    "",             // Default version
    nil,            // Default container
    postgres,       // PostgreSQL service
    "svix",         // Database name
    nil,            // Auto-generated JWT secret
    nil,            // Default config
)
if err != nil {
    return err
}
```

### Custom Configuration

```go
// Initialize with custom settings
svix, err := dag.Svix().New(ctx,
    "latest",           // Version
    customContainer,    // Custom container
    postgres,          // PostgreSQL service
    "custom-db",       // Custom database name
    jwtSecret,         // Custom JWT secret
    configFile,        // Custom config file
)
if err != nil {
    return err
}
```

## Configuration

### Constructor Options

The `New` function accepts:

- `version`: Version of Svix to use (default: "latest")
- `container`: Custom container to use as base (optional)
- `postgres`: PostgreSQL service integration (optional)
- `database`: Database name override (optional)
- `jwtSecret`: JWT secret for authentication (auto-generated if not provided)
- `config`: Custom Svix configuration file (optional)

### Default Settings

- Base image: `svix/svix-server:latest`
- Auto-generated JWT secret if not provided
- Default PostgreSQL database name: "svix"

## Examples

### Development Setup

```go
func DevSetup(ctx context.Context) error {
    // Initialize PostgreSQL
    postgres := dag.Postgres().New()

    // Initialize Svix with PostgreSQL
    svix, err := dag.Svix().New(ctx,
        "",         // Use default version
        nil,        // Use default container
        postgres,   // PostgreSQL service
        "",         // Use default database name
        nil,        // Auto-generate JWT secret
        nil,        // Use default config
    )
    if err != nil {
        return err
    }

    // Get the service
    service := svix.Service()

    return nil
}
```

### Custom Configuration

```go
func CustomSetup(ctx context.Context) error {
    // Create custom configuration
    config := dag.Directory().WithNewFile(
        "config.yaml",
        `
server:
  port: 8071
database:
  type: postgres
  dsn: ${DATABASE_URL}
jwt:
  secret: ${JWT_SECRET}
        `,
    ).File("config.yaml")

    // Initialize Svix with custom config
    svix, err := dag.Svix().New(ctx,
        "latest",
        nil,
        postgres,
        "custom-db",
        jwtSecret,
        config,
    )
    if err != nil {
        return err
    }

    return nil
}
```

### Production Setup

```go
func ProductionSetup(ctx context.Context) error {
    // Use custom container with production settings
    container := dag.Container().
        From("svix/svix-server:latest").
        WithEnvVariable("SVIX_ENV", "production").
        WithEnvVariable("SVIX_LOG_LEVEL", "info")

    // Initialize Svix with production settings
    svix, err := dag.Svix().New(ctx,
        "",
        container,
        postgres,
        "svix-prod",
        productionJwtSecret,
        productionConfig,
    )
    if err != nil {
        return err
    }

    return nil
}
```

## Dependencies

The module requires:

- Dagger SDK
- PostgreSQL module (optional, but recommended)
- Internet access to pull the Svix container image

## Testing

The module includes a test suite in the `tests` directory that demonstrates:

- Basic Svix server initialization
- PostgreSQL integration
- Custom configuration handling
- JWT authentication
- Service management

To run the tests:

```bash
dagger do test
```

## License

This module is part of the Dagger project and is subject to its licensing terms.

Copyright Camptocamp SA
SPDX-License-Identifier: AGPL-3.0-or-later
