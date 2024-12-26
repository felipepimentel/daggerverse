# Supergit Module

This Dagger module provides enhanced Git functionality for your Dagger pipelines, allowing you to work with Git repositories, remotes, tags, and branches.

## Usage

```typescript
import { supergit } from "@felipepimentel/daggerverse/supergit";

// Create a new instance
const git = supergit();

// Work with a remote repository
const remote = git.remote("https://github.com/user/repo.git");

// Get a specific tag
const tag = await remote.tag("v1.0.0");

// Get the files at that tag
const files = await tag.commit().tree();
```

## Functions

### Repository Operations

#### `repository()`

Creates a new Git repository.

#### `withGitCommand(args: string[])`

Executes a Git command in the repository.

#### `withRemote(name: string, url: string)`

Adds a remote to the repository.

### Remote Operations

#### `remote(url: string)`

Creates a new Git remote reference.

#### `tag(name: string)`

Looks up a tag in the remote.

#### `tags(filter?: string)`

Lists all tags in the remote, optionally filtered by a regular expression.

#### `branch(name: string)`

Looks up a branch in the remote.

#### `branches(filter?: string)`

Lists all branches in the remote, optionally filtered by a regular expression.

### Tag Operations

#### `tag(name: string)`

Gets a tag reference.

#### `tree()`

Gets the directory tree at the tag.

### Commit Operations

#### `commit(digest: string)`

Gets a commit reference.

#### `tree()`

Gets the directory tree at the commit.

## Example

```typescript
import { supergit } from "@felipepimentel/daggerverse/supergit";

export default async function example() {
  const git = supergit();

  // Create a new repository
  const repo = git.repository();

  // Add a remote
  const withRemote = repo.withRemote(
    "origin",
    "https://github.com/user/repo.git"
  );

  // Work with a remote directly
  const remote = git.remote("https://github.com/user/repo.git");

  // Get all tags matching a pattern
  const tags = await remote.tags("^v[0-9]");

  // Get files from the latest tag
  if (tags.length > 0) {
    const files = await tags[0].commit().tree();
  }
}
```

## License

See [LICENSE](../LICENSE) file in the root directory.
