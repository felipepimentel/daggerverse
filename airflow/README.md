# Airflow Module for Dagger

A Dagger module that provides integration with [Apache Airflow](https://airflow.apache.org/), the platform to programmatically author, schedule, and monitor workflows. This module allows you to run Airflow for development and testing purposes in your Dagger pipelines.

## Features

- Apache Airflow server deployment and configuration
- PostgreSQL database integration
- Redis broker integration
- Celery executor support
- DAGs directory mounting
- Plugins directory mounting
- Custom configuration support
- Additional Python requirements installation
- Automatic secret generation
- Service endpoint management

## Usage

### Basic Setup

```typescript
import { airflow } from "@felipepimentel/daggerverse/airflow";

// Initialize the Airflow module with default settings
const service = await airflow().serve();
```

### With Custom Configuration

```typescript
import { airflow } from "@felipepimentel/daggerverse/airflow";

// Initialize Airflow with custom settings
const service = await airflow({
  image: "apache/airflow", // Image name
  version: "2.9.3", // Version
  dags: dagsDir, // DAGs directory
  config: configDir, // Config directory
  plugins: pluginsDir, // Plugins directory
  requirements: "apache-airflow-providers-cncf-kubernetes", // Additional requirements
  databaseCacheName: "custom-cache", // Database cache name
}).serve();
```

## Configuration

### Constructor Options

The module accepts:

- `image`: Image name to use (default: "apache/airflow")
- `version`: Version of Apache Airflow to use (default: "2.9.3")
- `dags`: DAGs directory to mount (optional)
- `config`: Configuration directory to mount (optional)
- `plugins`: Plugins directory to mount (optional)
- `requirements`: Additional Python requirements to install (optional)
- `databaseCacheName`: Database cache name (default: "default")

### Default Settings

- Base image: `apache/airflow:2.9.3`
- Celery executor
- PostgreSQL database integration
- Redis broker integration
- Default credentials: username "airflow", password "airflow"

## Examples

### Development Setup

```typescript
import { airflow } from "@felipepimentel/daggerverse/airflow";

export async function devSetup() {
  // Initialize Airflow with default settings
  const service = await airflow().serve();
}
```

### Custom Configuration

```typescript
import { airflow } from "@felipepimentel/daggerverse/airflow";
import { Directory } from "@dagger.io/dagger";

export async function customSetup() {
  // Create custom DAGs directory
  const dags = dag.directory().withNewFile(
    "example_dag.py",
    `
from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime

def hello_world():
    print("Hello, World!")

dag = DAG(
    'example_dag',
    start_date=datetime(2024, 1, 1),
    schedule_interval='@daily'
)

hello_task = PythonOperator(
    task_id='hello_task',
    python_callable=hello_world,
    dag=dag
)
    `
  );

  // Initialize Airflow with custom DAGs
  const service = await airflow({
    image: "apache/airflow",
    version: "2.9.3",
    dags: dags,
    databaseCacheName: "dev",
  }).serve();
}
```

### Production Setup

```typescript
import { airflow } from "@felipepimentel/daggerverse/airflow";

export async function productionSetup() {
  // Create production configuration
  const config = dag.directory().withNewFile(
    "airflow.cfg",
    `
[core]
dags_folder = /opt/airflow/dags
load_examples = false
executor = CeleryExecutor

[webserver]
web_server_port = 8080
authenticate = True
auth_backend = airflow.api.auth.backend.basic_auth

[celery]
broker_url = redis://:@redis:6379/0
result_backend = db+postgresql://airflow:airflow@postgres/airflow
    `
  );

  // Initialize Airflow with production settings
  const service = await airflow({
    image: "apache/airflow",
    version: "2.9.3",
    dags: productionDags,
    config: config,
    plugins: productionPlugins,
    requirements:
      "apache-airflow-providers-amazon,apache-airflow-providers-google",
    databaseCacheName: "prod",
  }).serve();
}
```

## Dependencies

The module requires:

- Dagger SDK
- PostgreSQL module (automatically included)
- Redis module (automatically included)
- Internet access to pull the Apache Airflow container image

## Testing

The module includes tests that verify:

- Basic Airflow server initialization
- DAGs directory mounting
- Custom configuration
- Service management
- Database integration
- Redis integration

To run the tests:

```bash
dagger run test
```

## License

See [LICENSE](../LICENSE) file in the root directory.
