---
layout: default
title: Kafka Module
parent: Libraries
nav_order: 7
---

# Kafka Module

The Kafka module provides integration with [Apache Kafka](https://kafka.apache.org/), a distributed event streaming platform. This module allows you to manage Kafka brokers, topics, and messages in your Dagger pipelines.

## Features

- Kafka broker management
- Topic creation and management
- Message production and consumption
- ZooKeeper integration
- Multi-broker support
- Custom configuration
- Security settings
- Monitoring capabilities

## Installation

To use the Kafka module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/kafka"
)
```

## Usage Examples

### Basic Kafka Broker Setup

```go
func (m *MyModule) Example(ctx context.Context) (*Service, error) {
    kafka := dag.Kafka().New()
    
    // Start Kafka broker
    return kafka.Broker(
        ctx,
        "2.13-3.5.0",  // version
        9092,          // port
        nil,           // custom config
    )
}
```

### Topic Management

```go
func (m *MyModule) ManageTopics(ctx context.Context) error {
    kafka := dag.Kafka().New()
    
    // Create and configure topics
    return kafka.CreateTopic(
        ctx,
        "my-topic",
        3,              // partitions
        2,              // replication factor
        map[string]string{
            "cleanup.policy": "delete",
            "retention.ms": "604800000",
        },
    )
}
```

### Message Operations

```go
func (m *MyModule) HandleMessages(ctx context.Context) error {
    kafka := dag.Kafka().New()
    
    // Produce and consume messages
    return kafka.
        ProduceMessage(
            ctx,
            "my-topic",
            "key",
            "value",
        ).
        ConsumeMessages(
            ctx,
            "my-topic",
            "my-group",
            10,  // max messages
        )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Kafka Operations
on: [push]

jobs:
  kafka:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup Kafka
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/kafka
          args: |
            do -p '
              kafka := Kafka().New()
              broker := kafka.Broker(
                ctx,
                "2.13-3.5.0",
                9092,
                nil,
              )
              kafka.CreateTopic(
                ctx,
                "my-topic",
                3,
                2,
                nil,
              )
            '
```

## API Reference

### Kafka

Main module struct that provides access to Kafka functionality.

#### Constructor

- `New() *Kafka`
  - Creates a new Kafka instance
  - Default version: "2.13-3.5.0"
  - Default platform: "linux/amd64"

#### Methods

- `Broker(ctx context.Context, version string, port int, config map[string]string) (*Service, error)`
  - Starts a Kafka broker
  - Parameters:
    - `version`: Kafka version
    - `port`: Broker port
    - `config`: Custom broker configuration
  
- `CreateTopic(ctx context.Context, name string, partitions int, replicationFactor int, config map[string]string) error`
  - Creates a Kafka topic
  - Parameters:
    - `name`: Topic name
    - `partitions`: Number of partitions
    - `replicationFactor`: Replication factor
    - `config`: Topic configuration
  
- `ProduceMessage(ctx context.Context, topic string, key string, value string) error`
  - Produces a message to a topic
  - Parameters:
    - `topic`: Target topic
    - `key`: Message key
    - `value`: Message value
  
- `ConsumeMessages(ctx context.Context, topic string, group string, maxMessages int) error`
  - Consumes messages from a topic
  - Parameters:
    - `topic`: Source topic
    - `group`: Consumer group
    - `maxMessages`: Maximum messages to consume

## Best Practices

1. **Broker Configuration**
   - Use appropriate replication factors
   - Configure proper retention policies
   - Monitor broker health

2. **Topic Management**
   - Plan partitioning strategy
   - Set appropriate configs
   - Monitor topic performance

3. **Message Handling**
   - Use consistent serialization
   - Handle errors appropriately
   - Monitor consumer lag

4. **Security**
   - Enable authentication
   - Use encryption in transit
   - Follow security best practices

## Troubleshooting

Common issues and solutions:

1. **Connection Issues**
   ```
   Error: broker not available
   Solution: Check broker status and network connectivity
   ```

2. **Topic Creation Failures**
   ```
   Error: invalid replication factor
   Solution: Ensure replication factor <= number of brokers
   ```

3. **Message Problems**
   ```
   Error: message too large
   Solution: Check message size limits and broker configuration
   ```

## Configuration Example

```properties
# server.properties
broker.id=0
listeners=PLAINTEXT://:9092
log.dirs=/tmp/kafka-logs
zookeeper.connect=localhost:2181
num.partitions=3
default.replication.factor=2
log.retention.hours=168
```

## Advanced Usage

### Multi-Broker Setup

```go
func (m *MyModule) MultiNode(ctx context.Context) error {
    kafka := dag.Kafka().New()
    
    // Start multiple brokers
    broker1 := kafka.Broker(ctx, "2.13-3.5.0", 9092, map[string]string{
        "broker.id": "1",
    })
    
    broker2 := kafka.Broker(ctx, "2.13-3.5.0", 9093, map[string]string{
        "broker.id": "2",
    })
    
    broker3 := kafka.Broker(ctx, "2.13-3.5.0", 9094, map[string]string{
        "broker.id": "3",
    })
    
    return nil
}
```

### Custom ZooKeeper Configuration

```go
func (m *MyModule) CustomZK(ctx context.Context) (*Service, error) {
    kafka := dag.Kafka().New()
    
    // Configure with custom ZooKeeper
    return kafka.WithZooKeeper(
        ctx,
        "3.8.0",  // ZK version
        2181,     // ZK port
        map[string]string{
            "tickTime": "2000",
            "initLimit": "10",
        },
    ).Broker(ctx, "2.13-3.5.0", 9092, nil)
}
``` 