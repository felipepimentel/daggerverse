# OpenSSH Server Module for Dagger

A Dagger module that provides a configurable OpenSSH server for testing SSH connections and integrations. This module is built on top of Wolfi Linux and provides a lightweight, secure SSH server implementation.

## Features

- Easy setup of OpenSSH server instances
- Configurable port (default: 22)
- Custom container support
- Automatic host key generation
- Secure default configuration
- Integration with other Dagger modules
- Support for custom SSH configurations
- Built on Wolfi Linux base

## Usage

### Basic Setup

```typescript
import { opensshServer } from "@felipepimentel/daggerverse/openssh-server";

// Initialize a basic OpenSSH server
const client = opensshServer();

// Get the server as a service
const service = client.service();
```

### Custom Configuration

```typescript
// Initialize with custom port
const withPort = client.withPort(2222);

// Use a custom base container
const customContainer = dag
  .container()
  .from("custom-image:latest")
  .withExec(["apk", "add", "openssh-server"]);

const client = opensshServer({
  container: customContainer,
});
```

### Integration with SSH Keys

```typescript
// Generate SSH keys using the ssh-keygen module
const keys = dag.sshKeygen().ed25519();
const publicKey = await keys.publicKey();
const privateKey = await keys.privateKey();

// Configure the server with the public key
const withKey = client.withAuthorizedKey(publicKey);
```

## Configuration Options

### Port Configuration

```typescript
const withPort = client.withPort(2222);
```

### Custom Base Container

```typescript
const container = dag
  .container()
  .from("wolfi-base:latest")
  .withPackages(["openssh-server"]);

const client = opensshServer({
  container: container,
});
```

### SSH Configuration

The module uses a secure default configuration located in `etc/sshd_config`. You can customize it by:

```typescript
const withConfig = client.withConfig(
  dag
    .directory()
    .withFile("sshd_config", myConfig)
    .withDirectory("sshd_config.d", customConfigs)
);
```

## Examples

### Basic SSH Server Setup

```typescript
import { opensshServer } from "@felipepimentel/daggerverse/openssh-server";

export async function setupSSHServer() {
  // Initialize server
  const client = opensshServer();

  // Get the server as a service
  const service = client.service();
}
```

### SSH Server with Custom Configuration

```typescript
import { opensshServer } from "@felipepimentel/daggerverse/openssh-server";

export async function setupCustomSSHServer() {
  // Create custom configuration
  const config = dag
    .directory()
    .withFile(
      "custom.conf",
      dag.currentModule().source().file("configs/custom.conf")
    );

  // Initialize server with custom config
  const client = opensshServer().withPort(2222).withConfig(config);
}
```

### Integration with Git Server

```typescript
import { opensshServer } from "@felipepimentel/daggerverse/openssh-server";

export async function setupGitServer() {
  // Create a custom container with Git
  const container = dag
    .container()
    .from("wolfi-base:latest")
    .withPackages(["openssh-server", "git"]);

  // Initialize SSH server
  const client = opensshServer({
    container: container,
  }).withPort(22);
}
```

## Testing

The module includes a comprehensive test suite that can be run using:

```bash
dagger do test
```

The test suite includes:

- Basic server functionality tests
- Custom configuration tests
- SSH key authentication tests
- Service binding tests
- Integration tests with other modules

## Dependencies

The module requires:

- Dagger SDK
- Wolfi Linux base image (pulled automatically)
- OpenSSH server package
- Internet access for initial container pulling

## Implementation Details

- Uses Wolfi Linux as the base distribution
- Automatically generates host keys on startup
- Provides secure default SSH configuration
- Supports custom configuration through `sshd_config.d`
- Implements privilege separation
- Integrates with Dagger's service binding system

## License

See [LICENSE](../LICENSE) file in the root directory.
