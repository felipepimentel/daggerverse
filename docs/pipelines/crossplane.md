# Crossplane Pipeline

This module provides a comprehensive pipeline for building, packaging, and managing Crossplane configurations and packages.

## Features

- Package building and initialization
- Custom package templates
- Package publishing to container registries
- Support for custom Crossplane containers
- Template-based package generation
- Registry authentication
- Flexible configuration options

## Installation

```bash
dagger mod use github.com/felipepimentel/daggerverse/pipelines/crossplane@latest
```

## Usage

### Basic Example

```go
// Initialize the module
crossplane := dag.Crossplane()

// Create a new package
output := crossplane.InitPackage(ctx, "my-configuration")
```

### Configuration Options

The module supports the following configuration:

```go
type Crossplane struct {
    // Crossplane container configuration
    XplaneContainer *dagger.Container
}
```

## Package Management

### Creating a New Package

```go
// Initialize a basic package
output := crossplane.InitPackage(ctx, "my-configuration")

// Initialize a custom package with templates
output := crossplane.InitCustomPackage(ctx, "MyResource")
```

### Building a Package

```go
// Build a package from source
output := crossplane.Package(ctx, dag.Host().Directory("./my-package"))
```

### Publishing a Package

```go
// Push package to a registry
status := crossplane.Push(ctx,
    dag.Host().Directory("./my-package"),
    "ghcr.io",
    "username",
    dag.SetSecret("registry_password", "your-password"),
    "ghcr.io/org/package:tag")
```

## Custom Package Templates

The module supports custom package generation with predefined templates. The following data structure is used:

```go
data := map[string]interface{}{
    "namespace":           "crossplane-system",
    "claimName":          "incluster",
    "apiGroup":           "resources.stuttgart-things.com",
    "claimApiVersion":    "v1alpha1",
    "maintainer":         "your.email@domain.com",
    "source":             "github.com/org/repo",
    "license":            "Apache-2.0",
    "crossplaneVersion":  ">=v1.14.1-0",
    "kindLower":          "resourcename",
    "kindLowerX":         "xresourcename",
    "kind":               "XResourceName",
    "plural":             "xresourcenames",
    "claimKind":          "ResourceName",
    "claimPlural":        "resourcenames",
    "compositeApiVersion": "apiextensions.crossplane.io/v1",
}
```

## GitHub Actions Integration

Create a workflow file `.github/workflows/crossplane.yml`:

```yaml
name: Crossplane Package
on:
  push:
    branches: [main]
    paths:
      - 'packages/**'

jobs:
  package:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install Dagger CLI
        uses: dagger/dagger-for-github@v5
        with:
          version: "0.15.3"
      
      - name: Build and Push Package
        env:
          REGISTRY_USERNAME: ${{ secrets.REGISTRY_USERNAME }}
          REGISTRY_PASSWORD: ${{ secrets.REGISTRY_PASSWORD }}
        run: |
          dagger call --progress=plain \
            push \
            --source ./packages/my-package \
            --registry ghcr.io \
            --username "$REGISTRY_USERNAME" \
            --password "$REGISTRY_PASSWORD" \
            --destination "ghcr.io/org/package:latest"
```

## Examples

### Custom Container

```go
container := dag.Container().
    From("ghcr.io/stuttgart-things/crossplane-cli:v1.18.0")

crossplane := dag.Crossplane().
    WithXplaneContainer(container)
```

### Full Package Lifecycle

```go
// Initialize module
crossplane := dag.Crossplane()

// Create package
pkg := crossplane.InitCustomPackage(ctx, "MyResource")

// Build package
built := crossplane.Package(ctx, pkg)

// Push to registry
status := crossplane.Push(ctx,
    built,
    "ghcr.io",
    "username",
    dag.SetSecret("registry_password", "your-password"),
    "ghcr.io/org/myresource:latest")
```

## Best Practices

1. **Package Structure**:
   - Use consistent naming conventions
   - Follow Crossplane package guidelines
   - Include comprehensive documentation

2. **Version Control**:
   - Tag releases appropriately
   - Use semantic versioning
   - Document breaking changes

3. **Security**:
   - Use secrets for credentials
   - Implement proper RBAC
   - Follow least privilege principle

## Common Issues

1. **Build Failures**:
   - Check Crossplane CLI version
   - Verify package structure
   - Validate template syntax

2. **Push Errors**:
   - Verify registry credentials
   - Check network connectivity
   - Validate image naming

3. **Template Issues**:
   - Validate template syntax
   - Check variable substitution
   - Verify file paths

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](../CONTRIBUTING.md) for details on how to submit pull requests.

## License

This module is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details. 