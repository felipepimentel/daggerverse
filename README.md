# Daggerverse

Collection of reusable Dagger modules and GitHub Actions workflows.

## Usage Options

There are two ways to use this pipeline in your Python projects:

### Option 1: Using Reusable Workflow (Recommended)

This approach uses GitHub Actions' reusable workflows feature, providing better modularity and maintenance. The workflow handles everything automatically, including version management and git configuration:

```yaml
name: CI/CD

on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]
  release:
    types: [published]

permissions:
  contents: write
  id-token: write

jobs:
  pipeline:
    uses: felipepimentel/daggerverse/.github/workflows/reusable-python-ci.yml@main
    secrets:
      PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
    with:
      source-path: "."
```

The reusable workflow handles:

- Git checkout and configuration
- Version management using semantic versioning
- Python package building and testing
- PyPI publishing
- Error handling and reporting

### Option 2: Direct Module Usage

This approach calls Dagger modules directly, offering more flexibility and control when you need to customize the pipeline:

```yaml
name: CI/CD

on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]
  release:
    types: [published]

permissions:
  contents: write
  id-token: write

jobs:
  pipeline:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: Configure Git
        run: |
          git config --global user.name 'github-actions[bot]'
          git config --global user.email 'github-actions[bot]@users.noreply.github.com'

      - name: Version Management
        id: versioner
        uses: dagger/dagger-for-github@v7
        with:
          verb: call
          module: github.com/felipepimentel/daggerverse/versioner@main
          args: bump-version --source . --output-version

      - name: Run Python Pipeline
        uses: dagger/dagger-for-github@v7
        env:
          PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
          VERSION: ${{ steps.versioner.outputs.version }}
        with:
          verb: call
          module: github.com/felipepimentel/daggerverse/python-pipeline@main
          args: cicd --source . --token env:PYPI_TOKEN --version env:VERSION
```

## Required Secrets

- `PYPI_TOKEN`: Your PyPI token for publishing packages

Note: `GITHUB_TOKEN` is automatically provided by GitHub Actions and handled by the workflow.

## Available Modules

### Versioner

A module for managing semantic versioning of your projects.

#### Usage via CLI

```bash
dagger call -m github.com/felipepimentel/daggerverse/versioner@main bump-version --source . --output-version
```

### Python Pipeline

A module for handling Python project CI/CD pipelines.

#### Usage via CLI

```bash
dagger call -m github.com/felipepimentel/daggerverse/python-pipeline@main cicd --source . --token $PYPI_TOKEN --version $VERSION
```

## Development

### Module Structure

Each module should have:

- `dagger.json`: Module configuration
- `main.go`: Module implementation
- `README.md`: Module documentation
- Tests (when applicable)

### Contributing

1. Create a new branch for your changes
2. Follow the commit message format specified in `.cursorrules`
3. Submit a pull request

For more details on commit messages and other rules, see `.cursorrules`.

## Troubleshooting

### Common Issues

1. **Module Initialization Error**

   ```
   Error: module must be fully initialized
   ```

   Solution: Make sure to:

   - Use the correct module path (e.g., `python-pipeline` instead of `python`)
   - Include `@main` or specific version tag in module path
   - Verify that the module's `dagger.json` is valid

2. **Version Management Issues**
   If the version is not being passed correctly between steps, ensure:
   - The workflow is using the latest version of the reusable workflow
   - All required permissions are set correctly
   - The source path is correctly specified

### Best Practices

1. Always use tagged versions for stability
2. Use the reusable workflow when possible for standardization
3. Keep secrets secure using GitHub's secrets management
4. Test workflow changes in a feature branch before merging to main

## Workflow Features

The reusable workflow provides:

1. **Automated Version Management**

   - Semantic versioning based on commit messages
   - Automatic version bumping
   - Version output for downstream jobs

2. **Git Configuration**

   - Automatic repository checkout
   - Bot user configuration
   - Proper commit history handling

3. **Python Package Management**

   - Package building and testing
   - PyPI publishing
   - Dependency management

4. **Security**
   - Secure secrets handling
   - Token management
   - Permission controls
