# PostgreSQL Module for Dagger

A Dagger module that provides integration with PostgreSQL, allowing you to run PostgreSQL database servers for development and testing purposes in your Dagger pipelines.

## Features

- PostgreSQL server deployment and configuration
- Superuser management
- Database initialization
- Data persistence with cache volumes
- Custom initialization scripts
- Instance naming and labeling
- PSQL client integration
- Cross-platform support

## Usage

### Basic Setup

```typescript
import { postgres } from "@felipepimentel/daggerverse/postgres";

// Initialize PostgreSQL with default settings
const client = postgres();

// Get the PostgreSQL service
const service = client.service();
```

### Custom Configuration

```typescript
// Initialize with custom settings
const client = postgres({
  version: "16.3", // Version
  container: null, // Default container
  user: customUser, // Custom superuser
  password: customPassword, // Custom password
  database: "mydb", // Custom database name
  initdbArgs: ["-E", "UTF8"], // Custom initdb arguments
  dataVolume: cacheVolume, // Data persistence
  name: "my-instance", // Instance name
  initScripts: initScripts, // Initialization scripts
});
```

## Configuration

### Constructor Options

The module accepts:

- `version`: Version of PostgreSQL to use (default: "latest")
- `container`: Custom container to use as base (optional)
- `user`: Superuser name (default: "postgres")
- `password`: Superuser password (default: "postgres")
- `database`: Database name (default: "postgres")
- `initdbArgs`: Additional arguments for initdb (optional)
- `dataVolume`: Cache volume for data persistence (optional)
- `name`: Instance name for labeling (optional)
- `initScripts`: Directory containing initialization scripts (optional)

### Default Settings

- Base image: `postgres:latest`
- Port: 5432
- Environment variables:
  - `POSTGRES_USER`
  - `POSTGRES_PASSWORD`
  - `POSTGRES_DB`
  - `POSTGRES_INITDB_ARGS` (when specified)
- Data directory: `/var/lib/postgresql/data`

## Examples

### Development Setup

```typescript
import { postgres } from "@felipepimentel/daggerverse/postgres";

export async function devSetup() {
  // Create a cache volume for data persistence
  const dataVolume = dag.cacheVolume("postgres-data");

  // Initialize PostgreSQL
  const client = postgres({
    version: "16.3",
    database: "dev_db",
    dataVolume: dataVolume,
    name: "dev-instance",
  });

  // Get service
  const service = client.service();
}
```

### Custom Initialization

```typescript
import { postgres } from "@felipepimentel/daggerverse/postgres";

export async function customInit() {
  // Create initialization scripts
  const initScripts = dag.directory().withNewFile(
    "init.sql",
    `
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    );
    `
  );

  // Initialize PostgreSQL with custom scripts
  const client = postgres({
    version: "16.3",
    database: "app_db",
    initScripts: initScripts,
  });
}
```

### Production-like Setup

```typescript
import { postgres } from "@felipepimentel/daggerverse/postgres";

export async function productionSetup() {
  // Use custom container with production settings
  const container = dag
    .container()
    .from("postgres:16.3")
    .withEnvVariable("POSTGRES_MAX_CONNECTIONS", "100")
    .withEnvVariable("POSTGRES_SHARED_BUFFERS", "256MB");

  // Initialize PostgreSQL with production settings
  const client = postgres({
    container: container,
    user: productionUser,
    password: productionPassword,
    database: "prod_db",
    initdbArgs: ["-E", "UTF8", "--locale=en_US.UTF-8"],
    dataVolume: persistentVolume,
    name: "prod-instance",
    initScripts: productionInit,
  });
}
```

## Dependencies

The module requires:

- Dagger SDK
- Internet access to pull PostgreSQL images
- PSQL module (for client operations)

## Testing

The module includes a test suite in the `tests` directory that demonstrates:

- Basic PostgreSQL server initialization
- Database creation and management
- Custom user authentication
- Service operations
- Data persistence

To run the tests:

```bash
dagger do test
```

## License

See [LICENSE](../LICENSE) file in the root directory.
