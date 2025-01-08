# Python Poetry Module for Dagger

A Dagger module for managing Python projects with Poetry. This module provides essential Poetry operations in a containerized environment.

## Features

- Install project dependencies
- Build Python packages
- Run tests
- Update dependencies
- Manage lock files

## Usage

Import the module in your Dagger pipeline:

```go
poetry := dag.PythonPoetry()
```

### Installing Dependencies

```go
installed, err := poetry.With(source).Install(ctx)
if err != nil {
    // Handle error
}
```

### Building Package

```go
built, err := poetry.With(source).Build(ctx)
if err != nil {
    // Handle error
}
```

### Running Tests

```go
output, err := poetry.With(source).Test(ctx)
if err != nil {
    // Handle error
}
fmt.Println("Test output:", output)
```

### Updating Dependencies

```go
updated, err := poetry.With(source).Update(ctx)
if err != nil {
    // Handle error
}
```

### Managing Lock File

```go
locked, err := poetry.With(source).Lock(ctx)
if err != nil {
    // Handle error
}
```

## Requirements

- Dagger v0.15.1
- Go 1.23.4
- Python project with Poetry configuration

## Environment

The module uses `python:3.12-slim` as the base image and automatically installs Poetry.

## Example

Here's a complete example of using the module:

```go
func BuildProject(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
    poetry := dag.PythonPoetry()

    // Install dependencies
    installed, err := poetry.With(source).Install(ctx)
    if err != nil {
        return nil, err
    }

    // Run tests
    if _, err := poetry.With(installed).Test(ctx); err != nil {
        return nil, err
    }

    // Build package
    return poetry.With(installed).Build(ctx)
}
```
