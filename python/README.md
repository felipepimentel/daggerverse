# Python Module for Dagger

This module provides a complete CI/CD pipeline for Python projects using Poetry. It includes functionality for building, testing, versioning, and publishing Python packages.

## Features

- Automatic detection of `pyproject.toml` location
- Poetry-based dependency management
- Test execution with pytest
- Semantic versioning using semantic-release
- PyPI package publishing

## Usage

```go
import (
    "context"
    "fmt"
    "dagger.io/dagger"
)

func main() {
    ctx := context.Background()

    // Initialize Dagger client
    client, err := dagger.Connect(ctx)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Get source code directory
    source := client.Host().Directory(".")

    // Create Python module instance
    python := &Python{
        PackagePath: ".", // Path to your Python package
    }

    // Get PyPI token from environment
    token := client.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN"))

    // Run complete CI/CD pipeline
    version, err := python.CICD(ctx, source, token)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Successfully released version %s\n", version)
}
```

## Configuration

The module can be configured through the following fields:

- `PackagePath`: Path to your Python package within the source directory (default: ".")

## Environment Variables

The following environment variables are required:

- `PYPI_TOKEN`: PyPI authentication token for publishing packages
- `GITHUB_TOKEN`: GitHub token for semantic-release (when running in CI)
