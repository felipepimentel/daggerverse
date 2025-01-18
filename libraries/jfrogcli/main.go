// A Dagger Module to install and run the JFrog CLI.
//
// JFrog CLI is a command-line tool that allows you to interact with JFrog products,
// such as Artifactory and Xray.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/felipepimentel/daggerverse/libraries/jfrogcli/internal/dagger"
)

const (
	gitHubReleasesURL = "https://api.github.com/repos/jfrog/jfrog-cli/releases/latest"
	fallbackVersion   = "2.71.5" // from https://github.com/jfrog/jfrog-cli/releases
	binaryFileURLTpl  = "https://releases.jfrog.io/artifactory/jfrog-cli/v2-jf/%s/jfrog-cli-%s/jf"

	// use fixed base images for reproductible builds and improved caching
	// the base image: https://images.chainguard.dev/directory/image/wolfi-base/overview
	// retrieve the latest sha256 hash with: `crane digest cgr.dev/chainguard/wolfi-base:latest`
	// and to retrieve its creation time: `crane config cgr.dev/chainguard/wolfi-base:latest | jq .created`
	// This one is from 2024-11-12T16:53:09Z
	baseWolfiImage = "cgr.dev/chainguard/wolfi-base:latest@sha256:b3dd9cf08283b959c6a0a3c833e68b2882a50129930215060154b43ae6a3e81c"
)

// Jfrogcli is a Dagger Module to install and run the JFrog CLI.
type Jfrogcli struct {
	// Version of the JFrog CLI binary.
	Version string
}

func New(
	// version of the JFrog CLI to install. If empty, the latest version will be installed.
	// +optional
	// +default="2.71.5"
	version string,
) *Jfrogcli {
	return &Jfrogcli{
		Version: version,
	}
}

// GetLatestVersion returns the latest version of the JFrog CLI.
func (c *Jfrogcli) GetLatestVersion(ctx context.Context) (string, error) {
	body, err := dag.HTTP(gitHubReleasesURL).Contents(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get latest version from %s: %w", gitHubReleasesURL, err)
	}

	var release struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal([]byte(body), &release)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal release body: %w", err)
	}

	return release.Name, nil
}

// Install installs the JFrog CLI into the given container.
func (c *Jfrogcli) Install(
	ctx context.Context,
	// +optional
	base *dagger.Container,
) (*dagger.Container, error) {
	if c.Version == "" {
		var err error
		c.Version, err = c.GetLatestVersion(ctx)
		if err != nil || c.Version == "" {
			fmt.Println("failed to get latest version, using fallback version", fallbackVersion, err)
			c.Version = fallbackVersion
		}
	}

	ctr := base
	if ctr == nil {
		ctr = dag.Container().From(baseWolfiImage)
	}

	platform, err := ctr.Platform(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform: %w", err)
	}
	osAndArch := strings.ReplaceAll(string(platform), "/", "-")

	binURL := fmt.Sprintf(binaryFileURLTpl, c.Version, osAndArch)
	binFile := dag.HTTP(binURL)

	ctr = ctr.
		WithMountedFile("/usr/local/bin/jf", binFile).
		WithExec([]string{"chmod", "+x", "/usr/local/bin/jf"}).
		WithEnvVariable("PATH", "/usr/local/bin:$PATH", dagger.ContainerWithEnvVariableOpts{Expand: true}).
		WithEnvVariable("CI", "true").
		WithEnvVariable("JFROG_CLI_REPORT_USAGE", "false").
		WithEnvVariable("JFROG_CLI_AVOID_NEW_VERSION_WARNING", "true")

	return ctr, nil
}
