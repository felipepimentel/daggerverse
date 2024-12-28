# Daggerverse Modules

A collection of high-quality Dagger modules designed to streamline your CI/CD workflows. Each module is crafted with best practices, security, and flexibility in mind.

## Available Modules

### [Python Module](python/README.md)

A comprehensive Python module that streamlines Poetry-based Python development workflows. Features include:

- 🏗️ Poetry integration and dependency management
- 🧪 Advanced testing with coverage reporting
- 📦 Automated package building and publishing
- 🔄 Built-in Git operations and authentication
- 🔍 Code quality tools (linting, formatting)
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
   ```go
   // See individual module documentation for detailed usage examples
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

We welcome contributions! Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on:

- Code of Conduct
- Development process
- Commit message format
- Pull request process

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security

For security concerns, please see our [Security Policy](SECURITY.md).
