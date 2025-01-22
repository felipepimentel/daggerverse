package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/felipepimentel/daggerverse/pipelines/n8n/internal/dagger"
)

type N8N struct {
	Domain      string
	Subdomain   string
	N8nVersion  string
	Region      string
	DropletSize string
	dag         *dagger.Client
}

func New() *N8N {
	return &N8N{
		Domain:      "pepper88.com",
		Subdomain:   "n8n",
		N8nVersion:  "0.234.0",
		Region:      "syd1",
		DropletSize: "s-1vcpu-1gb",
	}
}

func (n *N8N) Connect() error {
	client := dagger.Connect()
	n.dag = client
	return nil
}

func (n *N8N) Deploy(ctx context.Context, doToken *dagger.Secret) (string, error) {
	if err := n.Connect(); err != nil {
		return "", err
	}

	// 1. Gerar par de chaves SSH
	sshKeys, err := n.createSSHKeys(ctx)
	if err != nil {
		return "", fmt.Errorf("SSH key generation failed: %w", err)
	}

	// 2. Provisionar infraestrutura
	dropletIP, err := n.provisionInfrastructure(ctx, doToken, sshKeys.publicKey)
	if err != nil {
		return "", fmt.Errorf("infrastructure setup failed: %w", err)
	}

	// 3. Implantar n8n
	if err := n.deployN8n(ctx, dropletIP, sshKeys.privateKey); err != nil {
		return "", fmt.Errorf("deployment failed: %w", err)
	}

	// 4. Verificar implantação
	if err := n.verifyDeployment(ctx, dropletIP); err != nil {
		return "", fmt.Errorf("verification failed: %w", err)
	}

	return fmt.Sprintf("https://%s.%s", n.Subdomain, n.Domain), nil
}

type sshKeyPair struct {
	publicKey  string
	privateKey *dagger.Secret
}

func (n *N8N) createSSHKeys(ctx context.Context) (sshKeyPair, error) {
	keygen := n.dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "openssh-keygen"}).
		WithExec([]string{
			"ssh-keygen",
			"-t", "ed25519",
			"-f", "/key",
			"-N", "",
			"-C", "n8n-deploy",
		})

	pubKey, err := keygen.File("/key.pub").Contents(ctx)
	if err != nil {
		return sshKeyPair{}, fmt.Errorf("failed to get public key: %w", err)
	}

	keyContent, err := keygen.File("/key").Contents(ctx)
	if err != nil {
		return sshKeyPair{}, fmt.Errorf("failed to get private key content: %w", err)
	}

	privateKey := n.dag.SetSecret("ssh-priv-key", keyContent)

	return sshKeyPair{
		publicKey:  strings.TrimSpace(pubKey),
		privateKey: privateKey,
	}, nil
}

func (n *N8N) provisionInfrastructure(
	ctx context.Context,
	doToken *dagger.Secret,
	publicKey string,
) (string, error) {
	doctl := n.dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", doToken)

	keyID, err := n.manageSSHKey(ctx, doctl, publicKey)
	if err != nil {
		return "", fmt.Errorf("SSH key management failed: %w", err)
	}

	dropletIP, err := n.createDroplet(ctx, doctl, keyID)
	if err != nil {
		return "", fmt.Errorf("droplet creation failed: %w", err)
	}

	if err := n.configureDNS(ctx, doctl, dropletIP); err != nil {
		return "", fmt.Errorf("DNS configuration failed: %w", err)
	}

	return dropletIP, nil
}

func (n *N8N) manageSSHKey(
	ctx context.Context,
	doctl *dagger.Container,
	publicKey string,
) (string, error) {
	output, err := doctl.
		WithExec([]string{"compute", "ssh-key", "list", "--format", "ID,Name"}).
		Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list SSH keys: %w", err)
	}

	for _, line := range strings.Split(output, "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == "n8n-key" {
			if _, err := doctl.
				WithExec([]string{"compute", "ssh-key", "delete", parts[0], "--force"}).
				Sync(ctx); err != nil {
				return "", fmt.Errorf("failed to delete SSH key: %w", err)
			}
			break
		}
	}

	keyID, err := doctl.
		WithExec([]string{
			"compute", "ssh-key", "create",
			"n8n-key",
			"--public-key", publicKey,
			"--format", "ID",
			"--no-header",
		}).
		Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create SSH key: %w", err)
	}

	return strings.TrimSpace(keyID), nil
}

func (n *N8N) createDroplet(
	ctx context.Context,
	doctl *dagger.Container,
	sshKeyID string,
) (string, error) {
	output, err := doctl.
		WithExec([]string{
			"compute", "droplet", "create",
			fmt.Sprintf("n8n-%s", n.Subdomain),
			"--region", n.Region,
			"--size", n.DropletSize,
			"--image", "docker-22-04",
			"--ssh-keys", sshKeyID,
			"--format", "ID,PublicIPv4",
			"--no-header",
			"--wait",
		}).
		Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create droplet: %w", err)
	}

	parts := strings.Fields(output)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid droplet output: %q", output)
	}

	fmt.Printf("✅ Droplet created:\nID: %s\nIP: %s\n", parts[0], parts[1])
	return parts[1], nil
}

func (n *N8N) configureDNS(
	ctx context.Context,
	doctl *dagger.Container,
	ip string,
) error {
	records, err := doctl.
		WithExec([]string{
			"compute", "domain", "records", "list",
			n.Domain,
			"--format", "ID,Type,Name,Data",
		}).
		Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to list DNS records: %w", err)
	}

	target := fmt.Sprintf("A\t%s\t%s", n.Subdomain, ip)
	if strings.Contains(records, target) {
		fmt.Println("ℹ️ DNS record already exists")
		return nil
	}

	_, err = doctl.
		WithExec([]string{
			"compute", "domain", "records", "create",
			n.Domain,
			"--record-type", "A",
			"--record-name", n.Subdomain,
			"--record-data", ip,
		}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to create DNS record: %w", err)
	}

	fmt.Println("✅ DNS record created")
	return nil
}

func (n *N8N) deployN8n(
	ctx context.Context,
	ip string,
	privateKey *dagger.Secret,
) error {
	sshClient := n.dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "openssh-client"}).
		WithMountedSecret("/root/.ssh/id_ed25519", privateKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_ed25519"})

	if err := n.waitForSSH(ctx, sshClient, ip); err != nil {
		return fmt.Errorf("SSH connection failed: %w", err)
	}

	if err := n.executeDeployment(ctx, sshClient, ip); err != nil {
		return fmt.Errorf("deployment script failed: %w", err)
	}

	return nil
}

func (n *N8N) waitForSSH(
	ctx context.Context,
	sshClient *dagger.Container,
	ip string,
) error {
	const maxAttempts = 10
	timeout := 5 * time.Second

	for i := 0; i < maxAttempts; i++ {
		fmt.Printf("🔌 SSH attempt %d/%d\n", i+1, maxAttempts)

		_, err := sshClient.
			WithExec([]string{
				"ssh",
				"-o", "ConnectTimeout=5",
				"-o", "StrictHostKeyChecking=no",
				fmt.Sprintf("root@%s", ip),
				"echo SSH_OK",
			}).
			Sync(ctx)

		if err == nil {
			fmt.Println("✅ SSH connection established")
			return nil
		}

		fmt.Printf("⚠️ Attempt failed: %v\n", err)
		time.Sleep(timeout)
		timeout *= 2
	}

	return fmt.Errorf("failed to establish SSH connection after %d attempts", maxAttempts)
}

func (n *N8N) executeDeployment(
	ctx context.Context,
	sshClient *dagger.Container,
	ip string,
) error {
	script := fmt.Sprintf(`#!/bin/bash
set -euo pipefail

# Instalar dependências
apt-get update
apt-get install -y docker.io docker-compose

# Configurar Docker
systemctl enable --now docker
usermod -aG docker $USER

# Configurar n8n
mkdir -p /opt/n8n
cat > /opt/n8n/docker-compose.yml <<EOL
version: '3.8'
services:
  n8n:
    image: n8nio/n8n:%s
    restart: unless-stopped
    ports:
      - "80:5678"
    environment:
      - N8N_HOST=%s.%s
      - N8N_PROTOCOL=https
    volumes:
      - n8n_data:/home/node/.n8n
volumes:
  n8n_data:
EOL

# Iniciar serviços
cd /opt/n8n && docker-compose up -d
`, n.N8nVersion, n.Subdomain, n.Domain)

	_, err := sshClient.
		WithNewFile("/deploy.sh", script, dagger.ContainerWithNewFileOpts{
			Permissions: 0755,
		}).
		WithExec([]string{
			"scp",
			"-o", "StrictHostKeyChecking=no",
			"/deploy.sh",
			fmt.Sprintf("root@%s:/tmp/deploy.sh", ip),
		}).
		WithExec([]string{
			"ssh",
			"-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("root@%s", ip),
			"bash /tmp/deploy.sh && rm /tmp/deploy.sh",
		}).
		Sync(ctx)

	return err
}

func (n *N8N) verifyDeployment(
	ctx context.Context,
	ip string,
) error {
	healthCheck := n.dag.Container().
		From("curlimages/curl:latest").
		WithExec([]string{
			"curl",
			"-sSf",
			"--retry", "5",
			"--retry-delay", "10",
			"--max-time", "30",
			fmt.Sprintf("http://%s:5678/healthz", ip),
		})

	_, err := healthCheck.Sync(ctx)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	fmt.Println("✅ n8n is operational")
	return nil
}
