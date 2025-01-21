package main

import (
	"context"
	"fmt"
	"os"

	n8n "github.com/felipepimentel/daggerverse/pipelines/n8n"
)

func main() {
	ctx := context.Background()

	// Create a new n8n deployment
	n8nModule := n8n.New(
		"pepper88.com",  // Domain
		"n8n",          // Subdomain
		"dagger",       // SSH key name
	).WithSSLEmail("admin@pepper88.com").
		WithSize("s-2vcpu-4gb").
		WithN8NVersion("0.234.0").
		WithPostgresConfig("15-alpine", "n8n", "n8n", "n8n").
		WithBackupConfig(true, "0 0 * * *", 7).
		WithMonitoringConfig(true, 8080)

	// Run CI checks
	fmt.Println("Running CI checks...")
	if err := n8nModule.CI(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "CI validation failed: %v\n", err)
		os.Exit(1)
	}

	// Get required secrets from environment
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "DO_TOKEN environment variable is required\n")
		os.Exit(1)
	}

	sshKey := os.Getenv("SSH_PRIVATE_KEY")
	if sshKey == "" {
		fmt.Fprintf(os.Stderr, "SSH_PRIVATE_KEY environment variable is required\n")
		os.Exit(1)
	}

	// Deploy n8n
	fmt.Println("Deploying n8n...")
	if err := n8nModule.Deploy(ctx, token, sshKey); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to deploy n8n: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully deployed n8n to %s\n", n8nModule.GetURL())
} 