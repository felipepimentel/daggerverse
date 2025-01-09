---
layout: default
title: Poetry Module
parent: Libraries
nav_order: 12
---

# Poetry Module

The Poetry module provides integration with [Poetry](https://python-poetry.org/), a modern dependency management and packaging tool for Python. This module allows you to manage Python dependencies and build packages in your Dagger pipelines.

## Features

- Dependency management
- Package building
- Virtual environment handling
- Lock file management
- Project initialization
- Package publishing
- Development dependencies
- Build isolation

## Installation

To use the Poetry module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/poetry"
)
```

## Usage Examples

### Basic Package Installation

```go
func (m *MyModule) Example(ctx context.Context) (*Container, error) {
    poetry := dag.Poetry().New()
    
    // Install dependencies
    return poetry.Install(
        ctx,
        dag.Directory("./python-project"),  // project directory
        false,                             // no dev dependencies
        nil,                              // default Python version
    )
}
```

### Package Building

```go
func (m *MyModule) BuildPackage(ctx context.Context) (*File, error) {
    poetry := dag.Poetry().New()
    
    // Build package
    return poetry.Build(
        ctx,
        dag.Directory("./python-project"),
        "wheel",  // format
        map[string]string{
            "python-version": "3.9",
        },
    )
}
```

### Development Environment

```go
func (m *MyModule) DevEnv(ctx context.Context) (*Container, error) {
    poetry := dag.Poetry().New()
    
    // Setup development environment
    return poetry.Install(
        ctx,
        dag.Directory("./python-project"),
        true,   // include dev dependencies
        "3.9",  // Python version
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Python Package
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build Package
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/poetry
          args: |
            do -p '
              poetry := Poetry().New()
              poetry.Build(
                ctx,
                dag.Directory("./python-project"),
                "wheel",
                map[string]string{
                  "python-version": "3.9",
                },
              )
            '
```

## API Reference

### Poetry

Main module struct that provides access to Poetry functionality.

#### Constructor

- `New() *Poetry`
  - Creates a new Poetry instance
  - Default version: "latest"
  - Default platform: "linux/amd64"

#### Methods

- `Install(ctx context.Context, src *Directory, dev bool, pythonVersion string) (*Container, error)`
  - Installs project dependencies
  - Parameters:
    - `src`: Project directory
    - `dev`: Include development dependencies
    - `pythonVersion`: Python version
  
- `Build(ctx context.Context, src *Directory, format string, config map[string]string) (*File, error)`
  - Builds Python package
  - Parameters:
    - `src`: Project directory
    - `format`: Build format (wheel, sdist)
    - `config`: Build configuration
  
- `Publish(ctx context.Context, src *Directory, repository string, token *Secret) error`
  - Publishes package to repository
  - Parameters:
    - `src`: Project directory
    - `repository`: Target repository
    - `token`: Authentication token

## Best Practices

1. **Dependency Management**
   - Use lock files
   - Pin versions
   - Separate dev dependencies

2. **Build Configuration**
   - Use build isolation
   - Configure Python versions
   - Manage build dependencies

3. **Package Publishing**
   - Use secure tokens
   - Verify package contents
   - Test before publishing

4. **Environment Management**
   - Use virtual environments
   - Handle Python versions
   - Manage system dependencies

## Troubleshooting

Common issues and solutions:

1. **Installation Issues**
   ```
   Error: dependency resolution failed
   Solution: Check dependency conflicts in pyproject.toml
   ```

2. **Build Problems**
   ```
   Error: build backend failed
   Solution: Verify build dependencies and configuration
   ```

3. **Publishing Errors**
   ```
   Error: authentication failed
   Solution: Check repository credentials and token
   ```

## Configuration Example

```toml
# pyproject.toml
[tool.poetry]
name = "my-package"
version = "0.1.0"
description = "My Python package"
authors = ["Author <author@example.com>"]

[tool.poetry.dependencies]
python = "^3.9"
requests = "^2.28.0"

[tool.poetry.dev-dependencies]
pytest = "^7.1.0"
black = "^22.3.0"

[build-system]
requires = ["poetry-core>=1.0.0"]
build-backend = "poetry.core.masonry.api"
```

## Advanced Usage

### Custom Build Configuration

```go
func (m *MyModule) CustomBuild(ctx context.Context) (*File, error) {
    poetry := dag.Poetry().New()
    
    // Build with custom configuration
    return poetry.BuildWithConfig(
        ctx,
        dag.Directory("./python-project"),
        dag.File("./build.toml"),
        map[string]string{
            "python-version": "3.9",
            "optimize": "2",
        },
    )
}
```

### Multi-Package Management

```go
func (m *MyModule) ManagePackages(ctx context.Context) error {
    poetry := dag.Poetry().New()
    
    // Manage multiple packages
    packages := []string{"package1", "package2", "package3"}
    for _, pkg := range packages {
        _, err := poetry.Install(
            ctx,
            dag.Directory(fmt.Sprintf("./%s", pkg)),
            false,
            "3.9",
        )
        if err != nil {
            return err
        }
    }
    
    return nil
} 