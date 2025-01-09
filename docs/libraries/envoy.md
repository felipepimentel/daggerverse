---
layout: default
title: Envoy Module
parent: Libraries
nav_order: 4
---

# Envoy Module

The Envoy module provides integration with [Envoy](https://www.envoyproxy.io/), a high-performance edge and service proxy. This module allows you to run and validate Envoy proxy configurations in your Dagger pipelines.

## Features

- Run Envoy proxy instances
- Validate Envoy configurations
- Multi-platform support
- Custom version selection
- Port exposure management
- Service mode operation

## Installation

To use the Envoy module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/envoy"
)
```

## Usage Examples

### Basic Envoy Proxy Service

```go
func (m *MyModule) Example(ctx context.Context) (*Service, error) {
    envoy := dag.Envoy().New()
    
    // Create Envoy proxy service
    return envoy.EnvoyProxyService(
        ctx,
        "v1.30-latest",      // version
        "linux/arm64",       // platform
        dag.File("./envoy.yaml"), // config
        []int{10000},        // ports to expose
    )
}
```

### Configuration Validation

```go
func (m *MyModule) ValidateConfig(ctx context.Context) (string, error) {
    envoy := dag.Envoy().New()
    
    // Validate Envoy configuration
    return envoy.ValidateConfig(
        ctx,
        "v1.30-latest",
        "linux/arm64",
        dag.File("./envoy.yaml"),
    )
}
```

### Custom Platform and Version

```go
func (m *MyModule) CustomSetup(ctx context.Context) (*Service, error) {
    envoy := dag.Envoy().New()
    
    // Use custom version and platform
    return envoy.EnvoyProxyService(
        ctx,
        "v1.29.1",
        "linux/amd64",
        dag.File("./envoy.yaml"),
        []int{10000, 9901}, // expose multiple ports
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Envoy Configuration
on: [push]

jobs:
  envoy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Validate Envoy Config
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/envoy
          args: |
            do -p '
              envoy := Envoy().New()
              envoy.ValidateConfig(
                ctx,
                "v1.30-latest",
                "linux/amd64",
                dag.File("./envoy.yaml"),
              )
            '
```

## API Reference

### Envoy

Main module struct that provides access to Envoy functionality.

#### Constructor

- `New() *Envoy`
  - Creates a new Envoy instance
  - Default version: "v1.30-latest"
  - Default platform: "linux/arm64"

#### Methods

- `EnvoyProxyService(ctx context.Context, version string, platform Platform, config *File, port []int) (*Service, error)`
  - Creates and runs an Envoy proxy service
  - Parameters:
    - `version`: Envoy version (optional, default: "v1.30-latest")
    - `platform`: Target platform (optional, default: "linux/arm64")
    - `config`: Envoy configuration file (required)
    - `port`: Ports to expose (required)
  
- `ValidateConfig(ctx context.Context, version string, platform Platform, config *File) (string, error)`
  - Validates an Envoy configuration file
  - Parameters:
    - `version`: Envoy version (optional, default: "v1.30-latest")
    - `platform`: Target platform (optional, default: "linux/arm64")
    - `config`: Envoy configuration file to validate (required)

## Configuration

### Example Envoy Configuration

```yaml
static_resources:
  listeners:
  - name: listener_0
    address:
      socket_address:
        address: 0.0.0.0
        port_value: 10000
    filter_chains:
    - filters:
      - name: envoy.filters.network.http_connection_manager
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
          stat_prefix: ingress_http
          route_config:
            name: local_route
            virtual_hosts:
            - name: local_service
              domains: ["*"]
              routes:
              - match:
                  prefix: "/"
                route:
                  host_rewrite_literal: www.envoyproxy.io
                  cluster: service_envoyproxy_io
          http_filters:
          - name: envoy.filters.http.router
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
  clusters:
  - name: service_envoyproxy_io
    type: LOGICAL_DNS
    dns_lookup_family: V4_ONLY
    load_assignment:
      cluster_name: service_envoyproxy_io
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: www.envoyproxy.io
                port_value: 443
    transport_socket:
      name: envoy.transport_sockets.tls
      typed_config:
        "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
```

## Best Practices

1. **Configuration Management**
   - Validate configurations before deployment
   - Use version control for configurations
   - Document configuration changes

2. **Version Control**
   - Use specific versions in production
   - Test upgrades in staging
   - Keep track of version compatibility

3. **Resource Management**
   - Monitor proxy performance
   - Configure appropriate resource limits
   - Use proper logging levels

4. **Security**
   - Follow Envoy security best practices
   - Keep Envoy version up to date
   - Configure TLS appropriately

## Troubleshooting

Common issues and solutions:

1. **Configuration Errors**
   ```
   Error: configuration is invalid
   Solution: Use ValidateConfig to check configuration syntax
   ```

2. **Port Conflicts**
   ```
   Error: address already in use
   Solution: Verify port availability and configurations
   ```

3. **Platform Issues**
   ```
   Error: no matching manifest
   Solution: Verify platform compatibility and availability
   ```

## Advanced Usage

### Custom Configuration Generation

```go
func (m *MyModule) GenerateConfig(ctx context.Context) (*Service, error) {
    // Generate configuration dynamically
    config := fmt.Sprintf(`
static_resources:
  listeners:
  - name: listener_0
    address:
      socket_address:
        address: 0.0.0.0
        port_value: %d
    ...
`, port)

    envoy := dag.Envoy().New()
    
    return envoy.EnvoyProxyService(
        ctx,
        "v1.30-latest",
        "linux/arm64",
        dag.File(config),
        []int{port},
    )
}
```

### Health Check Integration

```go
func (m *MyModule) WithHealthCheck(ctx context.Context) (*Service, error) {
    envoy := dag.Envoy().New()
    
    return envoy.EnvoyProxyService(
        ctx,
        "v1.30-latest",
        "linux/arm64",
        dag.File("./envoy.yaml"),
        []int{10000, 9901}, // 9901 for admin interface
    )
}
``` 