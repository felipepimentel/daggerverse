// Module to run k6 QA tests

package main

import (
	"context"
	"dagger/k-6/internal/dagger"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type K6 struct{}

// EnvironmentVariable represents a string that follows the pattern [a-zA-Z0-9]+=[a-zA-Z0-9]+
type EnvironmentVariable string

// IsValid checks if the EnvironmentVariable adheres to the specified format
func (cs EnvironmentVariable) IsValid() bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-]+=.+$`)
	return re.MatchString(string(cs))
}

// Parse returns the key and value of the EnvironmentVariable
func (cs EnvironmentVariable) Parse() (string, string) {
	parts := strings.Split(string(cs), "=")
	if len(parts) != 2 {
		panic("Error parsing environment variable")
	}
	return parts[0], parts[1]
}

// Runs the k6 QA tests
func (m *K6) Run(
	ctx context.Context,
	// The working directory containing the script
	//+required
	workingDir *dagger.Directory,
	// k6 Script file to execute
	//+required
	script string,
	// Environment variables to set
	//+optional
	env []EnvironmentVariable,
	// Virtual users to emulate
	//+optional
	//+default=1
	vus int,
	// Duration of the test
	//+optional
	//+default="1s"
	duration string,

) *dagger.Container {
	// We use Glob over Entries because it lists files recursively
	entries, _ := workingDir.Glob(ctx, "*")

	if !slices.Contains(entries, script) {
		panic("script not found")
	}

	workingDirMountPath := "/k6"
	outputDirMountPath := "/output"

	ctr := dag.Container().
		From("ghcr.io/grafana/xk6-dashboard")

	for _, e := range env {
		if !e.IsValid() {
			panic("invalid env variable")
		}

		key, value := e.Parse()

		ctr = ctr.WithEnvVariable(key, value)
	}

	command := []string{
		"k6",
		"run",
		"--vus", strconv.Itoa(vus),
		"--duration", duration,
		"--out", fmt.Sprintf("web-dashboard=export=%s", filepath.Join(outputDirMountPath, "report.html")),
		"--summary-export", filepath.Join(outputDirMountPath, "summary.json"),
		"--console-output", filepath.Join(outputDirMountPath, "errors.txt"),
		filepath.Join(workingDirMountPath, script),
	}

	ctr = ctr.WithDirectory(workingDirMountPath, workingDir).
		WithDirectory(outputDirMountPath, dag.Directory()).
		WithUser("root").
		WithExec([]string{
			"sh",
			"-c",
			fmt.Sprintf("%s || exit 0", strings.Join(command, " ")),
		}, dagger.ContainerWithExecOpts{
			UseEntrypoint: false,
		})

	ctr, _ = ctr.Sync(ctx)

	return ctr
}
