# Redis Module for Dagger

A Dagger module that provides integration with Redis, enabling you to run and manage Redis servers in your Dagger pipelines. This module is designed for development and CI purposes, offering a simple way to set up Redis instances.

## Features

- Single Redis server deployment
- Configurable port mapping
- Optional password protection
- Data persistence through cache volumes
- Version selection support
- Based on the official Redis image

## Usage

### Basic Setup

```go
// Initialize Redis with default settings
redis := dag.Redis().New()

// Get the Redis server
server := redis.Server()

// Get the service
service := server.AsService()
```

### Custom Configuration

```go
// Initialize Redis with custom settings
redis := dag.Redis().New(
    6380,                   // Custom port
    "7.2.4",               // Specific version
    dag.SetSecret("password", "mypassword"), // Password protection
    true,                  // Enable data persistence
)

// Get the server with custom configuration
server := redis.Server()
```

## Configuration Options

### Constructor Parameters

- `port` (optional, default: 6379): The port to expose Redis on
- `version` (optional, default: "7.2.4"): Redis version to use
- `password` (optional): Secret for Redis authentication
- `cache` (optional, default: false): Enable data persistence using cache volume

### Default Settings

- Port: 6379
- Version: 7.2.4
- No authentication
- No data persistence

## Implementation Details

### Server Configuration

The module configures Redis with:

- Network binding to all interfaces
- Optional password protection
- Optional persistent storage
- Customizable port mapping

### Data Persistence

When enabled, the module uses Dagger cache volumes to persist Redis data across pipeline runs.

### Security

- Optional password authentication
- Network access control through container configuration
- Secure secret management for passwords

## Dependencies

The module requires:

- Dagger SDK
- Internet access to pull Redis images
- Sufficient system resources for running Redis

## Examples

### With Password Protection

```go
// Initialize Redis with password
redis := dag.Redis().New(
    6379,
    "7.2.4",
    dag.SetSecret("redis-password", "secure123"),
    false,
)

// Get the secured server
server := redis.Server()
```

### With Data Persistence

```go
// Initialize Redis with data persistence
redis := dag.Redis().New(
    6379,
    "7.2.4",
    nil,
    true,
)

// Get the server with persistent storage
server := redis.Server()
```

## License

This module is part of the Dagger project and is subject to its licensing terms.

Copyright Camptocamp SA
SPDX-License-Identifier: AGPL-3.0-or-later
