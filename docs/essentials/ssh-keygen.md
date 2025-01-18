---
layout: default
title: SSH-Keygen Module
parent: Essentials
nav_order: 15
---

# SSH-Keygen Module

The SSH-Keygen module provides functionality for generating SSH keys in your Dagger pipelines. It allows you to create and manage SSH key pairs with various algorithms and configurations.

## Features

- SSH key generation
- Multiple algorithms
- Custom key sizes
- Passphrase support
- Key format selection
- Comment management
- Output formatting
- File permissions
- Error handling
- Secure operations

## Installation

To use the SSH-Keygen module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/ssh-keygen"
)
```

## Usage Examples

### Basic Key Generation

```go
func (m *MyModule) Example(ctx context.Context) (*Directory, error) {
    sshkeygen := dag.SshKeygen()
    
    // Generate key pair
    return sshkeygen.Generate(
        ctx,
        "",           // algorithm (default: rsa)
        "",           // bits (default: 2048)
        "",           // comment
        "",           // passphrase
        "",           // filename
    )
}
```

### Custom Algorithm and Size

```go
func (m *MyModule) CustomKey(ctx context.Context) (*Directory, error) {
    sshkeygen := dag.SshKeygen()
    
    // Generate ed25519 key
    return sshkeygen.Generate(
        ctx,
        "ed25519",   // algorithm
        "",          // bits (not used for ed25519)
        "deploy@example.com",  // comment
        "",          // passphrase
        "deploy_key", // filename
    )
}
```

### Protected Key

```go
func (m *MyModule) ProtectedKey(ctx context.Context) (*Directory, error) {
    sshkeygen := dag.SshKeygen()
    
    // Generate key with passphrase
    return sshkeygen.Generate(
        ctx,
        "rsa",       // algorithm
        "4096",      // bits
        "secure@example.com",  // comment
        dag.SetSecret("KEY_PASS", "your-passphrase"),  // passphrase
        "secure_key", // filename
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Generate Keys
on: [push]

jobs:
  keys:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Generate SSH Keys
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/ssh-keygen
          args: |
            do -p '
              sshkeygen := SshKeygen()
              sshkeygen.Generate(
                ctx,
                "ed25519",
                "",
                "github-action@example.com",
                "",
                "github_deploy",
              )
            '
```

## API Reference

### SshKeygen

Main module struct that provides access to SSH key generation functionality.

#### Methods

- `Generate(ctx context.Context, algorithm string, bits string, comment string, passphrase *Secret, filename string) (*Directory, error)`
  - Generates SSH key pair
  - Parameters:
    - `algorithm`: Key algorithm (rsa, dsa, ecdsa, ed25519)
    - `bits`: Key size in bits
    - `comment`: Key comment
    - `passphrase`: Key passphrase (optional)
    - `filename`: Output filename
  - Returns directory containing key files

## Best Practices

1. **Key Security**
   - Use strong algorithms
   - Protect private keys
   - Secure passphrases

2. **Algorithm Selection**
   - Choose modern algorithms
   - Use appropriate sizes
   - Consider compatibility

3. **Key Management**
   - Backup securely
   - Rotate regularly
   - Document usage

4. **File Handling**
   - Set proper permissions
   - Secure storage
   - Clear old keys

## Troubleshooting

Common issues and solutions:

1. **Generation Errors**
   ```
   Error: invalid algorithm
   Solution: Use supported algorithm
   ```

2. **Permission Issues**
   ```
   Error: permission denied
   Solution: Check file permissions
   ```

3. **Key Format**
   ```
   Error: invalid key format
   Solution: Verify algorithm/size
   ```

## Configuration Example

```yaml
# ssh-keygen-config.yaml
defaults:
  algorithm: ed25519
  comment: "automated@example.com"
  permissions:
    private: "0600"
    public: "0644"
  
algorithms:
  rsa:
    min_bits: 2048
    recommended_bits: 4096
  ed25519:
    recommended: true
```

## Advanced Usage

### Multi-Key Generation

```go
func (m *MyModule) MultiKeys(ctx context.Context) error {
    sshkeygen := dag.SshKeygen()
    
    // Define key configurations
    keys := []struct {
        algo     string
        bits     string
        comment  string
        filename string
    }{
        {"rsa", "4096", "rsa@example.com", "rsa_key"},
        {"ed25519", "", "ed25519@example.com", "ed_key"},
        {"ecdsa", "", "ecdsa@example.com", "ec_key"},
    }
    
    // Generate all keys
    for _, key := range keys {
        dir, err := sshkeygen.Generate(
            ctx,
            key.algo,
            key.bits,
            key.comment,
            "",  // no passphrase
            key.filename,
        )
        if err != nil {
            return err
        }
        
        // Process generated keys
        err = dag.Container().
            From("alpine:latest").
            WithMountedDirectory("/keys", dir).
            WithExec([]string{
                "sh", "-c",
                fmt.Sprintf(
                    "cp /keys/%s* /permanent/keys/",
                    key.filename,
                ),
            }).
            Sync(ctx)
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### Key Verification

```go
func (m *MyModule) VerifyKey(ctx context.Context) error {
    sshkeygen := dag.SshKeygen()
    
    // Generate key
    dir, err := sshkeygen.Generate(
        ctx,
        "rsa",
        "4096",
        "test@example.com",
        "",
        "test_key",
    )
    if err != nil {
        return err
    }
    
    // Verify key
    return dag.Container().
        From("alpine:latest").
        WithMountedDirectory("/keys", dir).
        WithExec([]string{"apk", "add", "openssh-client"}).
        WithExec([]string{
            "sh", "-c",
            `
            # Check private key
            ssh-keygen -l -f /keys/test_key
            
            # Verify key pair
            diff <(ssh-keygen -y -f /keys/test_key) /keys/test_key.pub
            
            echo "Key verification successful"
            `,
        }).
        Sync(ctx)
} 