---
layout: default
title: Poetry Module
parent: Libraries
nav_order: 12
---

# Poetry Library

This module provides comprehensive integration with Poetry for Python project management, including dependency management, building, testing, and package publishing.

## Features

- Poetry installation and configuration
- Dependency management
- Package building with version control
- Test execution
- Lock file management
- Dependency updates
- Custom Python base image support
- Virtual environment configuration

## Installation

```bash
dagger mod use github.com/felipepimentel/daggerverse/libraries/poetry@latest
```

## Usage

### Basic Example

```go
// Initialize the module
poetry := dag.Poetry().WithBaseImage("python:3.12-alpine")

// Install dependencies
output := poetry.Install(dag.Host().Directory("."))
```

### Configuration Options

```go
type Poetry struct {
    // Base image for Poetry operations
    BaseImage string // default: "python:3.12-alpine"
}
```

## Package Management

### Installing Dependencies

```go
// Install project dependencies
output := poetry.Install(dag.Host().Directory("."))
```

### Building Packages

```go
// Build package
dist := poetry.Build(dag.Host().Directory("."))

// Build with specific version
dist := poetry.BuildWithVersion(dag.Host().Directory("."), "1.0.0")
```

### Testing

```go
// Run tests
output, err := poetry.Test(ctx, dag.Host().Directory("."))
```

### Dependency Management

```go
// Update lock file
output := poetry.Lock(dag.Host().Directory("."))

// Update dependencies
output := poetry.Update(dag.Host().Directory("."))
```

## GitHub Actions Integration

Create a workflow file `.github/workflows/poetry.yml`:

```yaml
name: Python Package
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install Dagger CLI
        uses: dagger/dagger-for-github@v5
        with:
          version: "0.9.3"
      
      - name: Build and Test
        run: |
          dagger call --progress=plain \
            test \
            --source .

      - name: Build Package
        run: |
          dagger call --progress=plain \
            build \
            --source . \
            --version "1.0.0"
```

## Examples

### Custom Base Image

```go
poetry := dag.Poetry().
    WithBaseImage("python:3.11-slim")
```

### Complete Build Pipeline

```go
// Initialize with custom image
poetry := dag.Poetry().
    WithBaseImage("python:3.12-alpine")

// Install dependencies
installed := poetry.Install(dag.Host().Directory("."))

// Run tests
output, err := poetry.Test(ctx, installed)
if err != nil {
    return err
}

// Build with version
dist := poetry.BuildWithVersion(installed, "1.0.0")
```

### Development Workflow

```go
// Initialize Poetry
poetry := dag.Poetry()

// Update dependencies
updated := poetry.Update(dag.Host().Directory("."))

// Lock dependencies
locked := poetry.Lock(updated)

// Install and test
output, err := poetry.Test(ctx, locked)
```

## Best Practices

1. **Project Structure**:
   - Keep `pyproject.toml` in root directory
   - Use consistent dependency versions
   - Include comprehensive test suite

2. **Version Management**:
   - Use semantic versioning
   - Keep dependencies up to date
   - Lock dependencies for production

3. **Testing**:
   - Write comprehensive tests
   - Use pytest fixtures
   - Include coverage reports

## Common Issues

1. **Installation Problems**:
   - Check Python version compatibility
   - Verify dependency conflicts
   - Validate pyproject.toml syntax

2. **Build Issues**:
   - Check build dependencies
   - Verify package structure
   - Validate version format

3. **Test Failures**:
   - Check test dependencies
   - Verify test environment
   - Debug test output

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](../CONTRIBUTING.md) for details on how to submit pull requests.

## License

This module is licensed under the Apache License 2.0. See the [LICENSE](../LICENSE) file for details. 