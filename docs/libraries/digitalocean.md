# DigitalOcean Module Documentation

## Overview

The DigitalOcean module provides a reusable interface for managing DigitalOcean resources through Dagger pipelines. It abstracts common operations for managing droplets, DNS records, SSH keys, and container registries.

## Installation

Add the module as a dependency in your `dagger.json`:

```json
{
  "name": "your-module",
  "dependencies": [
    {
      "name": "digitalocean",
      "source": "github.com/felipepimentel/daggerverse/libraries/digitalocean"
    }
  ]
}
```

## Required Secrets

| Secret | Description | Required |
|--------|-------------|----------|
| `do_token` | DigitalOcean API token with appropriate permissions | Yes |

## Module Interface

### Initialization

```go
do := dag.DigitalOcean().
    WithToken(dag.SetSecret("do_token", os.Getenv("DIGITALOCEAN_TOKEN")))
```

### Configuration Types

#### SSHKeyConfig

Configuration for managing SSH keys:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Name` | string | Yes | Name of the SSH key |
| `PublicKey` | string | Yes | Content of the public key file |

#### RegistryConfig

Configuration for container registry operations:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Name` | string | Yes | Name of the registry |

#### DropletConfig

Configuration for creating a new droplet:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Name` | string | Yes | Name of the droplet |
| `Region` | string | Yes | DigitalOcean region (e.g., "nyc1", "sfo2") |
| `Size` | string | Yes | Droplet size (e.g., "s-1vcpu-1gb") |
| `Image` | string | Yes | Operating system image |
| `SSHKeyID` | string | Yes | SSH key identifier |
| `Monitoring` | bool | No | Enable monitoring (default: false) |
| `IPv6` | bool | No | Enable IPv6 (default: false) |
| `Tags` | []string | No | Array of tags to apply |

#### DNSConfig

Configuration for managing DNS records:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Domain` | string | Yes | Domain name |
| `Type` | string | Yes | Record type (A, AAAA, CNAME, etc.) |
| `Name` | string | Yes | Record name |
| `Value` | string | Yes | Record value |
| `TTL` | int | No | Time to live in seconds |
| `Priority` | int | No | Priority (for MX records) |
| `Weight` | int | No | Weight for load balancing |
| `Port` | int | No | Port (for SRV records) |
| `Flag` | int | No | Flag value |
| `Tag` | string | No | Tag value |

## Available Methods

### SSH Key Management

#### CreateSSHKey

Creates a new SSH key in DigitalOcean.

```go
key, err := do.CreateSSHKey(ctx, SSHKeyConfig{
    Name:      "deployment-key",
    PublicKey: "ssh-ed25519 AAAA...",
})
```

#### ListSSHKeys

Lists all SSH keys with optional format specification.

```go
// List all keys
keys, err := do.ListSSHKeys(ctx, "")

// Get only key IDs
keyIDs, err := do.ListSSHKeys(ctx, "ID")
```

### Registry Management

#### CreateRegistry

Creates a new container registry.

```go
registry, err := do.CreateRegistry(ctx, RegistryConfig{
    Name: "my-registry",
})
```

#### GetRegistry

Gets details about the container registry.

```go
details, err := do.GetRegistry(ctx)
```

#### ListRegistryTags

Lists all tags in a registry repository.

```go
tags, err := do.ListRegistryTags(ctx, "my-registry")
```

#### DeleteRegistry

Deletes a container registry.

```go
err := do.DeleteRegistry(ctx, "my-registry")
```

### Droplet Management

#### CreateDroplet

Creates a new droplet with the specified configuration.

```go
droplet, err := do.CreateDroplet(ctx, DropletConfig{
    Name:       "web-server",
    Region:     "nyc1",
    Size:       "s-1vcpu-1gb",
    Image:      "ubuntu-20-04-x64",
    SSHKeyID:   "12:23:34:45:56:67:78:89:90",
    Monitoring: true,
    IPv6:       true,
    Tags:       []string{"production", "web"},
})
```

#### GetDroplet

Retrieves information about a specific droplet.

```go
// Get all droplet information
info, err := do.GetDroplet(ctx, "web-server", "")

// Get only IP address
ip, err := do.GetDroplet(ctx, "web-server", "PublicIPv4")
```

#### DeleteDroplet

Deletes a droplet by name.

```go
err := do.DeleteDroplet(ctx, "web-server")
```

### DNS Management

#### CreateDNSRecord

Creates a new DNS record.

```go
err := do.CreateDNSRecord(ctx, DNSConfig{
    Domain:   "example.com",
    Type:     "A",
    Name:     "www",
    Value:    "1.2.3.4",
    TTL:      3600,
    Priority: 10,
})
```

#### ListDNSRecords

Lists all DNS records for a domain.

```go
records, err := do.ListDNSRecords(ctx, "example.com")
```

#### DeleteDNSRecord

Deletes a DNS record.

```go
err := do.DeleteDNSRecord(ctx, "example.com", "record-id")
```

### Utility Functions

#### WaitForDroplet

Waits for a droplet to reach a specific status.

```go
err := do.WaitForDroplet(ctx, "web-server", "active", 5*time.Minute)
```

#### ListDroplets

Lists all droplets in the account.

```go
droplets, err := do.ListDroplets(ctx)
```

## Error Handling

The module includes comprehensive error handling:

1. **Input Validation**: Methods validate required fields before making API calls
2. **Timeout Handling**: Operations that may take time include configurable timeouts
3. **Detailed Errors**: Error messages include context about what failed

Example error handling:

```go
droplet, err := do.CreateDroplet(ctx, config)
if err != nil {
    if strings.Contains(err.Error(), "missing required droplet configuration") {
        // Handle missing configuration
    }
    return fmt.Errorf("failed to create droplet: %w", err)
}
```

## Best Practices

1. **Resource Management**
   - Use meaningful names for resources
   - Apply tags for better organization
   - Clean up unused resources

2. **Security**
   - Store the API token securely using Dagger secrets
   - Use minimal required permissions for the API token
   - Regularly rotate API tokens

3. **Performance**
   - Set appropriate timeouts for operations
   - Use monitoring when needed
   - Consider resource costs when selecting droplet sizes

4. **DNS Management**
   - Use appropriate TTL values
   - Verify DNS propagation after changes
   - Document DNS record management

## Example Workflows

### Complete Infrastructure Setup

```go
func SetupInfrastructure(ctx context.Context) error {
    do := dag.DigitalOcean().
        WithToken(dag.SetSecret("do_token", os.Getenv("DIGITALOCEAN_TOKEN")))

    // Create SSH key
    key, err := do.CreateSSHKey(ctx, SSHKeyConfig{
        Name:      "deployment-key",
        PublicKey: os.Getenv("SSH_PUBLIC_KEY"),
    })
    if err != nil {
        return err
    }

    // Create registry
    registry, err := do.CreateRegistry(ctx, RegistryConfig{
        Name: "app-registry",
    })
    if err != nil {
        return err
    }

    // Create droplet
    droplet, err := do.CreateDroplet(ctx, DropletConfig{
        Name:       "app-server",
        Region:     "nyc1",
        Size:       "s-1vcpu-1gb",
        Image:      "docker-20-04",
        SSHKeyID:   key.ID,
        Monitoring: true,
        IPv6:       true,
        Tags:       []string{"production", "app"},
    })
    if err != nil {
        return err
    }

    // Wait for droplet to be ready
    err = do.WaitForDroplet(ctx, droplet.Name, "active", 5*time.Minute)
    if err != nil {
        return err
    }

    // Configure DNS
    err = do.CreateDNSRecord(ctx, DNSConfig{
        Domain: "example.com",
        Type:   "A",
        Name:   "app",
        Value:  droplet.IPv4,
        TTL:    3600,
    })
    if err != nil {
        return err
    }

    return nil
}
```

## Troubleshooting

Common issues and solutions:

1. **Authentication Failures**
   - Verify the API token is correct and has required permissions
   - Check if the token has expired
   - Ensure the token is properly set as a secret

2. **Resource Creation Failures**
   - Check resource name conflicts
   - Verify region availability
   - Confirm account resource limits

3. **DNS Issues**
   - Verify domain ownership
   - Check DNS record syntax
   - Allow time for DNS propagation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit changes following semantic commit messages
4. Submit a pull request

## License

MIT License - see LICENSE file for details 