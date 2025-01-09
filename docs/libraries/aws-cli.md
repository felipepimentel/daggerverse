---
layout: default
title: AWS CLI Module
parent: Libraries
nav_order: 1
---

# AWS CLI Module
{: .fs-9 }

Seamless integration with AWS Command Line Interface for Dagger pipelines.
{: .fs-6 .fw-300 }

[Get Started](#installation){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View on GitHub](https://github.com/felipepimentel/daggerverse/tree/main/aws-cli){: .btn .fs-5 .mb-4 .mb-md-0 }

---

## Overview

The AWS CLI module provides a seamless integration with AWS Command Line Interface, allowing you to interact with AWS services directly from your Dagger pipelines.

{: .note }
> This module requires AWS credentials to be configured properly. See the [Configuration](#configuration) section for details.

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
    awscli "github.com/felipepimentel/daggerverse/libraries/aws-cli"
)
```

## Basic Usage

Here's a simple example of using the AWS CLI module:

```go
func (m *MyModule) Example(ctx context.Context) (*Container, error) {
    awsCli := dag.AwsCli().New()
    
    // Execute a simple AWS command
    return awsCli.Exec([]string{"sts", "get-caller-identity"}), nil
}
```

## Configuration

### Custom Configuration

You can customize the AWS CLI configuration:

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

### Environment Variables

The module supports the following environment variables:

| Variable | Description |
|----------|-------------|
| `AWS_ACCESS_KEY_ID` | AWS access key |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key |
| `AWS_REGION` | AWS region |
| `AWS_PROFILE` | AWS profile name |

{: .important }
> Always use secrets for storing sensitive AWS credentials. Never hardcode them in your pipeline code.

## Examples

### List S3 Buckets

```go
awsCli.Exec([]string{"s3", "ls"})
```

### Deploy to ECS

```go
awsCli.Exec([]string{
    "ecs", "update-service",
    "--cluster", "my-cluster",
    "--service", "my-service",
    "--force-new-deployment"
})
```

### Using Temporary Credentials

```go
awsCli.WithTemporaryCredentials(
    dag.SetSecret("AWS_ACCESS_KEY_ID", tempKey),
    dag.SetSecret("AWS_SECRET_ACCESS_KEY", tempSecret),
    dag.SetSecret("AWS_SESSION_TOKEN", sessionToken),
)
```

## Best Practices

1. Always use secrets for credentials
2. Set specific AWS CLI versions for reproducibility
3. Configure region explicitly when needed
4. Use temporary credentials when possible
5. Follow the principle of least privilege

## Troubleshooting

Common issues and solutions:

{: .warning }
> If you encounter permission errors, ensure your AWS credentials have the necessary IAM permissions.

1. **Credentials not found**
   - Check if credentials are properly set
   - Verify secret names match exactly

2. **Region issues**
   - Explicitly set region in configuration
   - Check for environment variable conflicts

## API Reference

### Core Methods

| Method | Description |
|--------|-------------|
| `New()` | Create new AWS CLI instance |
| `Exec([]string)` | Execute AWS CLI command |
| `WithStaticCredentials()` | Configure static credentials |
| `WithTemporaryCredentials()` | Configure temporary credentials |

### Configuration Options

| Option | Type | Description |
|--------|------|-------------|
| `Version` | string | AWS CLI version |
| `Region` | string | AWS region |
| `Profile` | string | AWS profile name |

## Contributing

Please read our [Contributing Guidelines](https://github.com/felipepimentel/daggerverse/blob/main/CONTRIBUTING.md) for details on submitting pull requests. 