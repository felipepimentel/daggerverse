// Yq runtime

package main

import (
	"context"
	"fmt"
	"main/internal/dagger"
	"strings"
)

const (
	defaultPath = "/opt/"
)

type Yq struct {
	// +private
	Ctr *dagger.Container
}

func New(
	// The image to use for yq
	// +optional
	// +default="mikefarah/yq"
	image string,
	// The version of the image to use
	// +optional
	// +default="4.35.2"
	version string,
	// The source where yaml files are stored
	source *dagger.Directory,
) *Yq {
	return &Yq{
		Ctr: dag.Container().
			From(fmt.Sprintf("%s:%s", image, version)).
			WithEntrypoint([]string{"yq"}).
			WithDirectory(defaultPath, source, dagger.ContainerWithDirectoryOpts{Owner: "yq"}).
			WithWorkdir(defaultPath),
	}
}

// Edit a yaml file following the given expression
func (y *Yq) Set(
	// The yq expression to execute
	expr,
	// The yaml file path to edit
	yamlFilePath string,
) *Yq {
	y.Ctr = y.Ctr.
		WithExec(
			[]string{
				"-i",
				expr,
				defaultPath + yamlFilePath,
			},
			dagger.ContainerWithExecOpts{UseEntrypoint: true},
		)

	return y
}

// Fetch a value from a yaml file
func (y *Yq) Get(
	ctx context.Context,
	// The yq expression to execute
	expr,
	// The yaml file path to read
	yamlFilePath string,
) (string, error) {
	val, err := y.Ctr.
		WithExec(
			[]string{
				expr,
				defaultPath + yamlFilePath,
			},
			dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)

	return strings.TrimSuffix(val, "\n"), err
}

// Override the source directory
func (y *Yq) WithDirectory(
	// The source where yaml files are stored
	source *dagger.Directory,
) *Yq {
	y.Ctr = y.Ctr.WithDirectory(defaultPath, source)
	return y
}

// Get the directory given to Yq
func (y *Yq) State() *dagger.Directory {
	return y.Ctr.Directory(defaultPath)
}

// Get the yq container
func (y *Yq) Container() *dagger.Container {
	return y.Ctr
}

// Open a shell in the current container
func (y *Yq) Shell() *dagger.Container {
	return y.Ctr.Terminal()
}
