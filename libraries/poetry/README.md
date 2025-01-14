# Python Poetry Module for Dagger

A Dagger module for managing Python projects with Poetry. This module provides essential Poetry operations in a containerized environment.

## Features

- Install project dependencies
- Build Python packages
- Run tests
- Update dependencies
- Manage lock files
- Custom base image support

## Usage

Import the module in your Dagger pipeline:

```go
poetry := dag.Poetry()
```

### Installing Dependencies

```go
// Install dependencies
output := poetry.Install(dag.Host().Directory("."))
```

### Building Package

```go
// Build package
dist := poetry.Build(dag.Host().Directory("."))

// Build with specific version
dist := poetry.BuildWithVersion(dag.Host().Directory("."), "1.0.0")
```

### Running Tests

```go
// Run tests
output, err := poetry.Test(ctx, dag.Host().Directory("."))
if err != nil {
    // Handle error
}
fmt.Println("Test output:", output)
```

### Updating Dependencies

```go
// Update dependencies
updated := poetry.Update(dag.Host().Directory("."))

// Update lock file
locked := poetry.Lock(dag.Host().Directory("."))
```

## Requirements

- Dagger v0.15.1
- Go 1.23.4
- Python project with Poetry configuration

## Environment

The module uses `python:3.12-alpine` as the base image and automatically installs Poetry.

## Example

Here's a complete example of using the module:

```go
func BuildProject(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
    poetry := dag.Poetry()

    // Install dependencies
    installed := poetry.Install(source)

    // Run tests
    if _, err := poetry.Test(ctx, installed); err != nil {
        return nil, err
    }

    // Build package with version
    return poetry.BuildWithVersion(installed, "1.0.0"), nil
}
```

## License

This module is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
