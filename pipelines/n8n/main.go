package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger.io/dagger"
)

type N8N struct {
	Domain      string
	Subdomain   string
	N8nVersion  string
	Region      string
	DropletSize string
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

func (n *N8N) Deploy(ctx context.Context, doToken *dagger.Secret) (string, error) {
	// 1. Gerenciamento seguro de chaves SSH
	sshKeys, err := n.createSSHKeys(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create SSH keys: %w", err)
	}

	// 2. Provisionamento idempotente da infraestrutura
	dropletIP, err := n.provisionInfrastructure(ctx, doToken, sshKeys.publicKey)
	if err != nil {
		return "", fmt.Errorf("infrastructure provisioning failed: %w", err)
	}

	// 3. Implantação segura com verificação
	if err := n.deployN8n(ctx, dropletIP, sshKeys.privateKey); err != nil {
		return "", fmt.Errorf("deployment failed: %w", err)
	}

	// 4. Validação final
	if err := n.verifyDeployment(ctx, dropletIP); err != nil {
		return "", fmt.Errorf("post-deployment verification failed: %w", err)
	}

	return fmt.Sprintf("https://%s.%s", n.Subdomain, n.Domain), nil
}

type sshKeyPair struct {
	publicKey  string
	privateKey *dagger.Secret
}

func (n *N8N) createSSHKeys(ctx context.Context) (sshKeyPair, error) {
	keygen := dagger.Container().
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
		return sshKeyPair{}, err
	}

	return sshKeyPair{
		publicKey:  strings.TrimSpace(pubKey),
		privateKey: keygen.File("/key").Secret(),
	}, nil
}

func (n *N8N) provisionInfrastructure(
	ctx context.Context,
	doToken *dagger.Secret,
	publicKey string,
) (string, error) {
	doctl := dagger.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", doToken)

	// Gerenciamento idempotente de chaves SSH
	keyID, err := n.manageSSHKey(ctx, doctl, publicKey)
	if err != nil {
		return "", err
	}

	// Criação otimizada do droplet
	dropletIP, err := n.createDroplet(ctx, doctl, keyID)
	if err != nil {
		return "", err
	}

	// Configuração segura de DNS
	if err := n.configureDNS(ctx, doctl, dropletIP); err != nil {
		return "", err
	}

	return dropletIP, nil
}

func (n *N8N) manageSSHKey(
	ctx context.Context,
	doctl *dagger.Container,
	publicKey string,
) (string, error) {
	// Listar chaves existentes
	output, err := doctl.
		WithExec([]string{"compute", "ssh-key", "list", "--format", "ID,Name"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	// Procurar e remover chave existente
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == "n8n-key" {
			if _, err := doctl.
				WithExec([]string{"compute", "ssh-key", "delete", parts[0], "--force"}).
				Sync(ctx); err != nil {
				return "", err
			}
			break
		}
	}

	// Criar nova chave
	keyID, err := doctl.
		WithExec([]string{
			"compute", "ssh-key", "create",
			"n8n-key",
			"--public-key", publicKey,
			"--format", "ID",
			"--no-header",
		}).
		Stdout(ctx)

	return strings.TrimSpace(keyID), err
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
		return "", err
	}

	parts := strings.Fields(output)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid droplet output: %s", output)
	}

	fmt.Printf("Droplet created:\nID: %s\nIP: %s\n", parts[0], parts[1])
	return parts[1], nil
}

func (n *N8N) configureDNS(
	ctx context.Context,
	doctl *dagger.Container,
	ip string,
) error {
	// Verificar registros existentes
	records, err := doctl.
		WithExec([]string{
			"compute", "domain", "records", "list",
			n.Domain,
			"--format", "ID,Type,Name,Data",
		}).
		Stdout(ctx)
	if err != nil {
		return err
	}

	// Evitar duplicação
	target := fmt.Sprintf("A\t%s\t%s", n.Subdomain, ip)
	if strings.Contains(records, target) {
		fmt.Println("DNS record already exists")
		return nil
	}

	// Criar novo registro
	_, err = doctl.
		WithExec([]string{
			"compute", "domain", "records", "create",
			n.Domain,
			"--record-type", "A",
			"--record-name", n.Subdomain,
			"--record-data", ip,
		}).
		Sync(ctx)

	return err
}

func (n *N8N) deployN8n(
	ctx context.Context,
	ip string,
	privateKey *dagger.Secret,
) error {
	sshClient := dagger.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "openssh-client"}).
		WithSecretFile("/root/.ssh/id_ed25519", privateKey).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_ed25519"})

	// Conexão com backoff exponencial
	if err := n.waitForSSH(ctx, sshClient, ip); err != nil {
		return err
	}

	// Implantação segura
	return n.executeDeployment(ctx, sshClient, ip)
}

func (n *N8N) waitForSSH(
	ctx context.Context,
	sshClient *dagger.Container,
	ip string,
) error {
	const maxAttempts = 10
	timeout := 5 * time.Second

	for i := 0; i < maxAttempts; i++ {
		fmt.Printf("SSH attempt %d/%d\n", i+1, maxAttempts)

		_, err := sshClient.
			WithExec([]string{
				"ssh",
				"-o", "ConnectTimeout=5",
				"-o", "StrictHostKeyChecking=no",
				fmt.Sprintf("root@%s", ip),
				"echo ready",
			}).
			Sync(ctx)

		if err == nil {
			fmt.Println("SSH connection established")
			return nil
		}

		time.Sleep(timeout)
		timeout *= 2
	}

	return fmt.Errorf("SSH unreachable after %d attempts", maxAttempts)
}

func (n *N8N) executeDeployment(
	ctx context.Context,
	sshClient *dagger.Container,
	ip string,
) error {
	script := fmt.Sprintf(`#!/bin/bash
set -euo pipefail

# Criar diretório de trabalho
mkdir -p /opt/n8n && cd /opt/n8n

# Configurar docker-compose
cat > docker-compose.yml <<EOL
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
docker-compose up -d
`, n.N8nVersion, n.Subdomain, n.Domain)

	_, err := sshClient.
		WithNewFile("/deploy.sh", dagger.ContainerWithNewFileOpts{
			Contents:    script,
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
	healthCheck := dagger.Container().
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

	fmt.Println("n8n is operational")
	return nil
}
