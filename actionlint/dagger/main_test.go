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

func Test_Actionlint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_actionlint_error", func(t *testing.T) {
		t.Parallel()
		container := base("")
		require.NotNil(t, container)

		_, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"sh", "-c", "find . -type f -name '*.yml' -print0 | xargs -0 actionlint"}).
			Stderr(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 123")
	})
	t.Run("Test_actionlint_version", func(t *testing.T) {
		t.Parallel()
		container := base("")
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"actionlint", "-version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
	t.Run("Test_shellcheck_version", func(t *testing.T) {
		t.Parallel()
		container := base("")
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"shellcheck", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "version")
	})
}

func base(
	image string,
) *dagger.Container {

	defaultImageRepository := "rhysd/actionlint"
	var ctr *dagger.Container

	if image != "" {
		ctr = c.Container().From(image)
	} else {
		ctr = c.Container().From(defaultImageRepository)
	}

	return ctr
}
