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
TAG_PREFIX=""
if [ "$MODULE_NAME" = "." ]; then
    echo "Publishing root module"
    TAG_PREFIX="root"
else
    TAG_PREFIX="$MODULE_NAME"
fi

# Get the latest tag, trying different approaches
TAG=$(git describe --tags --abbrev=0 --match "$TAG_PREFIX/v*" 2>/dev/null || \
      git tag -l "$TAG_PREFIX/v*" | sort -V | tail -n1 || \
      echo "")

# If still no tag, something is wrong with the release process
if [ -z "$TAG" ]; then
    echo "::error::No semantic version tag found for module $MODULE_NAME. Release process failed."
    echo "::error::This should not happen as we initialize modules with v0.0.0."
    echo "::error::Please check the release step logs for more details."
    exit 1
fi

echo "Publishing module $MODULE_NAME with version $TAG"

# Change to module directory if not root
if [ "$MODULE_NAME" != "." ]; then
    cd "$MODULE_NAME"
fi

# Use --force flag if FORCE_PUBLISH is set to true
if [ "${FORCE_PUBLISH:-}" = "true" ]; then
    echo "Force publishing enabled"
    PUBLISH_CMD="dagger publish --force"
else
    PUBLISH_CMD="dagger publish"
fi

$PUBLISH_CMD || {
    echo "::error::Failed to publish module $MODULE_NAME"
    echo "::error::Please check if the module is properly configured and try again"
    exit 1
} 