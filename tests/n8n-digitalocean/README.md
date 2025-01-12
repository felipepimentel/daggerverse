# Testing n8n-digitalocean Module Locally

This directory contains scripts to test the n8n-digitalocean module locally. The tests simulate the CI/CD workflow that would normally be triggered by GitHub Actions.

## Prerequisites

1. Install Dagger CLI (version 0.15.1)
2. Have Docker running locally
3. Have all required secrets/environment variables:
   - `DIGITALOCEAN_ACCESS_TOKEN`: Your DigitalOcean API token
   - `DO_SSH_KEY_FINGERPRINT`: SSH key fingerprint registered in DigitalOcean
   - `DO_SSH_KEY_ID`: SSH key ID registered in DigitalOcean
   - `DO_SSH_PRIVATE_KEY`: SSH private key for DigitalOcean access
   - `N8N_BASIC_AUTH_PASSWORD`: Password for n8n basic auth
   - `N8N_DOMAIN`: Domain for n8n installation
   - `N8N_ENCRYPTION_KEY`: Encryption key for n8n

## Setup

1. Clone your n8n configuration repository:
   ```bash
   git clone <your-n8n-repo> tests/n8n-repo
   ```

2. Make the test script executable:
   ```bash
   chmod +x test.sh
   ```

3. Export the required environment variables:
   ```bash
   export DIGITALOCEAN_ACCESS_TOKEN="your-token"
   export DO_SSH_KEY_FINGERPRINT="your-fingerprint"
   export DO_SSH_KEY_ID="your-key-id"
   export DO_SSH_PRIVATE_KEY="your-private-key"
   export N8N_BASIC_AUTH_PASSWORD="your-password"
   export N8N_DOMAIN="your-domain"
   export N8N_ENCRYPTION_KEY="your-encryption-key"
   ```

## Running Tests

1. Run the test script:
   ```bash
   ./test.sh
   ```

The script will:
1. Verify all required environment variables are set
2. Set up SSH for DigitalOcean access
3. Run the CI pipeline
4. If CI passes and you're on the main branch, run the CD pipeline

## What's Being Tested

The test script simulates the same workflow that would be triggered by GitHub Actions:

1. **CI Phase**:
   - Checks the n8n configuration
   - Validates the setup
   - Runs any tests

2. **CD Phase** (only on main branch):
   - Builds the n8n container
   - Deploys to DigitalOcean
   - Sets up SSL/TLS with Caddy
   - Configures the domain

## Troubleshooting

If you encounter any issues:

1. Check that all environment variables are set correctly
2. Ensure your SSH key is valid and registered with DigitalOcean
3. Verify your DigitalOcean API token has the necessary permissions
4. Check that your domain's DNS is configured correctly
5. Make sure Docker is running and you have sufficient permissions

## Cleanup

After testing, you may want to:

1. Remove the test n8n deployment from DigitalOcean
2. Delete the SSH key from `~/.ssh/id_rsa`
3. Remove the n8n-repo directory:
   ```bash
   rm -rf tests/n8n-repo
   ``` 