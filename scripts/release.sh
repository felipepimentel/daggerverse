#!/usr/bin/env bash
set -euo pipefail

MODULE_NAME="${1:-}"
if [ -z "$MODULE_NAME" ]; then
    echo "Error: MODULE_NAME is required"
    exit 1
fi

# Check if module directory exists
if [ ! -d "$MODULE_NAME" ]; then
    echo "Error: Module directory $MODULE_NAME does not exist"
    exit 1
fi

# Handle root directory specially
if [ "$MODULE_NAME" = "." ]; then
    echo "Creating release for root module"
    export MODULE_NAME="root"
    export MODULE_PATH="."
else
    export MODULE_NAME
    export MODULE_PATH="$MODULE_NAME"
fi

# Ensure we're up to date with remote
echo "Syncing with remote repository..."
git pull --rebase origin main || {
    echo "Failed to rebase. Resolving conflicts..."
    git rebase --abort
    git reset --hard origin/main
    git clean -fd
}

# Create .releaserc.json for semantic-release configuration
cat > .releaserc.json << EOF
{
  "branches": ["main"],
  "plugins": [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/changelog",
    ["@semantic-release/exec", {
      "prepareCmd": "poetry version \${nextRelease.version}",
      "publishCmd": "git add pyproject.toml && git commit -m 'chore(release): bump version to \${nextRelease.version} [skip ci]' || true"
    }],
    ["@semantic-release/git", {
      "assets": ["CHANGELOG.md", "pyproject.toml"],
      "message": "chore(release): \${nextRelease.version} [skip ci]\n\n\${nextRelease.notes}"
    }],
    "@semantic-release/github"
  ]
}
EOF

echo "Running semantic-release dry-run for module $MODULE_NAME"

# Try dry-run first
if ! npx semantic-release --no-ci-skip --dry-run; then
    # Check if we have an initial version
    if ! git tag -l "$MODULE_NAME/v*" | grep -q .; then
        echo "::warning::No changes detected for a new version, but module has no version yet."
        echo "::warning::This should have been handled by init-module.sh"
        exit 1
    fi
    
    echo "::notice::No changes detected that would trigger a new version."
    echo "::notice::Using existing version for publishing."
    exit 0
fi

echo "Creating new release for module $MODULE_NAME"

# Actual release
if ! npx semantic-release; then
    echo "::error::Failed to create release for module $MODULE_NAME"
    echo "::error::Please check the semantic-release logs for more details"
    exit 1
fi 