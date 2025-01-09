---
layout: default
title: SSH Module
parent: Essentials
nav_order: 14
---

# SSH Module

The SSH module provides functionality for managing SSH connections and operations in your Dagger pipelines. It allows you to handle SSH keys, execute remote commands, and manage SSH configurations.

## Features

- SSH key management
- Remote command execution
- Key generation
- Configuration handling
- Known hosts management
- Agent forwarding
- Port forwarding
- Connection testing
- Error handling
- Secure operations

## Installation

To use the SSH module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/ssh"
)
```

## Usage Examples

### Basic SSH Connection

```go
func (m *MyModule) Example(ctx context.Context) error {
    ssh := dag.SSH()
    
    // Configure SSH connection
    ssh = ssh.WithPrivateKey(
        dag.SetSecret("SSH_KEY", "private-key-content"),
    ).WithKnownHosts(
        dag.SetSecret("KNOWN_HOSTS", "known-hosts-content"),
    )
    
    // Execute remote command
    return ssh.Run(
        ctx,
        "user@example.com",
        []string{"ls", "-la"},
    )
}
```

### Custom Port and Config

```go
func (m *MyModule) CustomConfig(ctx context.Context) error {
    ssh := dag.SSH()
    
    // Configure SSH with custom settings
    ssh = ssh.WithPrivateKey(
        dag.SetSecret("SSH_KEY", "private-key-content"),
    ).WithPort(2222).WithConfig(
        dag.Directory(".").File("ssh_config"),
    )
    
    // Execute command
    return ssh.Run(
        ctx,
        "user@example.com",
        []string{"whoami"},
    )
}
```

### Multiple Commands

```go
func (m *MyModule) MultiCommand(ctx context.Context) error {
    ssh := dag.SSH()
    
    // Configure SSH
    ssh = ssh.WithPrivateKey(
        dag.SetSecret("SSH_KEY", "private-key-content"),
    )
    
    // Execute multiple commands
    commands := [][]string{
        {"mkdir", "-p", "/tmp/test"},
        {"cd", "/tmp/test", "&&", "touch", "file.txt"},
        {"ls", "-la", "/tmp/test"},
    }
    
    for _, cmd := range commands {
        if err := ssh.Run(ctx, "user@example.com", cmd); err != nil {
            return err
        }
    }
    
    return nil
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: Remote Deploy
on: [push]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: SSH Command
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/ssh
          args: |
            do -p '
              ssh := SSH().
                WithPrivateKey(
                  SetSecret("SSH_KEY", "${{ secrets.SSH_PRIVATE_KEY }}"),
                ).
                WithKnownHosts(
                  SetSecret("KNOWN_HOSTS", "${{ secrets.KNOWN_HOSTS }}"),
                )
              ssh.Run(
                ctx,
                "user@example.com",
                []string{"deploy.sh"},
              )
            '
```

## API Reference

### SSH

Main module struct that provides access to SSH functionality.

#### Methods

- `WithPrivateKey(key *Secret) *SSH`
  - Sets SSH private key
  - Parameters:
    - `key`: Private key secret

- `WithKnownHosts(hosts *Secret) *SSH`
  - Sets known hosts file
  - Parameters:
    - `hosts`: Known hosts content

- `WithConfig(config *File) *SSH`
  - Sets SSH config file
  - Parameters:
    - `config`: SSH configuration file

- `WithPort(port int) *SSH`
  - Sets SSH port
  - Parameters:
    - `port`: Port number

- `Run(ctx context.Context, host string, command []string) error`
  - Executes remote command
  - Parameters:
    - `host`: Remote host
    - `command`: Command to execute

## Best Practices

1. **Key Management**
   - Secure storage
   - Regular rotation
   - Access control

2. **Configuration**
   - Use config files
   - Document settings
   - Version control

3. **Security**
   - Verify hosts
   - Limit access
   - Audit logs

4. **Operations**
   - Test connections
   - Handle timeouts
   - Error recovery

## Troubleshooting

Common issues and solutions:

1. **Connection Issues**
   ```
   Error: connection refused
   Solution: Check host and port
   ```

2. **Authentication Problems**
   ```
   Error: permission denied
   Solution: Verify key permissions
   ```

3. **Host Verification**
   ```
   Error: host key verification failed
   Solution: Update known hosts
   ```

## Configuration Example

```
# ssh_config
Host example
    HostName example.com
    User deploy
    Port 2222
    IdentityFile ~/.ssh/id_rsa
    StrictHostKeyChecking yes
    UserKnownHostsFile ~/.ssh/known_hosts
    ForwardAgent yes
```

## Advanced Usage

### Custom SSH Agent

```go
func (m *MyModule) CustomAgent(ctx context.Context) error {
    ssh := dag.SSH()
    
    // Configure SSH with agent
    ssh = ssh.WithPrivateKey(
        dag.SetSecret("SSH_KEY", "private-key-content"),
    ).WithConfig(
        dag.Directory(".").File("ssh_config"),
    )
    
    // Start agent and add key
    return dag.Container().
        From("alpine:latest").
        WithExec([]string{"apk", "add", "openssh-client"}).
        WithExec([]string{
            "sh", "-c",
            `
            # Start SSH agent
            eval $(ssh-agent)
            
            # Add private key
            echo "$SSH_KEY" > /tmp/key
            chmod 600 /tmp/key
            ssh-add /tmp/key
            
            # Run SSH command
            ssh -A user@example.com "ls -la"
            `,
        }).
        WithSecretVariable("SSH_KEY", dag.SetSecret("SSH_KEY", "private-key-content")).
        Sync(ctx)
}
```

### Secure File Transfer

```go
func (m *MyModule) SecureTransfer(ctx context.Context) error {
    ssh := dag.SSH()
    
    // Configure SSH
    ssh = ssh.WithPrivateKey(
        dag.SetSecret("SSH_KEY", "private-key-content"),
    ).WithKnownHosts(
        dag.SetSecret("KNOWN_HOSTS", "known-hosts-content"),
    )
    
    // Create test file
    file := dag.Container().
        From("alpine:latest").
        WithNewFile("/test.txt", dagger.ContainerWithNewFileOpts{
            Contents: "test content",
        }).
        File("/test.txt")
    
    // Transfer file using scp
    return dag.Container().
        From("alpine:latest").
        WithMountedFile("/source/test.txt", file).
        WithMountedSecret("/root/.ssh/id_rsa", dag.SetSecret("SSH_KEY", "private-key-content")).
        WithMountedSecret("/root/.ssh/known_hosts", dag.SetSecret("KNOWN_HOSTS", "known-hosts-content")).
        WithExec([]string{
            "sh", "-c",
            `
            # Set key permissions
            chmod 600 /root/.ssh/id_rsa
            
            # Transfer file
            scp -i /root/.ssh/id_rsa /source/test.txt user@example.com:/destination/
            `,
        }).
        Sync(ctx)
} 