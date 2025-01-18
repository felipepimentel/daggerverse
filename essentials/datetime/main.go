package main

import (
	"context"

	"github.com/felipepimentel/daggerverse/essentials/datetime/internal/dagger"
)

type Datetime struct {
	code string
	client *dagger.Client
}

func New() *Datetime {
	return &Datetime{
		code: "import datetime as dt; print(dt.datetime.now())",
		client: dag.InlinePython(),
	}
}

func (m *Datetime) Now(ctx context.Context) (string, error) {
	return m.client.Code(m.code).Stdout(ctx)
}
