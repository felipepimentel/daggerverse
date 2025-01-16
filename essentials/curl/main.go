package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipepimentel/daggerverse/essentials/curl/internal/dagger"
)

// Curl provides functionality for making HTTP requests
type Curl struct {
	client *dagger.Client
}

// Header represents an HTTP header
type Header struct {
	Key   string
	Value string
}

// RequestConfig holds configuration for making HTTP requests
type RequestConfig struct {
	URL            string
	Method         string
	Headers        []Header
	FollowRedirect bool
	Insecure       bool
	Output         string
	Timeout        int
	RetryAttempts  int
	RetryDelay     int
}

// New creates a new instance of the Curl module
func New() *Curl {
	return &Curl{}
}

// Request makes an HTTP request with the specified configuration
func (c *Curl) Request(ctx context.Context, config RequestConfig) (*dagger.Container, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}

	args := []string{"curl"}

	// Add method if specified
	if config.Method != "" {
		args = append(args, "-X", strings.ToUpper(config.Method))
	}

	// Add headers
	for _, header := range config.Headers {
		args = append(args, "-H", fmt.Sprintf("%s: %s", header.Key, header.Value))
	}

	// Add follow redirect option
	if config.FollowRedirect {
		args = append(args, "-L")
	}

	// Add insecure option
	if config.Insecure {
		args = append(args, "-k")
	}

	// Add output option
	if config.Output != "" {
		args = append(args, "-o", config.Output)
	}

	// Add timeout
	if config.Timeout > 0 {
		args = append(args, "--max-time", fmt.Sprintf("%d", config.Timeout))
	}

	// Add retry attempts
	if config.RetryAttempts > 0 {
		args = append(args, "--retry", fmt.Sprintf("%d", config.RetryAttempts))
	}

	// Add retry delay
	if config.RetryDelay > 0 {
		args = append(args, "--retry-delay", fmt.Sprintf("%d", config.RetryDelay))
	}

	// Add URL
	args = append(args, config.URL)

	return dag.Container().
		From("curlimages/curl:latest").
		WithExec(args), nil
}

// Head makes a HEAD request to check endpoint availability
func (c *Curl) Head(ctx context.Context, url string) (*dagger.Container, error) {
	return c.Request(ctx, RequestConfig{
		URL:            url,
		Method:         "HEAD",
		FollowRedirect: true,
		Timeout:        10,
		RetryAttempts:  3,
		RetryDelay:     5,
	})
}

// Get makes a GET request to retrieve content
func (c *Curl) Get(ctx context.Context, url string) (*dagger.Container, error) {
	return c.Request(ctx, RequestConfig{
		URL:            url,
		Method:         "GET",
		FollowRedirect: true,
		Timeout:        30,
		RetryAttempts:  3,
		RetryDelay:     5,
	})
}

// HealthCheck performs a health check on an endpoint
func (c *Curl) HealthCheck(ctx context.Context, url string) (*dagger.Container, error) {
	return c.Request(ctx, RequestConfig{
		URL:            url,
		Method:         "GET",
		Headers:        []Header{{Key: "Accept", Value: "application/json"}},
		FollowRedirect: true,
		Timeout:        5,
		RetryAttempts:  3,
		RetryDelay:     2,
	})
} 