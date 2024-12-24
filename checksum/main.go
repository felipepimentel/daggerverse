// Calculate and check the checksum of files.
package main

import (
	"dagger/checksum/internal/dagger"
	"fmt"
	"strings"
)

const alpineBaseImage = "alpine:latest"

// Calculate and check the checksum of files.
type Checksum struct{}

// Calculate the SHA-256 checksum of the given files.
func (m *Checksum) Sha256() *Sha256 {
	return &Sha256{}
}

type Sha256 struct{}

// Calculate the SHA-256 checksum of the given files.
func (m *Sha256) Calculate(
	// The files to calculate the checksum for.
	files []*dagger.File,
) *dagger.File {
	return calculate("sha256", files)
}

// Check the SHA-256 checksum of the given files.
func (m *Sha256) Check(
	// Checksum file.
	checksums *dagger.File,

	// The files to check the checksum if.
	files []*dagger.File,
) *dagger.Container {
	return check("sha256", checksums, files)
}

func calculate(algo string, files []*dagger.File) *dagger.File {
	return calculateDirectory(algo, dag.Directory().WithFiles("", files))
}

func calculateDirectory(algo string, dir *dagger.Directory) *dagger.File {
	const checksumFile = "/work/checksums.txt"

	cmd := []string{algo + "sum", "$(ls)", ">", checksumFile}

	return dag.Container().
		From(alpineBaseImage).
		WithWorkdir("/work/src").
		WithMountedDirectory("/work/src", dir).
		WithExec([]string{"sh", "-c", strings.Join(cmd, " ")}).
		File(checksumFile)
}

func check(algo string, checksums *dagger.File, files []*dagger.File) *dagger.Container {
	dir := dag.Directory()

	for _, file := range files {
		dir = dir.WithFile("", file)
	}

	return checkDirectory(algo, checksums, dir)
}

func checkDirectory(algo string, checksums *dagger.File, dir *dagger.Directory) *dagger.Container {
	dir = dir.WithFile("checksums.txt", checksums)

	return dag.Container().
		From(alpineBaseImage).
		WithWorkdir("/work").
		WithMountedDirectory("/work", dir).
		WithExec([]string{"sh", "-c", fmt.Sprintf("%ssum -w -c checksums.txt", algo)})
}
