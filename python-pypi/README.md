# Python PyPI Module for Dagger

A Dagger module for publishing Python packages to PyPI. This module handles the authentication and publishing process in a secure manner.

## Features

- Publish Python packages to PyPI
- Secure token handling
- Automated build and publish process

## Usage

Import the module in your Dagger pipeline:

```go
pypi := dag.PythonPypi()
```

### Publishing to PyPI

```go
err := pypi.With(source).WithSecret("PYPI_TOKEN", token).Publish(ctx)
if err != nil {
    // Handle error
}
```

## Requirements

- Dagger v0.15.1
- Go 1.23.4
- Python package ready for distribution
- PyPI token for authentication

## Environment

The module uses `python:3.11-slim` as the base image and automatically installs Poetry for package management.

## Security

- PyPI tokens are handled securely using Dagger secrets
- Tokens are never exposed in logs or container layers
- Authentication is handled automatically

## Example

Here's a complete example of publishing a package:

```go
func PublishPackage(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
    pypi := dag.PythonPypi()

    // Publish to PyPI with authentication
    return pypi.With(source).WithSecret("PYPI_TOKEN", token).Publish(ctx)
}
```

## Token Setup

1. Generate a PyPI token from your account settings
2. Pass the token as a secret to your pipeline:

```shell
export PYPI_TOKEN=your_token_here
dagger call publish --source=. --token=env:PYPI_TOKEN
```
