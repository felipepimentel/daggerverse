# Stainless Module for Dagger

A Dagger module that provides integration with Stainless, an API SDK generator tool. This module enables you to upload OpenAPI specifications and manage configurations for generating SDKs through the Stainless service.

## Features

- OpenAPI specification upload
- Configuration management
- Token-based authentication
- Multipart form data handling
- Conventional commit message support
- Project-based organization

## Usage

### Basic Setup

```go
// Initialize the Stainless module with your API token
stainless := dag.Stainless().New(dag.SetSecret("token", "your-stainless-token"))
```

### Upload OpenAPI Specification

```go
// Upload an OpenAPI spec file
result, err := stainless.UploadSpec(
    ctx,
    "my-project",                // Project name
    openApiFile,                 // OpenAPI spec file
    nil,                         // Optional: config file
    "",                         // Optional: commit message
)
if err != nil {
    return err
}
```

### Upload with Configuration

```go
// Upload with a custom configuration file
result, err := stainless.UploadSpec(
    ctx,
    "my-project",
    openApiFile,
    configFile,
    "feat: update API specification",
)
if err != nil {
    return err
}
```

## Configuration

### Constructor Options

The `New` function accepts:

- `token`: Stainless API token (\*dagger.Secret)

### Upload Options

The `UploadSpec` method accepts:

- `projectName`: Name of the Stainless project (required)
- `openapi`: OpenAPI specification file (\*dagger.File, required)
- `config`: Stainless configuration file (\*dagger.File, optional)
- `commitMessage`: Conventional commit message (string, optional)

## Examples

### Complete Example

```go
func UploadAPISpec(ctx context.Context) error {
    // Initialize Stainless with token
    stainless := dag.Stainless().New(
        dag.SetSecret("token", os.Getenv("STAINLESS_TOKEN")),
    )

    // Prepare OpenAPI spec file
    openApiFile := dag.Directory().
        WithNewFile("openapi.json", `{
            "openapi": "3.0.0",
            "info": {
                "title": "My API",
                "version": "1.0.0"
            }
        }`)

    // Prepare config file
    configFile := dag.Directory().
        WithNewFile("stainless.yaml", `
            name: my-api-sdk
            version: 1.0.0
            language: typescript
        `)

    // Upload specification
    result, err := stainless.UploadSpec(
        ctx,
        "my-project",
        openApiFile,
        configFile,
        "feat: update API specification with new endpoints",
    )
    if err != nil {
        return err
    }

    return nil
}
```

### Project Management

```go
func ManageProject(ctx context.Context) error {
    // Initialize Stainless
    stainless := dag.Stainless().New(
        dag.SetSecret("token", os.Getenv("STAINLESS_TOKEN")),
    )

    // Upload initial spec
    _, err := stainless.UploadSpec(
        ctx,
        "my-project",
        initialSpec,
        nil,
        "feat: initial API specification",
    )
    if err != nil {
        return err
    }

    // Update with new configuration
    _, err = stainless.UploadSpec(
        ctx,
        "my-project",
        updatedSpec,
        newConfig,
        "chore: update SDK configuration",
    )
    if err != nil {
        return err
    }

    return nil
}
```

## Dependencies

The module requires:

- Dagger SDK
- Stainless API token
- OpenAPI specification file
- Optional: Stainless configuration file

## Testing

The module includes tests that verify:

- OpenAPI specification upload
- Configuration file handling
- Token authentication
- Error handling

To run the tests:

```bash
dagger do test
```

## License

This module is part of the Dagger project and is subject to its licensing terms.

Copyright Camptocamp SA
SPDX-License-Identifier: AGPL-3.0-or-later
