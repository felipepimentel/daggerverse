---
layout: default
title: Vault Module
parent: Libraries
nav_order: 18
---

# Vault Module

The Vault module provides integration with [HashiCorp Vault](https://www.vaultproject.io/), a tool for secrets management, encryption as a service, and privileged access management. This module allows you to manage secrets and security in your Dagger pipelines.

## Features

- Secret management
- Dynamic secrets
- Encryption as a service
- Authentication methods
- Policy management
- Key rotation
- Audit logging
- High availability

## Installation

To use the Vault module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/vault"
)
```

## Usage Examples

### Basic Secret Management

```go
func (m *MyModule) Example(ctx context.Context) (*Secret, error) {
    vault := dag.Vault().New(
        dag.SetSecret("VAULT_TOKEN", "token"),
        "http://vault:8200",  // Vault address
    )
    
    // Read secret
    return vault.ReadSecret(
        ctx,
        "secret/data/myapp",  // path
        "password",           // key
    )
}
```

### Dynamic Database Credentials

```go
func (m *MyModule) GetDBCreds(ctx context.Context) (*Secret, error) {
    vault := dag.Vault().New(
        dag.SetSecret("VAULT_TOKEN", "token"),
        "http://vault:8200",
    )
    
    // Generate database credentials
    return vault.GetDynamicCreds(
        ctx,
        "database/creds/readonly",
        map[string]string{
            "ttl": "1h",
            "role": "readonly",
        },
    )
}
```

### Policy Management

```go
func (m *MyModule) ManagePolicy(ctx context.Context) error {
    vault := dag.Vault().New(
        dag.SetSecret("VAULT_TOKEN", "token"),
        "http://vault:8200",
    )
    
    // Create policy
    return vault.CreatePolicy(
        ctx,
        "app-policy",
        dag.File("./policy.hcl"),
        map[string]string{
            "description": "Application policy",
        },
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Secret Management
on: [push]

jobs:
  secrets:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Get Secret
        uses: dagger/dagger-action@v1
        env:
          VAULT_TOKEN: ${{ secrets.VAULT_TOKEN }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/vault
          args: |
            do -p '
              vault := Vault().New(
                dag.SetSecret("VAULT_TOKEN", VAULT_TOKEN),
                "http://vault:8200",
              )
              vault.ReadSecret(
                ctx,
                "secret/data/myapp",
                "password",
              )
            '
```

## API Reference

### Vault

Main module struct that provides access to Vault functionality.

#### Constructor

- `New(token *Secret, address string) *Vault`
  - Creates a new Vault instance
  - Parameters:
    - `token`: Vault authentication token
    - `address`: Vault server address

#### Methods

- `ReadSecret(ctx context.Context, path string, key string) (*Secret, error)`
  - Reads a secret value
  - Parameters:
    - `path`: Secret path
    - `key`: Secret key
  
- `GetDynamicCreds(ctx context.Context, path string, config map[string]string) (*Secret, error)`
  - Generates dynamic credentials
  - Parameters:
    - `path`: Credentials path
    - `config`: Configuration options
  
- `CreatePolicy(ctx context.Context, name string, policy *File, config map[string]string) error`
  - Creates access policy
  - Parameters:
    - `name`: Policy name
    - `policy`: Policy file
    - `config`: Policy configuration

## Best Practices

1. **Secret Management**
   - Use namespaces
   - Rotate secrets
   - Audit access

2. **Authentication**
   - Use appropriate auth methods
   - Limit token lifetimes
   - Enable MFA

3. **Policy**
   - Follow least privilege
   - Document policies
   - Review regularly

4. **Security**
   - Enable audit logging
   - Monitor access
   - Backup data

## Troubleshooting

Common issues and solutions:

1. **Authentication Issues**
   ```
   Error: permission denied
   Solution: Check token permissions
   ```

2. **Connection Problems**
   ```
   Error: connection refused
   Solution: Verify Vault address and status
   ```

3. **Policy Errors**
   ```
   Error: invalid policy syntax
   Solution: Validate policy HCL format
   ```

## Configuration Example

```hcl
# policy.hcl
path "secret/data/myapp/*" {
  capabilities = ["read", "list"]
}

path "database/creds/readonly" {
  capabilities = ["read"]
}

path "auth/token/create" {
  capabilities = ["create", "update"]
}

path "sys/policies/acl/*" {
  capabilities = ["read"]
}
```

## Advanced Usage

### High Availability Setup

```go
func (m *MyModule) HASetup(ctx context.Context) error {
    vault := dag.Vault().New(
        dag.SetSecret("VAULT_TOKEN", "token"),
        "http://vault:8200",
    )
    
    // Configure HA
    return vault.ConfigureHA(
        ctx,
        map[string]string{
            "api_addr": "https://vault.example.com",
            "cluster_addr": "https://vault.example.com:8201",
            "storage_type": "consul",
        },
    )
}
```

### Transit Encryption

```go
func (m *MyModule) TransitEncryption(ctx context.Context) error {
    vault := dag.Vault().New(
        dag.SetSecret("VAULT_TOKEN", "token"),
        "http://vault:8200",
    )
    
    // Use transit encryption
    return vault.TransitEncrypt(
        ctx,
        "mykey",
        dag.SetSecret("PLAINTEXT", "secret-data"),
        map[string]string{
            "type": "aes256-gcm96",
            "convergent_encryption": "true",
        },
    )
}
``` 