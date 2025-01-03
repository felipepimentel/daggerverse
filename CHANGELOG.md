## [1.11.2](https://github.com/felipepimentel/daggerverse/compare/v1.11.1...v1.11.2) (2025-01-03)


### Bug Fixes

* **global:** use existing semantic-release config ([ef37245](https://github.com/felipepimentel/daggerverse/commit/ef372456dfb5009f493053ab522ec1731fb9e916))

## [1.11.1](https://github.com/felipepimentel/daggerverse/compare/v1.11.0...v1.11.1) (2025-01-03)


### Bug Fixes

* **global:** add poetry installation to ci workflow ([d635815](https://github.com/felipepimentel/daggerverse/commit/d6358157abc99b17fd6491d3a61e416028445d27))
* **global:** improve version tag handling ([aedb999](https://github.com/felipepimentel/daggerverse/commit/aedb999c748004ee4b0a263dc0831a9f04dab1ba))
* **global:** remove poetry steps from ci workflow ([5d96422](https://github.com/felipepimentel/daggerverse/commit/5d96422aba656d079a213f1459924270652bb349))
* **global:** remove poetry-specific release steps ([cd8f0d7](https://github.com/felipepimentel/daggerverse/commit/cd8f0d7dfebbe69478b7d0a537dd673a3115107f))

# [1.11.0](https://github.com/felipepimentel/daggerverse/compare/v1.10.0...v1.11.0) (2025-01-02)


### Features

* implement automatic version bumping in CICD function ([1d95be0](https://github.com/felipepimentel/daggerverse/commit/1d95be088899719efd21bf1ac33b19493b639abf))

# [1.10.0](https://github.com/felipepimentel/daggerverse/compare/v1.9.0...v1.10.0) (2025-01-02)


### Features

* enhance CICD function to support publishing to PyPI ([7015f65](https://github.com/felipepimentel/daggerverse/commit/7015f658bdc6a0fbb78c558eb00d83602ad17c0a))

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
