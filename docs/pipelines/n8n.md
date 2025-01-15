# N8N Pipeline

This module provides a comprehensive CI/CD pipeline for n8n workflow automation platform, with support for building, testing, and deploying to various environments including DigitalOcean.

## Features

- Complete CI/CD pipeline for n8n
- Docker container build and publish
- Automated testing
- DigitalOcean App Platform deployment
- Environment variable management
- Registry authentication support
- Configurable instance sizing

## Installation

```bash
dagger mod use github.com/felipepimentel/daggerverse/pipelines/n8n@latest
```

## Usage

### Basic Example

```go
// Initialize the module
n8n := dag.N8N().
    WithSource(dag.Host().Directory(".")).
    WithRegistry("your-registry").
    WithTag("latest")

// Run CI/CD pipeline
container, err := n8n.CD(ctx)
```

### Configuration Options

The module supports the following configuration:

```go
type N8N struct {
    // Source directory containing n8n configuration
    Source *dagger.Directory
    // Environment variables for n8n
    EnvVars []EnvVar
    // Port to expose n8n on
    Port int
    // Registry to publish to
    Registry string
    // Image tag
    Tag string
    // Registry auth token
    RegistryAuth *dagger.Secret
    // DigitalOcean configuration
    DOConfig *DOConfig
}

type EnvVar struct {
    Name  string
    Value string
}

type DOConfig struct {
    Token        *dagger.Secret
    Region       string
    AppName      string
    InstanceSize string
}
```

### Building n8n

```go
container, err := n8n.Build(ctx)
```

### Testing

```go
err := n8n.Test(ctx)
```

### Publishing

```go
container, err := n8n.Publish(ctx)
```

## DigitalOcean Deployment

This module supports direct deployment to DigitalOcean App Platform:

```go
n8n := dag.N8N().
    WithSource(dag.Host().Directory(".")).
    WithRegistry("registry.digitalocean.com/your-registry").
    WithTag("latest").
    WithDOConfig(&DOConfig{
        Token:        dag.SetSecret("do_token", "your-token"),
        Region:       "nyc",
        AppName:      "my-n8n",
        InstanceSize: "basic-xxs",
    })

container, err := n8n.CD(ctx)
```

## GitHub Actions Integration

Create a workflow file `.github/workflows/n8n.yml`:

```yaml
name: N8N CI/CD
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install Dagger CLI
        uses: dagger/dagger-for-github@v5
        with:
          version: "0.15.2"
      
      - name: Deploy n8n
        env:
          DO_TOKEN: ${{ secrets.DO_TOKEN }}
        run: |
          dagger call --progress=plain \
            --source . \
            cd \
            --registry "registry.digitalocean.com/your-registry" \
            --tag "latest" \
            --do-token "$DO_TOKEN" \
            --do-region "nyc" \
            --do-app-name "my-n8n" \
            --do-instance-size "basic-xxs"
```

## Examples

### Custom Port

```go
n8n := dag.N8N().
    WithSource(dag.Host().Directory(".")).
    WithPort(8080)
```

### Environment Variables

```go
n8n := dag.N8N().
    WithSource(dag.Host().Directory(".")).
    WithEnvVars([]EnvVar{
        {Name: "N8N_HOST", Value: "0.0.0.0"},
        {Name: "N8N_PORT", Value: "5678"},
        {Name: "N8N_PROTOCOL", Value: "https"},
    })
```

### Registry Authentication

```go
n8n := dag.N8N().
    WithSource(dag.Host().Directory(".")).
    WithRegistry("your-registry").
    WithTag("latest").
    WithRegistryAuth(dag.SetSecret("registry_auth", "your-auth-token"))
```

## Best Practices

1. **Source Structure**:
   - Keep n8n configuration in a dedicated directory
   - Use `.env` files for environment variables
   - Include proper health check endpoints

2. **Configuration**:
   - Always specify a tag for container images
   - Use secrets for sensitive information
   - Configure appropriate instance sizes

3. **Development**:
   - Run tests before deployment
   - Use staging environments
   - Monitor resource usage

## Common Issues

1. **Build Failures**:
   - Check Node.js version compatibility
   - Verify all dependencies are installed
   - Ensure proper registry permissions

2. **Deployment Issues**:
   - Verify DigitalOcean token permissions
   - Check resource quotas
   - Validate health check configuration

3. **Runtime Problems**:
   - Monitor container logs
   - Check environment variables
   - Verify network connectivity

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](../CONTRIBUTING.md) for details on how to submit pull requests.

## License

This module is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details. 