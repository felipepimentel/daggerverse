# Daggerverse Modules

A collection of high-quality Dagger modules designed to streamline your CI/CD workflows. Each module is crafted with best practices, security, and flexibility in mind.

## Available Modules

### [Python Module](python/README.md)

A comprehensive Python module that streamlines Poetry-based Python development workflows. Features include:

- 🔄 Automatic version management with semantic-release
- 🔍 Automatic project structure detection
- 📦 Poetry integration and dependency management
- 🧪 Advanced testing with coverage reporting
- 🚀 Automated package building and publishing
- 🔐 Built-in Git operations and authentication
- 🎨 Code quality tools (linting, formatting)
- 📚 Documentation generation
- 💾 Optimized caching

## Getting Started

1. **Prerequisites**:

   - [Dagger CLI](https://docs.dagger.io/cli/465058/install)
   - Go 1.21 or later

2. **Installation**:

   ```bash
   # Import the desired module
   dagger mod use github.com/felipepimentel/daggerverse/python@latest
   ```

3. **Usage**:
   ```bash
   # Run complete CI/CD pipeline
   dagger -m github.com/felipepimentel/daggerverse/python call cicd --source . --token env:PYPI_TOKEN
   ```

## Module Structure

Each module follows a consistent structure:

```
module/
├── README.md           # Module documentation
├── dagger.json         # Module configuration
├── main.go            # Module implementation
└── examples/          # Usage examples
```

## Development

### Requirements

- Go 1.21+
- Dagger CLI
- Git

### Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/felipepimentel/daggerverse.git
   cd daggerverse
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   dagger test
   ```

## Contributing

We welcome contributions! Please follow our commit message format:

```
type(scope): subject

[optional body]

[optional footer(s)]
```

### Types

- `feat`: A new feature (minor version)
- `fix`: A bug fix (patch version)
- `perf`: A code change that improves performance (patch version)
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to our CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files

### Breaking Changes

Breaking changes must be indicated by:

1. `!` after the type/scope
2. `BREAKING CHANGE:` in the footer

Example:

```
feat(python)!: change module API interface

BREAKING CHANGE: The BuildConfig interface now requires explicit cache configuration.
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security

For security concerns, please see our [Security Policy](SECURITY.md).
