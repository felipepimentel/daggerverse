# Dagger Python Module

A powerful Dagger module for Python projects using Poetry, providing a streamlined workflow for building, testing, and publishing Python packages.

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
![Python Versions](https://img.shields.io/badge/Python-3.8%20%7C%203.9%20%7C%203.10%20%7C%203.11%20%7C%203.12-blue)
![Poetry](https://img.shields.io/badge/Poetry-1.7%2B-blue)

## Features

- 🚀 Automated Python package building with Poetry
- 🧪 Integrated testing with pytest and coverage reporting
- 📦 Secure PyPI publishing with token support
- 🔄 Caching for pip and Poetry dependencies
- 🐳 Containerized builds for consistent environments
- 🔧 Configurable Python versions and package paths

## Quick Start

```go
import (
    "context"
    "fmt"
)

func main() {
    ctx := context.Background()

    // Initialize Python module
    python := dag.Python().
        WithPythonVersion("3.12").
        WithPackagePath("my_package")

    // Run tests with coverage
    result, err := python.Test(ctx, dag.Host().Directory("."))
    if err != nil {
        panic(err)
    }
    fmt.Println(result)
}
```

## Installation

Add the Python module to your Dagger project:

```shell
dagger mod install github.com/daggerverse/python@latest
```

## Usage Examples

### Basic Testing

```go
// Run tests with default configuration
python := dag.Python()
result, err := python.Test(ctx, source)
```

### Building Package

```go
// Build Python package
container := python.Build(source)
```

### Publishing to PyPI

```go
// Configure PyPI deployment
pypiConfig := &PyPIConfig{
    Registry: "https://upload.pypi.org/legacy/",
    Token: dag.SetSecret("PYPI_TOKEN", "your-token"),
    SkipExisting: true,
}

// Initialize and publish
python := dag.Python().
    WithPyPIConfig(pypiConfig)
result, err := python.Publish(ctx, source)
```

## Configuration Options

### Python Module Configuration

| Option         | Method                               | Description               | Default |
| -------------- | ------------------------------------ | ------------------------- | ------- |
| Python Version | `WithPythonVersion(version string)`  | Set Python version        | "3.12"  |
| Package Path   | `WithPackagePath(path string)`       | Set package directory     | "."     |
| PyPI Config    | `WithPyPIConfig(config *PyPIConfig)` | Configure PyPI deployment | nil     |

### PyPI Configuration

| Field        | Type            | Description            | Default                           |
| ------------ | --------------- | ---------------------- | --------------------------------- |
| Registry     | string          | PyPI registry URL      | "https://upload.pypi.org/legacy/" |
| Token        | \*dagger.Secret | Authentication token   | nil                               |
| SkipExisting | bool            | Skip existing versions | false                             |
| AllowDirty   | bool            | Allow dirty versions   | false                             |

## Advanced Usage

### Custom Test Configuration

```go
python := dag.Python().
    WithPythonVersion("3.11").
    WithPackagePath("src/mypackage")

result, err := python.Test(ctx, source)
```

### Multi-stage Pipeline

```go
// Create a complete pipeline
python := dag.Python().
    WithPythonVersion("3.12").
    WithPackagePath("mypackage")

// 1. Run tests
if _, err := python.Test(ctx, source); err != nil {
    return err
}

// 2. Build package
container := python.Build(source)

// 3. Publish to PyPI
pypiConfig := &PyPIConfig{
    Token: dag.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN")),
    SkipExisting: true,
}

python = python.WithPyPIConfig(pypiConfig)
result, err := python.Publish(ctx, source)
```

## Development Environment

The module uses a containerized development environment with:

- Python (configurable version, default: 3.12)
- Poetry for dependency management
- pytest with coverage reporting
- Cached dependencies for faster builds

## Contributing

We welcome contributions! Here's how you can help:

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to your branch
5. Create a Pull Request

### Development Setup

1. Clone the repository:

```shell
git clone https://github.com/daggerverse/python.git
cd python
```

2. Install Dagger:

```shell
curl -L https://dl.dagger.io/dagger/install.sh | sh
```

3. Run tests:

```shell
dagger test
```

## Best Practices

- Always pin your Python version for reproducible builds
- Use Poetry's lock file for deterministic dependencies
- Store PyPI tokens securely using Dagger secrets
- Enable test coverage reporting for quality assurance
- Use caching to speed up builds

## Troubleshooting

### Common Issues

1. **Poetry Installation Fails**

   - Ensure your container has internet access
   - Check Python version compatibility

2. **PyPI Publishing Errors**

   - Verify token permissions
   - Check package version in pyproject.toml
   - Ensure package name is available on PyPI

3. **Test Coverage Issues**
   - Verify correct package path configuration
   - Check pytest-cov installation

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Dagger team for the amazing CI/CD framework
- Poetry team for the dependency management tool
- Python community for continuous support

## Support

- Create an issue for bug reports
- Start a discussion for feature requests
- Check existing issues before creating new ones

---

Built with ❤️ by the Dagger community
