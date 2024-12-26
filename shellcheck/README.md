# Shellcheck Module for Dagger

A Dagger module that provides integration with ShellCheck, a static analysis tool for shell scripts. This module allows you to check shell scripts for common errors, bugs, and potential issues in your Dagger pipelines.

## Features

- Static analysis of shell scripts
- Automatic script discovery in directories
- Integration with Alpine-based container
- Support for multiple shell script formats (sh, bash, dash, ksh)
- Comprehensive error reporting

## Usage

### Basic Setup

```go
// Initialize the Shellcheck module
shellcheck := dag.Shellcheck()

// Check scripts in a directory
result, err := shellcheck.Check(
    dag.Directory().WithFiles(map[string]string{
        "script.sh": "#!/bin/bash\necho 'Hello World'",
    }),
).Stdout(ctx)
```

### Container Integration

The module uses the official `koalaman/shellcheck-alpine` container image, providing a lightweight environment for script checking.

```go
// Check scripts in a mounted directory
container := shellcheck.Check(sourceDir)

// Execute the check and get results
output, err := container.Stdout(ctx)
```

## Configuration

The module provides a simple interface with the following main function:

### Check

The `Check` function accepts:

- `source`: A Dagger Directory containing shell scripts to check
  - Automatically finds all `.sh` files in the directory
  - Recursively checks subdirectories

## Error Handling

The module reports various types of issues, including:

- Syntax errors
- Common bugs and pitfalls
- Style issues
- POSIX compliance problems
- Performance considerations

Example error codes:

- `SC2086`: Double quote to prevent globbing and word splitting
- `SC2154`: Variable referenced but not assigned
- `SC2283`: Remove spaces around = in arithmetic expressions

## Testing

To test the shellcheck module:

1. Ensure you have Dagger installed
2. Run the test suite:

```bash
dagger do test
```

The test suite includes:

- Verification of shellcheck version
- Testing against sample scripts with known issues
- Directory mounting and scanning tests

## Examples

### Check a Single Script

```go
func CheckScript(ctx context.Context) error {
    shellcheck := dag.Shellcheck()

    // Create a directory with a test script
    dir := dag.Directory().WithNewFile(
        "test.sh",
        `#!/bin/bash
        echo "Hello World"`,
    )

    // Run shellcheck
    _, err := shellcheck.Check(dir).Stdout(ctx)
    return err
}
```

### Check Multiple Scripts in a Directory

```go
func CheckDirectory(ctx context.Context) error {
    shellcheck := dag.Shellcheck()

    // Get the source directory
    sourceDir := dag.Host().Directory("./scripts")

    // Run shellcheck on all scripts
    _, err := shellcheck.Check(sourceDir).Stdout(ctx)
    return err
}
```

## Dependencies

The module requires:

- Go 1.22 or later
- Dagger SDK
- Internet access to pull the shellcheck container image

## License

This module is part of the Dagger project and is subject to its licensing terms.

Copyright Camptocamp SA
SPDX-License-Identifier: AGPL-3.0-or-later
