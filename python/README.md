# Python Module for Dagger

A powerful Dagger module that streamlines Python development workflows using Poetry. This module automates building, testing, and publishing Python packages with best practices and security in mind.

## Features

- 🏗️ **Automated Package Building**

  - Poetry-based dependency management
  - Proper virtual environment handling
  - Configurable Python versions
  - Cache optimization for faster builds

- 🧪 **Advanced Testing**

  - Configurable test execution
  - Coverage reporting (XML, HTML, and Terminal)
  - Test result formatting options
  - Parallel test execution support

- 📦 **Secure Publishing**
  - PyPI and TestPyPI support
  - Secure token handling
  - Registry selection
  - Version management

## Quick Start

```go
// Initialize the Python module
python := dag.Python()

// Build your package
container := python.Build(ctx, source)

// Run tests with coverage
tested := python.Test(ctx, source)

// Publish to PyPI (with secure token handling)
published := python.WithPyPIToken(token).Publish(ctx, source)
```

## Configuration Options

### Python Version

```go
// Use a specific Python version
python := dag.Python().WithPythonVersion("3.11")
```

### Package Path

```go
// Specify custom package path
python := dag.Python().WithPackagePath("my_package")
```

### Testing Options

```go
// Configure test execution
tested := python.Test(ctx, source,
    WithCoverage(true),
    WithCoverageReport("xml"),
    WithVerbose(true))
```

### Publishing Options

```go
// Configure publishing
published := python.WithPyPIConfig(PyPIConfig{
    Registry: "https://upload.pypi.org/legacy/",
    SkipExisting: true,
    AllowDirty: false,
}).Publish(ctx, source)
```

## Security Best Practices

### Token Handling

Never hardcode PyPI tokens in your code. Instead, use environment variables or secrets management:

```go
token := dag.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN"))
python := dag.Python().WithPyPIToken(token)
```

## Examples

### Complete Build and Publish Workflow

```go
func BuildAndPublish(ctx context.Context, source *Directory) (*Container, error) {
    python := dag.Python().
        WithPythonVersion("3.11").
        WithPackagePath("my_package")

    // Build and test
    tested := python.Test(ctx, source)
    if tested.Error != nil {
        return nil, tested.Error
    }

    // Publish if tests pass
    token := dag.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN"))
    published := python.
        WithPyPIToken(token).
        Publish(ctx, source)

    return published, nil
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the Apache-2.0 License - see the LICENSE file for details.
