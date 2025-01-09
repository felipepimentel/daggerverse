---
layout: default
title: Spectral Module
parent: Libraries
nav_order: 16
---

# Spectral Module

The Spectral module provides integration with [Spectral](https://stoplight.io/open-source/spectral), a flexible JSON/YAML linter with out-of-the-box support for OpenAPI v2/v3 and AsyncAPI v2. This module allows you to validate API specifications in your Dagger pipelines.

## Features

- OpenAPI validation
- AsyncAPI validation
- Custom rulesets
- Format checking
- Style enforcement
- Error reporting
- Multiple formats
- CI/CD integration

## Installation

To use the Spectral module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/spectral"
)
```

## Usage Examples

### Basic API Validation

```go
func (m *MyModule) Example(ctx context.Context) error {
    spectral := dag.Spectral().New()
    
    // Validate OpenAPI spec
    return spectral.Lint(
        ctx,
        dag.File("./openapi.yaml"),  // API spec
        nil,                        // default ruleset
        map[string]string{
            "format": "stylish",
        },
    )
}
```

### Custom Ruleset

```go
func (m *MyModule) CustomRules(ctx context.Context) error {
    spectral := dag.Spectral().New()
    
    // Validate with custom rules
    return spectral.LintWithRuleset(
        ctx,
        dag.File("./openapi.yaml"),
        dag.File("./ruleset.yaml"),
        map[string]string{
            "failSeverity": "error",
            "displayFormat": "json",
        },
    )
}
```

### Multiple Files Validation

```go
func (m *MyModule) ValidateMultiple(ctx context.Context) error {
    spectral := dag.Spectral().New()
    
    // Validate multiple specs
    return spectral.LintDirectory(
        ctx,
        dag.Directory("./specs"),
        "**/*.yaml",
        nil,
        map[string]string{
            "ignore": "**/*.test.yaml",
        },
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: API Validation
on: [pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Lint OpenAPI Spec
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/spectral
          args: |
            do -p '
              spectral := Spectral().New()
              spectral.Lint(
                ctx,
                dag.File("./openapi.yaml"),
                nil,
                map[string]string{
                  "format": "stylish",
                },
              )
            '
```

## API Reference

### Spectral

Main module struct that provides access to Spectral functionality.

#### Constructor

- `New() *Spectral`
  - Creates a new Spectral instance
  - Default version: "latest"
  - Default platform: "linux/amd64"

#### Methods

- `Lint(ctx context.Context, spec *File, ruleset *File, config map[string]string) error`
  - Validates API specification
  - Parameters:
    - `spec`: API specification file
    - `ruleset`: Custom ruleset file (optional)
    - `config`: Linting configuration
  
- `LintWithRuleset(ctx context.Context, spec *File, ruleset *File, config map[string]string) error`
  - Validates with custom ruleset
  - Parameters:
    - `spec`: API specification file
    - `ruleset`: Custom ruleset file
    - `config`: Linting configuration
  
- `LintDirectory(ctx context.Context, dir *Directory, pattern string, ruleset *File, config map[string]string) error`
  - Validates multiple files
  - Parameters:
    - `dir`: Directory containing specs
    - `pattern`: File pattern to match
    - `ruleset`: Custom ruleset file (optional)
    - `config`: Linting configuration

## Best Practices

1. **Validation Strategy**
   - Use consistent rulesets
   - Define severity levels
   - Document exceptions

2. **Rule Management**
   - Customize for needs
   - Version control rules
   - Share across teams

3. **Integration**
   - Automate validation
   - Fail fast
   - Report clearly

4. **Maintenance**
   - Update rulesets
   - Monitor changes
   - Review exceptions

## Troubleshooting

Common issues and solutions:

1. **Validation Errors**
   ```
   Error: invalid specification format
   Solution: Check file syntax and format
   ```

2. **Ruleset Problems**
   ```
   Error: ruleset parsing failed
   Solution: Verify ruleset syntax
   ```

3. **Pattern Issues**
   ```
   Error: no files matched pattern
   Solution: Check file patterns and paths
   ```

## Configuration Example

```yaml
# .spectral.yaml
extends: spectral:oas
rules:
  operation-tags: error
  operation-description: warn
  no-$ref-siblings: off
  info-contact: error
  info-description: error
  info-license: warn
  oas3-api-servers: error
  operation-operationId: error
  path-params: error
  typed-enum: error
functions: []
```

## Advanced Usage

### Custom Function Rules

```go
func (m *MyModule) CustomFunctions(ctx context.Context) error {
    spectral := dag.Spectral().New()
    
    // Use custom function rules
    return spectral.LintWithFunctions(
        ctx,
        dag.File("./openapi.yaml"),
        dag.File("./functions.js"),
        map[string]string{
            "functionPath": "./functions",
            "extends": "spectral:oas",
        },
    )
}
```

### Continuous Validation

```go
func (m *MyModule) ContinuousValidation(ctx context.Context) error {
    spectral := dag.Spectral().New()
    
    // Watch and validate changes
    return spectral.Watch(
        ctx,
        dag.Directory("./specs"),
        "**/*.yaml",
        nil,
        map[string]string{
            "failOnChanges": "true",
            "debounceMs": "1000",
        },
    )
} 