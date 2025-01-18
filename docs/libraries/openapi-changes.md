---
layout: default
title: OpenAPI Changes Module
parent: Libraries
nav_order: 10
---

# OpenAPI Changes Module

The OpenAPI Changes module provides integration with OpenAPI diff tools to detect and analyze changes between different versions of OpenAPI specifications. This module helps you track API changes and maintain backward compatibility in your Dagger pipelines.

## Features

- OpenAPI spec comparison
- Breaking change detection
- Compatibility checking
- Change reporting
- Version tracking
- Custom rules support
- Multiple format support
- CI/CD integration

## Installation

To use the OpenAPI Changes module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/openapi-changes"
)
```

## Usage Examples

### Basic Spec Comparison

```go
func (m *MyModule) Example(ctx context.Context) (string, error) {
    changes := dag.OpenAPIChanges().New()
    
    // Compare OpenAPI specs
    return changes.Compare(
        ctx,
        dag.File("./api-v1.yaml"),  // old spec
        dag.File("./api-v2.yaml"),  // new spec
        "markdown",                 // output format
    )
}
```

### Breaking Changes Check

```go
func (m *MyModule) CheckBreaking(ctx context.Context) error {
    changes := dag.OpenAPIChanges().New()
    
    // Check for breaking changes
    return changes.CheckBreaking(
        ctx,
        dag.File("./api-v1.yaml"),
        dag.File("./api-v2.yaml"),
        map[string]string{
            "ignoreHeaderChanges": "true",
            "failOnIncompatible": "true",
        },
    )
}
```

### Custom Rules

```go
func (m *MyModule) CustomRules(ctx context.Context) (string, error) {
    changes := dag.OpenAPIChanges().New()
    
    // Apply custom rules
    return changes.CompareWithRules(
        ctx,
        dag.File("./api-v1.yaml"),
        dag.File("./api-v2.yaml"),
        dag.File("./rules.json"),
        "html",
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: API Changes
on: [pull_request]

jobs:
  api-diff:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Check API Changes
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/openapi-changes
          args: |
            do -p '
              changes := OpenAPIChanges().New()
              changes.Compare(
                ctx,
                dag.File("./api-v1.yaml"),
                dag.File("./api-v2.yaml"),
                "markdown",
              )
            '
```

## API Reference

### OpenAPIChanges

Main module struct that provides access to OpenAPI diff functionality.

#### Constructor

- `New() *OpenAPIChanges`
  - Creates a new OpenAPIChanges instance
  - Default version: "latest"
  - Default platform: "linux/amd64"

#### Methods

- `Compare(ctx context.Context, oldSpec *File, newSpec *File, format string) (string, error)`
  - Compares two OpenAPI specifications
  - Parameters:
    - `oldSpec`: Original specification file
    - `newSpec`: New specification file
    - `format`: Output format (markdown, html, json)
  
- `CheckBreaking(ctx context.Context, oldSpec *File, newSpec *File, config map[string]string) error`
  - Checks for breaking changes
  - Parameters:
    - `oldSpec`: Original specification file
    - `newSpec`: New specification file
    - `config`: Configuration options
  
- `CompareWithRules(ctx context.Context, oldSpec *File, newSpec *File, rules *File, format string) (string, error)`
  - Compares specs with custom rules
  - Parameters:
    - `oldSpec`: Original specification file
    - `newSpec`: New specification file
    - `rules`: Custom rules file
    - `format`: Output format

## Best Practices

1. **Version Control**
   - Track spec versions
   - Document changes
   - Use semantic versioning

2. **Change Management**
   - Review breaking changes
   - Plan deprecations
   - Communicate changes

3. **Compatibility**
   - Maintain backward compatibility
   - Use proper deprecation
   - Follow API guidelines

4. **Documentation**
   - Keep specs up to date
   - Document changes clearly
   - Include migration guides

## Troubleshooting

Common issues and solutions:

1. **Parse Errors**
   ```
   Error: invalid OpenAPI specification
   Solution: Validate spec syntax and format
   ```

2. **Rule Conflicts**
   ```
   Error: conflicting rules detected
   Solution: Review and update custom rules
   ```

3. **Format Issues**
   ```
   Error: unsupported output format
   Solution: Use supported format (markdown, html, json)
   ```

## Rules Example

```json
{
  "rules": {
    "breaking": {
      "path": {
        "delete": true,
        "rename": true
      },
      "parameter": {
        "required": true,
        "delete": true
      },
      "response": {
        "delete": true,
        "statusCode": true
      }
    },
    "ignore": {
      "paths": ["/internal/*"],
      "headers": ["X-Internal-*"]
    }
  }
}
```

## Advanced Usage

### Custom Report Generation

```go
func (m *MyModule) CustomReport(ctx context.Context) error {
    changes := dag.OpenAPIChanges().New()
    
    // Generate custom report
    report, err := changes.CompareWithTemplate(
        ctx,
        dag.File("./api-v1.yaml"),
        dag.File("./api-v2.yaml"),
        dag.File("./template.html"),
        map[string]string{
            "title": "API Changes Report",
            "version": "2.0.0",
        },
    )
    
    if err != nil {
        return err
    }
    
    // Process report
    return nil
}
```

### Batch Comparison

```go
func (m *MyModule) BatchCompare(ctx context.Context) error {
    changes := dag.OpenAPIChanges().New()
    
    // Compare multiple versions
    versions := []string{"v1", "v2", "v3"}
    for i := 0; i < len(versions)-1; i++ {
        oldSpec := fmt.Sprintf("./api-%s.yaml", versions[i])
        newSpec := fmt.Sprintf("./api-%s.yaml", versions[i+1])
        
        _, err := changes.Compare(
            ctx,
            dag.File(oldSpec),
            dag.File(newSpec),
            "markdown",
        )
        if err != nil {
            return err
        }
    }
    
    return nil
} 