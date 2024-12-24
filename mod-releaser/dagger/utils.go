package main

import (
	"context"
	"dagger/mod-releaser/internal/dagger"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type daggerManifest struct {
	Name          string `json:"name"`
	Sdk           string `json:"sdk"`
	Source        string `json:"source"`
	EngineVersion string `json:"engineVersion"`
}

func (m *ModReleaser) bumpVersion(major, minor, patch bool, customMsg string) (*ModReleaser, error) {
	var msg string
	var firstRelease bool
	prefixTag := m.Component + "/v"
	if len(m.Tags) == 0 {
		msg = "New component: " + m.Component
		m.Tag = prefixTag + "0.1.0"
		firstRelease = true
	}

	if !firstRelease {
		verMajor, verMinor, verPatch, err := m.parseSemver(strings.Split(strings.Split(m.Tags[len(m.Tags)-1], "/v")[1], "."))
		if err != nil {
			return nil, err
		}

		switch {
		case major:
			verMajor++
		case minor:
			verMinor++
		case patch:
			verPatch++
		default:
			return nil, fmt.Errorf("'major', 'minor', or 'patch' should be set to true")
		}

		m.Tag = prefixTag + fmt.Sprintf("%d.%d.%d", verMajor, verMinor, verPatch)
		msg = "New release " + m.Tag
	}

	if customMsg != "" {
		msg = customMsg
	}

	m.Tags = append(m.Tags, m.Tag)

	return m.WithContainer(m.Ctr.WithExec([]string{"git", "tag", "-a", m.Tag, "-m", msg})), nil
}

func (m *ModReleaser) parseSemver(semverParts []string) (int, int, int, error) {
	major, err := strconv.Atoi(semverParts[0])
	if err != nil {
		return 0, 0, 0, err
	}
	minor, err := strconv.Atoi(semverParts[1])
	if err != nil {
		return 0, 0, 0, err
	}
	patch, err := strconv.Atoi(semverParts[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return major, minor, patch, nil
}

func (m *ModReleaser) fetchTags(ctx context.Context) error {
	versionRegexp, err := regexp.Compile(m.Component + "/v\\d+\\.\\d+\\.\\d+")
	if err != nil {
		return err
	}

	output, err := m.Ctr.WithExec([]string{"git", "tag", "-l"}).Stdout(ctx)
	if err != nil {
		return err
	}

	for _, tag := range strings.Split(output, "\n") {
		if versionRegexp.MatchString(tag) {
			m.Tags = append(m.Tags, tag)
		}
	}

	slices.Sort(m.Tags)

	return nil
}

// Allow to override the current container
func (m *ModReleaser) WithContainer(ctr *dagger.Container) *ModReleaser {
	m.Ctr = ctr

	return m
}

// Open a shell
func (m *ModReleaser) Shell() *dagger.Container {
	return m.Ctr.WithDefaultTerminalCmd(nil).Terminal()
}
