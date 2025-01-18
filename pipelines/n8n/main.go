// Package main provides a Dagger module for n8n CI/CD operations
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/felipepimentel/daggerverse/pipelines/n8n/internal/dagger"
)

// N8N handles the deployment of n8n to DigitalOcean
type N8N struct {
	// Domain for the n8n instance
	Domain string
	// Subdomain for the n8n instance
	Subdomain string
	// Registry name in DigitalOcean
	RegistryName string
	// N8n version to deploy
	N8nVersion string
	// Region for the droplet
	Region string
	// Size for the droplet
	Size string
}

// New creates a new instance of the N8N module
func New() *N8N {
	return &N8N{
		Domain:       "pepper88.com",
		Subdomain:    "n8n",
		RegistryName: "n8n-registry",
		N8nVersion:   "0.234.0",
		Region:       "syd1",
		Size:         "s-1vcpu-1gb",
	}
}

// Deploy orchestrates the n8n deployment
func (n *N8N) Deploy(ctx context.Context, doToken *dagger.Secret) error {
	fmt.Println("Starting n8n deployment...")

	// Get DO token plaintext
	tokenValue, err := doToken.Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get DO token plaintext: %w", err)
	}

	// Pull and push image
	fmt.Println("\nPulling image...")
	container := dag.Container().From("n8nio/n8n:0.234.0")

	fmt.Println("\nAdding labels...")
	container = container.
		WithLabel("org.opencontainers.image.created", time.Now().UTC().Format(time.RFC3339)).
		WithLabel("org.opencontainers.image.source", "https://github.com/felipepimentel/daggerverse").
		WithLabel("org.opencontainers.image.version", "0.234.0")

	fmt.Println("\nPushing image...")
	targetRef := fmt.Sprintf("registry.digitalocean.com/n8n-registry/n8nio/n8n:%s", "0.234.0")
	container = container.WithRegistryAuth("registry.digitalocean.com", "", doToken)
	_, err = container.Publish(ctx, targetRef)
	if err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}
	fmt.Printf("\nSuccessfully pushed %s\n", targetRef)

	// Set up SSH key container
	fmt.Println("\nGenerating SSH key...")
	sshContainer := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "openssh"}).
		WithExec([]string{"mkdir", "-p", "/root/.ssh"}).
		WithExec([]string{
			"ssh-keygen",
			"-t", "ed25519",
			"-f", "/root/.ssh/id_ed25519",
			"-N", "",
			"-C", "n8n-deployment",
		})

	// Get public key contents
	pubKey, err := sshContainer.File("/root/.ssh/id_ed25519.pub").Contents(ctx)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	// Get private key contents
	privateKey, err := sshContainer.File("/root/.ssh/id_ed25519").Contents(ctx)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	// Save private key to file for debugging
	keyFile := "/tmp/n8n_id_ed25519"
	fmt.Printf("\nSaving private key to %s\n", keyFile)
	err = os.WriteFile(keyFile, []byte(privateKey), 0600)
	if err != nil {
		return fmt.Errorf("failed to save private key: %w", err)
	}

	fmt.Println("\nPrivate key for debugging:")
	fmt.Println(privateKey)

	// Set up DigitalOcean container
	fmt.Println("\nSetting up DigitalOcean client...")
	doContainer := dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", doToken)

	// Delete existing SSH key if it exists
	fmt.Println("\nCleaning up existing SSH keys...")
	_, _ = doContainer.WithExec([]string{
		"/app/doctl", "compute", "ssh-key", "delete", "n8n-key",
		"--force",
	}).Stdout(ctx)

	// Add SSH key to DigitalOcean
	fmt.Println("\nAdding SSH key to DigitalOcean...")
	sshKeyOutput, err := doContainer.WithExec([]string{
		"/app/doctl", "compute", "ssh-key", "create",
		"n8n-key",
		"--public-key", pubKey,
		"--format", "ID",
		"--no-header",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to add SSH key: %w", err)
	}
	sshKeyID := strings.TrimSpace(sshKeyOutput)

	// Create droplet with SSH key
	fmt.Println("\nCreating droplet...")
	dropletOutput, err := doContainer.WithExec([]string{
		"/app/doctl", "compute", "droplet", "create",
		fmt.Sprintf("n8n-%s", n.Subdomain),
		"--region", "syd1",
		"--size", "s-1vcpu-1gb",
		"--image", "ubuntu-22-04-x64",
		"--ssh-keys", sshKeyID,
		"--format", "ID,PublicIPv4",
		"--no-header",
		"--wait",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to create droplet: %w", err)
	}

	// Parse droplet output
	parts := strings.Fields(dropletOutput)
	if len(parts) != 2 {
		return fmt.Errorf("unexpected droplet output format: %s", dropletOutput)
	}
	dropletID := parts[0]
	dropletIP := parts[1]

	fmt.Printf("\nDroplet created successfully:\n")
	fmt.Printf("  ID: %s\n", dropletID)
	fmt.Printf("  IP: %s\n", dropletIP)

	// Configure DNS
	fmt.Println("\nConfiguring DNS...")
	_, err = doContainer.WithExec([]string{
		"/app/doctl", "compute", "domain", "records", "create",
		n.Domain,
		"--record-type", "A",
		"--record-name", n.Subdomain,
		"--record-data", dropletIP,
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to configure DNS: %w", err)
	}

	// Wait for droplet to be ready
	fmt.Println("\nWaiting for droplet to complete initialization (2 minutes)...")
	time.Sleep(2 * time.Minute)

	// Prepare setup script
	fmt.Println("\nPreparing setup script...")
	setupScript := fmt.Sprintf(`#!/bin/bash
	set -e

	# Update package list and install dependencies
	apt-get update
	apt-get install -y docker.io docker-compose

	# Start Docker service
	systemctl start docker
	systemctl enable docker

	# Log in to registry
	echo "%s" | docker login registry.digitalocean.com --username _ --password-stdin

	# Pull image
	docker pull %s

	# Create docker-compose.yml
	cat > /root/docker-compose.yml <<EOL
	version: '3.8'
	services:
	  n8n:
		image: %s
		restart: always
		ports:
		  - "80:5678"
		environment:
		  - N8N_HOST=%s.%s
		  - N8N_PROTOCOL=https
		  - NODE_ENV=production
		volumes:
		  - n8n_data:/home/node/.n8n
	volumes:
	  n8n_data:
	EOL

	# Start services
	cd /root && docker-compose up -d
	`, tokenValue, targetRef, targetRef, n.Subdomain, n.Domain)

	// Set up deployment container
	fmt.Println("\nPreparing deployment container...")
	deployContainer := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "openssh-client"}).
		WithNewFile("/root/.ssh/id_ed25519", privateKey).
		WithNewFile("/setup.sh", setupScript).
		WithExec([]string{"chmod", "600", "/root/.ssh/id_ed25519"}).
		WithExec([]string{"chmod", "755", "/setup.sh"})

	// Wait for SSH to be ready (try multiple times)
	fmt.Println("\nChecking SSH connectivity...")
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("\nAttempt %d/%d to establish SSH connection...\n", i+1, maxRetries)
		
		// Check droplet status
		status, err := doContainer.WithExec([]string{
			"/app/doctl", "compute", "droplet", "get",
			dropletID,
			"--format", "Status",
			"--no-header",
		}).Stdout(ctx)
		if err != nil {
			fmt.Printf("Failed to get droplet status: %v\n", err)
		} else {
			fmt.Printf("Droplet status: %s\n", strings.TrimSpace(status))
		}
		
		// Try to connect
		_, err = deployContainer.WithExec([]string{
			"ssh",
			"-o", "StrictHostKeyChecking=no",
			"-o", "ConnectTimeout=10",
			fmt.Sprintf("root@%s", dropletIP),
			"echo 'SSH connection successful'",
		}).Stdout(ctx)
		
		if err == nil {
			fmt.Println("SSH connection established successfully!")
			break
		}
		
		if i == maxRetries-1 {
			return fmt.Errorf("failed to establish SSH connection after %d attempts", maxRetries)
		}
		
		fmt.Printf("SSH connection failed: %v\n", err)
		fmt.Printf("Waiting 30 seconds before next attempt...\n")
		time.Sleep(30 * time.Second)
	}

	// Copy setup script
	fmt.Println("\nCopying setup script...")
	_, err = deployContainer.WithExec([]string{
		"scp",
		"-o", "StrictHostKeyChecking=no",
		"-o", "ConnectTimeout=10",
		"/setup.sh",
		fmt.Sprintf("root@%s:/root/setup.sh", dropletIP),
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to copy setup script: %w", err)
	}

	// Execute setup script
	fmt.Println("\nExecuting setup script...")
	_, err = deployContainer.WithExec([]string{
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "ConnectTimeout=10",
		fmt.Sprintf("root@%s", dropletIP),
		"bash /root/setup.sh",
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute setup script: %w", err)
	}

	fmt.Printf("\nDeployment complete! n8n will be available at https://%s.%s\n", n.Subdomain, n.Domain)
	fmt.Printf("Note: It may take a few minutes for DNS to propagate and for n8n to fully start up.\n")
	return nil
}
