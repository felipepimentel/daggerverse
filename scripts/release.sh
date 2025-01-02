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

# Export for semantic-release
export MODULE_NAME
export MODULE_PATH="$MODULE_NAME"

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