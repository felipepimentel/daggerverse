---
layout: default
title: GitHub Module
parent: Libraries
nav_order: 3
---

# GitHub Module

The GitHub module provides integration with GitHub through the GitHub CLI (`gh`), allowing you to interact with GitHub repositories and perform various GitHub operations directly from your Dagger pipelines.

## Features

- GitHub CLI command execution
- Repository cloning and management
- Git operations support
- GitHub token management
- Interactive terminal support
- Repository context management

## Installation

To use the GitHub module in your Dagger pipeline:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/gh"
)
```

## Usage Examples

### Basic GitHub CLI Setup

```go
func (m *MyModule) Example(ctx context.Context) (*Container, error) {
    gh, err := dag.Gh().New(
        dag.SetSecret("GITHUB_TOKEN", githubToken),
        "owner/repo",
        nil,
    )
    if err != nil {
        return nil, err
    }
    
    // Execute a GitHub CLI command
    return gh.Run("repo view"), nil
}
```

### Repository Operations

```go
func (m *MyModule) RepoOps(ctx context.Context) error {
    gh, err := dag.Gh().New(
        dag.SetSecret("GITHUB_TOKEN", githubToken),
        "",
        nil,
    )
    if err != nil {
        return err
    }
    
    // Clone a repository
    gh, err = gh.Clone("owner/repo")
    if err != nil {
        return err
    }
    
    // Execute git commands
    gh, err = gh.WithGitExec([]string{"checkout", "-b", "feature-branch"})
    if err != nil {
        return err
    }
    
    return nil
}
```

### Custom Commands

```go
func (m *MyModule) CustomCommands(ctx context.Context) (*Container, error) {
    gh, err := dag.Gh().New(
        dag.SetSecret("GITHUB_TOKEN", githubToken),
        "owner/repo",
        nil,
    )
    if err != nil {
        return nil, err
    }
    
    // Execute custom GitHub CLI commands
    return gh.Exec([]string{
        "issue", "create",
        "--title", "Bug Report",
        "--body", "Description of the bug",
    }), nil
}
```

## GitHub Actions Integration

You can use this module in your GitHub Actions workflows:

```yaml
name: GitHub Operations
on: [push]

jobs:
  github-ops:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: GitHub Operations with Dagger
        uses: dagger/dagger-action@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          module: github.com/felipepimentel/daggerverse/libraries/gh
          args: |
            do -p '
              gh := Gh().New(
                dag.SetSecret("GITHUB_TOKEN", GITHUB_TOKEN),
                "owner/repo",
                nil,
              )
              gh.Run("repo view")
            '
```

## API Reference

### Gh

Main module struct that provides access to GitHub CLI functionality.

#### Constructor

- `New(token *Secret, repo string, source *Directory) (*Gh, error)`
  - Creates a new GitHub CLI instance
  - Parameters:
    - `token`: GitHub token (optional)
    - `repo`: GitHub repository in "owner/repo" format (optional)
    - `source`: Git repository source directory (optional)

#### Methods

- `WithToken(token *Secret) *Gh`
  - Sets the GitHub token for authentication
  
- `WithRepo(repo string) (*Gh, error)`
  - Sets the GitHub repository context
  
- `WithSource(source *Directory) *Gh`
  - Sets the Git repository source directory
  
- `Clone(repo string) (*Gh, error)`
  - Clones a GitHub repository
  
- `Run(cmd string, token *Secret, repo string) *Container`
  - Runs a GitHub CLI command as a single string
  
- `Exec(args []string, token *Secret, repo string) *Container`
  - Executes a GitHub CLI command with arguments
  
- `WithGitExec(args []string) (*Gh, error)`
  - Executes a Git command in the repository
  
- `Terminal(token *Secret, repo string) *Container`
  - Opens an interactive terminal with GitHub CLI

## Best Practices

1. **Token Management**
   - Use secrets for GitHub tokens
   - Rotate tokens regularly
   - Use tokens with minimal required permissions

2. **Repository Operations**
   - Always check for errors after operations
   - Use appropriate repository context
   - Clean up temporary branches and clones

3. **Command Execution**
   - Prefer `Exec` over `Run` for complex commands
   - Handle command output appropriately
   - Use proper error handling

4. **Security**
   - Never commit tokens to source control
   - Use environment variables for sensitive data
   - Follow GitHub's security best practices

## Troubleshooting

Common issues and solutions:

1. **Authentication Issues**
   ```
   Error: HTTP 401: Bad credentials
   Solution: Verify GitHub token is valid and has required permissions
   ```

2. **Repository Access**
   ```
   Error: Repository not found
   Solution: Check repository name and access permissions
   ```

3. **Git Operations**
   ```
   Error: No git repository available
   Solution: Ensure repository is cloned before git operations
   ```

## Environment Variables

The module uses the following environment variables:

- `GITHUB_TOKEN`: For authentication
- `GH_REPO`: For repository context
- `GH_PROMPT_DISABLED`: Disabled by default
- `GH_NO_UPDATE_NOTIFIER`: Disabled by default

These can be set using the appropriate `With*` methods or through GitHub Actions secrets. 