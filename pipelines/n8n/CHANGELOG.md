# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2024-01-29

### Added
- Reusable module support for DigitalOcean, Docker, and SSH operations
- Configuration methods for region, size, and image customization
- Improved error handling and logging
- Comprehensive documentation in README.md

### Changed
- Refactored code to use reusable modules instead of direct implementations
- Simplified deployment process by removing unnecessary complexity
- Updated module structure to follow Daggerverse standards
- Improved SSH key handling with dedicated module

### Removed
- Unused configuration options and methods
- Direct Docker and SSH command execution
- Complex validation checks

## [0.1.0] - 2024-01-26

### Added
- Initial implementation of n8n deployment pipeline
- Basic DigitalOcean integration
- Docker and Docker Compose support
- Caddy reverse proxy configuration
- DNS management
- SSH key handling
- Configuration file generation 