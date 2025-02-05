package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/felipepimentel/daggerverse/essentials/ssh/internal/dagger"
)

type Ssh struct {
	BaseCtr     *dagger.Container
	Destination string
	Opts        []SshOpts
}

type SshOpts struct {
	IdentityFile *dagger.Secret
	Port         int
	Login        string
}

// SSHKeys represents the SSH keys used for deployment
type SSHKeys struct {
	name       string
	publicKey  string
	privateKey string
	keyID      string
}

func New(destination string, identityFile *dagger.Secret) (*Ssh, error) {
	baseCtr := dag.Container().From("alpine:3").WithExec([]string{"apk", "add", "--no-cache", "openssh-client"})

	opts := SshOpts{
		IdentityFile: identityFile,
	}

	ssh := &Ssh{
		Destination: destination,
		Opts:        []SshOpts{opts},
		BaseCtr:     baseCtr,
	}
	return ssh, nil
}

func (m *Ssh) GenerateKeys(ctx context.Context) (*SSHKeys, error) {
	fmt.Println("🔑 Generating SSH keys...")
	keyName := fmt.Sprintf("n8n-deploy-%d", time.Now().Unix())

	// Create SSH key pair in a container with proper permissions
	sshContainer := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "openssh"}).
		WithExec([]string{"mkdir", "-p", "/root/.ssh"}).
		WithExec([]string{"chmod", "700", "/root/.ssh"}).
		WithExec([]string{
			"ssh-keygen",
			"-t", "ed25519",
			"-f", "/root/.ssh/id_ed25519",
			"-C", keyName,
			"-N", "",
		})

	// Read the public key
	publicKey, err := sshContainer.WithExec([]string{
		"cat", "/root/.ssh/id_ed25519.pub",
	}).Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	// Read the private key
	privateKey, err := sshContainer.WithExec([]string{
		"cat", "/root/.ssh/id_ed25519",
	}).Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	fmt.Println("✅ SSH keys generated successfully")
	return &SSHKeys{
		name:       keyName,
		publicKey:  strings.TrimSpace(publicKey),
		privateKey: strings.TrimSpace(privateKey),
	}, nil
}

// example usage: "dagger call --destination USER@HOST --identity-file file:${HOME}/.ssh/id_ed25519 command --args whoami stdout"
func (m *Ssh) Command(args ...string) *dagger.Container {
	ctr := m.BaseCtr

	execArgs := []string{"/usr/bin/ssh", "-o", "StrictHostKeyChecking=no"}
	for i, o := range m.Opts {
		if o.IdentityFile != nil {
			// this allows to support several keys if many opts were passed
			keyPath := fmt.Sprintf("/key_%d", i)
			ctr = ctr.WithMountedSecret(keyPath, o.IdentityFile)
			execArgs = append(execArgs, "-i", keyPath)
		}
		if o.Login != "" {
			execArgs = append(execArgs, "-l", o.Login)
		}
		if o.Port > 0 {
			execArgs = append(execArgs, "-p", strconv.Itoa(o.Port))
		}
	}

	// add the destination address after the ssh args
	execArgs = append(execArgs, m.Destination)
	// add the command args
	execArgs = append(execArgs, args...)

	return ctr.WithExec(execArgs)
}
