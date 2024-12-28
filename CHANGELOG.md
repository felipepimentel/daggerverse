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
