# MkDocs Pipeline

This module provides a powerful and flexible way to build and deploy MkDocs documentation with Material theme support.

## Features

- Material theme with full feature support
- Custom Python version selection
- Requirements file support
- Git revision date integration
- HTML minification
- Strict mode validation
- Development server support

## Installation

```bash
dagger mod use github.com/felipepimentel/daggerverse/pipelines/mkdocs@latest
```

## Usage

### Basic Example

```go
// Initialize the module
mkdocs := dag.MkDocs()

// Build documentation
output, err := mkdocs.Build(ctx, &MkDocsConfig{
    Source: dag.Host().Directory("."),
    BaseURL: "https://yourdomain.github.io/project",
    Strict: true,
})
```

### Configuration Options

The module supports the following configuration options:

```go
type MkDocsConfig struct {
    // Source directory containing mkdocs.yml and docs/
    Source *dagger.Directory
    // Custom requirements file (optional)
    RequirementsFile *dagger.File
    // Output directory name (default: "site")
    OutputDir string
    // Base URL for the documentation
    BaseURL string
    // Whether to use strict mode
    Strict bool
    // Whether to minify HTML
    Minify bool
    // Whether to include git revision date
    GitRevisionDate bool
}
```

### Development Server

For local development, you can use the `Serve` function:

```go
container := mkdocs.Serve(&MkDocsConfig{
    Source: dag.Host().Directory("."),
})
```

This will start a development server on port 8000.

## GitHub Actions Integration

This module includes a reusable GitHub Actions workflow for easy integration:

```yaml
name: Documentation
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  docs:
    uses: felipepimentel/daggerverse/.github/workflows/reusable-mkdocs.yml@main
    with:
      source_dir: '.'
      base_url: 'https://yourdomain.github.io/project'
      environment: 'github-pages'
    secrets:
      github_token: ${{ secrets.GITHUB_TOKEN }}
```

### Workflow Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `python_version` | Python version to use | No | `3.11` |
| `source_dir` | Directory containing mkdocs.yml and docs/ | No | `.` |
| `requirements_file` | Path to custom requirements.txt file | No | `""` |
| `base_url` | Base URL for the documentation | No | `""` |
| `strict_mode` | Enable strict mode for MkDocs build | No | `true` |
| `minify` | Enable HTML minification | No | `true` |
| `git_revision_date` | Include git revision date | No | `true` |
| `environment` | GitHub environment to deploy to | No | `github-pages` |

## Examples

### Custom Python Version

```go
mkdocs := dag.MkDocs().WithPythonVersion("3.10")
```

### Custom Requirements

```go
output, err := mkdocs.Build(ctx, &MkDocsConfig{
    Source: dag.Host().Directory("."),
    RequirementsFile: dag.Host().File("requirements.txt"),
})
```

### Full Configuration

```go
output, err := mkdocs.Build(ctx, &MkDocsConfig{
    Source: dag.Host().Directory("."),
    RequirementsFile: dag.Host().File("requirements.txt"),
    OutputDir: "public",
    BaseURL: "https://docs.example.com",
    Strict: true,
    Minify: true,
    GitRevisionDate: true,
})
```

## Best Practices

1. **Source Structure**:
   - Keep your documentation in a `docs/` directory
   - Place `mkdocs.yml` in the root of your project
   - Use relative links in your markdown files

2. **Configuration**:
   - Enable strict mode during CI to catch issues early
   - Use git revision date for better version tracking
   - Enable minification for production builds

3. **Development**:
   - Use the `Serve` function for local development
   - Validate configuration before deployment
   - Keep requirements file up to date

## Common Issues

1. **Missing mkdocs.yml**:
   - Ensure `mkdocs.yml` exists in your source directory
   - Check file permissions

2. **Theme Issues**:
   - Verify Material theme is properly configured
   - Check custom theme requirements are installed

3. **Build Failures**:
   - Run with strict mode to identify issues
   - Check Python version compatibility
   - Verify all required dependencies are installed

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](../CONTRIBUTING.md) for details on how to submit pull requests.

## License

This module is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details. 