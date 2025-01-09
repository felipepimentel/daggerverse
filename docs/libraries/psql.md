---
layout: default
title: PSQL Module
parent: Libraries
nav_order: 5
---

# PSQL Module

The PSQL module provides integration with `psql`, the terminal-based front-end to PostgreSQL. It enables you to execute queries interactively, run commands from files, and automate database operations in your Dagger pipelines.

## Features

- Interactive PostgreSQL terminal
- Command execution (single or batch)
- File-based command execution
- Database listing and management
- SSL connection configuration
- Service binding support
- Secure credential handling

## Installation

To use the PSQL module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/psql"
)
```

## Usage Examples

### Basic PSQL Setup

```go
func (m *MyModule) Example(ctx context.Context) (*Container, error) {
    psql, err := dag.Psql().New(
        "15",           // version
        nil,           // custom container
        "db.host.com", // host
        nil,           // service
        5432,          // port
        dag.SetSecret("DB_USER", "myuser"),
        dag.SetSecret("DB_PASSWORD", "mypassword"),
        "mydb",        // database
        SSLModeVerifyFull,
    )
    if err != nil {
        return nil, err
    }
    
    return psql.Terminal(), nil
}
```

### Service Connection

```go
func (m *MyModule) WithService(ctx context.Context) (*Container, error) {
    // Create PostgreSQL service
    postgres, err := dag.Postgres().New(ctx)
    if err != nil {
        return nil, err
    }
    
    // Connect PSQL to the service
    psql, err := dag.Psql().New(
        "",    // default version
        nil,   // default container
        "",    // host will be set automatically
        postgres.Service(),
        5432,
        nil,   // default user
        nil,   // default password
        "",    // default database
        "",    // default SSL mode
    )
    if err != nil {
        return nil, err
    }
    
    return psql.Terminal(), nil
}
```

### Running Commands

```go
func (m *MyModule) Commands(ctx context.Context) error {
    psql, err := dag.Psql().New(
        "15",
        nil,
        "db.host.com",
        nil,
        5432,
        dag.SetSecret("DB_USER", "myuser"),
        dag.SetSecret("DB_PASSWORD", "mypassword"),
        "mydb",
        SSLModeVerifyFull,
    )
    if err != nil {
        return err
    }
    
    // Run a single command
    output, err := psql.RunCommand(ctx, "SELECT version();")
    if err != nil {
        return err
    }
    
    // Run commands from a file
    output, err = psql.RunFile(ctx, dag.File("./query.sql"))
    if err != nil {
        return err
    }
    
    return nil
}
```

### Database Operations

```go
func (m *MyModule) ListDatabases(ctx context.Context) ([]DatabaseListEntry, error) {
    psql, err := dag.Psql().New(
        "15",
        nil,
        "db.host.com",
        nil,
        5432,
        dag.SetSecret("DB_USER", "myuser"),
        dag.SetSecret("DB_PASSWORD", "mypassword"),
        "",
        SSLModeVerifyFull,
    )
    if err != nil {
        return nil, err
    }
    
    return psql.List(ctx)
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Database Operations
on: [push]

jobs:
  db-ops:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Database Commands
        uses: dagger/dagger-action@v1
        env:
          DB_USER: ${{ secrets.DB_USER }}
          DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/psql
          args: |
            do -p '
              psql := Psql().New(
                "15",
                nil,
                "db.host.com",
                nil,
                5432,
                dag.SetSecret("DB_USER", DB_USER),
                dag.SetSecret("DB_PASSWORD", DB_PASSWORD),
                "mydb",
                SSLModeVerifyFull,
              )
              psql.RunCommand(ctx, "SELECT version();")
            '
```

## API Reference

### Psql

Main module struct that provides access to PSQL functionality.

#### Constructor

- `New(version string, container *Container, host string, service *Service, port int, user *Secret, password *Secret, database string, sslmode SSLMode) (*Psql, error)`
  - Creates a new PSQL instance
  - Parameters:
    - `version`: PostgreSQL version (optional, defaults to "latest")
    - `container`: Custom base container (optional)
    - `host`: Database host (required if service not provided)
    - `service`: PostgreSQL service (required if host not provided)
    - `port`: Database port (optional, defaults to 5432)
    - `user`: Database user (optional, defaults to "postgres")
    - `password`: Database password (optional)
    - `database`: Database name (optional)
    - `sslmode`: SSL connection mode (optional)

#### Methods

- `Terminal() *Container`
  - Opens an interactive PSQL terminal
  
- `List(ctx context.Context) ([]DatabaseListEntry, error)`
  - Lists all available databases
  
- `RunCommand(ctx context.Context, command string) (string, error)`
  - Executes a single SQL command
  
- `RunFile(ctx context.Context, file *File) (string, error)`
  - Executes commands from a file

### SSLMode

Enum type for SSL connection modes:

- `SSLModeDisable`: Only try a non-SSL connection
- `SSLModeAllow`: Try non-SSL first, then SSL
- `SSLModePrefer`: Try SSL first, then non-SSL
- `SSLModeRequire`: Only try SSL connection
- `SSLModeVerifyCA`: SSL with CA verification
- `SSLModeVerifyFull`: SSL with full verification

## Best Practices

1. **Connection Management**
   - Use service binding for local development
   - Use SSL in production environments
   - Handle connection errors appropriately

2. **Security**
   - Use secrets for credentials
   - Enable appropriate SSL mode
   - Follow principle of least privilege

3. **Command Execution**
   - Use parameterized queries
   - Handle large result sets properly
   - Implement proper error handling

4. **Resource Management**
   - Close connections when done
   - Monitor query performance
   - Use appropriate timeouts

## Troubleshooting

Common issues and solutions:

1. **Connection Issues**
   ```
   Error: could not connect to server
   Solution: Verify host, port, and credentials
   ```

2. **SSL Problems**
   ```
   Error: SSL connection failed
   Solution: Check SSL mode and certificate configuration
   ```

3. **Authentication Failures**
   ```
   Error: password authentication failed
   Solution: Verify user credentials and permissions
   ```

## Environment Variables

The module uses standard PostgreSQL environment variables:

- `PGHOST`: Database host
- `PGPORT`: Database port
- `PGUSER`: Database user
- `PGPASSWORD`: Database password
- `PGDATABASE`: Database name
- `PGSSLMODE`: SSL mode

These can be set using the constructor parameters or through GitHub Actions secrets. 