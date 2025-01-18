---
layout: default
title: Docusaurus Module
parent: Libraries
nav_order: 3
---

# Docusaurus Module

The Docusaurus module provides integration with [Docusaurus](https://docusaurus.io/), a modern static website generator. This module allows you to build, serve, and develop Docusaurus documentation sites in your Dagger pipelines.

## Features

- Build production documentation
- Serve production builds
- Development server with hot reload
- NPM and Yarn support
- Cache management
- Custom working directory support
- Container and service modes

## Installation

To use the Docusaurus module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/docusaurus"
)
```

## Usage Examples

### Basic Documentation Build

```go
func (m *MyModule) Example(ctx context.Context) (*Directory, error) {
    docusaurus := dag.Docusaurus().New(
        dag.Directory("./docs"),  // source directory
        "/src",                  // working directory
        false,                   // disable cache
        "node-docusaurus-docs",  // cache volume name
        false,                   // use npm (not yarn)
    )
    
    return docusaurus.Build(), nil
}
```

### Development Server

```go
func (m *MyModule) DevServer(ctx context.Context) (*Service, error) {
    docusaurus := dag.Docusaurus().New(
        dag.Directory("./docs"),
        "/src",
        false,
        "node-docusaurus-docs",
        false,
    )
    
    return docusaurus.ServeDev(), nil
}
```

### Production Server

```go
func (m *MyModule) ProdServer(ctx context.Context) (*Service, error) {
    docusaurus := dag.Docusaurus().New(
        dag.Directory("./docs"),
        "/src",
        false,
        "node-docusaurus-docs",
        false,
    )
    
    return docusaurus.Serve(), nil
}
```

### Using Yarn

```go
func (m *MyModule) WithYarn(ctx context.Context) (*Directory, error) {
    docusaurus := dag.Docusaurus().New(
        dag.Directory("./docs"),
        "/src",
        false,
        "node-docusaurus-docs",
        true,  // use yarn
    )
    
    return docusaurus.Build(), nil
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Documentation
on: [push]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build Documentation
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/docusaurus
          args: |
            do -p '
              docusaurus := Docusaurus().New(
                dag.Directory("./docs"),
                "/src",
                false,
                "node-docusaurus-docs",
                false,
              )
              docusaurus.Build()
            '
```

## API Reference

### Docusaurus

Main module struct that provides access to Docusaurus functionality.

#### Constructor

- `New(src *Directory, dir string, disableCache bool, cacheVolumeName string, yarn bool) *Docusaurus`
  - Creates a new Docusaurus instance
  - Parameters:
    - `src`: Source directory containing Docusaurus site
    - `dir`: Working directory (optional, default: "/src")
    - `disableCache`: Disable caching (optional, default: false)
    - `cacheVolumeName`: Cache volume name (optional, default: "node-docusaurus-docs")
    - `yarn`: Use Yarn instead of NPM (optional, default: false)

#### Methods

- `Base() *Container`
  - Returns base container with Docusaurus dependencies installed
  
- `Build() *Directory`
  - Builds production documentation
  - Returns directory containing built site
  
- `Serve() *Service`
  - Serves production documentation
  - Returns service running on port 3000
  
- `ServeDev() *Service`
  - Serves development documentation with hot reload
  - Returns service running on port 3000

## Best Practices

1. **Cache Management**
   - Use cache for faster builds
   - Use unique cache volume names for multiple projects
   - Clear cache when dependencies change

2. **Development Workflow**
   - Use `ServeDev()` for local development
   - Use `Build()` for production builds
   - Use `Serve()` for testing production builds

3. **Resource Management**
   - Monitor memory usage with large sites
   - Clean up unused services
   - Use appropriate cache strategies

4. **Package Management**
   - Choose between NPM and Yarn based on project needs
   - Keep package.json up to date
   - Handle dependencies properly

## Troubleshooting

Common issues and solutions:

1. **Build Failures**
   ```
   Error: Failed to compile
   Solution: Check for syntax errors in MDX files
   ```

2. **Cache Issues**
   ```
   Error: npm ERR! Cannot read property 'matches' of undefined
   Solution: Try clearing the cache or using a different cache volume name
   ```

3. **Port Conflicts**
   ```
   Error: Port 3000 is already in use
   Solution: Stop other services using port 3000 or configure a different port
   ```

## Cache Volumes

The module uses several cache volumes for optimization:

1. `{cacheVolumeName}`: For node_modules
2. `{cacheVolumeName}-build`: For build output
3. `node-docusaurus-root`: For NPM cache
4. `node-docusaurus-root-yarn`: For Yarn cache

## Advanced Usage

### Custom Base Container

```go
container := docusaurus.Base().
    WithEnvVariable("NODE_ENV", "production").
    WithMountedDirectory("/custom", customDir)
```

### Custom Build Configuration

```go
func (m *MyModule) CustomBuild(ctx context.Context) (*Directory, error) {
    docusaurus := dag.Docusaurus().New(
        dag.Directory("./docs"),
        "/src",
        false,
        "node-docusaurus-docs-" + projectName,  // project-specific cache
        true,
    )
    
    // Add custom environment variables
    base := docusaurus.Base().
        WithEnvVariable("CUSTOM_VAR", "value")
    
    // Build with custom configuration
    return base.
        WithExec([]string{"yarn", "build", "--config", "custom.config.js"}).
        Directory("build"), nil
}
``` 