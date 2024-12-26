# Mod-Releaser Module for Dagger

A Dagger module that automates the release process for Dagger modules in your repository. This module helps manage versioning, tagging, and publishing of Dagger modules using semantic versioning.

## Features

- Automated semantic versioning
- Git tag management
- Module publishing to Dagger registry
- Git configuration support
- SSH key integration
- Branch selection
- Custom tag messages
- Version bumping (major, minor, patch)

## Usage

### Basic Setup

```go
// Initialize the module releaser
releaser, err := dag.ModReleaser().New(
    ctx,
    gitRepo,        // *dagger.Directory - Git repository
    "my-module",    // Module name to publish
)
```

### Version Management

```go
// Bump major version
releaser, err = releaser.Major("Release with breaking changes")

// Bump minor version
releaser, err = releaser.Minor("Release with new features")

// Bump patch version
releaser, err = releaser.Patch("Release with bug fixes")

// Publish the release
releaser = releaser.Publish(true)  // true to push git tag
```

### Custom Configuration

```go
// Configure git settings
releaser = releaser.WithGitConfig(
    gitConfig,              // Git config file
    "user@example.com",     // Git user email
    "User Name",            // Git user name
)

// Set up SSH keys
releaser = releaser.WithSshKeys(sshKeys)

// Select a branch
releaser = releaser.WithBranch("main")
```

## Configuration Options

### Constructor Parameters

- `gitRepo`: Git repository directory
- `component`: Name of the module to publish

### Git Configuration

- `cfg` (optional): Git config file
- `email` (optional): Git user email
- `name` (optional): Git user name

### Version Control

- Major version: Breaking changes
- Minor version: New features, backward compatible
- Patch version: Bug fixes, backward compatible

### Publishing Options

- `gitPush`: Whether to push git tags
- Custom tag messages for each version bump

## Implementation Details

### Version Management

The module:

- Automatically detects existing versions
- Follows semantic versioning rules
- Creates annotated git tags
- Manages version increments

### Git Integration

- Supports SSH key authentication
- Handles git configuration
- Manages git tags and branches
- Supports custom git messages

### Module Publishing

- Integrates with Dagger registry
- Handles module metadata
- Manages dependencies
- Supports version tagging

## Dependencies

The module requires:

- Dagger SDK
- Git
- SSH (for authentication)
- Internet access for publishing
- Valid Dagger module configuration

## Examples

### Complete Release Process

```go
// Initialize releaser
releaser, err := dag.ModReleaser().New(ctx, gitRepo, "my-module")
if err != nil {
    return err
}

// Configure git
releaser = releaser.WithGitConfig(
    nil,
    "ci@example.com",
    "CI Bot",
)

// Set up SSH keys
releaser = releaser.WithSshKeys(sshKeys)

// Create new minor version
releaser, err = releaser.Minor("New feature release")
if err != nil {
    return err
}

// Publish the module
releaser = releaser.Publish(true)

// Execute the release
_, err = releaser.Do(ctx)
return err
```

### Custom Branch Release

```go
// Initialize and configure
releaser, err := dag.ModReleaser().New(ctx, gitRepo, "my-module")
if err != nil {
    return err
}

// Switch to feature branch
releaser = releaser.WithBranch("feature/new-version")

// Create patch release
releaser, err = releaser.Patch("Bug fix in feature branch")
if err != nil {
    return err
}

// Publish without pushing git tag
releaser = releaser.Publish(false)
```

## License

This module is part of the Dagger project and is subject to its licensing terms.

Copyright Camptocamp SA
SPDX-License-Identifier: AGPL-3.0-or-later
