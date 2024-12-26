# Trivy Dagger Module

This Dagger module provides integration with [Trivy](https://github.com/aquasecurity/trivy), a comprehensive security scanner. Trivy can find vulnerabilities, misconfigurations, secrets, and generate SBOM in containers, Kubernetes, code repositories, clouds, and more.

## Features

- Container image scanning
- Filesystem scanning
- Helm chart scanning
- Binary scanning
- SBOM generation and scanning
- Multiple report formats support
- Cache persistence
- Custom configuration support

## Usage

### Basic Usage

```typescript
import { trivy } from "@felipepimentel/daggerverse/trivy";

// Initialize Trivy scanner
const scanner = trivy();

// Scan a container image
const results = await scanner.image("alpine:latest").output();
```

### Advanced Usage

#### Custom Configuration

```typescript
const scanner = trivy({
  version: "latest", // version
  container: null, // custom container
  config: myConfigFile, // config file
  cache: myCacheVolume, // cache volume
  databaseRepository: "my.registry/trivy-db", // database repository
  warmDatabaseCache: true, // warm cache
});
```

#### Different Scan Types

```typescript
// Image scanning
const imageScan = scanner.image("alpine:latest");

// Filesystem scanning
const fsScan = scanner.filesystem(myDirectory);

// Helm chart scanning
const helmScan = await scanner.helmChart(myChartFile);

// Binary scanning
const binaryScan = await scanner.binary(myBinaryFile);

// SBOM scanning
const sbomScan = await scanner.sbom(mySBOMFile);
```

#### Different Report Formats

```typescript
// Get results in different formats
const jsonReport = await scan.output({ format: "json" });
const sarifReport = await scan.output({ format: "sarif" });
const spdxReport = await scan.output({ format: "spdx" });

// Get report as a file
const reportFile = await scan.report({ format: "json" });
```

## API Reference

### Constructor

#### `trivy(options?: TrivyOptions)`

Creates a new instance of the Trivy scanner.

Parameters:

- `version` (optional): Version (image tag) to use from the official image repository
- `container` (optional): Custom container to use as base
- `config` (optional): Trivy configuration file
- `cache` (optional): Cache volume for persistence
- `databaseRepository` (optional): Custom trivy-db repository
- `warmDatabaseCache` (optional): Whether to warm the vulnerability database cache

### Scan Methods

#### `image(image: string, config?: File): Scan`

Scans a container image.

#### `filesystem(directory: Directory, target?: string, config?: File): Scan`

Scans a filesystem directory.

#### `helmChart(chart: File, ...): Promise<Scan>`

Scans a Helm chart.

#### `binary(binary: File, config?: File): Promise<Scan>`

Scans a binary file.

#### `sbom(sbom: File, config?: File): Promise<Scan>`

Scans an SBOM file.

### Report Formats

Available report formats:

- `table`
- `json`
- `template`
- `sarif`
- `cyclonedx`
- `spdx`
- `spdx-json`
- `github`
- `cosign-vuln`

## Examples

### Scanning a Container with Custom Configuration

```typescript
import { trivy } from "@felipepimentel/daggerverse/trivy";

export async function scanContainer() {
  const scanner = trivy();

  // Scan with custom config
  const scan = scanner.image("myapp:latest", myConfigFile);

  // Get JSON report
  const report = await scan.output({ format: "json" });
}
```

### Filesystem Scan with Cache

```typescript
import { trivy } from "@felipepimentel/daggerverse/trivy";

export async function scanFilesystem() {
  const cache = dag.cacheVolume("trivy-cache");

  const scanner = trivy({
    cache, // use cache
    warmDatabaseCache: true, // warm cache
  });

  const scan = scanner.filesystem(
    dag.host().directory("."),
    "src" // scan src directory
  );

  const report = await scan.output();
}
```

## License

See [LICENSE](../LICENSE) file in the root directory.
