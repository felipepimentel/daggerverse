---
layout: default
title: Make Module
parent: Essentials
nav_order: 10
---

# Make Module

The Make module provides integration with GNU Make in your Dagger pipelines. It allows you to execute Makefiles and build targets in a containerized environment.

## Features

- Makefile execution
- Custom target support
- Alpine-based container
- Directory mounting
- Argument passing
- Custom Makefile paths
- Build automation
- Error handling
- Output capture
- Directory modification

## Installation

To use the Make module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/make"
)
```

## Usage Examples

### Basic Make Execution

```go
func (m *MyModule) Example(ctx context.Context) (*Directory, error) {
    make := dag.Make()
    
    // Execute default target
    return make.Make(
        dag.Directory("."),  // source directory
        []string{},         // no additional arguments
        "",                 // use default Makefile
    ), nil
}
```

### Custom Target and Arguments

```go
func (m *MyModule) BuildTarget(ctx context.Context) (*Directory, error) {
    make := dag.Make()
    
    // Execute specific target with arguments
    return make.Make(
        dag.Directory("."),
        []string{
            "build",
            "ARGS=-v",
            "DEBUG=1",
        },
        "",
    ), nil
}
```

### Custom Makefile Path

```go
func (m *MyModule) CustomMakefile(ctx context.Context) (*Directory, error) {
    make := dag.Make()
    
    // Use custom Makefile
    return make.Make(
        dag.Directory("."),
        []string{"test"},
        "build/Makefile",
    ), nil
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Make Build
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build Project
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/make
          args: |
            do -p '
              make := Make()
              make.Make(
                Directory("."),
                []string{"build"},
                "",
              )
            '
```

## API Reference

### Make

Main module struct that provides access to Make functionality.

#### Methods

- `Make(dir *Directory, args []string, makefile string) *Directory`
  - Executes make command
  - Parameters:
    - `dir`: Source directory
    - `args`: Make arguments and targets
    - `makefile`: Custom Makefile path (optional, defaults to "Makefile")
  - Returns modified directory

## Best Practices

1. **Makefile Management**
   - Keep Makefiles simple
   - Document targets
   - Use variables

2. **Build Process**
   - Define dependencies
   - Handle errors
   - Clean targets

3. **Directory Structure**
   - Organize source files
   - Manage outputs
   - Clean artifacts

4. **Integration**
   - Automate builds
   - Version control
   - Document process

## Troubleshooting

Common issues and solutions:

1. **Makefile Issues**
   ```
   Error: Makefile not found
   Solution: Verify file path
   ```

2. **Target Problems**
   ```
   Error: target not found
   Solution: Check target name
   ```

3. **Build Errors**
   ```
   Error: build failed
   Solution: Check dependencies
   ```

## Configuration Example

```makefile
# Makefile
.PHONY: all build test clean

# Variables
BUILD_DIR = build
DEBUG = 0
ARGS =

# Default target
all: build test

# Build target
build:
    @echo "Building with args: $(ARGS)"
    @mkdir -p $(BUILD_DIR)
    @if [ "$(DEBUG)" = "1" ]; then \
        echo "Debug mode enabled"; \
    fi
    # Build commands here

# Test target
test:
    @echo "Running tests"
    # Test commands here

# Clean target
clean:
    @echo "Cleaning build directory"
    @rm -rf $(BUILD_DIR)
```

## Advanced Usage

### Multi-Stage Build

```go
func (m *MyModule) MultistageBuild(ctx context.Context) error {
    make := dag.Make()
    
    // Build stage
    buildDir := make.Make(
        dag.Directory("."),
        []string{
            "build",
            "BUILD_TYPE=release",
        },
        "",
    )
    
    // Test stage
    testDir := make.Make(
        buildDir,
        []string{
            "test",
            "TEST_ARGS=--verbose",
        },
        "",
    )
    
    // Deploy stage
    return make.Make(
        testDir,
        []string{
            "deploy",
            "ENV=production",
        },
        "",
    ).Sync(ctx)
}
```

### Parallel Builds

```go
func (m *MyModule) ParallelBuild(ctx context.Context) error {
    make := dag.Make()
    
    // Execute multiple targets in parallel
    return make.Make(
        dag.Directory("."),
        []string{
            "-j",           // enable parallel jobs
            "all",          // build all targets
            "-k",           // keep going on error
            "PARALLEL=1",   // enable parallel flag
        },
        "",
    ).Sync(ctx)
}
```

### Custom Environment

```go
func (m *MyModule) CustomEnv(ctx context.Context) (*Directory, error) {
    make := dag.Make()
    
    // Set custom environment variables
    env := map[string]string{
        "CC": "gcc",
        "CFLAGS": "-O2",
        "PREFIX": "/usr/local",
    }
    
    // Convert env to make arguments
    var args []string
    for k, v := range env {
        args = append(args, k+"="+v)
    }
    args = append(args, "install")
    
    // Execute make with custom environment
    return make.Make(
        dag.Directory("."),
        args,
        "",
    ), nil
}
``` 