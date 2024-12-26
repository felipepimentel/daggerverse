# Autodetection Module for Dagger

A Dagger module that provides runtime information detection capabilities for Node.js and OCI (Open Container Initiative) projects. This module helps analyze project structures, dependencies, and configurations automatically.

## Features

### Node.js Detection

- Package.json analysis
- Dependency management detection (npm/yarn)
- Test file detection
- Script detection
- Workspace configuration
- Engine version detection
- Package publishing configuration
- Repository information

### OCI Detection

- Dockerfile presence
- Containerfile presence
- Pattern-based detection

## Usage

### Node.js Analysis

```typescript
import { autodetection } from "@felipepimentel/daggerverse/autodetection";

// Initialize analyzer
const analyzer = await autodetection().node({
  src, // Source directory
  exclude: ["dist"], // Optional pattern exclusions
});

// Check for test files
const isTest = await analyzer.isTest();

// Check package manager
const isYarn = await analyzer.isYarn();
const isNpm = await analyzer.isNpm();

// Check for specific scripts
const hasTest = await analyzer.is("test");

// Get Node.js version
const nodeVersion = await analyzer.getEngineVersion();

// Get package information
const name = await analyzer.getName();
const version = await analyzer.getVersion();

// Get available scripts
const scripts = await analyzer.getScriptNames();

// Check workspaces
const workspaces = await analyzer.getWorkspaces();

// Check if it's a publishable package
const isPackage = await analyzer.isPackage();
```

### OCI Analysis

```typescript
import { autodetection } from "@felipepimentel/daggerverse/autodetection";

// Initialize analyzer
const analyzer = await autodetection().oci({
  src, // Source directory
  exclude: [".git"], // Optional pattern exclusions
});

// Check for OCI files
const isOci = await analyzer.isOci();
```

## Configuration Options

### Node.js Analyzer

#### Pattern Exclusions

Default exclusions:

- `node_modules`
- `.tsconfig`

#### Pattern Matching

Default patterns for detection:

- Test files:
  - `.+\.(test|spec)\.js`
  - `.+\.(test|spec)\.jsx`
  - `.+\.(test|spec)\.ts`
  - `.+\.(test|spec)\.tsx`
  - `(.+/)*(__)*tests*(__)*/.+`
- Package managers:
  - Yarn: `.*yarn.lock`
  - NPM: `.*package-lock.json`

### OCI Analyzer

#### Pattern Matching

Default patterns for detection:

- `.*Dockerfile`
- `.*Containerfile`

## Implementation Details

### Node.js Analysis

The module analyzes:

- Package.json structure and content
- Project file patterns
- Dependency management files
- Test file patterns
- Script configurations
- Engine requirements
- Publishing configurations

### OCI Analysis

The module detects:

- Container definition files
- Build configurations
- Container-related patterns

## Examples

### Complete Node.js Analysis

```typescript
import { autodetection } from "@felipepimentel/daggerverse/autodetection";

export async function analyze() {
  // Initialize analyzer
  const analyzer = await autodetection().node({ src });

  // Get project information
  const name = await analyzer.getName();
  const version = await analyzer.getVersion();
  const nodeVersion = await analyzer.getEngineVersion();

  // Check development setup
  if (await analyzer.isYarn()) {
    // Use Yarn commands
  } else if (await analyzer.isNpm()) {
    // Use NPM commands
  }

  // Check test configuration
  if (await analyzer.isTest()) {
    const hasTestScript = await analyzer.is("test");
    if (hasTestScript) {
      // Run tests
    }
  }

  // Check workspaces
  const workspaces = await analyzer.getWorkspaces();
  if (workspaces.length > 0) {
    // Handle monorepo setup
  }
}
```

### Basic OCI Analysis

```typescript
import { autodetection } from "@felipepimentel/daggerverse/autodetection";

export async function analyzeOci() {
  // Initialize analyzer
  const analyzer = await autodetection().oci({ src });

  // Check for container configuration
  if (await analyzer.isOci()) {
    // Handle container build
  }
}
```

## Dependencies

The module requires:

- Dagger SDK
- Go 1.22 or later
- Access to project filesystem

## License

See [LICENSE](../LICENSE) file in the root directory.
