# OpenAPI Changes Module for Dagger

A Dagger module that provides functionality to detect and analyze breaking changes between OpenAPI specifications. This module is built on top of the `pb33f/openapi-changes` tool and supports OpenAPI 3.1, 3.0, and Swagger specifications.

## Features

- Compare two OpenAPI specifications for breaking changes
- Track changes in a single specification over time
- Support for OpenAPI 3.1, 3.0, and Swagger
- Git integration for historical analysis
- Configurable output styling
- Base URL/directory support for relative references
- Limit and filtering options for change analysis
- CI/CD friendly output options

## Usage

### Basic Comparison

```typescript
import { openapiChanges } from "@felipepimentel/daggerverse/openapi-changes";

// Initialize the OpenAPI Changes module
const client = openapiChanges();

// Compare two OpenAPI specs
const result = await client.diff({
  old: oldSpec,
  new: newSpec,
});
```

### Git History Analysis

```typescript
// Initialize with styling disabled (for CI)
const client = openapiChanges({
  noStyle: true,
});

// Analyze changes in a Git repository
const result = await client.git({
  url: "https://github.com/org/repo",
  top: true,
  limit: 5,
});
```

## Configuration Options

### Version Selection

```typescript
// Use a specific version of the tool
const client = openapiChanges({
  version: "1.0.0",
});
```

### Custom Container

```typescript
// Use a custom container
const container = dag
  .container()
  .from("custom-image:latest")
  .withEnvVariable("CUSTOM_VAR", "value");

const client = openapiChanges({
  container: container,
});
```

### Output Styling

```typescript
// Disable color output and styling
const client = openapiChanges({
  noStyle: true,
});
```

## Examples

### Compare Local Specifications

```typescript
import { openapiChanges } from "@felipepimentel/daggerverse/openapi-changes";

export async function compareSpecs() {
  // Initialize module
  const client = openapiChanges();

  // Get spec files
  const oldSpec = dag.host().file("old-api.yaml");
  const newSpec = dag.host().file("new-api.yaml");

  // Compare specs
  const result = await client.diff({
    old: oldSpec,
    new: newSpec,
    noStyle: true,
    base: "http://api.example.com",
  });
}
```

### Analyze Git Repository Changes

```typescript
import { openapiChanges } from "@felipepimentel/daggerverse/openapi-changes";

export async function analyzeGitHistory() {
  // Initialize module
  const client = openapiChanges();

  // Analyze recent changes
  const result = await client.git({
    url: "https://github.com/org/api-specs",
    top: true,
    limit: 10,
    base: "/api/v1",
  });
}
```

### CI/CD Integration

```typescript
import { openapiChanges } from "@felipepimentel/daggerverse/openapi-changes";

export async function ciAnalysis() {
  // Initialize with CI-friendly options
  const client = openapiChanges({
    noStyle: true,
    version: "latest",
  });

  // Compare specs
  const result = await client.diff({
    old: oldSpec,
    new: newSpec,
    source: sourceDir,
    base: "http://api.example.com",
  });
}
```

## API Reference

### Constructor Options

The module accepts:

- `version`: Version of the tool to use (optional)
- `container`: Custom container to use (optional)
- `noStyle`: Disable color output and styling (optional)

### Methods

#### `diff(options: DiffOptions)`

Compare two OpenAPI specs.

Parameters:

- `old`: Old specification file
- `new`: New specification file
- `source`: Source directory (optional)
- `noStyle`: Disable styling (optional)
- `base`: Base URL for relative references (optional)

#### `git(options: GitOptions)`

Analyze changes in a Git repository.

Parameters:

- `url`: Repository URL
- `top`: Show only top-level changes (optional)
- `limit`: Maximum number of changes to show (optional)
- `base`: Base path for relative references (optional)

## Testing

The module includes a test suite that can be run using:

```bash
dagger do test
```

The test suite includes:

- Basic comparison tests
- Git integration tests
- Style configuration tests
- Error handling verification

## Dependencies

The module requires:

- Dagger SDK
- pb33f/openapi-changes image (pulled automatically)
- Internet access for initial container pulling
- Git (for repository analysis)

## Implementation Details

- Uses pb33f/openapi-changes as the base tool
- Supports multiple OpenAPI specification versions
- Provides both synchronous and asynchronous operation support
- Implements comprehensive error handling
- Supports various output formats and styling options

## License

See [LICENSE](../LICENSE) file in the root directory.
