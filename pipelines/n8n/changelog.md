# Changelog

## [0.0.1] - 2024-03-26

### Added
- Initial setup of the n8n deployment pipeline
- Basic structure for deploying to DigitalOcean using Dagger
- Environment variables configuration for n8n, PostgreSQL, and Caddy
- Docker Compose configuration for running the services
- Caddy configuration for SSL/TLS and reverse proxy
- Added base64 decoding support for SSH private key
- Added wait for system updates before Docker installation

### Changed
- Removed previous deployment artifacts (droplet and DNS records)
- Reset deployment state for fresh start
- Updated SSH key handling to use proper container-based approach
- Added SSH private key parameter to Deploy function
- Improved SSH command execution with proper key mounting and permissions
- Fixed SSH known_hosts handling
- Changed SSH key handling to use WithNewFile instead of WithMountedSecret to handle read-only filesystem
- Updated SSH key permissions handling to use file creation options
- Fixed WithNewFile usage to match the correct API signature
- Removed unused sshKey variable and simplified SSH setup
- Fixed ssh-keyscan command to properly redirect output
- Added proper trimming of SSH private key content
- Enhanced SSH key formatting with proper line endings and base64 decoding
- Fixed docker-compose.yml version string format
- Added wait for apt-get lock before Docker installation

### Fixed
- SSH authentication issues by properly mounting and configuring SSH keys
- Command execution issues by using direct SSH commands instead of high-level SSH client
- Read-only filesystem issues when setting SSH key permissions
- Compilation errors with WithNewFile usage
- Unused variable warning
- SSH known_hosts file creation issue
- SSH key formatting issues
- Base64 decoding issues with SSH private key
- Docker Compose version string format issue
- Docker installation failing due to apt-get lock 