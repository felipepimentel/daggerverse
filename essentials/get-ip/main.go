package main

import (
	"context"
	"time"

	"github.com/felipepimentel/daggerverse/essentials/get-ip/internal/dagger"
)

type GetIp struct {
	Ctr *dagger.Container
}

// Get IP Address of Current Container
// example usage: dagger call run
func (m *GetIp) Run(ctx context.Context) (string, error) {
	return m.Ctr.
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"curl", "https://api.ipify.org/"}).
		Stdout(ctx)
}
