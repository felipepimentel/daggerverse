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

# Get the latest tag, trying different approaches
TAG=$(git describe --tags --abbrev=0 --match "$MODULE_NAME/v*" 2>/dev/null || \
      git tag -l "$MODULE_NAME/v*" | sort -V | tail -n1 || \
      echo "")

# If still no tag, something is wrong with the release process
if [ -z "$TAG" ]; then
    echo "::error::No semantic version tag found for module $MODULE_NAME. Release process failed."
    echo "::error::This should not happen as we initialize modules with v0.0.0."
    echo "::error::Please check the release step logs for more details."
    exit 1
fi

echo "Publishing module $MODULE_NAME with version $TAG"

# Publish to Daggerverse
cd "$MODULE_NAME"
dagger publish || {
    echo "::error::Failed to publish module $MODULE_NAME"
    echo "::error::Please check if the module is properly configured and try again"
    exit 1
} 