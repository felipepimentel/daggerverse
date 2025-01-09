---
layout: default
title: PyPI Module
parent: Libraries
nav_order: 14
---

# PyPI Module

The PyPI module provides integration with [PyPI](https://pypi.org/) (Python Package Index), the official repository for Python packages. This module allows you to publish and manage Python packages in your Dagger pipelines.

## Features

- Package publishing
- Package downloading
- Version management
- Authentication handling
- Package verification
- Distribution formats
- Repository selection
- Metadata management

## Installation

To use the PyPI module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/pypi"
)
```

## Usage Examples

### Basic Package Publishing

```go
func (m *MyModule) Example(ctx context.Context) error {
    pypi := dag.PyPI().New()
    
    // Publish package to PyPI
    return pypi.Publish(
        ctx,
        dag.Directory("./dist"),  // distribution directory
        dag.SetSecret("PYPI_TOKEN", "your-token"),
        "pypi",                  // repository (or "testpypi")
    )
}
```

### Package Download

```go
func (m *MyModule) DownloadPackage(ctx context.Context) (*Directory, error) {
    pypi := dag.PyPI().New()
    
    // Download package
    return pypi.Download(
        ctx,
        "requests",    // package name
        "2.28.0",     // version
        "wheel",      // format
    )
}
```

### Package Verification

```go
func (m *MyModule) VerifyPackage(ctx context.Context) error {
    pypi := dag.PyPI().New()
    
    // Verify package integrity
    return pypi.Verify(
        ctx,
        dag.File("./dist/package-1.0.0-py3-none-any.whl"),
        map[string]string{
            "checkHashAlgorithm": "sha256",
            "expectedHash": "abc123...",
        },
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: PyPI Release
on: [push]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Publish to PyPI
        uses: dagger/dagger-action@v1
        env:
          PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/pypi
          args: |
            do -p '
              pypi := PyPI().New()
              pypi.Publish(
                ctx,
                dag.Directory("./dist"),
                dag.SetSecret("PYPI_TOKEN", PYPI_TOKEN),
                "pypi",
              )
            '
```

## API Reference

### PyPI

Main module struct that provides access to PyPI functionality.

#### Constructor

- `New() *PyPI`
  - Creates a new PyPI instance
  - Default version: "latest"
  - Default platform: "linux/amd64"

#### Methods

- `Publish(ctx context.Context, dist *Directory, token *Secret, repository string) error`
  - Publishes package to PyPI
  - Parameters:
    - `dist`: Distribution directory
    - `token`: PyPI API token
    - `repository`: Target repository
  
- `Download(ctx context.Context, package string, version string, format string) (*Directory, error)`
  - Downloads package from PyPI
  - Parameters:
    - `package`: Package name
    - `version`: Package version
    - `format`: Distribution format
  
- `Verify(ctx context.Context, package *File, config map[string]string) error`
  - Verifies package integrity
  - Parameters:
    - `package`: Package file
    - `config`: Verification configuration

## Best Practices

1. **Package Management**
   - Use semantic versioning
   - Include proper metadata
   - Validate distributions

2. **Security**
   - Use API tokens
   - Verify package hashes
   - Secure credentials

3. **Distribution**
   - Support multiple formats
   - Include documentation
   - Test installations

4. **Repository Selection**
   - Use TestPyPI for testing
   - Verify package visibility
   - Manage permissions

## Troubleshooting

Common issues and solutions:

1. **Upload Failures**
   ```
   Error: invalid distribution format
   Solution: Check package build configuration
   ```

2. **Authentication Issues**
   ```
   Error: unauthorized access
   Solution: Verify PyPI token and permissions
   ```

3. **Version Conflicts**
   ```
   Error: version already exists
   Solution: Update package version number
   ```

## Configuration Example

```toml
# setup.cfg
[metadata]
name = my-package
version = 1.0.0
author = Author Name
author_email = author@example.com
description = Package description
long_description = file: README.md
long_description_content_type = text/markdown
url = https://github.com/username/my-package
classifiers =
    Programming Language :: Python :: 3
    License :: OSI Approved :: MIT License
    Operating System :: OS Independent

[options]
package_dir =
    = src
packages = find:
python_requires = >=3.6

[options.packages.find]
where = src
```

## Advanced Usage

### Custom Repository Configuration

```go
func (m *MyModule) CustomRepo(ctx context.Context) error {
    pypi := dag.PyPI().New()
    
    // Publish to custom repository
    return pypi.PublishToRepository(
        ctx,
        dag.Directory("./dist"),
        dag.SetSecret("REPO_TOKEN", "token"),
        "https://custom.repo.com/simple",
        map[string]string{
            "verify_ssl": "true",
            "repository_name": "custom",
        },
    )
}
```

### Multi-Format Publishing

```go
func (m *MyModule) MultiFormat(ctx context.Context) error {
    pypi := dag.PyPI().New()
    token := dag.SetSecret("PYPI_TOKEN", "your-token")
    
    // Publish multiple formats
    formats := []string{"wheel", "sdist"}
    for _, format := range formats {
        dist := fmt.Sprintf("./dist/%s", format)
        err := pypi.Publish(
            ctx,
            dag.Directory(dist),
            token,
            "pypi",
        )
        if err != nil {
            return err
        }
    }
    
    return nil
}
``` 