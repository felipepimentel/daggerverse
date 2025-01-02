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
    echo "Initializing root module"
    TAG_PREFIX="root"
else
    TAG_PREFIX="$MODULE_NAME"
fi

# Check if module has any tags
if ! git tag -l "$TAG_PREFIX/v*" | grep -q .; then
    echo "Initializing module $MODULE_NAME with v0.0.0"
    # Create an initial tag if none exists
    git tag -a "$TAG_PREFIX/v0.0.0" -m "Initial version"
    git push origin "$TAG_PREFIX/v0.0.0"
fi 