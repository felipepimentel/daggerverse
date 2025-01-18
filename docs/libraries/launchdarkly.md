---
layout: default
title: LaunchDarkly Module
parent: Libraries
nav_order: 8
---

# LaunchDarkly Module

The LaunchDarkly module provides integration with [LaunchDarkly](https://launchdarkly.com/), a feature management platform. This module allows you to manage feature flags and experiments in your Dagger pipelines.

## Features

- Feature flag management
- Environment configuration
- Flag targeting rules
- A/B testing setup
- Rollout management
- SDK key handling
- Multi-environment support
- Audit logging

## Installation

To use the LaunchDarkly module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/launchdarkly"
)
```

## Usage Examples

### Basic Flag Management

```go
func (m *MyModule) Example(ctx context.Context) error {
    ld := dag.LaunchDarkly().New(
        dag.SetSecret("LD_API_KEY", "api-key"),
        "production",  // environment
    )
    
    // Create or update a feature flag
    return ld.CreateFlag(
        ctx,
        "new-feature",
        "A new feature toggle",
        true,           // temporary
        map[string]interface{}{
            "defaultValue": false,
        },
    )
}
```

### Environment Management

```go
func (m *MyModule) ManageEnvironments(ctx context.Context) error {
    ld := dag.LaunchDarkly().New(
        dag.SetSecret("LD_API_KEY", "api-key"),
        "production",
    )
    
    // Create a new environment
    return ld.CreateEnvironment(
        ctx,
        "staging",
        "Staging Environment",
        map[string]interface{}{
            "color": "blue",
            "tags": []string{"staging", "test"},
        },
    )
}
```

### Flag Targeting

```go
func (m *MyModule) ConfigureTargeting(ctx context.Context) error {
    ld := dag.LaunchDarkly().New(
        dag.SetSecret("LD_API_KEY", "api-key"),
        "production",
    )
    
    // Update targeting rules
    return ld.UpdateTargeting(
        ctx,
        "new-feature",
        map[string]interface{}{
            "rules": []map[string]interface{}{
                {
                    "clauses": []map[string]interface{}{
                        {
                            "attribute": "email",
                            "op": "endsWith",
                            "values": []string{"@company.com"},
                        },
                    },
                    "variation": 1,
                },
            },
        },
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: LaunchDarkly Operations
on: [push]

jobs:
  launchdarkly:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Manage Feature Flags
        uses: dagger/dagger-action@v1
        env:
          LD_API_KEY: ${{ secrets.LD_API_KEY }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/launchdarkly
          args: |
            do -p '
              ld := LaunchDarkly().New(
                dag.SetSecret("LD_API_KEY", LD_API_KEY),
                "production",
              )
              ld.CreateFlag(
                ctx,
                "new-feature",
                "A new feature toggle",
                true,
                map[string]interface{}{
                  "defaultValue": false,
                },
              )
            '
```

## API Reference

### LaunchDarkly

Main module struct that provides access to LaunchDarkly functionality.

#### Constructor

- `New(apiKey *Secret, environment string) *LaunchDarkly`
  - Creates a new LaunchDarkly instance
  - Parameters:
    - `apiKey`: LaunchDarkly API key
    - `environment`: Target environment

#### Methods

- `CreateFlag(ctx context.Context, key string, description string, temporary bool, config map[string]interface{}) error`
  - Creates a new feature flag
  - Parameters:
    - `key`: Flag key
    - `description`: Flag description
    - `temporary`: Whether the flag is temporary
    - `config`: Flag configuration
  
- `CreateEnvironment(ctx context.Context, key string, name string, config map[string]interface{}) error`
  - Creates a new environment
  - Parameters:
    - `key`: Environment key
    - `name`: Environment name
    - `config`: Environment configuration
  
- `UpdateTargeting(ctx context.Context, flagKey string, rules map[string]interface{}) error`
  - Updates targeting rules for a flag
  - Parameters:
    - `flagKey`: Flag key
    - `rules`: Targeting rules

## Best Practices

1. **Flag Management**
   - Use descriptive flag keys
   - Document flag purposes
   - Clean up temporary flags

2. **Environment Strategy**
   - Use consistent naming
   - Separate production concerns
   - Manage access control

3. **Targeting Rules**
   - Keep rules simple
   - Test targeting logic
   - Document rule changes

4. **Security**
   - Protect API keys
   - Audit access regularly
   - Monitor flag changes

## Troubleshooting

Common issues and solutions:

1. **Authentication Issues**
   ```
   Error: invalid API key
   Solution: Verify API key and permissions
   ```

2. **Flag Creation Failures**
   ```
   Error: flag key already exists
   Solution: Use unique flag keys or update existing flags
   ```

3. **Targeting Problems**
   ```
   Error: invalid targeting rule
   Solution: Check rule syntax and attributes
   ```

## Configuration Example

```json
{
  "flag": {
    "key": "new-feature",
    "description": "A new feature toggle",
    "temporary": true,
    "variations": [
      {
        "value": false,
        "name": "Off"
      },
      {
        "value": true,
        "name": "On"
      }
    ],
    "rules": [
      {
        "clauses": [
          {
            "attribute": "email",
            "op": "endsWith",
            "values": ["@company.com"]
          }
        ],
        "variation": 1
      }
    ]
  }
}
```

## Advanced Usage

### Experimentation Setup

```go
func (m *MyModule) SetupExperiment(ctx context.Context) error {
    ld := dag.LaunchDarkly().New(
        dag.SetSecret("LD_API_KEY", "api-key"),
        "production",
    )
    
    // Configure an A/B test
    return ld.CreateExperiment(
        ctx,
        "button-color",
        "Button Color Test",
        map[string]interface{}{
            "variations": []map[string]interface{}{
                {"value": "blue", "weight": 50},
                {"value": "green", "weight": 50},
            },
            "metrics": []string{"clicks", "conversions"},
        },
    )
}
```

### Rollout Management

```go
func (m *MyModule) ManageRollout(ctx context.Context) error {
    ld := dag.LaunchDarkly().New(
        dag.SetSecret("LD_API_KEY", "api-key"),
        "production",
    )
    
    // Configure gradual rollout
    return ld.UpdateRollout(
        ctx,
        "new-feature",
        map[string]interface{}{
            "percentage": 25,
            "bucketBy": "userId",
        },
    )
}
``` 