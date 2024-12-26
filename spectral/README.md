# Spectral Module for Dagger

A Dagger module that integrates [Spectral](https://github.com/stoplightio/spectral), an open-source API style guide enforcer and linter. This module helps validate JSON/YAML OpenAPI documents against custom rulesets.

## Features

- JSON/YAML OpenAPI document linting
- Custom ruleset support
- Configurable severity levels
- JSON reference resolver support
- Multiple document validation
- Flexible output options
- Custom container support

## Usage

### Basic Setup

```typescript
import { spectral } from "@felipepimentel/daggerverse/spectral";

const client = spectral({
  version: "latest", // Version (optional)
  image: "", // Custom image (optional)
  container: null, // Custom container (optional)
});
```

### Document Linting

```typescript
// Lint documents with custom ruleset
const container = await client.lint({
  documents, // File[] - OpenAPI documents
  ruleset, // File - Ruleset file
  failSeverity: "error", // Fail severity (optional)
  displayOnlyFailures: false, // Display only failures (optional)
  resolver: null, // Custom resolver (optional)
  encoding: "utf8", // Text encoding (optional)
  verbose: false, // Verbose output (optional)
  quiet: false, // Quiet mode (optional)
});
```

## Configuration Options

### Constructor Parameters

- `version`: Version (image tag) from the official image repository
- `image`: Custom image reference in "repository:tag" format
- `container`: Custom container to use as base

### Linting Options

- `documents`: Array of JSON/YAML OpenAPI documents to validate
- `ruleset`: Spectral ruleset file defining validation rules
- `failSeverity`: Minimum severity level for failures (error, warn, info, hint)
- `displayOnlyFailures`: Show only results equal to or above fail severity
- `resolver`: Custom JSON reference resolver
- `encoding`: Text encoding (utf8, ascii, utf-8, utf16le, ucs2, ucs-2, base64, latin1)
- `verbose`: Enable verbose output
- `quiet`: Suppress logging, output only results

## Implementation Details

### Base Container

The module uses:

- Default image: `stoplight/spectral`
- Configurable version/tag
- Support for custom images
- Custom container integration

### Document Processing

- Mounts documents in container workspace
- Supports multiple document validation
- Handles file encoding
- Processes JSON references

### Ruleset Management

- Custom ruleset file support
- Flexible rule configuration
- Severity level management
- Failure criteria configuration

## Examples

### Basic Document Validation

```typescript
import { spectral } from "@felipepimentel/daggerverse/spectral";

// Initialize Spectral
const client = spectral();

// Prepare documents and ruleset
const documents = [dag.directory().withFile("openapi.yaml", openApiContent)];
const ruleset = dag.directory().withFile(".spectral.yaml", rulesetContent);

// Lint documents
const result = await client.lint({
  documents,
  ruleset,
  failSeverity: "error", // Fail on errors
  displayOnlyFailures: false, // Show all results
  resolver: null, // No custom resolver
  encoding: "utf8", // UTF-8 encoding
  verbose: false, // Normal output
  quiet: false, // Show logging
});
```

### Custom Configuration

```typescript
import { spectral } from "@felipepimentel/daggerverse/spectral";

// Use custom container
const customContainer = dag
  .container()
  .from("alpine:latest")
  .withExec(["apk", "add", "nodejs", "npm"])
  .withExec(["npm", "install", "-g", "@stoplight/spectral-cli"]);

const client = spectral({
  container: customContainer,
});

// Lint with custom options
const result = await client.lint({
  documents,
  ruleset,
  failSeverity: "warn", // Fail on warnings
  displayOnlyFailures: true, // Show only failures
  resolver, // Custom resolver
  encoding: "utf8",
  verbose: true, // Verbose output
  quiet: false,
});
```

## Dependencies

The module requires:

- Dagger SDK
- Node.js (provided in container)
- Spectral CLI (provided in container)
- Access to OpenAPI documents
- Valid ruleset file

## License

See [LICENSE](../LICENSE) file in the root directory.
