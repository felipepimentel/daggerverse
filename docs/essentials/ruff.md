---
layout: default
title: Ruff Module
parent: Essentials
nav_order: 13
---

# Ruff Module

The Ruff module provides integration with Ruff, an extremely fast Python linter written in Rust. This module allows you to perform code quality checks and formatting in your Dagger pipelines.

## Features

- Python code linting
- Code formatting
- Fast execution
- Configuration management
- Error reporting
- Fix suggestions
- Directory scanning
- Rule customization
- Output formatting
- Cache management

## Installation

To use the Ruff module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/ruff"
)
```

## Usage Examples

### Basic Linting

```go
func (m *MyModule) Example(ctx context.Context) error {
    ruff := dag.Ruff()
    
    // Lint Python files
    return ruff.Check(
        ctx,
        dag.Directory("."),  // source directory
        nil,                // config file (optional)
        nil,                // ignore file (optional)
    )
}
```

### Format Code

```go
func (m *MyModule) FormatCode(ctx context.Context) (*Directory, error) {
    ruff := dag.Ruff()
    
    // Format Python files
    return ruff.Format(
        ctx,
        dag.Directory("."),  // source directory
        nil,                // config file (optional)
        nil,                // ignore file (optional)
    )
}
```

### Custom Configuration

```go
func (m *MyModule) CustomConfig(ctx context.Context) error {
    ruff := dag.Ruff()
    
    // Use custom configuration
    return ruff.Check(
        ctx,
        dag.Directory("."),
        dag.Directory(".").File("ruff.toml"),  // config file
        dag.Directory(".").File(".ruffignore"), // ignore file
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Python Lint
on: [push]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Lint Python
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/ruff
          args: |
            do -p '
              ruff := Ruff()
              ruff.Check(
                ctx,
                Directory("."),
                nil,
                nil,
              )
            '
```

## API Reference

### Ruff

Main module struct that provides access to Ruff functionality.

#### Methods

- `Check(ctx context.Context, dir *Directory, config *File, ignore *File) error`
  - Lints Python code
  - Parameters:
    - `dir`: Source directory
    - `config`: Configuration file (optional)
    - `ignore`: Ignore file (optional)

- `Format(ctx context.Context, dir *Directory, config *File, ignore *File) (*Directory, error)`
  - Formats Python code
  - Parameters:
    - `dir`: Source directory
    - `config`: Configuration file (optional)
    - `ignore`: Ignore file (optional)
  - Returns modified directory

## Best Practices

1. **Configuration Management**
   - Use project config
   - Document rules
   - Version control

2. **Code Quality**
   - Regular checks
   - Fix suggestions
   - Style consistency

3. **Integration**
   - Pre-commit hooks
   - CI/CD pipeline
   - Team standards

4. **Performance**
   - Cache results
   - Ignore patterns
   - Parallel checks

## Troubleshooting

Common issues and solutions:

1. **Linting Errors**
   ```
   Error: code style error
   Solution: Apply suggested fixes
   ```

2. **Configuration Issues**
   ```
   Error: invalid config
   Solution: Verify TOML syntax
   ```

3. **Path Problems**
   ```
   Error: file not found
   Solution: Check file paths
   ```

## Configuration Example

```toml
# ruff.toml
line-length = 88
target-version = "py39"

[lint]
select = [
    "E",  # pycodestyle errors
    "F",  # pyflakes
    "I",  # isort
]
ignore = ["E501"]  # line too long

[format]
quote-style = "double"
indent-style = "space"
line-ending = "lf"
```

## Advanced Usage

### Custom Linting Rules

```go
func (m *MyModule) CustomLint(ctx context.Context) error {
    ruff := dag.Ruff()
    
    // Create custom config
    config := `
[lint]
select = [
    "E",   # pycodestyle errors
    "F",   # pyflakes
    "I",   # isort
    "N",   # pep8-naming
    "UP",  # pyupgrade
]
ignore = []

[lint.per-file-ignores]
"__init__.py" = ["F401"]  # unused imports

[lint.isort]
known-first-party = ["myproject"]
`
    
    // Create config file
    configFile := dag.Container().
        From("alpine:latest").
        WithNewFile("/ruff.toml", dagger.ContainerWithNewFileOpts{
            Contents: config,
        }).
        File("/ruff.toml")
    
    // Run linting
    return ruff.Check(
        ctx,
        dag.Directory("."),
        configFile,
        nil,
    )
}
```

### Format and Verify

```go
func (m *MyModule) FormatAndVerify(ctx context.Context) error {
    ruff := dag.Ruff()
    
    // Format code
    formatted, err := ruff.Format(
        ctx,
        dag.Directory("."),
        nil,
        nil,
    )
    if err != nil {
        return err
    }
    
    // Verify formatting
    return dag.Container().
        From("alpine:latest").
        WithMountedDirectory("/src", formatted).
        WithWorkdir("/src").
        WithExec([]string{
            "sh", "-c",
            `
            # Check for any remaining style issues
            if [ -n "$(find . -name '*.py' -print0 | xargs -0 grep -l 'TODO: format')" ]; then
                echo "Found files with formatting markers"
                exit 1
            fi
            
            echo "All files properly formatted"
            `,
        }).
        Sync(ctx)
} 