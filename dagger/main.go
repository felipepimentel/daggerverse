package main

import (
	"context"
	"fmt"

	"github.com/felipepimentel/daggerverse/dagger/internal/dagger"
)

type Daggerverse struct {
	Source   *dagger.Directory
	RepoName string
}

func New(
	// +optional
	// +defaultPath="/"
	source *dagger.Directory,
	// +optional
	// +default="vbehar/daggerverse"
	repoName string,
) *Daggerverse {
	return &Daggerverse{
		Source:   source,
		RepoName: repoName,
	}
}

func (d *Daggerverse) Release(
	ctx context.Context,
	gitToken *dagger.Secret,
) (string, error) {
	nextVersion, err := dag.JxReleaseVersion().NextVersion(ctx,
		d.Source.Directory(".git"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to determine next version: %w", err)
	}

	tagName := "v" + nextVersion

	err = dag.Gh(dagger.GhOpts{
		Token:  gitToken,
		Source: d.Source,
		Repo:   d.RepoName,
	}).Release().Create(ctx, tagName, tagName, dagger.GhReleaseCreateOpts{
		GenerateNotes: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create release %q: %w", tagName, err)
	}

	return tagName, nil
}

func (d *Daggerverse) Publish(
	ctx context.Context,
	// +optional
	// +default=false
	dryRun bool,
) ([]string, error) {
	return dag.DaggerverseCockpit().Publish(ctx, d.Source, dagger.DaggerverseCockpitPublishOpts{
		DryRun: dryRun,
		Exclude: []string{
			"dagger.json", // don't include our own CI
			"artifactory/examples/go",
			"jfrogcli/examples/go",
		},
	})
}
