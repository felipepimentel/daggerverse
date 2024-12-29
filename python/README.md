# Python Module for Dagger

A Dagger module for Python projects that provides CI/CD capabilities with Poetry support.

## Features

- 🔄 Automatic version management with semantic-release
- 🔍 Automatic project structure detection (finds `pyproject.toml` automatically)
- 📦 Poetry-based dependency management
- 🧪 Test execution with pytest
- 🎨 Code quality tools (linting and formatting)
- 📚 Documentation generation
- 🚀 PyPI publishing
- ⚙️ Flexible configuration options

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
- `CICD`: Runs the complete pipeline:
  1. Bumps version using semantic-release
  2. Runs tests with coverage
  3. Performs linting and formatting
  4. Builds the package
  5. Updates version in pyproject.toml
  6. Publishes to PyPI

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
- `GITHUB_TOKEN`: GitHub token (required for semantic-release)
  - Used for version management and release creation

## Examples

### Running Tests with Coverage

```bash
dagger -m github.com/felipepimentel/daggerverse/python call test --source .
```

### Publishing to PyPI

```bash
dagger -m github.com/felipepimentel/daggerverse/python call cd --source . --token env:PYPI_TOKEN
```

### Complete CI/CD Pipeline

```bash
dagger -m github.com/felipepimentel/daggerverse/python call cicd --source . --token env:PYPI_TOKEN
```

## GitHub Actions Integration

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main]

jobs:
  pipeline:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Run Pipeline
        uses: dagger/dagger-for-github@v7
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
        with:
          verb: call
          module: github.com/felipepimentel/daggerverse/python
          args: cicd --source . --token env:PYPI_TOKEN
```

## Version Management

The module uses semantic-release to automatically manage versions based on commit messages. The version is determined by analyzing commit messages since the last release:

- `feat:` commits trigger a minor version bump
- `fix:` commits trigger a patch version bump
- `perf:` commits trigger a patch version bump
- Breaking changes (indicated by `!` or `BREAKING CHANGE:`) trigger a major version bump

The version is automatically updated in:

- Git tags
- GitHub releases
- pyproject.toml
- Release notes/changelog
