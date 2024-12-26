# Editorconfig Module

This Dagger module provides integration with [EditorConfig](https://editorconfig.org/), a tool that helps maintain consistent coding styles across different editors and IDEs.

## Usage

```typescript
import { editorconfig } from "@felipepimentel/daggerverse/editorconfig";

// Create a new instance
const checker = editorconfig();

// Run editorconfig check on your source code
await checker.check(dag.host().directory("."), ".git");
```

## Functions

### `new(image?: string)`

Creates a new instance of the Editorconfig checker.

- `image`: Optional. Custom image reference in "repository:tag" format to use as a base container. Defaults to "mstruebing/editorconfig-checker:latest".

### `check(source: Directory, excludeDirectoryPattern?: string)`

Runs the editorconfig-checker command on the specified source directory.

- `source`: The directory to check.
- `excludeDirectoryPattern`: Optional. Pattern to exclude directories. Defaults to ".git".

## Example

```typescript
import { editorconfig } from "@felipepimentel/daggerverse/editorconfig";

export default async function check() {
  const checker = editorconfig();

  // Check the current directory, excluding .git
  await checker.check(dag.host().directory("."));
}
```

## License

See [LICENSE](../LICENSE) file in the root directory.
