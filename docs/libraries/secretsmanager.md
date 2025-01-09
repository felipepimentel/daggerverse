---
layout: default
title: Secrets Manager Module
parent: Libraries
nav_order: 15
---

# Secrets Manager Module

The Secrets Manager module provides integration with [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/), a service that helps you protect access to your applications, services, and IT resources. This module allows you to manage secrets in your Dagger pipelines.

## Features

- Secret management
- Secret rotation
- Version control
- Access control
- Encryption management
- Cross-region replication
- Resource tagging
- Secret recovery

## Installation

To use the Secrets Manager module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/secretsmanager"
)
```

## Usage Examples

### Basic Secret Retrieval

```go
func (m *MyModule) Example(ctx context.Context) (*Secret, error) {
    sm := dag.SecretsManager().New(
        dag.SetSecret("AWS_ACCESS_KEY_ID", "key"),
        dag.SetSecret("AWS_SECRET_ACCESS_KEY", "secret"),
    )
    
    // Get secret value
    return sm.GetSecret(
        ctx,
        "my-secret-name",
        "us-west-2",  // AWS region
    )
}
```

### Secret Creation

```go
func (m *MyModule) CreateSecret(ctx context.Context) error {
    sm := dag.SecretsManager().New(
        dag.SetSecret("AWS_ACCESS_KEY_ID", "key"),
        dag.SetSecret("AWS_SECRET_ACCESS_KEY", "secret"),
    )
    
    // Create new secret
    return sm.CreateSecret(
        ctx,
        "new-secret",
        dag.SetSecret("SECRET_VALUE", "mysecretvalue"),
        "us-west-2",
        map[string]string{
            "Environment": "Production",
            "Project": "MyApp",
        },
    )
}
```

### Secret Rotation

```go
func (m *MyModule) RotateSecret(ctx context.Context) error {
    sm := dag.SecretsManager().New(
        dag.SetSecret("AWS_ACCESS_KEY_ID", "key"),
        dag.SetSecret("AWS_SECRET_ACCESS_KEY", "secret"),
    )
    
    // Rotate secret value
    return sm.RotateSecret(
        ctx,
        "database-credentials",
        dag.SetSecret("NEW_SECRET", "newvalue"),
        "us-west-2",
        map[string]string{
            "RotationPeriod": "30d",
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
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/secretsmanager
          args: |
            do -p '
              sm := SecretsManager().New(
                dag.SetSecret("AWS_ACCESS_KEY_ID", AWS_ACCESS_KEY_ID),
                dag.SetSecret("AWS_SECRET_ACCESS_KEY", AWS_SECRET_ACCESS_KEY),
              )
              sm.GetSecret(
                ctx,
                "my-secret-name",
                "us-west-2",
              )
            '
```

## API Reference

### SecretsManager

Main module struct that provides access to AWS Secrets Manager functionality.

#### Constructor

- `New(accessKey *Secret, secretKey *Secret) *SecretsManager`
  - Creates a new Secrets Manager instance
  - Parameters:
    - `accessKey`: AWS access key ID
    - `secretKey`: AWS secret access key

#### Methods

- `GetSecret(ctx context.Context, name string, region string) (*Secret, error)`
  - Retrieves a secret value
  - Parameters:
    - `name`: Secret name
    - `region`: AWS region
  
- `CreateSecret(ctx context.Context, name string, value *Secret, region string, tags map[string]string) error`
  - Creates a new secret
  - Parameters:
    - `name`: Secret name
    - `value`: Secret value
    - `region`: AWS region
    - `tags`: Resource tags
  
- `RotateSecret(ctx context.Context, name string, newValue *Secret, region string, config map[string]string) error`
  - Rotates a secret value
  - Parameters:
    - `name`: Secret name
    - `newValue`: New secret value
    - `region`: AWS region
    - `config`: Rotation configuration

## Best Practices

1. **Secret Management**
   - Use descriptive names
   - Implement rotation
   - Tag resources properly

2. **Security**
   - Limit access rights
   - Enable encryption
   - Monitor usage

3. **Organization**
   - Use naming conventions
   - Group related secrets
   - Document purpose

4. **Recovery**
   - Enable backups
   - Plan recovery
   - Test procedures

## Troubleshooting

Common issues and solutions:

1. **Authentication Issues**
   ```
   Error: invalid credentials
   Solution: Verify AWS credentials and permissions
   ```

2. **Access Problems**
   ```
   Error: access denied
   Solution: Check IAM roles and policies
   ```

3. **Region Issues**
   ```
   Error: region not found
   Solution: Verify AWS region configuration
   ```

## Configuration Example

```json
{
  "SecretString": "{\n  \"username\":\"admin\",\n  \"password\":\"secret123\"\n}",
  "Tags": [
    {
      "Key": "Environment",
      "Value": "Production"
    },
    {
      "Key": "Project",
      "Value": "MyApp"
    }
  ],
  "RotationRules": {
    "AutomaticallyAfterDays": 30
  }
}
```

## Advanced Usage

### Multi-Region Management

```go
func (m *MyModule) MultiRegion(ctx context.Context) error {
    sm := dag.SecretsManager().New(
        dag.SetSecret("AWS_ACCESS_KEY_ID", "key"),
        dag.SetSecret("AWS_SECRET_ACCESS_KEY", "secret"),
    )
    
    // Manage secret across regions
    regions := []string{"us-west-2", "us-east-1", "eu-west-1"}
    for _, region := range regions {
        err := sm.ReplicateSecret(
            ctx,
            "global-secret",
            region,
            map[string]string{
                "KmsKeyId": "alias/aws/secretsmanager",
            },
        )
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### Secret Version Management

```go
func (m *MyModule) VersionManagement(ctx context.Context) error {
    sm := dag.SecretsManager().New(
        dag.SetSecret("AWS_ACCESS_KEY_ID", "key"),
        dag.SetSecret("AWS_SECRET_ACCESS_KEY", "secret"),
    )
    
    // Manage secret versions
    return sm.ManageVersions(
        ctx,
        "my-secret",
        "us-west-2",
        map[string]string{
            "RetentionPeriod": "7d",
            "MaxVersions": "5",
        },
    )
} 