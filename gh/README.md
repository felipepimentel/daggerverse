# GitHub CLI Module for Dagger

A Dagger module that provides integration with GitHub CLI (`gh`), allowing you to interact with GitHub repositories and perform various GitHub operations in your Dagger pipelines.

## Features

- GitHub repository cloning and management
- GitHub CLI command execution
- Git operations support
- Token-based authentication
- Interactive terminal access
- Repository context management
- Source directory mounting

## Usage

### Basic Setup

```typescript
import { gh } from "@felipepimentel/daggerverse/gh";

// Initialize the GitHub CLI module
const client = gh({
  token, // Optional: GitHub token
  repo: "owner/repo", // Optional: GitHub repository
  source: sourceDir, // Optional: Git repository source
});
```

### Authentication

```typescript
// Set GitHub token
const withToken = client.withToken(githubToken);
```

### Repository Management

```typescript
// Set repository context
const withRepo = await client.withRepo("owner/repo");

// Clone a repository
const withClone = await client.clone("owner/repo");

// Load existing repository
const withSource = client.withSource(repoDir);
```

### Command Execution

```typescript
// Run a GitHub CLI command
const container = client.run("pr list", token, "owner/repo");

// Execute with custom arguments
const container = client.exec(
  ["issue", "list", "--state", "open"],
  token,
  "owner/repo"
);

// Execute Git commands
const withGit = await client.withGitExec(["status"]);
```

### Interactive Terminal

```typescript
// Open an interactive terminal
const terminal = client.terminal(token, "owner/repo");
```

## Configuration

### Constructor Options

The module accepts:

- `token`: GitHub token for authentication (optional)
- `repo`: GitHub repository in "owner/repo" format (optional)
- `source`: Git repository source directory (optional)

### Default Settings

- Base image: Wolfi container with `gh` and `git` packages
- Working directory: `/work/repo` when source is provided
- Environment variables:
  - `GH_PROMPT_DISABLED=true`
  - `GH_NO_UPDATE_NOTIFIER=true`
  - `GH_REPO` (when repository is specified)
  - `GITHUB_TOKEN` (when token is provided)

## Examples

### Create a Pull Request

```typescript
import { gh } from "@felipepimentel/daggerverse/gh";

export async function createPR() {
  const client = gh({
    token: githubToken,
    repo: "owner/repo",
  });

  // Clone repository
  const withClone = await client.clone("");

  // Create and checkout branch
  const withBranch = await withClone.withGitExec([
    "checkout",
    "-b",
    "feature-branch",
  ]);

  // Make changes and commit
  const withAdd = await withBranch.withGitExec(["add", "."]);
  const withCommit = await withAdd.withGitExec([
    "commit",
    "-m",
    "Add new feature",
  ]);

  // Create PR
  await withCommit.run("pr create --title 'New Feature' --body 'Description'");
}
```

### List Issues

```typescript
import { gh } from "@felipepimentel/daggerverse/gh";

export async function listIssues() {
  const client = gh({
    token: githubToken,
    repo: "owner/repo",
  });

  // List open issues
  const container = await client.exec([
    "issue",
    "list",
    "--state",
    "open",
    "--limit",
    "10",
  ]);
}
```

### Work with Git Repository

```typescript
import { gh } from "@felipepimentel/daggerverse/gh";

export async function gitOperations() {
  const client = gh({
    source: sourceDir,
  });

  // Check repository status
  const withStatus = await client.withGitExec(["status"]);

  // Pull latest changes
  const withPull = await withStatus.withGitExec(["pull", "origin", "main"]);
}
```

## Dependencies

The module requires:

- Dagger SDK
- Internet access to pull the Wolfi container image
- GitHub token for authenticated operations
- Git repository (for source-based operations)

## License

See [LICENSE](../LICENSE) file in the root directory.
