package main_test

import (
	"context"
	"flag"
	"os"
	"path/filepath"
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

func Test_Editorconfig(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_mounted_host_directory_without_git_directory", func(t *testing.T) {
		t.Parallel()
		container := base("")
		require.NotNil(t, container)

		dir, err := os.Getwd()
		require.NoError(t, err)
		require.NotEmpty(t, dir)

		out, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory(filepath.Join(dir, ".."), dagger.HostDirectoryOpts{
				Exclude: []string{".git"},
			})).
			WithWorkdir("/tmp").
			WithExec([]string{"ec", "-dry-run"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotNil(t, out)
		require.NotRegexp(t, `\.git/`, out)
	})
	t.Run("Test_mounted_host_directory", func(t *testing.T) {
		t.Parallel()
		container := base("")
		require.NotNil(t, container)

		out, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"ec", "-dry-run"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `t\.txt`, out)
	})
	t.Run("Test_editorconfig-checker_help", func(t *testing.T) {
		t.Parallel()
		container := base("")
		require.NotNil(t, container)

		_, err := container.
			WithExec([]string{"ec", "-help"}).
			Stderr(ctx)
		require.NoError(t, err)
	})
}

func base(
	image string,
) *dagger.Container {

	defaultImageRepository := "mstruebing/editorconfig-checker"
	var ctr *dagger.Container

	if image != "" {
		ctr = c.Container().From(image)
	} else {
		ctr = c.Container().From(defaultImageRepository)
	}

	return ctr
}
