#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Check if required environment variables are set
required_vars=(
    "DIGITALOCEAN_ACCESS_TOKEN"
    "DO_SSH_KEY_FINGERPRINT"
    "DO_SSH_KEY_ID"
    "DO_SSH_PRIVATE_KEY"
    "N8N_BASIC_AUTH_PASSWORD"
    "N8N_DOMAIN"
    "N8N_ENCRYPTION_KEY"
)

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo -e "${RED}Error: $var is not set${NC}"
        echo "Please set all required environment variables:"
        printf '%s\n' "${required_vars[@]}"
        exit 1
    fi
done

echo -e "${GREEN}All required environment variables are set${NC}"

# Set up SSH key
echo "Setting up SSH key..."
mkdir -p ~/.ssh
chmod 700 ~/.ssh
echo "$DO_SSH_PRIVATE_KEY" > ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa
ssh-keyscan -H digitalocean.com >> ~/.ssh/known_hosts 2>/dev/null
ssh-keyscan -H registry.digitalocean.com >> ~/.ssh/known_hosts 2>/dev/null
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_rsa

# Change to the n8n-digitalocean module directory
cd ../../pipelines/n8n-digitalocean || exit 1

# Run CI
echo -e "\n${GREEN}Running CI...${NC}"
dagger call ci \
    --n8n-repo "../../../tests/n8n-repo"

# If CI passes and we're on main branch, run CD
if [ $? -eq 0 ] && [ "$(git rev-parse --abbrev-ref HEAD)" = "main" ]; then
    echo -e "\n${GREEN}Running CD...${NC}"
    dagger call deploy \
        --token "$DIGITALOCEAN_ACCESS_TOKEN" \
        --region "nyc" \
        --app-name "n8n" \
        --ssh-key "$DO_SSH_PRIVATE_KEY" \
        --ssh-key-fingerprint "$DO_SSH_KEY_FINGERPRINT" \
        --ssh-key-id "$DO_SSH_KEY_ID" \
        --domain "$N8N_DOMAIN" \
        --basic-auth-password "$N8N_BASIC_AUTH_PASSWORD" \
        --encryption-key "$N8N_ENCRYPTION_KEY" \
        --n8n-repo "../../../tests/n8n-repo"
fi 