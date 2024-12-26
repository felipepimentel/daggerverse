# SLSA Verifier Module for Dagger

A Dagger module that provides integration with SLSA (Supply-chain Levels for Software Artifacts) verifier, enabling verification of provenance for software artifacts. This module helps ensure the integrity and authenticity of your software supply chain by verifying SLSA provenance statements.

## Features

- Automated SLSA verifier binary management
- Automatic version detection and download
- Provenance verification for multiple artifacts
- Support for various verification options:
  - Source repository validation
  - Builder ID verification
  - Branch and tag validation
  - Semantic version matching

## Usage

### Basic Verification

```typescript
import { slsaVerifier } from "@felipepimentel/daggerverse/slsa-verifier";

// Initialize SLSA verifier with default settings (latest version)
const verifier = slsaVerifier("");

// Verify artifacts with provenance
const result = await verifier.verifyArtifact(
  artifacts, // List of artifact files to verify
  provenanceFile, // Provenance file
  "github.com/org/repo" // Source repository URI
);
```

### Custom Version and Options

```typescript
// Initialize with specific version
const verifier = slsaVerifier("1.2.3");

// Verify with additional options
const result = await verifier.verifyArtifact(
  artifacts,
  provenanceFile,
  "github.com/org/repo",
  "github-actions", // Builder ID
  "main", // Source branch
  "v1.0.0", // Source tag
  "1.0.0" // Source versioned tag
);
```

## Configuration Options

### Version

- Specifies the version of SLSA verifier to use
- Optional: Defaults to latest version from GitHub releases
- Format: Semantic version (e.g., "1.2.3")

### Verification Parameters

#### Required

- `artifacts`: List of artifact files to verify
- `provenance`: Provenance file containing SLSA metadata
- `sourceURI`: Expected source repository URI

#### Optional

- `builderID`: The unique builder ID who created the provenance
- `sourceBranch`: Expected branch the binary was compiled from
- `sourceTag`: Expected tag the binary was compiled from
- `sourceVersionedTag`: Expected version using semantic version matching

## Dependencies

The module requires:

- Dagger SDK
- Internet access to download SLSA verifier binary
- Alpine-based container runtime

## Implementation Details

The module:

1. Downloads the appropriate SLSA verifier binary
2. Sets up a container environment
3. Mounts artifacts and provenance
4. Executes verification with specified parameters
5. Returns verification results

## Error Handling

The module includes error handling for:

- Missing artifacts
- Binary download failures
- Version resolution issues
- Verification failures

## License

See [LICENSE](../LICENSE) file in the root directory.
