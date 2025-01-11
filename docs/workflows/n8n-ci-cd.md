# n8n CI/CD Reusable Workflow

This reusable workflow provides CI/CD capabilities for n8n deployments to DigitalOcean. It leverages our Dagger modules to provide a consistent and reliable deployment process.

## Usage

To use this workflow in your repository, create a workflow file (e.g., `.github/workflows/n8n-deploy.yml`) with the following content:

```yaml
name: Deploy n8n

on:
  push:
    branches:
      - main  # or your preferred branch
  pull_request:  # if you want CI checks on PRs

jobs:
  deploy-n8n:
    uses: felipepimentel/daggerverse/.github/workflows/reusable-n8n-ci-cd.yml@main
    with:
      dagger_version: "0.9.3"  # optional, defaults to 0.9.3
      region: "nyc"            # optional, defaults to nyc
      app_name: "my-n8n"      # optional, defaults to n8n
      env_vars: '{"N8N_BASIC_AUTH_ACTIVE": "true", "N8N_BASIC_AUTH_USER": "admin"}'  # optional
    secrets:
      digitalocean_token: ${{ secrets.DIGITALOCEAN_TOKEN }}  # required
```

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|---------|
| `dagger_version` | Version of Dagger to use | No | `0.9.3` |
| `region` | DigitalOcean region for deployment | No | `nyc` |
| `app_name` | Name of the n8n application | No | `n8n` |
| `env_vars` | JSON string of environment variables | No | `{"N8N_BASIC_AUTH_ACTIVE": "true", "N8N_BASIC_AUTH_USER": "admin"}` |

## Secrets

| Name | Description | Required |
|------|-------------|----------|
| `digitalocean_token` | DigitalOcean API token | Yes |

## Workflow Steps

1. **CI Job**:
   - Checks out the code
   - Sets up Dagger with specified version
   - Runs CI checks using our Dagger modules

2. **CD Job** (only runs on main branch):
   - Checks out the code
   - Sets up Dagger with specified version
   - Deploys to DigitalOcean using our Dagger modules

## Related Modules

This workflow uses the following Dagger modules from our repository:

- `/pipelines/n8n`: Core n8n pipeline module
- `/pipelines/n8n-digitalocean`: DigitalOcean-specific deployment module
- `/libraries/digitalocean`: DigitalOcean API integration module

## Example Environment Variables

Here's an example of environment variables you might want to set:

```json
{
  "N8N_BASIC_AUTH_ACTIVE": "true",
  "N8N_BASIC_AUTH_USER": "admin",
  "N8N_PROTOCOL": "https",
  "N8N_HOST": "your-domain.com",
  "N8N_PORT": "5678"
}
```

## Security Considerations

- Never commit your DigitalOcean token. Always use GitHub Secrets.
- Consider encrypting sensitive environment variables.
- Review the environment variables you expose to n8n.

## Troubleshooting

If you encounter issues:

1. Ensure your DigitalOcean token has the necessary permissions
2. Check that your environment variables are properly formatted JSON
3. Verify your DigitalOcean region is valid
4. Make sure you're using a compatible Dagger version 