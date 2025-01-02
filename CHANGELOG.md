# [1.9.0](https://github.com/felipepimentel/daggerverse/compare/v1.8.0...v1.9.0) (2025-01-02)


### Features

* implement CICD function for Python projects in main.go ([8aebac7](https://github.com/felipepimentel/daggerverse/commit/8aebac7c0d9d270985000c95d7834deb79f4ff85))

# [1.8.0](https://github.com/felipepimentel/daggerverse/compare/v1.7.1...v1.8.0) (2025-01-01)


### Features

* **python:** enhance CI/CD functionality and project structure detection ([4052cf9](https://github.com/felipepimentel/daggerverse/commit/4052cf983aaf7d13d92c0964457acd89c86bc97e))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Support for code linting with ruff
- Support for code formatting with black
- Support for documentation generation with Sphinx/MkDocs
- Comprehensive test configuration options
- Advanced coverage reporting features
- Poetry dependency management features
- Build environment customization options

### Fixed

- Replace `${process.env.MODULE_PATH}` with proper environment variable handling
- Map type replacement with KeyValue struct for Dagger compatibility
- Float64 type replacement with int for coverage percentage

### Changed

- Update to use latest Dagger commands
- Improve configuration defaults for better usability
- Enhanced error handling and logging

## [0.1.0] - 2024-01-17

### Added

- Initial release with basic Python module functionality
- Poetry integration for package management
- PyPI publishing support
- Basic test execution with pytest
- Simple build environment setup
