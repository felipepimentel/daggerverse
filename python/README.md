# Python Module for Dagger

A Dagger module for Python projects that provides CI/CD capabilities with Poetry support.

## Features

- Automatic project structure detection (finds `pyproject.toml` automatically)
- Poetry-based dependency management
- Test execution with pytest
- Code quality tools (linting and formatting)
- Documentation generation
- PyPI publishing
- Flexible configuration options

## Usage

### Basic Example

```go
// Initialize client
client, err := dagger.Connect(context.Background())
if err != nil {
    panic(err)
}
defer client.Close()

// Create Python module
python := dag.Python()

// Run CI/CD pipeline
output, err := python.CICD(ctx, dag.Host().Directory("."), dag.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN")))
if err != nil {
    panic(err)
}
fmt.Println(output)
```

### Pipeline Functions

- `CI`: Runs tests and builds the package
- `CD`: Publishes the package to PyPI
- `CICD`: Runs the complete pipeline (CI + CD)

### Individual Functions

- `Test`: Runs pytest with configurable options
- `Build`: Creates a Poetry environment and builds the package
- `Publish`: Publishes the package to PyPI
- `Lint`: Runs code linting with ruff
- `Format`: Formats code with black
- `BuildDocs`: Generates documentation with Sphinx or MkDocs

## Configuration

### Test Configuration

```go
python := dag.Python().WithTestConfig(&TestConfig{
    Verbose: true,
    Workers: 4,
    Coverage: &CoverageConfig{
        Enabled: true,
        Formats: []string{"xml", "html"},
        MinCoverage: 80,
    },
})
```

### Build Configuration

```go
python := dag.Python().WithBuildConfig(&BuildConfig{
    DependencyGroups: []string{"dev", "test"},
    Cache: &CacheConfig{
        PipCache: true,
        PoetryCache: true,
    },
})
```

### PyPI Configuration

```go
python := dag.Python().WithPyPIConfig(&PyPIConfig{
    Registry: "https://upload.pypi.org/legacy/",
    SkipExisting: true,
})
```

## Environment Variables

- `PYPI_TOKEN`: PyPI authentication token (required for publishing)
  - Can be provided via environment variable
  - Can be stored in `.env` file
  - Can be passed directly as a secret

## Examples

### Running Tests with Coverage

```go
output, err := python.Test(ctx, dag.Host().Directory("."))
```

### Publishing to PyPI

```go
token := dag.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN"))
output, err := python.Publish(ctx, dag.Host().Directory("."), token)
```

### Complete CI/CD Pipeline

```go
token := dag.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN"))
output, err := python.CICD(ctx, dag.Host().Directory("."), token)
```

## GitHub Actions Integration

```yaml
name: CI/CD

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  pipeline:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dagger/dagger-action@v1
        with:
          module: github.com/dagger/daggerverse/python@main
        env:
          PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
        run: |
          # Run only CI for pull requests
          if [ "${{ github.event_name }}" = "pull_request" ]; then
            dagger call ci --source .
          else
            # Run complete CI/CD pipeline for pushes to main
            dagger call cicd --source . --token env:PYPI_TOKEN
          fi
```

## Command Line Usage

```bash
# Run tests
dagger -m github.com/felipepimentel/daggerverse/python call test --source .

# Run linting
dagger -m github.com/felipepimentel/daggerverse/python call lint --source .

# Run formatting
dagger -m github.com/felipepimentel/daggerverse/python call format --source .

# Run CI pipeline
dagger -m github.com/felipepimentel/daggerverse/python call ci --source .

# Run CD pipeline
dagger -m github.com/felipepimentel/daggerverse/python call cd --source . --token env:PYPI_TOKEN

# Run complete CI/CD pipeline
dagger -m github.com/felipepimentel/daggerverse/python call cicd --source . --token env:PYPI_TOKEN
```
