// Kubernetes native configuration management.
package main

import (
	"fmt"
	"path"
	"strings"
)

const (
	// defaultImageRepository is used when no image is specified.
	defaultImageRepository = "registry.k8s.io/kustomize/kustomize"

	// defaultVersion is used when no version is specified.
	//
	// (there is no latest tag published in the default image repository)
	defaultVersion = "v5.0.1"
)

type Kustomize struct {
	Container *Container
}

func New(
	// Version (image tag) to use from the official image repository as a base container.
	//
	// +optional
	version string,

	// Custom container to use as a base container.
	//
	// +optional
	container *Container,
) *Kustomize {
	if container == nil {
		if version == "" {
			version = defaultVersion
		}

		container = dag.Container().From(fmt.Sprintf("%s:%s", defaultImageRepository, version))
	}

	return &Kustomize{container}
}

func cleanPath(s string) string {
	s = path.Clean(s)

	for strings.HasPrefix(s, "../") {
		s = strings.TrimPrefix(s, "../")
	}

	return s
}

// Build a kustomization target from a directory or URL.
func (m *Kustomize) Build(
	source *Directory,

	// Subdirectory within the source to use as the target.
	//
	// +optional
	dir string,
) *File {
	sourcePath := "/work/src"
	output := "/work/output.yaml"

	args := []string{"build", "--output", output}

	if dir != "" {
		args = append(args, cleanPath(dir))
	}

	return m.Container.
		WithWorkdir(sourcePath).
		WithMountedDirectory(sourcePath, source).
		WithExec(args).
		File(output)
}

// Edit a kustomization file.
func (m *Kustomize) Edit(
	source *Directory,

	// Subdirectory within the source to use as the target.
	//
	// +optional
	dir string,
) *Edit {
	workdir := "/work"

	if dir != "" {
		workdir = path.Join(workdir, cleanPath(dir))
	}

	return &Edit{m.Container.WithMountedDirectory("/work", source).WithWorkdir(workdir)}
}

// Edit a kustomization file.
type Edit struct {
	// +private
	Container *Container
}

// Retrieve the source containing the modifications.
func (m *Edit) Directory() *Directory {
	return m.Container.Directory("/work")
}

// Set the value of different fields in kustomization file.
func (m *Edit) Set() *Set {
	return &Set{m.Container}
}

// Set the value of different fields in kustomization file.
type Set struct {
	// +private
	Container *Container
}

// Sets one or more commonAnnotations in kustomization.yaml.
func (m *Set) Annotation(key string, value string) *Edit {
	return &Edit{m.Container.WithExec([]string{"edit", "set", "annotation", fmt.Sprintf("%s:%s", key, value)})}
}

// Set images and their new names, new tags or digests in the kustomization file.
func (m *Set) Image(image string) *Edit {
	return &Edit{m.Container.WithExec([]string{"edit", "set", "image", image})}
}

// Set the value of the namespace field in the kustomization file.
func (m *Set) Namespace(namespace string) *Edit {
	return &Edit{m.Container.WithExec([]string{"edit", "set", "namespace", namespace})}
}

// Set the value of the nameSuffix field in the kustomization file.
func (m *Set) Namesuffix(nameSuffix string) *Edit {
	return &Edit{m.Container.WithExec([]string{"edit", "set", "namesuffix", nameSuffix})}
}
