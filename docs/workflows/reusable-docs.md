# Reusable Documentation Workflow

This workflow provides a standardized way to build and deploy documentation using various documentation generators.

## Supported Generators

- **Just The Docs**: A modern, highly customizable, and responsive Jekyll theme for documentation
- **MkDocs**: A fast, simple, and downright gorgeous static site generator that's geared towards building project documentation
- **Sphinx**: A powerful documentation generator that supports multiple output formats and has excellent Python integration

## Usage

### Basic Example

```yaml
name: "Deploy Documentation"
on:
  push:
    branches: [main]
    paths:
      - "docs/**"
      - ".github/workflows/docs.yml"

jobs:
  docs:
    uses: felipepimentel/daggerverse/.github/workflows/reusable-docs.yml@main
    with:
      generator: just-the-docs  # or 'mkdocs' or 'sphinx'
      source_dir: docs
      environment: github-pages
    secrets:
      token: ${{ secrets.GITHUB_TOKEN }}
```

### Generator-Specific Examples

#### Just The Docs

```yaml
jobs:
  docs:
    uses: felipepimentel/daggerverse/.github/workflows/reusable-docs.yml@main
    with:
      generator: just-the-docs
      source_dir: docs
      ruby_version: "3.3"  # optional
      environment: github-pages
    secrets:
      token: ${{ secrets.GITHUB_TOKEN }}
```

#### MkDocs

```yaml
jobs:
  docs:
    uses: felipepimentel/daggerverse/.github/workflows/reusable-docs.yml@main
    with:
      generator: mkdocs
      source_dir: docs
      python_version: "3.12"  # optional
      strict_mode: true  # optional
      minify: true  # optional
      git_revision_date: true  # optional
      environment: github-pages
    secrets:
      token: ${{ secrets.GITHUB_TOKEN }}
```

#### Sphinx with Poetry

```yaml
jobs:
  docs:
    uses: felipepimentel/daggerverse/.github/workflows/reusable-docs.yml@main
    with:
      generator: sphinx
      source_dir: docs
      python_version: "3.12"  # optional
      poetry_version: "latest"  # optional
      poetry_virtualenvs_in_project: true  # optional
      poetry_groups: "dev,docs"  # optional
      environment: github-pages
    secrets:
      token: ${{ secrets.GITHUB_TOKEN }}
```

## Configuration Options

### Common Options

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `generator` | Documentation generator to use (`just-the-docs`, `mkdocs`, `sphinx`) | Yes | - |
| `source_dir` | Directory containing documentation source | No | `docs` |
| `environment` | GitHub environment to deploy to | No | `github-pages` |
| `base_url` | Base URL for the documentation | No | `''` |

### Just The Docs Options

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `ruby_version` | Ruby version to use | No | `3.3` |

### MkDocs Options

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `python_version` | Python version to use | No | `3.12` |
| `strict_mode` | Enable strict mode | No | `true` |
| `minify` | Enable HTML minification | No | `true` |
| `git_revision_date` | Include git revision date | No | `true` |

### Sphinx Options

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `python_version` | Python version to use | No | `3.12` |
| `poetry_version` | Poetry version to use | No | `latest` |
| `poetry_virtualenvs_in_project` | Create virtualenvs in project directory | No | `true` |
| `poetry_groups` | Poetry dependency groups to install | No | `dev,docs` |

## Required Repository Setup

### Just The Docs
- `Gemfile` with Just The Docs theme
- Jekyll configuration in `_config.yml`

### MkDocs
- `mkdocs.yml` configuration file
- Documentation in `docs/` directory

### Sphinx
- `pyproject.toml` with Poetry configuration
- Sphinx configuration in `docs/conf.py`
- Makefile in `docs/` directory

## Security

The workflow requires the following permissions:
- `contents: read`
- `pages: write`
- `id-token: write`

These permissions are necessary for checking out code, deploying to GitHub Pages, and handling authentication. 