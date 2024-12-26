# SSH-Keygen Module for Dagger

A Dagger module that provides functionality to generate SSH key pairs using various algorithms (Ed25519, ECDSA, RSA). This module enables secure key generation with optional passphrase protection and customizable key parameters.

## Features

- Multiple key generation algorithms:
  - Ed25519 (modern, recommended)
  - ECDSA (with 256, 384, or 521 bits)
  - RSA (customizable key size)
- Passphrase protection support
- Customizable key names
- OpenSSH format compatibility
- Secure random number generation
- PEM encoding for private keys
- Authorized keys format for public keys

## Usage

### Ed25519 Keys (Recommended)

```typescript
import { sshKeygen } from "@felipepimentel/daggerverse/ssh-keygen";

// Initialize the SSH-Keygen module
const keygen = sshKeygen();

// Generate an Ed25519 key pair
const keyPair = await keygen.ed25519().generate({
  name: "id_ed25519", // Optional: key name
  passphrase: null, // Optional: passphrase
});

// Access the keys
const publicKey = keyPair.publicKey; // File
const privateKey = keyPair.privateKey; // Secret
```

### ECDSA Keys

```typescript
// Initialize with specific bit size (256, 384, or 521)
const ecdsa = await keygen.ecdsa(256);

// Generate an ECDSA key pair
const keyPair = await ecdsa.generate({
  name: "id_ecdsa", // Optional: key name
  passphrase, // Optional: passphrase
});
```

### RSA Keys

```typescript
// Initialize with specific bit size (default: 4096)
const rsa = keygen.rsa(4096);

// Generate an RSA key pair
const keyPair = await rsa.generate({
  name: "id_rsa", // Optional: key name
  passphrase, // Optional: passphrase
});
```

## Configuration

### Ed25519 Options

The `ed25519().generate()` method accepts:

- `name`: Key name (optional, default: "id_ed25519")
- `passphrase`: Secret for private key encryption (optional)

### ECDSA Options

The `ecdsa()` constructor accepts:

- `bits`: Key size in bits (optional, default: 256)
  - Valid values: 256, 384, 521

The `generate()` method accepts:

- `name`: Key name (optional, default: "id_ecdsa")
- `passphrase`: Secret for private key encryption (optional)

### RSA Options

The `rsa()` constructor accepts:

- `bits`: Key size in bits (optional, default: 4096)

The `generate()` method accepts:

- `name`: Key name (optional, default: "id_rsa")
- `passphrase`: Secret for private key encryption (optional)

## Examples

### Generate Protected Key Pair

```typescript
import { sshKeygen } from "@felipepimentel/daggerverse/ssh-keygen";

export async function generateProtectedKey() {
  // Initialize module
  const keygen = sshKeygen();

  // Create passphrase
  const passphrase = dag.setSecret("passphrase", "your-secure-passphrase");

  // Generate Ed25519 key pair with passphrase
  const keyPair = await keygen.ed25519().generate({
    name: "deploy_key",
    passphrase,
  });

  // Use the keys
  const container = dag
    .container()
    .from("alpine:latest")
    .withMountedFile("/root/.ssh/id_ed25519.pub", keyPair.publicKey)
    .withMountedSecret("/root/.ssh/id_ed25519", keyPair.privateKey);
}
```

### Multiple Key Types

```typescript
import { sshKeygen } from "@felipepimentel/daggerverse/ssh-keygen";

export async function generateMultipleKeys() {
  const keygen = sshKeygen();

  // Ed25519 key
  const ed25519Key = await keygen.ed25519().generate({
    name: "id_ed25519",
  });

  // ECDSA key
  const ecdsaKey = await keygen.ecdsa(384).generate({
    name: "id_ecdsa",
  });

  // RSA key
  const rsaKey = await keygen.rsa(4096).generate({
    name: "id_rsa",
  });

  // Use the keys...
}
```

## Dependencies

The module requires:

- Dagger SDK
- Go crypto libraries
- OpenSSH compatibility

## Testing

The module includes tests that verify:

- Key pair generation for all algorithms
- Passphrase protection
- OpenSSH compatibility
- Key format validation

To run the tests:

```bash
dagger do test
```

## License

See [LICENSE](../LICENSE) file in the root directory.
