# N8N Deployment Guide

## Prerequisites
- DigitalOcean account with API token
- Domain name with DNS managed by DigitalOcean
- Docker and Docker Compose installed locally

## Deployment Steps

### 1. Initial Setup
```bash
# Set your DigitalOcean token
export DIGITALOCEAN_TOKEN=your_digitalocean_token

# Navigate to the n8n pipeline directory
cd pipelines/n8n

# Run the deployment
dagger call deploy --do-token env:DIGITALOCEAN_TOKEN
```

### 2. Automated Deployment Process
The deployment will automatically:
1. Clean up any existing resources (droplets, DNS records)
2. Create SSH keys for secure access
3. Create a DigitalOcean droplet
4. Configure DNS records
5. Install and configure Docker
6. Deploy n8n with Caddy as reverse proxy
7. Set up SSL/TLS certificates automatically

### 3. Default Configuration

#### Environment Variables
```env
DATA_FOLDER=/home/n8n-user/n8n-docker-caddy
DOMAIN_NAME=pepper88.com
SUBDOMAIN=n8n
GENERIC_TIMEZONE=America/Sao_Paulo
N8N_BASIC_AUTH_ACTIVE=true
N8N_BASIC_AUTH_USER=admin
N8N_BASIC_AUTH_PASSWORD=admin123
N8N_HOST=n8n.pepper88.com
N8N_PROTOCOL=https
N8N_PORT=5678
N8N_ENCRYPTION_KEY=your-random-encryption-key
```

#### Caddy Configuration
```caddyfile
n8n.pepper88.com {
    reverse_proxy n8n:5678 {
        flush_interval -1
    }
}
```

### 4. Verification Steps

1. **Check DNS Propagation**
```bash
dig n8n.pepper88.com +short
```

2. **Verify HTTPS Access**
```bash
curl -I https://n8n.pepper88.com
```

3. **Check Container Status**
```bash
doctl compute ssh n8n --ssh-command "docker ps"
```

### 5. Default Access

#### N8N Web Interface
- URL: https://n8n.pepper88.com
- Username: admin
- Password: admin123

### 6. Troubleshooting

#### DNS Issues
- Wait for DNS propagation (can take up to 48 hours)
- Verify DNS record:
  ```bash
  doctl compute domain records list yourdomain.com
  ```

#### Container Issues
- Check container logs:
  ```bash
  doctl compute ssh n8n --ssh-command "docker logs n8n-docker-caddy_n8n_1"
  ```
- Check Caddy logs:
  ```bash
  doctl compute ssh n8n --ssh-command "docker logs n8n-docker-caddy_caddy_1"
  ```

#### SSL/TLS Issues
- Caddy handles certificates automatically
- Check Caddy logs for certificate issues
- Ensure ports 80 and 443 are open

### 7. Security Considerations

1. **Firewall Rules**
- Port 80 (HTTP) - Required for initial SSL setup
- Port 443 (HTTPS) - Required for secure access
- Port 22 (SSH) - Required for management
- Port 5678 (n8n) - Internal only, not exposed

2. **Authentication**
- Change default n8n credentials after first login
- Use strong passwords
- Enable basic authentication (enabled by default)

### 8. Maintenance

1. **Backup Considerations**
- n8n data is stored in Docker volumes
- Regular backups recommended:
  ```bash
  docker volume ls | grep n8n
  ```

2. **Updates**
To update n8n:
```bash
doctl compute ssh n8n --ssh-command "cd /home/n8n-user/n8n-docker-caddy && docker-compose pull && docker-compose up -d"
```

### 9. Clean Up
To remove all resources:
```bash
# Delete droplet
doctl compute droplet delete n8n -f

# Delete DNS record
doctl compute domain records list yourdomain.com | grep n8n
doctl compute domain records delete yourdomain.com <record-id> -f
```

## Common Issues and Solutions

### 1. Deployment Fails
**Problem**: Initial deployment fails
**Solution**: 
1. Clean up existing resources
2. Verify DIGITALOCEAN_TOKEN is set correctly
3. Ensure domain is managed by DigitalOcean

### 2. Cannot Access n8n
**Problem**: n8n URL not accessible
**Solutions**:
1. Wait for DNS propagation
2. Verify droplet is running
3. Check container logs
4. Verify Caddy configuration

### 3. SSL Certificate Issues
**Problem**: SSL certificate not working
**Solutions**:
1. Ensure DNS is properly configured
2. Check Caddy logs for certificate errors
3. Verify domain ownership

## Next Steps

1. **Custom Domain**
- Update domain in configuration
- Update SSL certificate

2. **Production Hardening**
- Change default credentials
- Configure backup solution
- Set up monitoring

3. **Performance Optimization**
- Adjust container resources
- Monitor system metrics
- Configure caching if needed

## References

1. [N8N Documentation](https://docs.n8n.io/)
2. [Nginx Documentation](https://nginx.org/en/docs/)
3. [Docker Documentation](https://docs.docker.com/)
4. [Node Exporter Documentation](https://github.com/prometheus/node_exporter) 