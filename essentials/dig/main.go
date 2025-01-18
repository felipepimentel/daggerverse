package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/essentials/dig/internal/dagger"
)

// Dig provides functionality for DNS queries
type Dig struct {
	client *dagger.Client
}

// QueryConfig holds configuration for DNS queries
type QueryConfig struct {
	Domain    string
	Type      string
	Server    string
	Port      int
	Short     bool
	Trace     bool
	Timeout   int
	Retries   int
	NoRecurse bool
}

// New creates a new instance of the Dig module
func New() *Dig {
	return &Dig{}
}

// Query performs a DNS query with the specified configuration
func (d *Dig) Query(ctx context.Context, config QueryConfig) (*dagger.Container, error) {
	if config.Domain == "" {
		return nil, fmt.Errorf("domain is required")
	}

	args := []string{"dig"}

	// Add query type if specified
	if config.Type != "" {
		args = append(args, config.Type)
	}

	// Add domain
	args = append(args, config.Domain)

	// Add server if specified
	if config.Server != "" {
		args = append(args, fmt.Sprintf("@%s", config.Server))
		if config.Port > 0 {
			args = append(args, fmt.Sprintf("-p %d", config.Port))
		}
	}

	// Add short output option
	if config.Short {
		args = append(args, "+short")
	}

	// Add trace option
	if config.Trace {
		args = append(args, "+trace")
	}

	// Add timeout
	if config.Timeout > 0 {
		args = append(args, fmt.Sprintf("+time=%d", config.Timeout))
	}

	// Add retries
	if config.Retries > 0 {
		args = append(args, fmt.Sprintf("+tries=%d", config.Retries))
	}

	// Add no recurse option
	if config.NoRecurse {
		args = append(args, "+norecurse")
	}

	return dag.Container().
		From("alpine/bind-tools:latest").
		WithExec(args), nil
}

// A performs an A record lookup
func (d *Dig) A(ctx context.Context, domain string) (*dagger.Container, error) {
	return d.Query(ctx, QueryConfig{
		Domain:  domain,
		Type:    "A",
		Short:   true,
		Timeout: 5,
		Retries: 3,
	})
}

// AAAA performs an AAAA record lookup
func (d *Dig) AAAA(ctx context.Context, domain string) (*dagger.Container, error) {
	return d.Query(ctx, QueryConfig{
		Domain:  domain,
		Type:    "AAAA",
		Short:   true,
		Timeout: 5,
		Retries: 3,
	})
}

// MX performs an MX record lookup
func (d *Dig) MX(ctx context.Context, domain string) (*dagger.Container, error) {
	return d.Query(ctx, QueryConfig{
		Domain:  domain,
		Type:    "MX",
		Short:   true,
		Timeout: 5,
		Retries: 3,
	})
}

// NS performs an NS record lookup
func (d *Dig) NS(ctx context.Context, domain string) (*dagger.Container, error) {
	return d.Query(ctx, QueryConfig{
		Domain:  domain,
		Type:    "NS",
		Short:   true,
		Timeout: 5,
		Retries: 3,
	})
}

// TXT performs a TXT record lookup
func (d *Dig) TXT(ctx context.Context, domain string) (*dagger.Container, error) {
	return d.Query(ctx, QueryConfig{
		Domain:  domain,
		Type:    "TXT",
		Short:   true,
		Timeout: 5,
		Retries: 3,
	})
}

// Trace performs a trace of the DNS resolution
func (d *Dig) Trace(ctx context.Context, domain string) (*dagger.Container, error) {
	return d.Query(ctx, QueryConfig{
		Domain:  domain,
		Trace:   true,
		Timeout: 10,
		Retries: 3,
	})
}

// VerifyPropagation checks if a DNS record has propagated
func (d *Dig) VerifyPropagation(ctx context.Context, domain string, expectedValue string) (*dagger.Container, error) {
	return d.Query(ctx, QueryConfig{
		Domain:  domain,
		Type:    "A",
		Short:   true,
		Timeout: 5,
		Retries: 5,
	})
} 