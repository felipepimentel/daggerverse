# DigitalOcean Module

A reusable Dagger module for managing DigitalOcean resources. This module provides a clean abstraction over DigitalOcean's API, allowing you to manage droplets, DNS records, and other resources programmatically.

## Features

- Droplet management (create, delete, list, get status)
- DNS record management (create, delete, list)
- Resource monitoring and status checks
- Secure token handling

## Prerequisites

- DigitalOcean account
- DigitalOcean API token
- Dagger CLI installed

## Installation

Add this module as a dependency in your `dagger.json`:

```json
{
  "dependencies": [
    {
      "name": "digitalocean",
      "source": "github.com/felipepimentel/daggerverse/libraries/digitalocean"
    }
  ]
}
```

## Usage

### Creating a New Droplet

```go
do := dag.DigitalOcean().
    WithToken(dag.SetSecret("do_token", os.Getenv("DIGITALOCEAN_TOKEN")))

droplet, err := do.CreateDroplet(ctx, DropletConfig{
    Name:       "my-server",
    Region:     "nyc1",
    Size:       "s-1vcpu-1gb",
    Image:      "ubuntu-20-04-x64",
    SSHKeyID:   "your-ssh-key-id",
    Monitoring: true,
    IPv6:       true,
    Tags:       []string{"production", "web"},
})
```

### Managing DNS Records

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

### Listing Resources

```go
// List all droplets
droplets, err := do.ListDroplets(ctx)

// List DNS records for a domain
records, err := do.ListDNSRecords(ctx, "example.com")
```

## Configuration

### Droplet Configuration

The `DropletConfig` struct allows you to specify:

- `Name`: Droplet name
- `Region`: DigitalOcean region (e.g., "nyc1", "sfo2")
- `Size`: Droplet size (e.g., "s-1vcpu-1gb")
- `Image`: Operating system image
- `SSHKeyID`: SSH key identifier
- `Monitoring`: Enable monitoring
- `IPv6`: Enable IPv6
- `Tags`: Array of tags

### DNS Configuration

The `DNSConfig` struct allows you to specify:

- `Domain`: Domain name
- `Type`: Record type (A, AAAA, CNAME, etc.)
- `Name`: Record name
- `Value`: Record value
- `TTL`: Time to live
- `Priority`: Record priority (for MX records)
- Additional fields for specific record types

## Error Handling

The module includes comprehensive error handling:

- Input validation for required fields
- Timeout handling for long-running operations
- Detailed error messages for troubleshooting

## Best Practices

1. Always use secrets for API tokens
2. Set appropriate timeouts for operations
3. Use tags to organize resources
4. Monitor resource status after creation
5. Clean up resources when no longer needed

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - see LICENSE file for details 