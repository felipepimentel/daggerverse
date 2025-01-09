---
layout: default
title: AWS CLI Module
parent: Libraries
nav_order: 2
---

# AWS CLI Module

The AWS CLI module provides a seamless integration with AWS Command Line Interface, allowing you to interact with AWS services directly from your Dagger pipelines.

## Features

- Execute AWS CLI commands in Dagger pipelines
- Support for AWS credentials management
- Region configuration
- Profile management
- Support for temporary credentials
- Custom AWS CLI version selection

## Installation

To use the AWS CLI module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/aws-cli"
)
```

## Usage Examples

### Basic AWS CLI Setup

```go
func (m *MyModule) Example(ctx context.Context) (*Container, error) {
    awsCli := dag.AwsCli().New()
    
    // Execute a simple AWS command
    return awsCli.Exec([]string{"sts", "get-caller-identity"}), nil
}
```

### Custom Configuration

```go
func (m *MyModule) CustomConfig(ctx context.Context) (*Container, error) {
    awsCli := dag.AwsCli().New(
        dagger.AwsCliOpts{
            Version: "2.13.0",
            Region: "us-west-2",
        },
    )
    
    // Add AWS credentials
    awsCli = awsCli.WithStaticCredentials(
        dag.SetSecret("AWS_ACCESS_KEY_ID", accessKey),
        dag.SetSecret("AWS_SECRET_ACCESS_KEY", secretKey),
    )
    
    return awsCli.Exec([]string{"s3", "ls"}), nil
}
```

### Using AWS Profiles

```go
func (m *MyModule) WithProfile(ctx context.Context) (*Container, error) {
    awsCli := dag.AwsCli().New()
    
    // Mount AWS config and credentials
    awsCli = awsCli.
        WithConfig(dag.File("~/.aws/config")).
        WithCredentials(dag.SetSecret("credentials", credentialsFile)).
        WithProfile("production")
    
    return awsCli.Exec([]string{"eks", "list-clusters"}), nil
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: AWS Operations
on: [push]

jobs:
  aws-ops:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: AWS Operations with Dagger
        uses: dagger/dagger-action@v1
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/aws-cli
          args: |
            do -p '
              aws := AwsCli().New()
              aws.WithStaticCredentials(
                dag.SetSecret("AWS_ACCESS_KEY_ID", dag.SetSecret("key", AWS_ACCESS_KEY_ID)),
                dag.SetSecret("AWS_SECRET_ACCESS_KEY", dag.SetSecret("secret", AWS_SECRET_ACCESS_KEY)),
              ).Exec(["s3", "ls"])
            '
```

## API Reference

### AwsCli

Main module struct that provides access to AWS CLI functionality.

#### Constructor

- `New(version string, container *Container, region string) AwsCli`
  - Creates a new AWS CLI instance
  - Parameters:
    - `version`: AWS CLI version (optional, defaults to "latest")
    - `container`: Custom base container (optional)
    - `region`: Default AWS region (optional)

#### Methods

- `WithRegion(region string) AwsCli`
  - Sets the AWS region for all commands
  
- `WithConfig(source *File) AwsCli`
  - Mounts an AWS CLI config file
  
- `WithCredentials(source *Secret) AwsCli`
  - Mounts an AWS CLI credentials file
  
- `WithProfile(profile string) AwsCli`
  - Sets the AWS profile for all commands
  
- `WithStaticCredentials(accessKeyId *Secret, secretAccessKey *Secret, sessionToken *Secret) AwsCli`
  - Sets static AWS credentials
  - `sessionToken` is optional for permanent credentials
  
- `WithTemporaryCredentials(accessKeyId *Secret, secretAccessKey *Secret, sessionToken *Secret) AwsCli`
  - Sets temporary AWS credentials (requires session token)
  
- `WithoutStaticCredentials() AwsCli`
  - Removes previously set static credentials
  
- `Exec(args []string) *Container`
  - Executes an AWS CLI command

## Best Practices

1. **Credential Management**
   - Use secrets for storing credentials
   - Prefer IAM roles when possible
   - Rotate credentials regularly

2. **Region Selection**
   - Set explicit regions for predictability
   - Use environment-specific regions

3. **Error Handling**
   - Check command output for errors
   - Implement proper error handling in pipelines

4. **Security**
   - Never commit credentials to source control
   - Use temporary credentials when possible
   - Follow the principle of least privilege

## Troubleshooting

Common issues and solutions:

1. **Authentication Failures**
   ```
   Error: Unable to locate credentials
   Solution: Ensure credentials are properly mounted or environment variables are set
   ```

2. **Region Issues**
   ```
   Error: Region is required
   Solution: Set region using WithRegion() or in AWS config
   ```

3. **Profile Problems**
   ```
   Error: Profile not found
   Solution: Verify profile exists in mounted config and credentials files
   ```

## Environment Variables

The module respects standard AWS environment variables:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_SESSION_TOKEN`
- `AWS_REGION`
- `AWS_PROFILE`

These can be set using the appropriate `With*` methods or through GitHub Actions secrets. 