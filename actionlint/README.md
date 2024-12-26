# Actionlint Dagger Module

This Dagger module provides integration with [actionlint](https://github.com/rhysd/actionlint), a static checker for GitHub Actions workflow files.

## Features

- Run actionlint checks on GitHub Actions workflow files
- Customizable base container image
- Simple integration with your Dagger pipelines

## Usage

### Basic Usage

```go
// Import the module in your Dagger code
import "dagger.io/dagger"
import "github.com/yourorg/daggerverse/actionlint"

// Create a new instance
actionlint := actionlint.New()

// Run checks on your workflow files
result := actionlint.Check(dag.Host().Directory("./github/workflows"))
```

### Custom Image

You can specify a custom actionlint image:

```go
actionlint := actionlint.New(
    "custom/actionlint:v1.0.0"
)
```

## API Reference

### Constructor

#### `New(image string) *Actionlint`

Creates a new instance of the Actionlint module.

Parameters:

- `image` (optional): Custom image reference in "repository:tag" format. Defaults to "rhysd/actionlint:latest"

### Methods

#### `Container() *dagger.Container`

Returns the underlying Dagger container used by the module.

#### `Check(source *dagger.Directory) *dagger.Container`

Runs actionlint checks on the provided source directory.

Parameters:

- `source`: Directory containing GitHub Actions workflow files to check

Returns a container with the check results.

## Example

```go
func Pipeline(ctx context.Context) (*dagger.Container, error) {
    d := dagger.Connect(ctx)
    defer d.Close()

    // Initialize actionlint
    linter := actionlint.New()

    // Run checks on workflows directory
    result := linter.Check(
        d.Host().Directory("./.github/workflows"),
    )

    return result, nil
}
```

## Testing

To test the module, you can create a simple workflow file and run the checks:

1. Create a test workflow file:

```yaml
# test.yml
name: Test Workflow
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
```

2. Run the check:

```shell
dagger call check --source ./path/to/workflows
```

## License

This module is available under the [Apache License 2.0](LICENSE).
