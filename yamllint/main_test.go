package main_test

import (
	"context"
	"flag"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

var c *dagger.Client

func TestMain(m *testing.M) {
	flag.Parse()

	ctx := context.Background()

	c, _ = dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	defer c.Close()

	code := m.Run()
	defer c.Close()
	os.Exit(code)
}

func Test_yamllint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_yamllint_directory_with_host_directory", func(t *testing.T) {
		t.Parallel()
		container := base("")
		require.NotNil(t, container)

		_, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithExec([]string{"yamllint",
				"--config-data",
				"{extends: default, rules: {line-length: {level: warning}}}",
				"--no-warnings",
				"."}).
			Stderr(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 1")
	})
	t.Run("Test_yamllint_error", func(t *testing.T) {
		t.Parallel()
		container := base("")
		require.NotNil(t, container)

		// editorconfig-checker-disable
		badYAML := `
	foo: "bar"
	  baz: "qux"
	- item1
	- item2
	key : value
	`
		// editorconfig-checker-enable

		container = container.WithNewFile("/tmp/bad.yaml", badYAML)
		_, err := container.
			WithExec([]string{"yamllint",
				"--config-data",
				"{extends: default, rules: {line-length: {level: warning}}}",
				"--no-warnings",
				"/tmp/bad.yaml"}).
			Stderr(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 1")
	})
	t.Run("Test_yamllint_version", func(t *testing.T) {
		t.Parallel()
		container := base("")
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"yamllint", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
}

func base(
	image string,
) *dagger.Container {

	defaultImageRepository := "pipelinecomponents/yamllint"
	var ctr *dagger.Container

	if image != "" {
		ctr = c.Container().From(image)
	} else {
		ctr = c.Container().From(defaultImageRepository)
	}

	return ctr
}
