# Python Poetry Dagger Module

A Dagger module for building, testing, and publishing Python projects that use Poetry for dependency management and Ruff for linting.

## Features

- Poetry dependency installation
- Ruff linting
- Project building
- Package publishing

## Usage

### Command Line

```bash
# Install dependencies
dagger call python-poetry install-deps --src . --python-version 3.11

# Run linter
dagger call python-poetry lint --src . --python-version 3.11

# Build project
dagger call python-poetry build --src . --python-version 3.11

# Publish package
dagger call python-poetry publish --src . --python-version 3.11 --repository https://pypi.org/simple
```

### Go Integration

```go
import (
    "context"
    poetry "python-poetry"
)

func main() {
    ctx := context.Background()
    client, err := dagger.Connect(ctx)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    poetryMod := poetry.New(client)

    // Install dependencies
    container, err := poetryMod.InstallDeps(ctx, client.Host().Directory("."), "3.11")
    if err != nil {
        panic(err)
    }

    // Run linter
    container, err = poetryMod.Lint(ctx, client.Host().Directory("."), "3.11")
    if err != nil {
        panic(err)
    }

    // Build project
    container, err = poetryMod.Build(ctx, client.Host().Directory("."), "3.11")
    if err != nil {
        panic(err)
    }

    // Publish package
    container, err = poetryMod.Publish(ctx, client.Host().Directory("."), "3.11", "https://pypi.org/simple")
    if err != nil {
        panic(err)
    }
}
```

## Requirements

Your Python project should have:

- A valid `pyproject.toml` file
- Poetry configuration
- Ruff configuration (if using the linting feature)

## Example pyproject.toml

```toml
[tool.poetry]
name = "your-project"
version = "0.1.0"
description = "Your project description"

[tool.poetry.dependencies]
python = "^3.11"

[tool.poetry.group.dev.dependencies]
ruff = "^0.1.0"

[tool.ruff]
line-length = 88
target-version = "py311"
```
