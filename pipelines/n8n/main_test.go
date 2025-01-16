package main

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestN8N(t *testing.T) {
	// Skip if not running in CI or explicitly enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=1 to run")
	}

	// Test configuration
	domain := os.Getenv("TEST_DOMAIN")
	if domain == "" {
		t.Fatal("TEST_DOMAIN environment variable is required")
	}

	sshKeyName := os.Getenv("TEST_SSH_KEY_NAME")
	if sshKeyName == "" {
		t.Fatal("TEST_SSH_KEY_NAME environment variable is required")
	}

	// Initialize n8n module
	n8n := New(domain, "test-n8n", sshKeyName)

	// Configure test instance
	n8n.WithSize("s-1vcpu-1gb").
		WithN8NVersion("0.234.0").
		WithSSLEmail("test@example.com").
		WithRegistryConfig("test-n8n-registry")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Run CI checks
	t.Log("Running CI checks...")
	if err := n8n.CI(ctx); err != nil {
		t.Fatalf("CI checks failed: %v", err)
	}

	// Deploy n8n
	t.Log("Deploying n8n...")
	if err := n8n.Deploy(ctx); err != nil {
		t.Fatalf("Deploy failed: %v", err)
	}

	// Get deployment status
	t.Log("Getting deployment status...")
	status, err := n8n.GetStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	t.Logf("Deployment status: %s", status)

	// Clean up resources
	t.Log("Cleaning up resources...")
	if err := n8n.Cleanup(ctx); err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}
}

func TestN8NValidation(t *testing.T) {
	// Test invalid configuration
	tests := []struct {
		name        string
		setupN8N    func() *N8N
		expectError bool
	}{
		{
			name: "missing domain",
			setupN8N: func() *N8N {
				return New("", "test", "key")
			},
			expectError: true,
		},
		{
			name: "missing subdomain",
			setupN8N: func() *N8N {
				return New("example.com", "", "key")
			},
			expectError: true,
		},
		{
			name: "missing ssh key",
			setupN8N: func() *N8N {
				return New("example.com", "test", "")
			},
			expectError: true,
		},
		{
			name: "valid configuration",
			setupN8N: func() *N8N {
				return New("example.com", "test", "key").
					WithSSLEmail("test@example.com").
					WithRegistryConfig("test-registry")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			n8n := tt.setupN8N()
			err := n8n.CI(ctx)
			
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
} 