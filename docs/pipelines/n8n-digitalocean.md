# N8N DigitalOcean Pipeline (Deprecated)

This module has been deprecated and its functionality has been incorporated into the main n8n module.

Please use the n8n module with the DigitalOceanProvider instead. See [N8N Pipeline](n8n.md) for details.

## Migration Guide

### Old Way

```go
n8n := dag.N8NDigitalocean().
    WithSource(source).
    WithRegistry("registry.digitalocean.com/your-registry").
    WithTag("latest").
    WithDOConfig(&DOConfig{
        Token:        token,
        Region:       "nyc",
        AppName:      "my-n8n",
        InstanceSize: "basic-xxs",
    })
```

### New Way

```go
// Create a DigitalOcean provider
provider := &DigitalOceanProvider{
    Token:        token,
    Region:       "nyc1",
    AppName:      "my-n8n",
    InstanceSize: "basic-xxs",
    Domain:       "n8n.example.com",  // Optional, adds Caddy reverse proxy if specified
}

// Configure n8n with the provider
n8n := dag.N8N().
    WithSource(source).
    WithRegistry("registry.digitalocean.com/myregistry/n8n").
    WithTag("latest").
    WithProvider(provider)

// Deploy
container, err := n8n.Deploy(ctx)
``` 