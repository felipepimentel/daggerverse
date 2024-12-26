# Kind Module for Dagger

A Dagger module that provides integration with Kind (Kubernetes in Docker), allowing you to create and manage local Kubernetes clusters within your Dagger pipelines.

## Features

- Create and manage Kind clusters
- Load container images into Kind clusters
- Interactive K9s terminal for cluster management
- Docker socket integration
- Custom cluster configuration
- Kubernetes tools included (kubectl, k9s)

## Usage

### Basic Setup

```typescript
import { kind } from "@felipepimentel/daggerverse/kind";

// Initialize the Kind module with Docker socket and service
const client = kind({
  dockerSocket: dag.hostSocket("/var/run/docker.sock"),
  kindSvc: dag.service().withEndpoint("tcp://127.0.0.1:3000"),
  clusterName: "my-cluster",
});
```

### Loading Images into Kind

```typescript
// Load a container image into the Kind cluster
const container = dag.container().from("nginx:latest");
await client.loadContainerOnKind(container, "my-nginx");
```

### Interactive K9s Terminal

```typescript
// Launch K9s terminal for cluster management
const terminal = await client.knines();
```

## Configuration Options

### Constructor Options

The module accepts:

- `dockerSocket`: Docker socket path (required)
  - Example: `/var/run/docker.sock`
- `kindSvc`: Kind service endpoint (required)
  - Format: `tcp://127.0.0.1:<port>`
  - Port must be between 1024 and 65535
- `clusterName`: Name of the Kind cluster (optional)
  - Default: "dagger-kubernetes-cluster"

### Kind Configuration

The module creates a Kind cluster with the following default configuration:

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerPort: <specified-port>
```

## Prerequisites

Before using this module:

1. Ensure Docker is installed and running on the host machine
2. Configure `/etc/hosts` with an entry for `127.0.0.1` as localhost
3. Choose an available port for the Kind API server (between 1024 and 65535)

## Testing

To test the Kind module:

1. Ensure you have Dagger installed
2. Run the test suite:

```bash
dagger do test
```

## Examples

### Complete Example: Create Cluster and Deploy Application

```typescript
import { kind } from "@felipepimentel/daggerverse/kind";

export async function deployToKind() {
  // Initialize Kind module
  const client = kind({
    dockerSocket: dag.hostSocket("/var/run/docker.sock"),
    kindSvc: dag.service().withEndpoint("tcp://127.0.0.1:3000"),
    clusterName: "test-cluster",
  });

  // Build and load application container
  const appContainer = dag
    .container()
    .from("golang:1.22")
    .withDirectory(
      "/app",
      dag.directory().withNewFile(
        "main.go",
        `package main
      import "fmt"
      func main() {
          fmt.Println("Hello from Kind!")
      }`
      )
    )
    .withWorkdir("/app")
    .withExec(["go", "build", "-o", "app"]);

  // Load the container into Kind
  await client.loadContainerOnKind(appContainer, "my-app:latest");
}
```

### Managing Multiple Clusters

```typescript
import { kind } from "@felipepimentel/daggerverse/kind";

export async function manageClusters() {
  // Create development cluster
  const devCluster = kind({
    dockerSocket: dag.hostSocket("/var/run/docker.sock"),
    kindSvc: dag.service().withEndpoint("tcp://127.0.0.1:3001"),
    clusterName: "dev-cluster",
  });

  // Create testing cluster
  const testCluster = kind({
    dockerSocket: dag.hostSocket("/var/run/docker.sock"),
    kindSvc: dag.service().withEndpoint("tcp://127.0.0.1:3002"),
    clusterName: "test-cluster",
  });

  // Load different versions of an application
  const devContainer = dag.container().from("my-app:dev");
  const testContainer = dag.container().from("my-app:test");

  await devCluster.loadContainerOnKind(devContainer, "my-app");
  await testCluster.loadContainerOnKind(testContainer, "my-app");
}
```

## Command Line Usage

The module can be used directly from the command line:

```bash
# Create a cluster and launch K9s
dagger call --docker-sock=/var/run/docker.sock --kind-svc=tcp://127.0.0.1:3000 knines

# Inspect the container configuration
dagger call --docker-sock=/var/run/docker.sock --kind-svc=tcp://127.0.0.1:3000 inspect
```

## Dependencies

The module requires:

- Dagger SDK
- Docker Engine
- Kind v0.25.0 (installed automatically)
- kubectl (installed automatically)
- k9s (installed automatically)

## License

See [LICENSE](../LICENSE) file in the root directory.
