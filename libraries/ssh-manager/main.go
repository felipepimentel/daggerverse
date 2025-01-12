// Package main provides a Dagger module for managing SSH keys
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/digitalocean/godo"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
)

// SSHManager represents the SSH key management module
type SSHManager struct {
	Token string // DigitalOcean API token
}

// SSHKey represents an SSH key pair
type SSHKey struct {
	Name        string
	PublicKey   string
	PrivateKey  string
	Fingerprint string
	ID          int
}

// New creates a new SSHManager instance
func New() *SSHManager {
	return &SSHManager{}
}

// WithToken sets the DigitalOcean API token
func (s *SSHManager) WithToken(token string) *SSHManager {
	s.Token = token
	return s
}

// getClient creates a new DigitalOcean API client
func (s *SSHManager) getClient() (*godo.Client, error) {
	if s.Token == "" {
		return nil, fmt.Errorf("DigitalOcean API token is required")
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: s.Token})
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	return godo.NewClient(oauthClient), nil
}

// GenerateKey generates a new SSH key pair
func (s *SSHManager) GenerateKey(ctx context.Context, name string) (*SSHKey, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Convert private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyStr := string(pem.EncodeToMemory(privateKeyPEM))

	// Generate public key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate public key: %w", err)
	}
	publicKeyStr := string(ssh.MarshalAuthorizedKey(publicKey))

	// Create key in DigitalOcean
	client, err := s.getClient()
	if err != nil {
		return nil, err
	}

	createRequest := &godo.KeyCreateRequest{
		Name:      name,
		PublicKey: publicKeyStr,
	}

	key, _, err := client.Keys.Create(ctx, createRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create key in DigitalOcean: %w", err)
	}

	return &SSHKey{
		Name:        name,
		PublicKey:   publicKeyStr,
		PrivateKey:  privateKeyStr,
		Fingerprint: key.Fingerprint,
		ID:          key.ID,
	}, nil
}

// DeleteKey deletes an SSH key from DigitalOcean
func (s *SSHManager) DeleteKey(ctx context.Context, fingerprint string) error {
	client, err := s.getClient()
	if err != nil {
		return err
	}

	_, err = client.Keys.DeleteByFingerprint(ctx, fingerprint)
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// ListKeys lists all SSH keys in the DigitalOcean account
func (s *SSHManager) ListKeys(ctx context.Context) ([]*SSHKey, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, err
	}

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	keys, _, err := client.Keys.List(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	var sshKeys []*SSHKey
	for _, key := range keys {
		sshKeys = append(sshKeys, &SSHKey{
			Name:        key.Name,
			PublicKey:   key.PublicKey,
			Fingerprint: key.Fingerprint,
			ID:          key.ID,
		})
	}

	return sshKeys, nil
}

// GetKey gets an SSH key by fingerprint
func (s *SSHManager) GetKey(ctx context.Context, fingerprint string) (*SSHKey, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, err
	}

	key, _, err := client.Keys.GetByFingerprint(ctx, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	return &SSHKey{
		Name:        key.Name,
		PublicKey:   key.PublicKey,
		Fingerprint: key.Fingerprint,
		ID:          key.ID,
	}, nil
} 