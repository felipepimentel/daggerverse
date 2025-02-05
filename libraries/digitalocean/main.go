// Package main provides a Dagger module for DigitalOcean operations
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/felipepimentel/daggerverse/libraries/digitalocean/internal/dagger"
)

// DigitalOcean provides functionality for managing DigitalOcean resources
type DigitalOcean struct {
	token *dagger.Secret
}

// SSHKeyConfig holds configuration for SSH key operations
type SSHKeyConfig struct {
	Name      string
	PublicKey string
}

// RegistryConfig holds configuration for container registry operations
type RegistryConfig struct {
	Name string
}

// DropletConfig holds configuration for creating a droplet
type DropletConfig struct {
	Name       string
	Region     string
	Size       string
	Image      string
	SSHKeyID   string
	Monitoring bool
	IPv6       bool
	Tags       []string
}

// DNSConfig holds configuration for managing DNS records
type DNSConfig struct {
	Domain   string
	Type     string
	Name     string
	Value    string
	TTL      int
	Weight   int
	Port     int
	Flag     int
	Tag      string
	Priority int
}

// New creates a new instance of the DigitalOcean module
func New(token *dagger.Secret) *DigitalOcean {
	return &DigitalOcean{
		token: token,
	}
}

// SSH Key Management

// CreateSSHKey creates a new SSH key
func (do *DigitalOcean) CreateSSHKey(ctx context.Context, config SSHKeyConfig) (*dagger.Container, error) {
	fmt.Printf("🔑 Creating SSH key: %s\n", config.Name)
	return dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"compute",
			"ssh-key",
			"create",
			config.Name,
			"--public-key", config.PublicKey,
			"--format", "ID",
			"--no-header",
		}), nil
}

// ListSSHKeys lists all SSH keys
func (do *DigitalOcean) ListSSHKeys(ctx context.Context, format string) (*dagger.Container, error) {
	fmt.Println("🔍 Listing SSH keys...")
	args := []string{"compute", "ssh-key", "list"}
	if format != "" {
		args = append(args, "--format", format)
	}
	if format == "ID" {
		args = append(args, "--no-header")
	}

	return dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec(args), nil
}

// Registry Management

// CreateRegistry creates a new container registry
func (do *DigitalOcean) CreateRegistry(ctx context.Context, config RegistryConfig) (*dagger.Container, error) {
	fmt.Printf("🔧 Creating registry: %s\n", config.Name)
	return dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"registry",
			"create",
			config.Name,
		}), nil
}

// GetRegistry gets registry details
func (do *DigitalOcean) GetRegistry(ctx context.Context) (*dagger.Container, error) {
	fmt.Println("🔍 Getting registry details...")
	return dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"registry",
			"get",
		}), nil
}

// ListRegistryTags lists all tags in a registry repository
func (do *DigitalOcean) ListRegistryTags(ctx context.Context, registry string) (*dagger.Container, error) {
	fmt.Printf("🔍 Listing tags for registry: %s\n", registry)
	return dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"registry",
			"repository",
			"list-tags",
			registry,
		}), nil
}

// DeleteRegistry deletes a container registry
func (do *DigitalOcean) DeleteRegistry(ctx context.Context, name string) error {
	fmt.Printf("🗑️ Deleting registry: %s\n", name)
	_, err := dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"registry",
			"delete",
			name,
			"--force",
		}).
		Stdout(ctx)
	return err
}

// Droplet Management

// CreateDroplet creates a new droplet with the given configuration
func (do *DigitalOcean) CreateDroplet(ctx context.Context, config DropletConfig) (*dagger.Container, error) {
	if config.Name == "" || config.Region == "" || config.Size == "" || config.Image == "" {
		return nil, fmt.Errorf("missing required droplet configuration")
	}

	fmt.Printf("🚀 Creating droplet: %s\n", config.Name)
	fmt.Printf("  Region: %s\n", config.Region)
	fmt.Printf("  Size: %s\n", config.Size)
	fmt.Printf("  Image: %s\n", config.Image)
	fmt.Printf("  SSH Keys: %s\n", config.SSHKeyID)

	args := []string{
		"compute",
		"droplet",
		"create",
		config.Name,
		"--region", config.Region,
		"--size", config.Size,
		"--image", config.Image,
		"--ssh-keys", config.SSHKeyID,
		"--wait",
		"--format", "ID,Name,PublicIPv4",
		"--no-header",
	}

	if config.Monitoring {
		args = append(args, "--enable-monitoring")
	}

	if config.IPv6 {
		args = append(args, "--enable-ipv6")
	}

	if len(config.Tags) > 0 {
		args = append(args, "--tag-names", fmt.Sprintf("[%s]", config.Tags[0]))
		for _, tag := range config.Tags[1:] {
			args[len(args)-1] = fmt.Sprintf("%s,%s", args[len(args)-1], tag)
		}
	}

	return dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec(args), nil
}

// GetDroplet retrieves information about a droplet by name
func (do *DigitalOcean) GetDroplet(ctx context.Context, name string, format string) (*dagger.Container, error) {
	fmt.Printf("🔍 Getting droplet: %s\n", name)
	args := []string{
		"compute",
		"droplet",
		"get",
		name,
	}

	if format != "" {
		args = append(args, "--format", format)
		if format == "PublicIPv4" {
			args = append(args, "--no-header")
		}
	}

	return dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec(args), nil
}

// DeleteDroplet deletes a droplet by name
func (do *DigitalOcean) DeleteDroplet(ctx context.Context, name string) error {
	fmt.Printf("🗑️ Deleting droplet: %s\n", name)
	_, err := dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"compute",
			"droplet",
			"delete",
			name,
			"--force",
		}).
		Stdout(ctx)
	return err
}

// DNS Management

// CreateDNSRecord creates a new DNS record
func (do *DigitalOcean) CreateDNSRecord(ctx context.Context, config DNSConfig) error {
	fmt.Printf("🌐 Creating DNS record: %s.%s -> %s\n", config.Name, config.Domain, config.Value)
	args := []string{
		"compute",
		"domain",
		"records",
		"create",
		config.Domain,
		"--record-type", config.Type,
		"--record-name", config.Name,
		"--record-data", config.Value,
	}

	if config.TTL > 0 {
		args = append(args, "--record-ttl", fmt.Sprintf("%d", config.TTL))
	}

	if config.Priority > 0 {
		args = append(args, "--record-priority", fmt.Sprintf("%d", config.Priority))
	}

	_, err := dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec(args).
		Stdout(ctx)
	return err
}

// ListDNSRecords lists all DNS records for a domain
func (do *DigitalOcean) ListDNSRecords(ctx context.Context, domain string) (*dagger.Container, error) {
	fmt.Printf("🔍 Listing DNS records for domain: %s\n", domain)
	return dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"compute",
			"domain",
			"records",
			"list",
			domain,
			"--format", "ID,Type,Name,Data",
		}), nil
}

// DeleteDNSRecord deletes a DNS record
func (do *DigitalOcean) DeleteDNSRecord(ctx context.Context, domain string, recordID string) error {
	fmt.Printf("🗑️ Deleting DNS record: %s (ID: %s)\n", domain, recordID)
	_, err := dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"compute",
			"domain",
			"records",
			"delete",
			domain,
			recordID,
			"--force",
		}).
		Stdout(ctx)
	return err
}

// Utility Functions

// WaitForDroplet waits for a droplet to reach the desired status
func (do *DigitalOcean) WaitForDroplet(ctx context.Context, name string, status string, timeout time.Duration) error {
	fmt.Printf("⏳ Waiting for droplet %s to reach status %s (timeout: %s)\n", name, status, timeout)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		container, err := do.GetDroplet(ctx, name, "Status")
		if err != nil {
			return err
		}

		output, err := container.Stdout(ctx)
		if err != nil {
			return err
		}

		if output == status {
			fmt.Printf("✅ Droplet %s reached status %s\n", name, status)
			return nil
		}

		time.Sleep(10 * time.Second)
	}

	return fmt.Errorf("timeout waiting for droplet %s to reach status %s", name, status)
}

// ListDroplets lists all droplets in the account
func (do *DigitalOcean) ListDroplets(ctx context.Context) (*dagger.Container, error) {
	fmt.Println("🔍 Listing all droplets...")
	return dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"compute",
			"droplet",
			"list",
			"--format", "ID,Name,PublicIPv4,Status",
		}), nil
}

// DeleteSSHKey deletes an SSH key by ID
func (do *DigitalOcean) DeleteSSHKey(ctx context.Context, keyID string) error {
	fmt.Printf("🗑️ Deleting SSH key: %s\n", keyID)
	_, err := dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"compute",
			"ssh-key",
			"delete",
			keyID,
			"--force",
		}).
		Stdout(ctx)
	return err
}

// RegisterSSHKey registers an SSH key with DigitalOcean
func (do *DigitalOcean) RegisterSSHKey(ctx context.Context, name string, publicKey string) error {
	fmt.Printf("📝 Registering SSH key: %s\n", name)
	_, err := dag.Container().
		From("digitalocean/doctl:1.101.0").
		WithSecretVariable("DIGITALOCEAN_ACCESS_TOKEN", do.token).
		WithExec([]string{
			"compute",
			"ssh-key",
			"create",
			name,
			"--public-key", publicKey,
			"--format", "ID",
			"--no-header",
		}).
		Stdout(ctx)
	return err
}
