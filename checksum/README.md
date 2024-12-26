# Checksum Module for Dagger

A Dagger module that provides functionality to calculate and verify file checksums. This module supports SHA-256 hash algorithm for ensuring file integrity and verification.

## Features

- Calculate SHA-256 checksums for single or multiple files
- Verify file integrity using checksum files
- Simple and intuitive API
- Support for multiple files in a single operation
- Alpine-based lightweight implementation
- Parallel processing support in test suite

## Usage

### Basic Checksum Calculation

```typescript
import { checksum } from "@felipepimentel/daggerverse/checksum";

// Initialize the Checksum module
const client = checksum();

// Calculate SHA-256 checksums for files
const files = [
  dag.host().file("path/to/file1"),
  dag.host().file("path/to/file2"),
];
const checksums = await client.sha256().calculate(files);
```

### Checksum Verification

```typescript
import { checksum } from "@felipepimentel/daggerverse/checksum";

// Initialize the Checksum module
const client = checksum();

// Verify files against their checksums
const files = [
  dag.host().file("path/to/file1"),
  dag.host().file("path/to/file2"),
];
const result = await client.sha256().check(checksumFile, files);
```

## Examples

### Calculate and Verify Checksums

```typescript
import { checksum } from "@felipepimentel/daggerverse/checksum";

export async function calculateAndVerifyChecksums() {
  // Initialize files to check
  const files = [dag.host().file("file1.txt"), dag.host().file("file2.txt")];

  const client = checksum();

  // Calculate checksums
  const checksums = await client.sha256().calculate(files);

  // Verify checksums
  await client.sha256().check(checksums, files);
}
```

### Batch Processing

```typescript
import { checksum } from "@felipepimentel/daggerverse/checksum";

export async function batchChecksumProcessing() {
  // Get all files from a directory
  const sourceDir = dag.host().directory("./files");
  const files = await sourceDir.files();

  const client = checksum();

  // Calculate checksums for all files
  const checksums = await client.sha256().calculate(files);

  // Verify all files
  await client.sha256().check(checksums, files);
}
```

## API Reference

### Checksum

The main module interface:

```typescript
interface Checksum {
  // Initialize SHA-256 operations
  sha256(): Sha256;
}
```

### Sha256

SHA-256 specific operations:

```typescript
interface Sha256 {
  // Calculate SHA-256 checksums for given files
  calculate(files: File[]): Promise<File>;

  // Check files against their checksums
  check(checksums: File, files: File[]): Promise<Container>;
}
```

## Testing

The module includes a comprehensive test suite that can be run using:

```bash
dagger run test
```

The test suite includes:

- Basic checksum calculation verification
- File integrity checking
- Parallel processing tests
- Error handling verification

## Dependencies

The module requires:

- Go 1.22 or later
- Dagger SDK
- Alpine-based container (pulled automatically)
- Internet access for initial container pulling

## Implementation Details

- Uses Alpine Linux base image for lightweight operation
- Implements standard SHA-256 algorithm
- Supports parallel processing in test suite using the `sourcegraph/conc` package
- Provides both synchronous and asynchronous operation support

## License

See [LICENSE](../LICENSE) file in the root directory.
