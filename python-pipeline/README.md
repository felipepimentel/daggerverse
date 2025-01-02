# Python Pipeline Module for Dagger

A complete pipeline for Python projects that combines Poetry for dependency management and PyPI for package publishing. This module provides an end-to-end workflow for Python package development and distribution.

## Features

- Complete build and publish pipeline
- Dependency management with Poetry
- Automated testing
- PyPI publishing
- Dependency updates

## Usage

Import the module in your Dagger pipeline:

```go
pipeline := dag.PythonPipeline()
```

### Building and Publishing

```go
err := pipeline.BuildAndPublish(ctx, source, token)
if err != nil {
    // Handle error
}
```

This will:

1. Install dependencies using Poetry
2. Run tests
3. Build the package
4. Publish to PyPI

### Updating Dependencies

```go
updated, err := pipeline.UpdateDependencies(ctx, source)
if err != nil {
    // Handle error
}
```

## Requirements

- Dagger v0.15.1
- Go 1.23.4
- Python project with Poetry configuration
- PyPI token for publishing

## Dependencies

This module integrates:

- `python-poetry` module for Poetry operations
- `python-pypi` module for PyPI publishing

## Example

Here's a complete example of using the pipeline:

```go
func DeployPackage(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
    pipeline := dag.PythonPipeline()

    // First, update dependencies
    updated, err := pipeline.UpdateDependencies(ctx, source)
    if err != nil {
        return err
    }

    // Then build and publish
    return pipeline.BuildAndPublish(ctx, updated, token)
}
```

## CLI Usage

```shell
# Update dependencies
dagger call update-dependencies --source=.

# Build and publish
export PYPI_TOKEN=your_token_here
dagger call build-and-publish --source=. --token=env:PYPI_TOKEN
```

## Environment

The pipeline uses:

- `python:3.11-slim` as base image
- Poetry for package management
- Automated PyPI authentication

## Functions

### CICD

Runs the complete CI/CD pipeline for a Python project. This includes:

1. Installing dependencies with Poetry
2. Running tests with pytest
3. Running linting with pylint (if configured)
4. Building the package

```bash
# Run the CI/CD pipeline
dagger call cicd --source .
```

### BuildAndPublish

Builds and publishes a Python package to PyPI. This includes:
