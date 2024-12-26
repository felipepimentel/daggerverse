# CI Module for Dagger

A Dagger module that provides comprehensive CI/CD pipeline functionality with support for Node.js applications, infrastructure management with Terragrunt, YAML manipulation with YQ, and automated technology detection.

## Features

- Node.js pipeline automation with NPM and Yarn support
- Infrastructure management with Terragrunt
- YAML file manipulation with YQ
- Automated technology stack detection
- OCI image building and publishing
- Package management and publishing
- Configurable pipeline options
- Dry-run capabilities

## Components

### Node.js Pipeline

```typescript
import { node } from "@felipepimentel/daggerverse/ci";

// Lazy mode with OCI build
const result = await node().withAutoSetup("my-app", source).pipeline({
  dryRun: true,
  ttl: "5m",
  isOci: true,
});

// Explicit mode with package build
const result = await node()
  .withPipelineID("my-app")
  .withVersion("20.9.0")
  .withSource(source)
  .withNpm()
  .install()
  .test()
  .build()
  .publish({
    dryRun: true,
    devTag: "beta",
  });
```

### Infrastructure Management

```typescript
import { infrabox } from "@felipepimentel/daggerverse/ci";

// Terragrunt operations
const result = await infrabox()
  .terragrunt()
  .withSource("/terraform", source)
  .disableColor()
  .plan("/terraform/stacks/dev/region/env/stack")
  .apply("/terraform/stacks/dev/region/env/stack");
```

### Technology Detection

```typescript
import { autodetection } from "@felipepimentel/daggerverse/ci";

// Analyze Node.js project
const analyzer = await autodetection().node(source);

// Check for specific features
const isTest = await analyzer.isTest();
const isPackage = await analyzer.isPackage();
const isNpm = await analyzer.isNpm();
const isYarn = await analyzer.isYarn();
```

### YAML Manipulation

```typescript
import { yq } from "@felipepimentel/daggerverse/ci";

// Read YAML values
const value = await yq(source).get(".path.to.key", "config.yaml");

// Modify YAML values
const result = await yq(source).set('.path.to.key="new_value"', "config.yaml");
```

## Examples

### Complete Node.js Pipeline

```typescript
import { node } from "@felipepimentel/daggerverse/ci";

export async function nodePipeline() {
  // Initialize Node.js pipeline
  const pipeline = node().withAutoSetup("my-app", source);

  // Configure and run pipeline
  const refs = await pipeline.pipeline({
    dryRun: false,
    ttl: "1h",
    isOci: true,
    packageDevTag: "beta",
  });

  console.log("Built images:", refs);
}
```

### Infrastructure Management Pipeline

```typescript
import { infrabox } from "@felipepimentel/daggerverse/ci";

export async function infrastructurePipeline() {
  // Initialize Terragrunt pipeline
  const infra = infrabox()
    .terragrunt()
    .withSource("/terraform", source)
    .disableColor();

  // Plan and apply changes
  await infra
    .plan("/terraform/stacks/dev/region/env/stack")
    .apply("/terraform/stacks/dev/region/env/stack")
    .plan("/terraform/stacks/dev/region/env/stack", {
      detailedExitCode: true,
    });
}
```

### Technology Analysis

```typescript
import { autodetection } from "@felipepimentel/daggerverse/ci";

export async function analyzeProject() {
  // Initialize analyzer
  const analyzer = await autodetection().node(source);

  // Perform analysis
  const isTest = await analyzer.isTest();
  const isNpm = await analyzer.isNpm();
  const isYarn = await analyzer.isYarn();

  console.log(`Project analysis:
- Has tests: ${isTest}
- Uses NPM: ${isNpm}
- Uses Yarn: ${isYarn}`);
}
```

## API Reference

### Node Pipeline

```typescript
interface Node {
  // Auto-setup mode
  withAutoSetup(name: string, source: Directory): Node;

  // Explicit mode
  withPipelineID(id: string): Node;
  withVersion(version: string): Node;
  withSource(source: Directory): Node;
  withNpm(): Node;

  // Pipeline operations
  install(): Node;
  test(): Node;
  build(): Node;
  ociBuild(opts?: NodeOciBuildOpts): Promise<string[]>;
  publish(opts: NodePublishOpts): Node;
  pipeline(opts: NodePipelineOpts): Promise<string>;
}
```

### Infrastructure Management

```typescript
interface Infrabox {
  terragrunt(): Terragrunt;
}

interface Terragrunt {
  withSource(path: string, source: Directory): Terragrunt;
  disableColor(): Terragrunt;
  plan(path: string, opts?: InfraboxTfPlanOpts): Terragrunt;
  apply(path: string): Terragrunt;
}
```

### Technology Detection

```typescript
interface Autodetection {
  node(source: Directory): NodeAnalyzer;
}

interface NodeAnalyzer {
  isTest(): Promise<boolean>;
  isPackage(): Promise<boolean>;
  isNpm(): Promise<boolean>;
  isYarn(): Promise<boolean>;
  is(feature: string): Promise<boolean>;
}
```

## Configuration Options

### Node Pipeline Options

```typescript
interface NodePipelineOpts {
  dryRun?: boolean; // Enable dry-run mode
  ttl?: string; // Time-to-live for builds
  isOci?: boolean; // Enable OCI image building
  packageDevTag?: string; // Development tag for packages
}
```

### Infrastructure Options

```typescript
interface InfraboxTfPlanOpts {
  detailedExitCode?: boolean; // Enable detailed exit codes
}
```

## Testing

The module includes a comprehensive test suite that can be run using:

```bash
dagger run test
```

The test suite verifies:

- Node.js pipeline functionality
- Infrastructure management operations
- Technology detection accuracy
- YAML manipulation capabilities
- Error handling and edge cases
- Configuration options
- Integration scenarios

## License

See [LICENSE](../LICENSE) file in the root directory.
