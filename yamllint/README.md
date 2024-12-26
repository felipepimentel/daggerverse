# Yamllint Dagger Module

This Dagger module provides integration with [yamllint](https://github.com/adrienverge/yamllint), a linter for YAML files. It helps you detect syntax errors, key repetition, and enforce best practices in your YAML files.

## Features

- YAML syntax validation
- Configurable linting rules
- Custom base image support
- Default configuration with relaxed line length warnings
- Easy integration with your Dagger pipelines

## Usage

### Basic Usage

```typescript
import { yamllint } from "@felipepimentel/daggerverse/yamllint";

// Initialize yamllint
const linter = yamllint();

// Run checks on a directory
const container = linter.check(
  dag.host().directory("./config") // directory containing YAML files
);
```

### Custom Image

You can specify a custom yamllint image:

```typescript
const linter = yamllint({
  image: "custom/yamllint:1.0.0",
});
```

## API Reference

### Constructor

#### `yamllint(options?: YamllintOptions)`

Creates a new instance of the Yamllint module.

Parameters:

- `image` (optional): Custom image reference in "repository:tag" format. Defaults to "pipelinecomponents/yamllint:latest"

### Methods

#### `container(): Container`

Returns the underlying Dagger container used by the module.

#### `check(source: Directory): Container`

Runs yamllint checks on the provided source directory.

Parameters:

- `source`: Directory containing YAML files to check

Returns a container with the check results.

## Default Configuration

The module uses a default configuration that:

- Extends the default yamllint rules
- Sets line-length rule to warning level
- Suppresses warnings in the output

The configuration is equivalent to:

```yaml
extends: default
rules:
  line-length:
    level: warning
```

## Example

### Basic Linting Pipeline

```typescript
import { yamllint } from "@felipepimentel/daggerverse/yamllint";

export async function lint() {
  // Initialize yamllint
  const linter = yamllint();

  // Run linting checks
  const container = linter.check(
    dag.host().directory("./kubernetes") // directory with YAML files
  );

  // Get the results
  const output = await container.stdout();

  // Process the results
  if (output) {
    console.log("YAML Linting issues found:");
    console.log(output);
  } else {
    console.log("No YAML issues found!");
  }
}
```

## Testing

To test the module:

1. Create a test YAML file:

```yaml
# test.yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  ports:
    - port: 80
```

2. Run the linter:

```typescript
const container = linter.check(dag.host().directory("."));
```

3. Check the output for any linting issues

## License

See [LICENSE](../LICENSE) file in the root directory.
