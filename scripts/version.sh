#!/bin/bash

# Function to determine commit type from message
get_commit_type() {
    local commit_msg="$1"
    if [[ $commit_msg == *"BREAKING CHANGE"* ]]; then
        echo "BREAKING CHANGE"
    elif [[ $commit_msg == "feat"* ]]; then
        echo "feat"
    elif [[ $commit_msg == "fix"* ]] || [[ $commit_msg == "perf"* ]]; then
        echo "fix"
    else
        echo "patch"
    fi
}

# Function to bump version for a module
bump_version() {
    local module="$1"
    local commit_type="$2"
    
    echo "Bumping version for $module based on commit type: $commit_type"
    dagger call --module versioner bump-version \
        --source . \
        --module "$module" \
        --commit-type "$commit_type"
}

# Get the last commit message
COMMIT_MSG=$(git log -1 --pretty=%B)
COMMIT_TYPE=$(get_commit_type "$COMMIT_MSG")

# List of modules to version
MODULES=(
    "python-poetry"
    "python-pypi"
    "python-pipeline"
)

# Bump version for each module
for module in "${MODULES[@]}"; do
    if [ -d "$module" ]; then
        bump_version "$module" "$COMMIT_TYPE"
    fi
done 