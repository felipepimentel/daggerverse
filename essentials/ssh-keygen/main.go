package main

import (
	"context"
	"crypto"
	"encoding/pem"

	"golang.org/x/crypto/ssh"

	"github.com/felipepimentel/daggerverse/essentials/ssh-keygen/internal/dagger"
)

// Generate a new SSH key pair.
type SshKeygen struct{}

// An SSH key pair.
type KeyPair struct {
	PublicKey  *dagger.File
	PrivateKey *dagger.Secret
}

func formatKeys(
	ctx context.Context,
	name string,
	publicKey crypto.PublicKey,
	privateKey crypto.PrivateKey,
	passphrase *dagger.Secret,
) (*KeyPair, error) {
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	var sshPrivateKey *pem.Block
	{
		var err error

		if passphrase != nil {
			pass, err := passphrase.Plaintext(ctx)
			if err != nil {
				return nil, err
			}

			sshPrivateKey, err = ssh.MarshalPrivateKeyWithPassphrase(privateKey, "", []byte(pass))
			if err != nil {
				return nil, err
			}
		} else {
			sshPrivateKey, err = ssh.MarshalPrivateKey(privateKey, "")
			if err != nil {
				return nil, err
			}
		}
	}

	return &KeyPair{
		PublicKey:  dag.Directory().WithNewFile(name+".pub", string(ssh.MarshalAuthorizedKey(sshPublicKey))).File(name + ".pub"),
		PrivateKey: dag.SetSecret(name, string(pem.EncodeToMemory(sshPrivateKey))),
	}, nil
}
