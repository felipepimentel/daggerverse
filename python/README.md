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

## CLI Usage

### Secure PyPI Token Handling

You can securely pass your PyPI token when using the module via command line in several ways:

1. **Using environment variable**:

```bash
# Set the token as environment variable
export PYPI_TOKEN="your_token_here"

# Use the token from environment
dagger call publish --source . --token env:PYPI_TOKEN
```

2. **Using a secret file** (recommended for local development):

```bash
# Store token in a file (make sure to add it to .gitignore)
echo "your_token_here" > .pypi_token

# Use the token from file
dagger call publish --source . --token file:.pypi_token
```

3. **Using Dagger's secret management** (recommended for CI/CD):

```bash
# Store the secret in Dagger
dagger secret create pypi-token "your_token_here"

# Use the stored secret
dagger call publish --source . --token secret:pypi-token
```

⚠️ **Security Best Practices**:

- Never commit tokens to version control
- Add `.pypi_token` to your `.gitignore`
- Use environment variables in CI/CD pipelines
- Rotate tokens periodically
- Use PyPI's trusted publishing when possible (see [PyPI's documentation](https://docs.pypi.org/trusted-publishers/))

## Publishing

This module supports two methods of publishing to the Daggerverse:

### Automatic Publishing

The module will be automatically published to the Daggerverse whenever someone uses it via `dagger call`. This is the recommended approach as it ensures your module is always available when needed.

### Manual Publishing

If you need to manually publish:

```bash
# 1. Push your changes
git add .
git commit -m "feat: your changes"
git push origin main

# 2. Tag with semantic version
git tag v1.0.0  # For major version 1
# or for module in monorepo
git tag python/v1.0.0

# 3. Push the tag
git push origin v1.0.0
```

Then visit the [Daggerverse](https://daggerverse.dev) and click "Publish".

### Versioning

This module follows semantic versioning (`vMAJOR.MINOR.PATCH`) with automated release management using semantic-release.

#### Commit Convention

For monorepo management, use scoped conventional commits:

```bash
# Format
type(scope): description

# Examples for Python module
feat(python): add new testing feature
fix(python): resolve token handling issue
docs(python): update installation guide
perf(python): improve build cache
```

Commit types that trigger version updates (when scoped to `python`):

- `feat(python): new feature` -> MINOR version bump
- `fix(python): bug fix` -> PATCH version bump
- `feat(python)!:` or `fix(python)!:` -> MAJOR version bump
- `chore(python):`, `docs(python):`, `style(python):`, etc -> No version bump

Commits without the `python` scope or with different scopes won't trigger version updates for this module.

Examples:

```bash
# Patch release (python/v1.0.0 -> python/v1.0.1)
git commit -m "fix(python): correct PyPI token handling"

# Minor release (python/v1.0.1 -> python/v1.1.0)
git commit -m "feat(python): add support for custom test commands"

# Major release (python/v1.1.0 -> python/v2.0.0)
git commit -m "feat(python)!: change publish API interface"

# No release (commit to another module)
git commit -m "feat(nodejs): add new feature"
```

The version will be automatically bumped and tagged (with prefix `python/`) when merging to main branch.

Examples:

- `v1.0.0` - Initial stable release
- `v1.1.0` - New features added
- `v1.1.1` - Bug fixes
- `python/v1.0.0` - When module is in a monorepo

To use a specific version:

```bash
dagger mod install github.com/daggerverse/python@v1.0.0
```

---

Built with ❤️ by the Dagger community
