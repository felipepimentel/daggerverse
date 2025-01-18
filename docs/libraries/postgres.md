---
layout: default
title: PostgreSQL Module
parent: Libraries
nav_order: 13
---

# PostgreSQL Module

The PostgreSQL module provides integration with [PostgreSQL](https://www.postgresql.org/), a powerful open-source relational database. This module allows you to manage PostgreSQL databases and run queries in your Dagger pipelines.

## Features

- Database management
- Query execution
- Backup and restore
- User management
- Database initialization
- Custom configurations
- SSL/TLS support
- Extension management

## Installation

To use the PostgreSQL module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/postgres"
)
```

## Usage Examples

### Basic Database Setup

```go
func (m *MyModule) Example(ctx context.Context) (*Service, error) {
    postgres := dag.Postgres().New()
    
    // Start PostgreSQL server
    return postgres.Server(
        ctx,
        "15",           // version
        5432,          // port
        "mydb",        // database name
        dag.SetSecret("POSTGRES_PASSWORD", "secret"),
    )
}
```

### Database Initialization

```go
func (m *MyModule) InitDB(ctx context.Context) (*Service, error) {
    postgres := dag.Postgres().New()
    
    // Initialize with schema
    return postgres.WithSchema(
        ctx,
        dag.File("./schema.sql"),
        map[string]string{
            "POSTGRES_DB": "myapp",
            "POSTGRES_USER": "app",
        },
    ).Server(ctx, "15", 5432, "myapp", nil)
}
```

### Custom Configuration

```go
func (m *MyModule) CustomConfig(ctx context.Context) (*Service, error) {
    postgres := dag.Postgres().New()
    
    // Configure PostgreSQL settings
    return postgres.WithConfig(
        ctx,
        map[string]string{
            "max_connections": "100",
            "shared_buffers": "256MB",
            "work_mem": "4MB",
        },
    ).Server(ctx, "15", 5432, "mydb", nil)
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Database Operations
on: [push]

jobs:
  postgres:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup PostgreSQL
        uses: dagger/dagger-action@v1
        env:
          POSTGRES_PASSWORD: ${{ secrets.POSTGRES_PASSWORD }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/postgres
          args: |
            do -p '
              postgres := Postgres().New()
              postgres.Server(
                ctx,
                "15",
                5432,
                "mydb",
                dag.SetSecret("POSTGRES_PASSWORD", POSTGRES_PASSWORD),
              )
            '
```

## API Reference

### Postgres

Main module struct that provides access to PostgreSQL functionality.

#### Constructor

- `New() *Postgres`
  - Creates a new PostgreSQL instance
  - Default version: "15"
  - Default platform: "linux/amd64"

#### Methods

- `Server(ctx context.Context, version string, port int, database string, password *Secret) (*Service, error)`
  - Starts a PostgreSQL server
  - Parameters:
    - `version`: PostgreSQL version
    - `port`: Server port
    - `database`: Initial database name
    - `password`: Admin password
  
- `WithSchema(ctx context.Context, schema *File, env map[string]string) *Postgres`
  - Configures database schema
  - Parameters:
    - `schema`: SQL schema file
    - `env`: Environment variables
  
- `WithConfig(ctx context.Context, config map[string]string) *Postgres`
  - Sets PostgreSQL configuration
  - Parameters:
    - `config`: Configuration parameters

## Best Practices

1. **Database Configuration**
   - Tune for workload
   - Configure memory properly
   - Set appropriate limits

2. **Security**
   - Use strong passwords
   - Enable SSL/TLS
   - Manage access control

3. **Performance**
   - Monitor connections
   - Configure buffers
   - Optimize queries

4. **Backup Strategy**
   - Regular backups
   - Test restores
   - Version control schemas

## Troubleshooting

Common issues and solutions:

1. **Connection Issues**
   ```
   Error: could not connect to server
   Solution: Check port and credentials
   ```

2. **Memory Problems**
   ```
   Error: out of memory
   Solution: Adjust memory configuration
   ```

3. **Schema Errors**
   ```
   Error: relation does not exist
   Solution: Verify schema initialization
   ```

## Configuration Example

```conf
# postgresql.conf
max_connections = 100
shared_buffers = 256MB
work_mem = 4MB
maintenance_work_mem = 64MB
effective_cache_size = 1GB
max_wal_size = 1GB
checkpoint_completion_target = 0.9
random_page_cost = 1.1
effective_io_concurrency = 200
```

## Advanced Usage

### Replication Setup

```go
func (m *MyModule) SetupReplication(ctx context.Context) error {
    postgres := dag.Postgres().New()
    
    // Configure primary server
    primary := postgres.WithConfig(
        ctx,
        map[string]string{
            "wal_level": "replica",
            "max_wal_senders": "3",
        },
    ).Server(ctx, "15", 5432, "mydb", nil)
    
    // Configure replica
    replica := postgres.WithReplica(
        ctx,
        primary,
        map[string]string{
            "hot_standby": "on",
        },
    ).Server(ctx, "15", 5433, "mydb", nil)
    
    return nil
}
```

### Backup Management

```go
func (m *MyModule) ManageBackups(ctx context.Context) error {
    postgres := dag.Postgres().New()
    
    // Create backup
    backup, err := postgres.Backup(
        ctx,
        "mydb",
        map[string]string{
            "format": "custom",
            "compress": "9",
        },
    )
    if err != nil {
        return err
    }
    
    // Restore from backup
    return postgres.Restore(
        ctx,
        backup,
        "mydb_restored",
        nil,
    )
}
``` 