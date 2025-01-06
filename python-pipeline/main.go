// Package main provides a complete pipeline for Python projects using Poetry and PyPI.
package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/python-pipeline/internal/dagger"
)

// PythonPipeline orchestrates Python project workflows using Poetry and PyPI.
type PythonPipeline struct{}

// New creates a new instance of PythonPipeline.
func New() *PythonPipeline {
	return &PythonPipeline{}
}

// CICD runs the complete CI/CD pipeline for a Python project.
// This includes:
// 1. Installing dependencies
// 2. Running tests
// 3. Running linting (if configured)
// 4. Building the package
// 5. Publishing to PyPI (if token is provided)
//
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// - token: Optional PyPI token for publishing. If provided, the package will be published
//
// Returns:
// - error: Any error that occurred during the process
func (m *PythonPipeline) CICD(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	// Setup Python container with Poetry
	container := dag.Container().
		From("python:3.12-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "curl", "ca-certificates"}).
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Configure git
	container = container.
		WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"})

	// Install dependencies
	container = container.WithExec([]string{"poetry", "install", "--no-interaction"})

	// Run tests
	_, err := container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running tests: %v", err)
	}

	// Run black check
	_, err = container.WithExec([]string{"poetry", "run", "black", ".", "--check"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running black check: %v", err)
	}

	// Run ruff check
	_, err = container.WithExec([]string{"poetry", "run", "ruff", "check", "."}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running ruff check: %v", err)
	}

	// If token is provided, publish to PyPI
	if token != nil {
		container = container.WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token)

		// Install python-semantic-release
		container = container.WithExec([]string{"pip", "install", "python-semantic-release"})

		// Get GitHub token from environment
		githubToken := dag.SetSecret("GITHUB_TOKEN", "")  // Primary token
		container = container.WithSecretVariable("GITHUB_TOKEN", githubToken)
		container = container.WithSecretVariable("GH_TOKEN", githubToken)  // Backwards compatibility

		// Configure git with token and get repository info
		container = container.WithExec([]string{"bash", "-c", `
			# Get repository info from git remote
			REPO_URL=$(git remote get-url origin)
			REPO_URL=${REPO_URL#*github.com[/:]}  # Remove everything before github.com
			REPO_URL=${REPO_URL%.git}             # Remove .git suffix
			REPO_OWNER=${REPO_URL%/*}             # Get owner
			REPO_NAME=${REPO_URL#*/}              # Get repo name
			
			# Configure git with token
			git config --global user.email "github-actions[bot]@users.noreply.github.com"
			git config --global user.name "github-actions[bot]"
			git config --global --add safe.directory '*'
			git config --global credential.helper store
			
			# Set up git remote with token
			git remote remove origin || true
			git remote add origin "https://$GITHUB_TOKEN@github.com/$REPO_OWNER/$REPO_NAME.git"
			git fetch origin main
			git reset --hard origin/main

			# Create a backup of the original pyproject.toml
			cp pyproject.toml pyproject.toml.bak

			# Check if semantic-release config already exists
			if ! grep -q "\[tool.semantic_release\]" pyproject.toml; then
				# Add semantic-release config only if it doesn't exist
				cat >> pyproject.toml << EOF

[tool.semantic_release]
version_variables = ["pyproject.toml:version"]
commit_author = "github-actions[bot] <github-actions[bot]@users.noreply.github.com>"
commit_parser = "angular"
branch = "main"
upload_to_pypi = true
build_command = "poetry build"
repository = "$REPO_NAME"
repository_owner = "$REPO_OWNER"

[tool.semantic_release.remote]
type = "github"
token = "\${GITHUB_TOKEN}"

[tool.semantic_release.publish]
dist_glob_patterns = ["dist/*"]
upload_to_vcs_release = true
upload_to_repository = true

[tool.semantic_release.branches.main]
match = "main"
prerelease_token = "rc"
prerelease = false

[tool.semantic_release.publish.pypi]
build = true
remove_dist = true
token = "\${POETRY_PYPI_TOKEN_PYPI}"
EOF
			fi

			# Ensure the token is available in the environment
			export GITHUB_TOKEN
			export GH_TOKEN="$GITHUB_TOKEN"
			export POETRY_PYPI_TOKEN_PYPI

			# Configure git credentials and test API access
			echo "https://$GITHUB_TOKEN@github.com" > ~/.git-credentials
			chmod 600 ~/.git-credentials
			
			# Test GitHub API access with token header
			curl -H "Authorization: Bearer $GITHUB_TOKEN" https://api.github.com/repos/$REPO_OWNER/$REPO_NAME
		`})

		// Run semantic-release version to determine and update version
		_, err = container.WithExec([]string{
			"semantic-release",
			"version",
			"--patch",  // Force a patch version increment
		}).Stdout(ctx)
		if err != nil {
			// Restore original pyproject.toml if semantic-release fails
			container = container.WithExec([]string{"mv", "pyproject.toml.bak", "pyproject.toml"})
			return fmt.Errorf("error running semantic-release version: %v", err)
		}

		// Clean up backup file
		container = container.WithExec([]string{"rm", "-f", "pyproject.toml.bak"})

		// Run semantic-release publish to handle publishing
		_, err = container.WithExec([]string{
			"semantic-release",
			"publish",
		}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("error running semantic-release publish: %v", err)
		}
	}

	return nil
}

// BuildAndPublish builds a Python package and publishes it to PyPI.
// The process includes:
// 1. Installing dependencies with Poetry
// 2. Running tests
// 3. Building the package
// 4. Publishing to PyPI
//
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
// - token: The PyPI authentication token
//
// Returns:
// - error: Any error that occurred during the process
func (m *PythonPipeline) BuildAndPublish(ctx context.Context, source *dagger.Directory, token *dagger.Secret) error {
	// Setup Python container with Poetry
	container := dag.Container().
		From("python:3.12-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Install dependencies
	container = container.WithExec([]string{"poetry", "install", "--no-interaction"})

	// Run tests
	_, err := container.WithExec([]string{"poetry", "run", "pytest"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running tests: %v", err)
	}

	// Build package
	container = container.WithExec([]string{"poetry", "build"})

	// Configure PyPI credentials and publish
	container = container.WithSecretVariable("POETRY_PYPI_TOKEN_PYPI", token)
	_, err = container.WithExec([]string{"poetry", "publish", "--no-interaction"}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error publishing to PyPI: %v", err)
	}

	return nil
}

// UpdateDependencies updates project dependencies and lock file.
// Parameters:
// - ctx: The context for the operation
// - source: The source directory containing the Python project
//
// Returns:
// - *dagger.Directory: The directory with updated dependencies
// - error: Any error that occurred during the update
func (m *PythonPipeline) UpdateDependencies(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	// Setup Python container with Poetry
	container := dag.Container().
		From("python:3.12-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"pip", "install", "--no-cache-dir", "poetry"})

	// Update dependencies
	container = container.WithExec([]string{"poetry", "update", "--no-interaction"})

	// Export the updated directory
	return container.Directory("/src"), nil
} 