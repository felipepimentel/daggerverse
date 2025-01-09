---
layout: default
title: TTLSH Module
parent: Essentials
nav_order: 16
---

# TTLSH Module

The TTLSH module provides functionality for running shell commands with time-to-live (TTL) constraints in your Dagger pipelines. It allows you to execute commands that automatically terminate after a specified duration.

## Features

- Command execution with TTL
- Timeout management
- Signal handling
- Process cleanup
- Exit code handling
- Output capture
- Error reporting
- Alpine-based execution
- Resource cleanup
- Graceful termination

## Installation

To use the TTLSH module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/ttlsh"
)
```

## Usage Examples

### Basic Command Execution

```go
func (m *MyModule) Example(ctx context.Context) error {
    ttlsh := dag.Ttlsh()
    
    // Run command with default TTL
    return ttlsh.Run(
        ctx,
        []string{"echo", "Hello, World!"},
        "",           // default TTL
    )
}
```

### Custom TTL Duration

```go
func (m *MyModule) CustomTTL(ctx context.Context) error {
    ttlsh := dag.Ttlsh()
    
    // Run command with 30s TTL
    return ttlsh.Run(
        ctx,
        []string{"sleep", "10"},
        "30s",        // TTL duration
    )
}
```

### Long-Running Process

```go
func (m *MyModule) LongProcess(ctx context.Context) error {
    ttlsh := dag.Ttlsh()
    
    // Run long process with timeout
    return ttlsh.Run(
        ctx,
        []string{
            "sh", "-c",
            "for i in $(seq 1 10); do echo $i; sleep 1; done",
        },
        "5s",         // TTL duration
    )
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: TTL Command
on: [push]

jobs:
  run:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Command
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/ttlsh
          args: |
            do -p '
              ttlsh := Ttlsh()
              ttlsh.Run(
                ctx,
                []string{"echo", "Running in CI"},
                "10s",
              )
            '
```

## API Reference

### Ttlsh

Main module struct that provides access to TTL shell functionality.

#### Methods

- `Run(ctx context.Context, command []string, ttl string) error`
  - Executes command with TTL
  - Parameters:
    - `command`: Command and arguments to execute
    - `ttl`: Time-to-live duration (e.g., "30s", "1m")
  - Returns error if command fails or times out

## Best Practices

1. **TTL Management**
   - Set appropriate timeouts
   - Handle termination
   - Clean up resources

2. **Command Execution**
   - Use proper signals
   - Handle exit codes
   - Capture output

3. **Error Handling**
   - Check timeouts
   - Handle failures
   - Log issues

4. **Resource Usage**
   - Monitor processes
   - Clean temporary files
   - Manage memory

## Troubleshooting

Common issues and solutions:

1. **Timeout Issues**
   ```
   Error: command timed out
   Solution: Increase TTL duration
   ```

2. **Process Problems**
   ```
   Error: process terminated
   Solution: Check command validity
   ```

3. **Resource Errors**
   ```
   Error: resource exhausted
   Solution: Monitor resource usage
   ```

## Configuration Example

```yaml
# ttlsh-config.yaml
defaults:
  ttl: "60s"
  signal: "SIGTERM"
  grace_period: "5s"
  
logging:
  level: info
  format: json
  
cleanup:
  enabled: true
  temp_files: true
  processes: true
```

## Advanced Usage

### Process Monitoring

```go
func (m *MyModule) MonitorProcess(ctx context.Context) error {
    ttlsh := dag.Ttlsh()
    
    // Run process with monitoring
    return ttlsh.Run(
        ctx,
        []string{
            "sh", "-c",
            `
            # Start main process
            while true; do
                echo "Running..."
                sleep 1
            done &
            
            # Store PID
            echo $! > /tmp/process.pid
            
            # Monitor process
            while kill -0 $(cat /tmp/process.pid) 2>/dev/null; do
                echo "Process is alive"
                sleep 1
            done
            `,
        },
        "15s",
    )
}
```

### Cleanup Handler

```go
func (m *MyModule) CleanupHandler(ctx context.Context) error {
    ttlsh := dag.Ttlsh()
    
    // Run with cleanup
    return ttlsh.Run(
        ctx,
        []string{
            "sh", "-c",
            `
            # Create temporary files
            echo "Creating temp files..."
            mkdir -p /tmp/work
            touch /tmp/work/file1.txt
            touch /tmp/work/file2.txt
            
            # Set trap for cleanup
            trap 'echo "Cleaning up..."; rm -rf /tmp/work' EXIT
            
            # Main process
            echo "Working..."
            sleep 10
            
            # Normal exit triggers cleanup
            echo "Done"
            `,
        },
        "30s",
    )
}
```

### Resource Monitoring

```go
func (m *MyModule) ResourceMonitor(ctx context.Context) error {
    ttlsh := dag.Ttlsh()
    
    // Run with resource monitoring
    return ttlsh.Run(
        ctx,
        []string{
            "sh", "-c",
            `
            # Start resource monitoring
            (
                while true; do
                    echo "Memory usage:"
                    free -h
                    echo "CPU usage:"
                    top -bn1 | head -n 3
                    sleep 5
                done
            ) &
            
            # Main process
            echo "Starting work..."
            for i in $(seq 1 5); do
                echo "Processing batch $i"
                # Simulate CPU work
                yes > /dev/null &
                sleep 2
                killall yes
            done
            
            echo "Work completed"
            `,
        },
        "1m",
    )
} 