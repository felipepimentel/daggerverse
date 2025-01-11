# n8n CI/CD Reusable Workflow

This reusable workflow provides CI/CD capabilities for n8n deployments to DigitalOcean. It leverages our Dagger modules to provide a consistent and reliable deployment process.

## Features

- 🔄 Automated CI/CD pipeline for n8n
- 🌊 DigitalOcean App Platform deployment
- 🔐 Secure secrets management
- ⚙️ Configurable environment variables
- 🎯 Region-specific deployment
- 🏷️ Custom app naming

## Prerequisites

- GitHub repository with n8n configuration
- DigitalOcean account and API token
- GitHub Actions enabled in your repository

## Usage

### Basic Example

Create `.github/workflows/n8n-deploy.yml` in your repository:

```yaml
name: Deploy n8n

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  deploy-n8n:
    uses: felipepimentel/daggerverse/.github/workflows/reusable-n8n-ci-cd.yml@main
    secrets:
      digitalocean_token: ${{ secrets.DIGITALOCEAN_TOKEN }}
```

### Advanced Example

```yaml
name: Deploy n8n with Custom Configuration

on:
  push:
    branches: [main, staging]
    paths:
      - 'n8n/**'
      - '.github/workflows/**'
  pull_request:
    branches: [main, staging]
    paths:
      - 'n8n/**'
      - '.github/workflows/**'

jobs:
  deploy-n8n:
    uses: felipepimentel/daggerverse/.github/workflows/reusable-n8n-ci-cd.yml@main
    with:
      dagger_version: "0.15.1"
      region: "fra1"
      app_name: "n8n-production"
      env_vars: |
        {
          "N8N_BASIC_AUTH_ACTIVE": "true",
          "N8N_BASIC_AUTH_USER": "admin",
          "N8N_PROTOCOL": "https",
          "N8N_HOST": "workflow.example.com",
          "N8N_PORT": "5678",
          "N8N_ENCRYPTION_KEY": "${{ secrets.N8N_ENCRYPTION_KEY }}",
          "WEBHOOK_TUNNEL_URL": "https://workflow.example.com/"
        }
    secrets:
      digitalocean_token: ${{ secrets.DIGITALOCEAN_TOKEN }}
```

## Inputs

| Name | Description | Required | Default | Example Values |
|------|-------------|----------|---------|----------------|
| `dagger_version` | Version of Dagger to use | No | `0.15.1` | `0.15.1`, `0.15.0` |
| `region` | DigitalOcean region for deployment | No | `nyc` | `nyc`, `fra1`, `sgp1` |
| `app_name` | Name of the n8n application | No | `n8n` | `n8n-prod`, `n8n-staging` |
| `env_vars` | JSON string of environment variables | No | Basic auth config | See environment variables section |

## Secrets

| Name | Description | Required | How to Obtain |
|------|-------------|----------|---------------|
| `digitalocean_token` | DigitalOcean API token | Yes | [Create in DigitalOcean](https://cloud.digitalocean.com/account/api/tokens) |

## Environment Variables

### Required Variables

```json
{
  "N8N_BASIC_AUTH_ACTIVE": "true",
  "N8N_BASIC_AUTH_USER": "admin"
}
```

### Recommended Variables

```json
{
  "N8N_PROTOCOL": "https",
  "N8N_HOST": "your-domain.com",
  "N8N_PORT": "5678",
  "N8N_ENCRYPTION_KEY": "your-secure-key",
  "WEBHOOK_TUNNEL_URL": "https://your-domain.com/",
  "N8N_EMAIL_MODE": "smtp",
  "N8N_SMTP_HOST": "smtp.example.com",
  "N8N_SMTP_PORT": "587",
  "N8N_SMTP_USER": "your-smtp-user",
  "N8N_SMTP_PASS": "your-smtp-password"
}
```

## Workflow Details

### CI Process

The CI job performs the following steps:

1. Code checkout
2. Dagger setup with specified version
3. Runs tests and validations
4. Checks n8n configuration
5. Validates environment variables

### CD Process

The CD job (runs only on main branch) performs:

1. Code checkout
2. Dagger setup
3. DigitalOcean authentication
4. App deployment configuration
5. Environment variables setup
6. Application deployment
7. Health check verification

## Best Practices

1. **Environment Variables**
   - Store sensitive data in GitHub Secrets
   - Use descriptive names for variables
   - Document all custom variables

2. **Security**
   - Rotate DigitalOcean tokens regularly
   - Use HTTPS for production deployments
   - Enable basic authentication
   - Set strong encryption keys

3. **Deployment**
   - Use different app names for staging/production
   - Set appropriate resource limits
   - Configure automatic backups
   - Monitor deployment logs

## Troubleshooting

### Common Issues

1. **Authentication Failures**
   ```
   Error: Unable to authenticate with DigitalOcean
   ```
   - Verify token permissions
   - Check token expiration
   - Ensure token is properly set in secrets

2. **Environment Variables**
   ```
   Error: Invalid environment variables format
   ```
   - Validate JSON syntax
   - Check for missing quotes
   - Ensure all values are strings

3. **Region Issues**
   ```
   Error: Invalid region specified
   ```
   - Use valid DigitalOcean region codes
   - Check region availability
   - Verify resource availability in region

### Debug Steps

1. Check workflow run logs
2. Verify environment variables
3. Validate DigitalOcean configuration
4. Test n8n configuration locally
5. Review deployment status in DigitalOcean dashboard

## Related Resources

- [n8n Documentation](https://docs.n8n.io/)
- [DigitalOcean App Platform](https://www.digitalocean.com/products/app-platform)
- [Dagger Documentation](https://docs.dagger.io/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

## Support

For issues and feature requests:
- Open an issue in the [repository](https://github.com/felipepimentel/daggerverse)
- Provide workflow run logs
- Include error messages
- Describe expected vs actual behavior 