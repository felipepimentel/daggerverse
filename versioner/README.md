# Versioner Module for Dagger

A Dagger module for automated semantic versioning and changelog generation.

## Features

- Semantic versioning automation
- Conventional commits validation
- Automatic version bumping based on commit types
- Changelog generation
- Git tag management
- Multi-file version updating
- Flexible configuration options

## Usage

### Basic Example

```go
// Initialize client
client, err := dagger.Connect(context.Background())
if err != nil {
    panic(err)
}
defer client.Close()

// Create versioner module
versioner := dag.Versioner()

// Bump version based on commits
newVersion, err := versioner.BumpVersion(ctx, dag.Host().Directory("."))
if err != nil {
    panic(err)
}
fmt.Println("New version:", newVersion)
```

### Command Line Usage

```bash
# Get current version
dagger -m github.com/felipepimentel/daggerverse/versioner call get-current-version --source .

# Bump version based on commits
dagger -m github.com/felipepimentel/daggerverse/versioner call bump-version --source .
```

## Configuration

### Version Configuration

```go
versioner := dag.Versioner().WithConfig(&VersionConfig{
    Strategy: "semver",
    TagPrefix: "v",
    VersionFiles: []string{"pyproject.toml", "package.json"},
    VersionPattern: `version\s*=\s*["']?(\d+\.\d+\.\d+)["']?`,
    ValidateCommits: true,
    GenerateChangelog: true,
    CreateTags: true,
})
```

### Branch Configuration

```go
versioner := dag.Versioner().WithConfig(&VersionConfig{
    BranchConfig: &BranchConfig{
        MainBranch: "main",
        DevelopBranch: "develop",
        ReleaseBranchPrefix: "release/",
        HotfixBranchPrefix: "hotfix/",
        FeatureBranchPrefix: "feature/",
    },
})
```

### Commit Configuration

```go
versioner := dag.Versioner().WithConfig(&VersionConfig{
    CommitConfig: &CommitConfig{
        Types: map[string]VersionIncrement{
            "feat":     Minor,
            "fix":      Patch,
            "perf":     Patch,
            "refactor": None,
        },
        Scopes: []string{"core", "deps", "docs"},
        RequireScope: false,
        AllowCustomScopes: true,
        BreakingChangeIndicators: []string{"BREAKING CHANGE:", "!"},
    },
})
```

## Commit Message Format

The module follows the Conventional Commits specification:

```
type(scope): subject

[optional body]

[optional footer(s)]
```

### Types and Version Increments

- `feat`: Minor version bump (0.1.0 -> 0.2.0)
  - New features
  - New functionality
  - New APIs
- `fix`: Patch version bump (0.1.0 -> 0.1.1)
  - Bug fixes
  - Error corrections
  - Security patches
- `perf`: Patch version bump
  - Performance improvements
  - Optimizations
- `refactor`: No version bump
  - Code restructuring
  - Moving/renaming
- `style`: No version bump
  - Formatting
  - White-space
- `docs`: No version bump
  - Documentation
  - Comments
- `test`: No version bump
  - Adding/updating tests
- `build`: No version bump
  - Build process
  - Dependencies
- `ci`: No version bump
  - CI configuration
  - Pipeline changes
- `chore`: No version bump
  - Maintenance
  - Cleanup

### Breaking Changes

Breaking changes can be indicated by:

1. Adding `!` after type/scope
2. Including `BREAKING CHANGE:` in footer

Example:

```
feat(api)!: change authentication interface

BREAKING CHANGE: The auth interface now requires explicit token configuration.
Previous token configuration will need to be updated.
```

## Integration with Other Modules

### Python Module

```go
// Get current version
version, err := versioner.GetCurrentVersion(ctx, source)
if err != nil {
    return err
}

// Update version in pyproject.toml
python := dag.Python().
    WithPackagePath(".")

// Run CI/CD with new version
output, err := python.CICD(ctx, source, token)
```

### GitHub Actions Integration

```yaml
name: Version and Release

on:
  push:
    branches: [main]

jobs:
  version:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/versioner@main
        run: |
          # Bump version based on commits
          dagger call bump-version --source .
```

## Best Practices

1. **Commit Messages**

   - Use conventional commit format
   - Be descriptive in commit messages
   - Mark breaking changes appropriately

2. **Version Files**

   - Keep version in standard locations
   - Use consistent version format
   - Update all relevant files

3. **Changelog**

   - Review generated changelog
   - Add manual entries if needed
   - Keep entries clear and concise

4. **Git Tags**

   - Use consistent tag prefix
   - Include version in tag message
   - Push tags after creation

5. **Breaking Changes**
   - Document migration steps
   - Update documentation
   - Provide upgrade guide
