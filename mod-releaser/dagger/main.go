// A module for handling module release in daggerverse

package main

import (
	"context"
	"dagger/mod-releaser/internal/dagger"
	"encoding/json"
	"fmt"
	"strings"
)

const workingDir = "/opt/repo/"

func New(
	ctx context.Context,
	// A git repository where the release process will be applied
	gitRepo *dagger.Directory,
	// The module name to publish
	component string,
) (*ModReleaser, error) {
	dagManifest := daggerManifest{}
	c, err := gitRepo.Directory(component).File("dagger.json").Contents(ctx)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(c), &dagManifest)
	if err != nil {
		return nil, err
	}

	releaser := &ModReleaser{
		Component: component,
		Ctr: dag.
			Container().
			From("alpine:latest").
			WithExec([]string{"apk", "add", "--no-cache", "git", "curl", "openssh-client"}).
			WithExec([]string{
				"sh", "-c",
				fmt.Sprintf(
					"curl -L https://dl.dagger.io/dagger/install.sh | DAGGER_VERSION=%s sh",
					strings.Split(dagManifest.EngineVersion, "v")[1]),
			}).
			WithDirectory(workingDir, gitRepo).
			WithWorkdir(workingDir),
	}

	err = releaser.fetchTags(ctx)
	if err != nil {
		return nil, err
	}

	return releaser, nil
}

type ModReleaser struct {
	Tags []string
	Tag  string
	// +private
	Ctr *dagger.Container
	// +private
	Component string
}

// Setup global git config, it won't affect the git config of the local repository
func (m *ModReleaser) WithGitConfig(
	// A path to a git config file to use
	// +optional
	cfg *dagger.File,
	// the email to use in the git config
	// +optional
	email string,
	// the username to use in the git config
	// +optional
	name string,
) *ModReleaser {
	if cfg != nil {
		m.WithContainer(m.Ctr.WithFile("/etc/gitconfig", cfg))
	}

	if email != "" {
		m.WithContainer(m.Ctr.WithExec([]string{"git", "config", "--global", "user.email", email}))
	}

	if name != "" {
		m.WithContainer(m.Ctr.WithExec([]string{"git", "config", "--global", "user.name", name}))
	}

	return m
}

// Mount ssh keys from the host
func (m *ModReleaser) WithSshKeys(
	// The directory with ssh keys to mount
	src *dagger.Directory,
) *ModReleaser {
	return m.WithContainer(m.Ctr.WithDirectory("/root/.ssh", src))
}

// Select a specific git branch
func (m *ModReleaser) WithBranch(
	// Define the branch from where to publish
	// +optional
	// +default="main"
	branch string,
) *ModReleaser {
	return m.WithContainer(m.Ctr.WithExec([]string{"git", "checkout", branch}))
}

// Increase the major version
func (m *ModReleaser) Major(
	// Define a custom message for the git tag otherwise it will be the default from the function
	// +optional
	msg string,
) (*ModReleaser, error) {
	return m.bumpVersion(true, false, false, msg)
}

// Increase the minor version
func (m *ModReleaser) Minor(
	// Define a custom message for the git tag otherwise it will be the default from the function
	// +optional
	msg string,
) (*ModReleaser, error) {
	return m.bumpVersion(false, true, false, msg)
}

// Increase the patch version
func (m *ModReleaser) Patch(
	// Define a custom message for the git tag otherwise it will be the default from the function
	// +optional
	msg string,
) (*ModReleaser, error) {
	return m.bumpVersion(false, false, true, msg)
}

// Publish the git tag and the module
func (m *ModReleaser) Publish(
	// Indicate if the publish process should git push the tag
	// +optional
	gitPush bool,
) *ModReleaser {
	if gitPush {
		m.WithContainer(m.Ctr.WithExec([]string{"git", "push", "origin", m.Tag}))
	}

	return m.WithContainer(m.Ctr.WithExec([]string{"dagger", "publish", "-m", m.Component}, dagger.ContainerWithExecOpts{ExperimentalPrivilegedNesting: true}))
}

// Return the git repository
func (m *ModReleaser) Repository() *dagger.Directory {
	return m.Ctr.Directory(workingDir)
}

// Execute all commands
func (m *ModReleaser) Do(ctx context.Context) (string, error) {
	return m.Ctr.Stdout(ctx)
}
