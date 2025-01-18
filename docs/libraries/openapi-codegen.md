---
layout: default
title: OpenAPI Codegen Module
parent: Libraries
nav_order: 11
---

# OpenAPI Codegen Module

The OpenAPI Codegen module provides integration with [OpenAPI Generator](https://openapi-generator.tech/), a powerful tool for generating clients, servers, and documentation from OpenAPI (Swagger) definitions. This module allows you to automate code generation in your Dagger pipelines.

## Features

- Client code generation
- Server stub generation
- Documentation generation
- Multiple language support
- Template customization
- Configuration management
- Validation options
- Generator selection

## Installation

To use the OpenAPI Codegen module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/openapi-codegen"
)
```

## Usage Examples

### Basic Client Generation

```go
func (m *MyModule) Example(ctx context.Context) (*Directory, error) {
    codegen := dag.OpenAPICodegen().New()
    
    // Generate client code
    return codegen.GenerateClient(
        ctx,
        dag.File("./api.yaml"),  // spec file
        "go",                    // language
        map[string]string{
            "packageName": "client",
            "apiPackage": "api",
        },
    )
}
```

### Server Stub Generation

```go
func (m *MyModule) GenerateServer(ctx context.Context) (*Directory, error) {
    codegen := dag.OpenAPICodegen().New()
    
    // Generate server stubs
    return codegen.GenerateServer(
        ctx,
        dag.File("./api.yaml"),
        "python-flask",
        map[string]string{
            "packageName": "myapi",
            "serverPort": "8080",
        },
    )
}
```

### Documentation Generation

```go
func (m *MyModule) GenerateDocs(ctx context.Context) (*Directory, error) {
    codegen := dag.OpenAPICodegen().New()
    
    // Generate documentation
    return codegen.GenerateDocs(
        ctx,
        dag.File("./api.yaml"),
        "markdown",
        map[string]string{
            "outputFile": "API.md",
            "includeExamples": "true",
        },
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: API Code Generation
on: [push]

jobs:
  codegen:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Generate API Client
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/openapi-codegen
          args: |
            do -p '
              codegen := OpenAPICodegen().New()
              codegen.GenerateClient(
                ctx,
                dag.File("./api.yaml"),
                "typescript-axios",
                map[string]string{
                  "npmName": "@myorg/api-client",
                  "supportsES6": "true",
                },
              )
            '
```

## API Reference

### OpenAPICodegen

Main module struct that provides access to OpenAPI Generator functionality.

#### Constructor

- `New() *OpenAPICodegen`
  - Creates a new OpenAPICodegen instance
  - Default version: "latest"
  - Default platform: "linux/amd64"

#### Methods

- `GenerateClient(ctx context.Context, spec *File, language string, config map[string]string) (*Directory, error)`
  - Generates API client code
  - Parameters:
    - `spec`: OpenAPI specification file
    - `language`: Target language/framework
    - `config`: Generator configuration
  
- `GenerateServer(ctx context.Context, spec *File, language string, config map[string]string) (*Directory, error)`
  - Generates server stubs
  - Parameters:
    - `spec`: OpenAPI specification file
    - `language`: Target language/framework
    - `config`: Generator configuration
  
- `GenerateDocs(ctx context.Context, spec *File, format string, config map[string]string) (*Directory, error)`
  - Generates documentation
  - Parameters:
    - `spec`: OpenAPI specification file
    - `format`: Documentation format
    - `config`: Generator configuration

## Best Practices

1. **Code Generation**
   - Use consistent naming
   - Configure proper packages
   - Handle generated code

2. **Language Selection**
   - Choose appropriate generators
   - Configure language features
   - Test generated code

3. **Documentation**
   - Include examples
   - Generate multiple formats
   - Keep docs in sync

4. **Configuration**
   - Use version control
   - Document options
   - Validate settings

## Troubleshooting

Common issues and solutions:

1. **Generation Errors**
   ```
   Error: invalid specification
   Solution: Validate OpenAPI spec format
   ```

2. **Language Issues**
   ```
   Error: unsupported generator
   Solution: Check available generators and versions
   ```

3. **Configuration Problems**
   ```
   Error: invalid configuration option
   Solution: Verify generator-specific options
   ```

## Configuration Example

```yaml
# config.yaml
inputSpec: ./api.yaml
generatorName: go
output: ./generated
additionalProperties:
  packageName: myapi
  apiPackage: api
  modelPackage: models
  generateInterfaces: true
  enumClassPrefix: true
  structPrefix: true
```

## Advanced Usage

### Custom Templates

```go
func (m *MyModule) CustomTemplates(ctx context.Context) (*Directory, error) {
    codegen := dag.OpenAPICodegen().New()
    
    // Generate with custom templates
    return codegen.GenerateWithTemplates(
        ctx,
        dag.File("./api.yaml"),
        "java",
        dag.Directory("./templates"),
        map[string]string{
            "apiPackage": "com.example.api",
            "modelPackage": "com.example.model",
        },
    )
}
```

### Multi-Language Generation

```go
func (m *MyModule) MultiLanguage(ctx context.Context) error {
    codegen := dag.OpenAPICodegen().New()
    spec := dag.File("./api.yaml")
    
    // Generate for multiple languages
    languages := []string{"go", "typescript", "python"}
    for _, lang := range languages {
        _, err := codegen.GenerateClient(
            ctx,
            spec,
            lang,
            map[string]string{
                "packageName": fmt.Sprintf("api-%s", lang),
            },
        )
        if err != nil {
            return err
        }
    }
    
    return nil
} 