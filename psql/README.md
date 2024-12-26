# PSQL Module for Dagger

A Dagger module that provides integration with `psql`, the PostgreSQL interactive terminal. This module enables you to interact with PostgreSQL databases through commands and scripts in your Dagger pipelines.

## Features

- Interactive terminal access to PostgreSQL databases
- Command execution support
- File-based script execution
- Database listing and management
- SSL/TLS connection configuration
- Service binding for container-based PostgreSQL instances
- CSV output parsing for structured data
- Environment variable-based configuration

## Usage

### Basic Setup

```typescript
import { psql } from "@felipepimentel/daggerverse/psql";

// Initialize PSQL with a host
const client = psql({
  host: "db.example.com",
  port: 5432,
});

// Or with a PostgreSQL service
const client = psql({
  service: postgresService,
});
```

### Interactive Terminal

```typescript
// Open an interactive psql terminal
const terminal = client.terminal();
```

### Execute Commands

```typescript
// Run a single command
const result = await client.run().withCommand("SELECT version();").execute();

// Run multiple commands
const result = await client
  .run()
  .withCommand("CREATE DATABASE myapp;")
  .withCommand("\\l")
  .execute();

// Execute a SQL script file
const result = await client.run().withFile(sqlFile).execute();
```

### List Databases

```typescript
// Get a list of databases
const databases = await client.list();
for (const db of databases) {
  console.log(`Database: ${db.name}, Owner: ${db.owner}`);
}
```

## Configuration

### Connection Options

- `host`: PostgreSQL server hostname
- `port`: Server port (default: 5432)
- `user`: Username for authentication
- `password`: Password for authentication
- `database`: Database name to connect to
- `sslmode`: SSL mode for connections
- `service`: Optional PostgreSQL service container

### SSL Modes

- `disable`: Only try a non-SSL connection
- `allow`: First try non-SSL, then SSL if that fails
- `prefer`: First try SSL, then non-SSL if that fails
- `require`: Only try an SSL connection
- `verify-ca`: Verify server certificate
- `verify-full`: Verify server certificate and hostname

## Examples

### Complete Example

```typescript
import { psql } from "@felipepimentel/daggerverse/psql";
import { postgres } from "@felipepimentel/daggerverse/postgres";

export async function postgresInteraction() {
  // Initialize PostgreSQL service
  const postgres = postgres({
    version: "16.3",
    user: dag.setSecret("postgres-user", "admin"),
    password: dag.setSecret("postgres-password", "secret"),
    database: "myapp",
  });

  // Initialize PSQL client
  const client = psql({
    service: postgres.service(),
    user: dag.setSecret("postgres-user", "admin"),
    password: dag.setSecret("postgres-password", "secret"),
  });

  // Create tables and load data
  const result = await client
    .run()
    .withFile(
      dag
        .directory()
        .withNewFile(
          "schema.sql",
          `
        CREATE TABLE users (
          id SERIAL PRIMARY KEY,
          name TEXT NOT NULL
        );
      `
        )
        .file("schema.sql")
    )
    .withCommand("INSERT INTO users (name) VALUES ('Alice'), ('Bob');")
    .execute();
}
```

### Database Management

```typescript
import { psql } from "@felipepimentel/daggerverse/psql";

export async function manageDatabases() {
  const client = psql({
    host: "localhost",
    user: dag.setSecret("user", "postgres"),
    password: dag.setSecret("password", "postgres"),
  });

  // List databases
  const databases = await client.list();

  // Create a new database
  const result = await client
    .run()
    .withCommand("CREATE DATABASE newapp;")
    .execute();
}
```

## Testing

The module includes comprehensive tests that verify:

- Database connection and authentication
- Command execution
- Script file processing
- Database listing functionality
- SSL/TLS configuration

To run the tests:

```bash
dagger do test
```

## Dependencies

The module requires:

- Dagger SDK
- PostgreSQL client tools
- Optional: PostgreSQL server or service

## License

See [LICENSE](../LICENSE) file in the root directory.
