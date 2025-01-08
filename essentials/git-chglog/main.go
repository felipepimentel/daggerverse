// A Dagger Module to install and run the git-chglog CLI.
//
// Git-Changelog is a command-line tool that allows you to generate a changelog from a git repository.
// See https://github.com/git-chglog/git-chglog for more information.
package main

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/felipepimentel/daggerverse/essentials/git-chglog/internal/dagger"
)

// GitChglog is a Dagger Module to interact with the git-chglog CLI.
type GitChglog struct {
	ImageRepository string
	ImageTag        string
	ChglogDir       *dagger.Directory
}

func New(
	// image repository of git-chglog.
	// +optional
	// +default="quay.io/git-chglog/git-chglog"
	imageRepository string,
	// image tag of git-chglog.
	// See https://github.com/git-chglog/git-chglog/releases for available versions.
	// Run `crane digest quay.io/git-chglog/git-chglog:0.15.4` to get the digest.
	// +optional
	// +default="0.15.4@sha256:c791b1e8264387690cce4ce32e18b4f59ca3ffd8d55cb4093dc6de74529493f4"
	imageTag string,
	// directory containing the chglog config and template.
	// See https://github.com/git-chglog/git-chglog/#configuration
	// and https://github.com/git-chglog/git-chglog/#templates
	// If empty, a default config and template will be used.
	// +optional
	chglogDir *dagger.Directory,
) *GitChglog {
	return &GitChglog{
		ImageRepository: imageRepository,
		ImageTag:        imageTag,
		ChglogDir:       chglogDir,
	}
}

// Container returns a container with git-chglog installed and configured.
func (g *GitChglog) Container(
	ctx context.Context,
	// git directory to include in the container.
	// +optional
	gitDirectory *dagger.Directory,
) *dagger.Container {
	ctr := dag.Container().
		From(g.ImageRepository + ":" + g.ImageTag)

	if g.ChglogDir == nil {
		g.ChglogDir = dag.CurrentModule().Source().Directory("templates")
	}

	if gitDirectory != nil {
		// It has to be "/workdir"
		// See https://github.com/git-chglog/git-chglog/blob/master/Dockerfile
		ctr = ctr.WithMountedDirectory("/workdir/.git", gitDirectory)

		// dynamically set the repository URL in the config file
		tplDir, err := replaceRepositoryURL(ctx, ctr, g.ChglogDir)
		if err == nil {
			g.ChglogDir = tplDir
		}
	}

	ctr = ctr.WithMountedDirectory("/workdir/.chglog", g.ChglogDir)

	return ctr
}

// Changelog generates a changelog for the given git repository and version.
func (g *GitChglog) Changelog(
	ctx context.Context,
	// git directory
	gitDirectory *dagger.Directory,
	// version to generate the changelog for.
	// See https://github.com/git-chglog/git-chglog/#cli-usage for the supported formats.
	version string,
	// container to use for the command, instead of the default container
	// you can use this to customize the container
	// +optional
	ctr *dagger.Container,
) *dagger.File {
	if ctr == nil {
		ctr = g.Container(ctx, gitDirectory)
	}

	changelogFilePath := "CHANGELOG.md"

	return ctr.
		WithExec([]string{
			"--output", changelogFilePath,
			"--no-color",
			"--no-emoji",
			version,
		},
			dagger.ContainerWithExecOpts{
				UseEntrypoint: true,
			}).
		File(changelogFilePath)
}

func replaceRepositoryURL(ctx context.Context, ctr *dagger.Container, tplDir *dagger.Directory) (*dagger.Directory, error) {
	repoURL, err := retrieveRepositoryURL(ctx, ctr)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repository URL: %w", err)
	}

	content, err := tplDir.File("config.yml").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yml: %w", err)
	}

	content = strings.Replace(content, "REPOSITORY_URL", repoURL, -1)

	return tplDir.
		WithoutFile("config.yml").
		WithNewFile("config.yml", content), nil
}

var (
	gitHostnameRegexp = regexp.MustCompile(`git@([^:]+):`)
)

func retrieveRepositoryURL(ctx context.Context, ctr *dagger.Container) (string, error) {
	gitURL, err := ctr.WithExec([]string{
		"git", "config", "--get", "remote.origin.url",
	}).Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get git remote URL: %w", err)
	}

	gitURL = strings.TrimSpace(gitURL)
	gitURL = strings.TrimSuffix(gitURL, ".git")

	repoURL := gitHostnameRegexp.ReplaceAllStringFunc(gitURL, func(match string) string {
		hostname := gitHostnameRegexp.FindStringSubmatch(match)[1]
		return "https://" + hostname + "/"
	})

	u, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("invalid repository URL %q: %w", repoURL, err)
	}

	// ensure we won't leak the user info
	u.User = nil

	return u.String(), nil
}
