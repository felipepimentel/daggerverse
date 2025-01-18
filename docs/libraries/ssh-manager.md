# SSH Manager Module

The SSH Manager module provides functionality for managing SSH keys in DigitalOcean. It can generate new SSH key pairs, register them with DigitalOcean, and manage existing keys.

## Features

- Generate new SSH key pairs (4096-bit RSA)
- Register keys with DigitalOcean
- List existing keys
- Delete keys
- Get key details by fingerprint

## Usage

### Import

```go
import sshmanager "github.com/felipepimentel/daggerverse/libraries/ssh-manager"
```

### Initialize

```go
sshManager := sshmanager.New().WithToken("your-digitalocean-token")
```

### Generate and Register a New Key

```go
key, err := sshManager.GenerateKey(ctx, "my-key-name")
if err != nil {
    return err
}

fmt.Printf("Key ID: %d\n", key.ID)
fmt.Printf("Fingerprint: %s\n", key.Fingerprint)
fmt.Printf("Private Key:\n%s\n", key.PrivateKey)
fmt.Printf("Public Key:\n%s\n", key.PublicKey)
```

### List Keys

```go
keys, err := sshManager.ListKeys(ctx)
if err != nil {
    return err
}

for _, key := range keys {
    fmt.Printf("Name: %s, Fingerprint: %s\n", key.Name, key.Fingerprint)
}
```

### Get Key Details

```go
key, err := sshManager.GetKey(ctx, "fingerprint")
if err != nil {
    return err
}

fmt.Printf("Name: %s\n", key.Name)
fmt.Printf("Public Key: %s\n", key.PublicKey)
```

### Delete Key

```go
err := sshManager.DeleteKey(ctx, "fingerprint")
if err != nil {
    return err
}
```

## Integration with n8n-digitalocean Pipeline

The SSH Manager module can be used to automate SSH key management in the n8n-digitalocean pipeline. Here's an example:

```go
func (n *N8NDigitalOcean) Deploy(ctx context.Context) (*dagger.Container, error) {
    // Create SSH manager
    sshManager := sshmanager.New().WithToken(n.DigitalOcean.Token)

    // Generate ephemeral SSH key
    key, err := sshManager.GenerateKey(ctx, fmt.Sprintf("n8n-%s", n.AppName))
    if err != nil {
        return nil, fmt.Errorf("failed to generate SSH key: %w", err)
    }

    // Use the key for deployment
    // ...

    // Clean up the key after deployment
    defer func() {
        if err := sshManager.DeleteKey(ctx, key.Fingerprint); err != nil {
            fmt.Printf("Warning: failed to delete SSH key: %v\n", err)
        }
    }()

    // Continue with deployment
    // ...
}
```

## Error Handling

The module returns detailed error messages for common failure scenarios:

- Missing DigitalOcean API token
- Failed key generation
- Failed key registration
- Key not found
- Network errors
- API rate limiting

## Security Considerations

1. The module generates 4096-bit RSA keys for maximum security
2. Private keys are never stored, only returned to the caller
3. Keys can be ephemeral (created for a single deployment and then deleted)
4. All communication with DigitalOcean API is done over HTTPS

## Dependencies

- `github.com/digitalocean/godo` - DigitalOcean API client
- `golang.org/x/crypto/ssh` - SSH key generation
- `golang.org/x/oauth2` - OAuth2 authentication for DigitalOcean API 