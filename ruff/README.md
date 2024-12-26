# Ruff Dagger Module

This Dagger module provides integration with [Ruff](https://github.com/astral-sh/ruff), an extremely fast Python linter written in Rust. Ruff can replace Flake8, isort, pydocstyle, yesqa, eradicate, pyupgrade, and autoflake.

## Features

- Fast Python code linting
- Support for custom configuration files
- Uses official Ruff container image
- Simple integration with your Dagger pipelines
- Configurable through `.ruff.toml`

## Usage

### Basic Usage

```typescript
import { ruff } from "@felipepimentel/daggerverse/ruff";

// Initialize ruff
const linter = ruff();

// Run checks on a directory
const container = linter.check(
  dag.host().directory("./src") // directory containing Python files
);
```

### Using Custom Configuration

```typescript
// Run with custom configuration file
const container = linter.checkWithConfig(
  dag.host().directory("./src"),
  dag.host().file(".ruff.toml")
);
```

## API Reference

### Constructor

#### `ruff()`

Creates a new instance of the Ruff module.

### Methods

#### `check(source: Directory): Container`

Runs Ruff checks on the provided source directory using default configuration.

Parameters:

- `source`: Directory containing Python files to check

Returns a container with the check results.

#### `checkWithConfig(source: Directory, file: File): Container`

Runs Ruff checks on the provided source directory using a custom configuration file.

Parameters:

- `source`: Directory containing Python files to check
- `file`: Ruff configuration file (`.ruff.toml`)

Returns a container with the check results.

## Configuration

### Default Configuration

When using the `check` method, Ruff will use its default configuration settings. For most projects, this provides a good starting point.

### Custom Configuration

When using `checkWithConfig`, you can provide a `.ruff.toml` file with your preferred settings. Example configuration:

```toml
# .ruff.toml
line-length = 88
target-version = "py37"
select = ["E", "F", "I"]
ignore = ["E501"]

[per-file-ignores]
"__init__.py" = ["F401"]

[mccabe]
max-complexity = 10
```

## Example

### Complete Linting Pipeline

```typescript
import { ruff } from "@felipepimentel/daggerverse/ruff";

export async function lint() {
  // Initialize ruff
  const linter = ruff();

  // Run linting checks with custom config
  const container = linter.checkWithConfig(
    dag.host().directory("./python_project"),
    dag.host().file("./python_project/.ruff.toml")
  );

  // Get the results
  const output = await container.stdout();

  // Process the results
  if (output) {
    console.log("Python Linting issues found:");
    console.log(output);
  } else {
    console.log("No Python issues found!");
  }
}
```

## Testing

To test the module:

1. Create a test Python file:

```python
# test.py
def hello_world():
    print ('Hello, World!')  # Extra space in parentheses
```

2. Run the linter:

```typescript
const container = linter.check(dag.host().directory("."));
```

3. Check the output for any linting issues

## License

See [LICENSE](../LICENSE) file in the root directory.
