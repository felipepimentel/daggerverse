# Python Module for Dagger

A Dagger module for automating Python development workflows using Poetry. This module provides functions for building packages, running tests, and publishing to PyPI with proper dependency management and caching.

## Usage

### Import the module in your Dagger pipeline

```python
import dagger

async def build():
    # Connect to Dagger
    async with dagger.Connection() as client:
        # Get reference to the Python module
        python = client.container().from_("registry.dagger.io/engine")

        # Import the Python module from Daggerverse
        mod = await client.host().module("github.com/felipepimentel/daggerverse/python")

        # Use the module with your source code
        source = client.host().directory(".")

        # Run tests
        await mod.test(source)

        # Build package
        await mod.build(source)

        # Publish to PyPI (if you have a token)
        token = client.set_secret("PYPI_TOKEN", os.getenv("PYPI_TOKEN"))
        await mod.publish(source, token)

if __name__ == "__main__":
    asyncio.run(build())
```

### Available Functions

#### `test(source: Directory) -> str`

Run tests for your Python package.

```python
# Basic usage
await mod.test(source)

# With custom configuration
await mod.test(
    source,
    verbosity=2,
    parallel_workers=4,
    enable_coverage=True,
    coverage_threshold=80
)
```

#### `build(source: Directory) -> Container`

Build your Python package using Poetry.

```python
# Basic usage
await mod.build(source)

# With custom configuration
await mod.build(
    source,
    python_version="3.11",
    extra_dependencies=["pytest", "pytest-cov"],
    cache_dependencies=True
)
```

#### `publish(source: Directory, token: Secret) -> str`

Publish your package to PyPI.

```python
# Basic usage
token = client.set_secret("PYPI_TOKEN", os.getenv("PYPI_TOKEN"))
await mod.publish(source, token)

# With custom configuration
await mod.publish(
    source,
    token,
    repository="https://test.pypi.org/legacy/",
    skip_existing=True
)
```

### Using with GitHub Actions

```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Dagger
        run: pip install dagger-io

      - name: Test and Build
        run: |
          python3 << 'EOF'
          import dagger
          import asyncio
          import os

          async def pipeline():
              async with dagger.Connection() as client:
                  mod = await client.host().module("github.com/felipepimentel/daggerverse/python")
                  source = client.host().directory(".")
                  
                  # Run tests
                  print("Running tests...")
                  await mod.test(source)
                  
                  # Build package
                  print("Building package...")
                  await mod.build(source)
                  
                  if os.getenv("PYPI_TOKEN"):
                      print("Publishing to PyPI...")
                      token = client.set_secret("PYPI_TOKEN", os.getenv("PYPI_TOKEN"))
                      await mod.publish(source, token)

          asyncio.run(pipeline())
          EOF
        env:
          PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
```

### Command Line Usage

You can also use the module directly from the command line:

```bash
# Run tests
dagger call -m github.com/felipepimentel/daggerverse/python test --source .

# Build package
dagger call -m github.com/felipepimentel/daggerverse/python build --source .

# Publish to PyPI
dagger call -m github.com/felipepimentel/daggerverse/python publish \
  --source . \
  --token $PYPI_TOKEN
```

## Configuration Options

### Test Configuration

- `verbosity`: Test output verbosity level (default: 2)
- `parallel_workers`: Number of parallel test workers (default: 1)
- `enable_coverage`: Enable coverage reporting (default: true)
- `coverage_threshold`: Minimum coverage percentage (default: 80)
- `coverage_formats`: Coverage report formats (default: ["xml", "html"])

### Build Configuration

- `python_version`: Python version to use (default: "3.11")
- `poetry_version`: Poetry version to use (default: "1.7.1")
- `extra_dependencies`: Additional dependencies to install
- `cache_dependencies`: Enable dependency caching (default: true)

### Publish Configuration

- `repository`: PyPI repository URL (default: "pypi")
- `skip_existing`: Skip if version exists (default: false)
- `verify_ssl`: Verify SSL certificates (default: true)

## Examples

### Basic Package Development

```python
async with dagger.Connection() as client:
    mod = await client.host().module("github.com/felipepimentel/daggerverse/python")
    source = client.host().directory(".")

    # Run tests with coverage
    await mod.test(
        source,
        enable_coverage=True,
        coverage_threshold=90
    )

    # Build with caching
    await mod.build(
        source,
        cache_dependencies=True
    )
```

### Publishing to Test PyPI First

```python
async with dagger.Connection() as client:
    mod = await client.host().module("github.com/felipepimentel/daggerverse/python")
    source = client.host().directory(".")
    token = client.set_secret("PYPI_TOKEN", os.getenv("PYPI_TOKEN"))

    # Publish to Test PyPI
    await mod.publish(
        source,
        token,
        repository="https://test.pypi.org/legacy/"
    )

    # If successful, publish to production PyPI
    await mod.publish(source, token)
```

## Best Practices

- Always use secrets for sensitive data like PyPI tokens
- Enable caching for faster builds in CI
- Set appropriate coverage thresholds for your project
- Use Test PyPI for testing package publishing
