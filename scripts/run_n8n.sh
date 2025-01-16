#!/usr/bin/env bash

set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function to print colored messages
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Check if DO_TOKEN is set
if [ -z "${DO_TOKEN:-}" ]; then
    warn "DO_TOKEN environment variable is not set"
    exit 1
fi

# Change to the workspace root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Generate SSH key for deployment
SSH_KEY_PATH="${WORKSPACE_DIR}/.ssh/n8n_ed25519"
SSH_KEY_DIR="$(dirname "${SSH_KEY_PATH}")"

log "Generating SSH key for deployment..."
mkdir -p "${SSH_KEY_DIR}"
if [ ! -f "${SSH_KEY_PATH}" ]; then
    ssh-keygen -t ed25519 -f "${SSH_KEY_PATH}" -C "n8n-deployment-key" -N ""
fi

# Read the private and public keys
PRIVATE_KEY="$(cat "${SSH_KEY_PATH}")"
PUBLIC_KEY="$(cat "${SSH_KEY_PATH}.pub")"

# Build and run dependent modules first
MODULES=(
    "libraries/digitalocean"
    "essentials/ssh"
    "essentials/dig"
    "essentials/curl"
    "libraries/docker"
    "pipelines/n8n"
)

for module in "${MODULES[@]}"; do
    log "Building module: ${module}"
    (cd "${WORKSPACE_DIR}/${module}" && dagger develop)
done

# Run the n8n module
log "Running n8n module..."
cd "${WORKSPACE_DIR}/pipelines/n8n"

# Export required secrets
export DAGGER_SESSION_TOKEN="$(dagger session)"

# Run the module with the generated SSH keys
dagger call deploy \
    --do-token "${DO_TOKEN}" \
    --ssh-key "${PRIVATE_KEY}" \
    --ssh-pub-key "${PUBLIC_KEY}"

log "n8n module execution completed successfully!" 