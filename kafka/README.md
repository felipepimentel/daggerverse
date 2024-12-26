# Kafka Module for Dagger

A Dagger module that provides integration with Apache Kafka, enabling you to run and manage Kafka clusters in your Dagger pipelines. This module supports both single-node and multi-node Kafka clusters using KRaft (Kafka Raft) for metadata management.

## Features

- Single-node Kafka deployment
- Multi-node Kafka cluster deployment
- KRaft (Kafka Raft) metadata management
- Configurable node roles (controller, broker)
- Customizable environment variables
- Service discovery and networking
- Port exposure management
- Based on Bitnami's Kafka image

## Usage

### Single Node Setup

```typescript
import { kafka } from "@felipepimentel/daggerverse/kafka";

// Initialize Kafka with default settings
const client = kafka();

// Create a single-node Kafka instance
const node = client.singleNode("kafka");

// Get the service
const service = node.service();
```

### Cluster Setup

```typescript
// Initialize Kafka with a specific version
const client = kafka("3.5.1");

// Create a 3-node Kafka cluster
const cluster = client.cluster("kafka", 3);

// Get all nodes
const nodes = await cluster.nodes();

// Access individual nodes
const controller = nodes[0]; // First node is the controller
const broker1 = nodes[1];
const broker2 = nodes[2];
```

### Custom Configuration

```typescript
// Initialize with custom settings
const client = kafka({
  version: "", // Version (optional)
  container: customContainer, // Custom container (optional)
});

// Add environment variables
const withEnv = client.withEnvVariable("KAFKA_HEAP_OPTS", "-Xmx1G -Xms1G");

// Create a cluster with custom settings
const cluster = withEnv.cluster(
  "kafka", // Service name prefix
  5 // Number of nodes
);
```

## Configuration Options

### Version

- Specifies the version of Kafka to use
- Optional: Defaults to "latest"
- Uses Bitnami's Kafka image tags

### Node Configuration

#### Controller Node (Node 0)

- Process Roles: controller, broker
- Ports: 9092 (PLAINTEXT), 9093 (CONTROLLER)
- Listeners: PLAINTEXT, CONTROLLER
- Security Protocol: PLAINTEXT

#### Broker Nodes (Node 1+)

- Process Roles: broker
- Ports: 9092 (PLAINTEXT)
- Listeners: PLAINTEXT
- Security Protocol: PLAINTEXT

### Environment Variables

The module supports all Kafka configuration options through environment variables:

- `KAFKA_CFG_*`: Kafka configuration options
- `KAFKA_HEAP_OPTS`: JVM heap settings
- `KAFKA_OPTS`: Additional JVM options
- Custom environment variables via `withEnvVariable`

## Implementation Details

### KRaft Configuration

The module uses KRaft (Kafka Raft) for metadata management:

- Node 0 serves as both controller and broker
- Other nodes operate as brokers only
- Controller quorum voters configuration is handled automatically
- Cluster ID is pre-configured (TODO: make configurable)

### Networking

- PLAINTEXT listener on port 9092 for all nodes
- CONTROLLER listener on port 9093 for node 0
- Automatic service discovery between nodes
- Advertised listeners configured for service names

## Dependencies

The module requires:

- Dagger SDK
- Internet access to pull Kafka images
- Sufficient system resources for running Kafka nodes

## License

See [LICENSE](../LICENSE) file in the root directory.
