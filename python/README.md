# Python Module for Dagger

A comprehensive Python module for Dagger that streamlines Poetry-based Python development workflows. Features include automated package building with proper dependency management, configurable test execution with coverage reporting, secure PyPI publishing with registry selection, and integrated Git operations.

## Features

- 🏗️ **Poetry Integration**: Full support for Poetry package management
- 🧪 **Advanced Testing**: Configurable pytest execution with coverage reporting
- 📦 **Package Building**: Automated build process with dependency management
- 🚀 **PyPI Publishing**: Secure package publishing with registry selection
- 🔍 **Code Quality**: Integrated linting and formatting tools
- 📚 **Documentation**: Automated documentation generation
- 💾 **Caching**: Optimized build and dependency caching
- 🔄 **Git Operations**: Built-in repository checkout and authentication

## Quick Start

```go
import (
    "context"
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/python"
)

func main() {
    ctx := context.Background()

    // Initialize client
    client, err := dagger.Connect(ctx)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Configure Python module with Git checkout
    python := dag.Python().
        WithPythonVersion("3.12").
        WithGitConfig(&GitConfig{
            Repository: "https://github.com/username/repo",
            Ref: "main",
            Token: client.SetSecret("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN")),
        })

    // Checkout repository
    source, err := python.Checkout(ctx)
    if err != nil {
        panic(err)
    }

    // Run tests with coverage
    output, err := python.Test(ctx, source)
    if err != nil {
        panic(err)
    }
    fmt.Println(output)

    // Build and publish package
    token := client.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN"))
    output, err = python.Publish(ctx, source, token)
    if err != nil {
        panic(err)
    }
    fmt.Println(output)
}
```

## Configuration Options

### Git Configuration

```go
gitConfig := &GitConfig{
    // Repository URL to clone
    Repository: "https://github.com/username/repo",

    // Branch or tag to checkout (default: main)
    Ref: "main",

    // Depth of git history to clone (default: 1)
    Depth: 1,

    // Git clone options
    FetchAll: false,
    FetchTags: false,
    Submodules: false,

    // Authentication (choose one)
    Token: client.SetSecret("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN")),
    SSHKey: client.SetSecret("SSH_KEY", os.Getenv("SSH_KEY")),
    KnownHosts: "github.com ssh-rsa AAAA...",

    // Additional configuration
    Config: []KeyValue{
        {Key: "user.name", Value: "CI Bot"},
        {Key: "user.email", Value: "ci@example.com"},
    },

    // Environment variables
    Env: []KeyValue{
        {Key: "GIT_SSL_NO_VERIFY", Value: "true"},
    },
}

python = python.WithGitConfig(gitConfig)

// Checkout repository
source, err := python.Checkout(ctx)
```

### Python Environment

```go
python := dag.Python().
    WithPythonVersion("3.12").  // Python version to use
    WithPackagePath("src")      // Path to package within source
```

### Build Configuration

```go
buildConfig := &BuildConfig{
    // Additional build arguments
    BuildArgs: []string{"--no-dev"},

    // Extra dependencies to install
    ExtraDependencies: []string{"pytest-cov"},

    // Poetry configuration
    PoetryConfig: []KeyValue{
        {Key: "virtualenvs.in-project", Value: "true"},
    },

    // Environment variables
    Env: []KeyValue{
        {Key: "POETRY_HOME", Value: "/opt/poetry"},
    },

    // Cache configuration
    Cache: &CacheConfig{
        PipCache: true,
        PoetryCache: true,
        PipCacheVolume: "pip-cache",
        PoetryCacheVolume: "poetry-cache",
    },

    // Poetry dependency groups to install
    DependencyGroups: []string{"dev", "test"},

    // Optional dependencies
    OptionalDependencies: []string{"docs"},

    // Installation options
    SkipDependencies: false,
    OnlyGroups: false,
    SkipRoot: true,
}

python = python.WithBuildConfig(buildConfig)
```

### Test Configuration

```go
testConfig := &TestConfig{
    // Test execution options
    Verbose: true,
    Workers: 4,  // Number of parallel workers

    // Coverage configuration
    Coverage: &CoverageConfig{
        Enabled: true,
        Formats: []string{"term", "xml", "html"},
        MinCoverage: 80,
        OutputDir: "coverage",
        Include: []string{"src/*"},
        Exclude: []string{"tests/*"},
        ShowMissing: true,
        Branch: true,
        Context: 3,
    },

    // Environment variables
    Env: []KeyValue{
        {Key: "PYTHONPATH", Value: "src"},
    },

    // Test selection
    Markers: []string{"unit", "integration"},
    TestPaths: []string{"tests"},

    // Test execution controls
    SkipInstall: false,
    JUnitXML: "test-results.xml",
    MaxTestTime: 300,
    FailFast: true,
}

python = python.WithTestConfig(testConfig)
```

### PyPI Publishing Configuration

```go
pypiConfig := &PyPIConfig{
    // Registry configuration
    Registry: "https://test.pypi.org/legacy/",
    Token: client.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN")),
    RepositoryName: "testpypi",

    // Publishing options
    SkipExisting: true,
    AllowDirty: false,
    SkipBuild: false,
    SkipVerify: false,

    // Additional configuration
    ExtraArgs: []string{"--verbose"},
    Env: []KeyValue{
        {Key: "POETRY_HTTP_TIMEOUT", Value: "60"},
    },
}

python = python.WithPyPIConfig(pypiConfig)
```

## Available Functions

### Core Functions

- `Checkout(ctx)`: Clones a Git repository and returns its directory
- `Build(source)`: Creates a Python package using Poetry
- `BuildEnv(source)`: Prepares a Python development environment
- `Publish(ctx, source, token)`: Publishes package to PyPI
- `Test(ctx, source)`: Runs test suite with coverage reporting

### Quality Tools

- `Lint(ctx, source)`: Runs code linting using ruff
- `Format(ctx, source)`: Formats code using black
- `BuildDocs(ctx, source)`: Generates documentation

## Common Use Cases

### CI/CD Pipeline

```go
func main() {
    ctx := context.Background()
    client, err := dagger.Connect(ctx)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Initialize Python module
    python := dag.Python().
        WithPythonVersion("3.12").
        WithGitConfig(&GitConfig{
            Repository: "https://github.com/username/repo",
            Ref: "main",
            Token: client.SetSecret("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN")),
        })

    // Checkout code
    source, err := python.Checkout(ctx)
    if err != nil {
        panic(err)
    }

    // Run linting
    if _, err := python.Lint(ctx, source); err != nil {
        panic(err)
    }

    // Run tests with coverage
    if _, err := python.Test(ctx, source); err != nil {
        panic(err)
    }

    // Build documentation
    if _, err := python.BuildDocs(ctx, source); err != nil {
        panic(err)
    }

    // Publish to PyPI
    if _, err := python.Publish(ctx, source, client.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN"))); err != nil {
        panic(err)
    }
}
```

### Development Environment

```go
func main() {
    ctx := context.Background()
    client, err := dagger.Connect(ctx)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Initialize Python module with development configuration
    python := dag.Python().
        WithPythonVersion("3.12").
        WithBuildConfig(&BuildConfig{
            DependencyGroups: []string{"dev", "test"},
            Cache: &CacheConfig{
                PipCache: true,
                PoetryCache: true,
            },
        })

    // Get source code
    source := client.Host().Directory(".")

    // Create development environment
    container := python.BuildEnv(source)

    // Run development shell
    if _, err := container.WithExec([]string{"poetry", "shell"}).Stdout(ctx); err != nil {
        panic(err)
    }
}
```

## Best Practices

1. **Git Operations**:

   - Use shallow clones (depth: 1) for CI/CD pipelines
   - Enable submodules only when needed
   - Use tokens or SSH keys securely
   - Configure known hosts for SSH authentication

2. **Cache Management**:

   - Enable both pip and Poetry caches for faster builds
   - Use custom volume names for specific projects
   - Consider cache cleanup for large projects

3. **Testing Strategy**:

   - Configure appropriate test markers
   - Set reasonable coverage thresholds
   - Use parallel testing for large test suites

4. **Code Quality**:

   - Enable automatic fixes when appropriate
   - Configure line length consistently
   - Use consistent Python version targets

5. **Documentation**:

   - Choose appropriate documentation tool
   - Enable relevant extensions
   - Configure theme for better readability

6. **Publishing**:
   - Always use secure token handling
   - Consider using test PyPI first
   - Enable appropriate verifications

## Security Considerations

1. **Authentication**:

   - Never hardcode tokens or SSH keys
   - Use environment variables or secrets management
   - Rotate tokens regularly
   - Use read-only tokens when possible

2. **Git Security**:

   - Verify SSL certificates
   - Use known hosts for SSH
   - Limit repository access
   - Use shallow clones when possible

3. **Package Security**:
   - Sign releases with GPG
   - Use secure registries
   - Verify dependencies
   - Keep dependencies updated

## Contributing

Please read our [Contributing Guidelines](../CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
