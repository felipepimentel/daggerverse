---
layout: default
title: Nginx Module
parent: Libraries
nav_order: 9
---

# Nginx Module

The Nginx module provides integration with [Nginx](https://nginx.org/), a high-performance HTTP server and reverse proxy. This module allows you to configure and manage Nginx servers in your Dagger pipelines.

## Features

- Server configuration
- Virtual host management
- SSL/TLS support
- Reverse proxy setup
- Load balancing
- Static file serving
- Custom configuration
- Health checks

## Installation

To use the Nginx module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/nginx"
)
```

## Usage Examples

### Basic Server Setup

```go
func (m *MyModule) Example(ctx context.Context) (*Service, error) {
    nginx := dag.Nginx().New()
    
    // Start Nginx server
    return nginx.Server(
        ctx,
        "1.24",         // version
        80,            // port
        dag.Directory("./static"), // content
    )
}
```

### Reverse Proxy Configuration

```go
func (m *MyModule) ReverseProxy(ctx context.Context) (*Service, error) {
    nginx := dag.Nginx().New()
    
    // Configure reverse proxy
    return nginx.WithProxy(
        ctx,
        "api",
        "http://backend:8080",
        map[string]string{
            "proxy_set_header": "Host $host",
            "proxy_ssl": "off",
        },
    ).Server(ctx, "1.24", 80, nil)
}
```

### SSL Configuration

```go
func (m *MyModule) WithSSL(ctx context.Context) (*Service, error) {
    nginx := dag.Nginx().New()
    
    // Configure SSL
    return nginx.WithSSL(
        ctx,
        dag.File("./cert.pem"),
        dag.File("./key.pem"),
        443,
    ).Server(ctx, "1.24", 80, dag.Directory("./static"))
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Nginx Operations
on: [push]

jobs:
  nginx:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup Nginx Server
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/nginx
          args: |
            do -p '
              nginx := Nginx().New()
              nginx.Server(
                ctx,
                "1.24",
                80,
                dag.Directory("./static"),
              )
            '
```

## API Reference

### Nginx

Main module struct that provides access to Nginx functionality.

#### Constructor

- `New() *Nginx`
  - Creates a new Nginx instance
  - Default version: "1.24"
  - Default platform: "linux/amd64"

#### Methods

- `Server(ctx context.Context, version string, port int, content *Directory) (*Service, error)`
  - Starts an Nginx server
  - Parameters:
    - `version`: Nginx version
    - `port`: Server port
    - `content`: Static content directory
  
- `WithProxy(ctx context.Context, location string, upstream string, config map[string]string) *Nginx`
  - Configures reverse proxy
  - Parameters:
    - `location`: URL path
    - `upstream`: Upstream server URL
    - `config`: Proxy configuration
  
- `WithSSL(ctx context.Context, cert *File, key *File, port int) *Nginx`
  - Configures SSL/TLS
  - Parameters:
    - `cert`: SSL certificate file
    - `key`: SSL private key file
    - `port`: HTTPS port

## Best Practices

1. **Server Configuration**
   - Use appropriate worker processes
   - Configure buffer sizes
   - Enable compression

2. **Security**
   - Enable HTTPS
   - Configure security headers
   - Follow security best practices

3. **Performance**
   - Enable caching
   - Optimize static file serving
   - Monitor resource usage

4. **Logging**
   - Configure access logs
   - Enable error logging
   - Monitor log rotation

## Troubleshooting

Common issues and solutions:

1. **Configuration Issues**
   ```
   Error: invalid configuration
   Solution: Verify nginx.conf syntax
   ```

2. **Port Conflicts**
   ```
   Error: address already in use
   Solution: Change port or stop conflicting service
   ```

3. **SSL Problems**
   ```
   Error: SSL certificate error
   Solution: Check certificate and key files
   ```

## Configuration Example

```nginx
# nginx.conf
worker_processes auto;
events {
    worker_connections 1024;
}

http {
    include mime.types;
    default_type application/octet-stream;

    # Logging
    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;

    # Gzip
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;

    # Virtual Host
    server {
        listen 80;
        server_name example.com;
        root /usr/share/nginx/html;

        location / {
            try_files $uri $uri/ /index.html;
        }

        # Reverse Proxy
        location /api {
            proxy_pass http://backend:8080;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
```

## Advanced Usage

### Load Balancing

```go
func (m *MyModule) LoadBalancer(ctx context.Context) (*Service, error) {
    nginx := dag.Nginx().New()
    
    // Configure load balancing
    return nginx.WithUpstream(
        ctx,
        "backend",
        []string{
            "server1:8080",
            "server2:8080",
            "server3:8080",
        },
        map[string]string{
            "least_conn": "",
            "keepalive": "32",
        },
    ).Server(ctx, "1.24", 80, nil)
}
```

### Custom Configuration

```go
func (m *MyModule) CustomConfig(ctx context.Context) (*Service, error) {
    nginx := dag.Nginx().New()
    
    // Use custom configuration
    return nginx.WithConfig(
        ctx,
        dag.File("./nginx.conf"),
        map[string]string{
            "client_max_body_size": "100M",
            "keepalive_timeout": "65",
        },
    ).Server(ctx, "1.24", 80, nil)
}
``` 