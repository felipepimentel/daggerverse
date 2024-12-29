# Python CI Pipeline

This is a reusable GitHub Actions workflow that encapsulates the Python module from Daggerverse for CI/CD.

## Usage

To use this workflow in your Python project, create a file `.github/workflows/ci.yml` with the following content:

```yaml
name: CI/CD

on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]
  release:
    types: [published]

# Required permissions
permissions:
  contents: read
  id-token: write

jobs:
  python-pipeline:
    uses: felipepimentel/daggerverse/.github/workflows/reusable/python-ci.yml@v1
    secrets:
      PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
```

## Configuration

### Required Permissions

The workflow requires the following permissions:

- `contents: read` - To access source code
- `id-token: write` - For OIDC authentication with PyPI

### Inputs

| Name                    | Description           | Default | Required |
| ----------------------- | --------------------- | ------- | -------- |
| `python-module-version` | Python module version | `main`  | No       |
| `source-path`           | Source code path      | `.`     | No       |

### Secrets

| Name         | Description               | Required |
| ------------ | ------------------------- | -------- |
| `PYPI_TOKEN` | PyPI token for publishing | Yes      |

## Examples

### Basic Configuration

```yaml
name: CI/CD

on:
  push:
    branches: ["main"]

permissions:
  contents: read
  id-token: write

jobs:
  python-pipeline:
    uses: felipepimentel/daggerverse/.github/workflows/reusable/python-ci.yml@v1
    secrets:
      PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
```

### Custom Configuration

```yaml
name: CI/CD

on:
  push:
    branches: ["main"]

permissions:
  contents: read
  id-token: write

jobs:
  python-pipeline:
    uses: felipepimentel/daggerverse/.github/workflows/reusable/python-ci.yml@v1
    with:
      python-module-version: "v1.2.3"
      source-path: "./src"
    secrets:
      PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
```
