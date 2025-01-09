---
layout: default
title: Terraform Module
parent: Libraries
nav_order: 17
---

# Terraform Module

The Terraform module provides integration with [Terraform](https://www.terraform.io/), a tool for building, changing, and versioning infrastructure safely and efficiently. This module allows you to manage infrastructure as code in your Dagger pipelines.

## Features

- Infrastructure provisioning
- State management
- Resource planning
- Configuration validation
- Module management
- Workspace handling
- Variable management
- Output retrieval

## Installation

To use the Terraform module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/terraform"
)
```

## Usage Examples

### Basic Infrastructure Deployment

```go
func (m *MyModule) Example(ctx context.Context) error {
    tf := dag.Terraform().New()
    
    // Deploy infrastructure
    return tf.Apply(
        ctx,
        dag.Directory("./terraform"),  // configuration directory
        map[string]string{
            "region": "us-west-2",
            "environment": "production",
        },
    )
}
```

### Plan Infrastructure Changes

```go
func (m *MyModule) PlanChanges(ctx context.Context) (*File, error) {
    tf := dag.Terraform().New()
    
    // Generate plan
    return tf.Plan(
        ctx,
        dag.Directory("./terraform"),
        map[string]string{
            "region": "us-west-2",
            "instance_type": "t3.micro",
        },
        true,  // detailed output
    )
}
```

### Workspace Management

```go
func (m *MyModule) ManageWorkspace(ctx context.Context) error {
    tf := dag.Terraform().New()
    
    // Create and use workspace
    return tf.WithWorkspace(
        ctx,
        "staging",
        dag.Directory("./terraform"),
        map[string]string{
            "environment": "staging",
            "region": "us-east-1",
        },
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Infrastructure Deployment
on: [push]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Deploy Infrastructure
        uses: dagger/dagger-action@v1
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/terraform
          args: |
            do -p '
              tf := Terraform().New()
              tf.Apply(
                ctx,
                dag.Directory("./terraform"),
                map[string]string{
                  "region": "us-west-2",
                  "environment": "production",
                },
              )
            '
```

## API Reference

### Terraform

Main module struct that provides access to Terraform functionality.

#### Constructor

- `New() *Terraform`
  - Creates a new Terraform instance
  - Default version: "latest"
  - Default platform: "linux/amd64"

#### Methods

- `Apply(ctx context.Context, config *Directory, vars map[string]string) error`
  - Applies infrastructure changes
  - Parameters:
    - `config`: Terraform configuration directory
    - `vars`: Variable values
  
- `Plan(ctx context.Context, config *Directory, vars map[string]string, detailed bool) (*File, error)`
  - Generates execution plan
  - Parameters:
    - `config`: Terraform configuration directory
    - `vars`: Variable values
    - `detailed`: Include detailed output
  
- `WithWorkspace(ctx context.Context, name string, config *Directory, vars map[string]string) error`
  - Manages Terraform workspace
  - Parameters:
    - `name`: Workspace name
    - `config`: Terraform configuration directory
    - `vars`: Variable values

## Best Practices

1. **State Management**
   - Use remote state
   - Enable state locking
   - Backup state files

2. **Configuration**
   - Use modules
   - Version providers
   - Document variables

3. **Security**
   - Encrypt sensitive data
   - Use least privilege
   - Rotate credentials

4. **Organization**
   - Structure workspaces
   - Tag resources
   - Maintain consistency

## Troubleshooting

Common issues and solutions:

1. **State Issues**
   ```
   Error: state lock acquisition failed
   Solution: Check state locking mechanism
   ```

2. **Provider Problems**
   ```
   Error: provider not found
   Solution: Verify provider configuration
   ```

3. **Variable Errors**
   ```
   Error: required variable not set
   Solution: Check variable definitions
   ```

## Configuration Example

```hcl
# main.tf
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
  
  backend "s3" {
    bucket = "terraform-state"
    key    = "prod/terraform.tfstate"
    region = "us-west-2"
  }
}

provider "aws" {
  region = var.region
}

variable "region" {
  type = string
}

variable "environment" {
  type = string
}

resource "aws_instance" "example" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  
  tags = {
    Environment = var.environment
  }
}
```

## Advanced Usage

### Remote State Management

```go
func (m *MyModule) RemoteState(ctx context.Context) error {
    tf := dag.Terraform().New()
    
    // Configure remote state
    return tf.WithBackend(
        ctx,
        "s3",
        map[string]string{
            "bucket": "terraform-state",
            "key": "prod/terraform.tfstate",
            "region": "us-west-2",
            "encrypt": "true",
        },
    ).Apply(ctx, dag.Directory("./terraform"), nil)
}
```

### Module Management

```go
func (m *MyModule) ManageModules(ctx context.Context) error {
    tf := dag.Terraform().New()
    
    // Initialize and use modules
    return tf.WithModules(
        ctx,
        map[string]string{
            "vpc": "terraform-aws-modules/vpc/aws",
            "eks": "terraform-aws-modules/eks/aws",
        },
        map[string]string{
            "vpc_version": "3.0.0",
            "eks_version": "18.0.0",
        },
    ).Apply(ctx, dag.Directory("./terraform"), nil)
}
``` 