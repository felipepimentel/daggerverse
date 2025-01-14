/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipepimentel/daggerverse/pipelines/crossplane/internal/dagger"
	reg "github.com/felipepimentel/daggerverse/pipelines/crossplane/registry"
	"github.com/felipepimentel/daggerverse/pipelines/crossplane/templates"
)

type Crossplane struct {
	XplaneContainer *dagger.Container
}

// Package Crossplane Package
func (m *Crossplane) Package(ctx context.Context, src *dagger.Directory) *dagger.Directory {

	xplane := m.XplaneContainer.
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"crossplane", "xpkg", "build"})

	buildArtifact, err := xplane.WithExec(
		[]string{"find", "-maxdepth", "1", "-name", "*.xpkg", "-exec", "basename", "{}", ";"}).
		Stdout(ctx)

	if err != nil {
		fmt.Println("ERROR GETTING BUILD ARTIFACT: ", err)
	}

	fmt.Println("BUILD PACKAGE: ", buildArtifact)

	return xplane.Directory("/src")
}

// Push Crossplane Package
func (m *Crossplane) Push(
	ctx context.Context,
	src *dagger.Directory,
	// +default="ghcr.io"
	registry string,
	username string,
	password *dagger.Secret,
	destination string) string {

	dirWithPackage := m.Package(ctx, src)

	passwordPlaintext, err := password.Plaintext(ctx)

	configJSON, err := reg.CreateDockerConfigJSON(username, passwordPlaintext, registry)
	if err != nil {
		fmt.Printf("ERROR CREATING DOCKER config.json: %v\n", err)
	}

	status, err := m.XplaneContainer.
		WithNewFile("/root/.docker/config.json", configJSON).
		WithDirectory("/src", dirWithPackage).
		WithWorkdir("/src").
		WithExec([]string{"crossplane", "xpkg", "push", destination}).
		Stdout(ctx)

	if err != nil {
		fmt.Println("ERROR PUSHING PACKAGE: ", err)
	}

	fmt.Println("PACKAGE STATUS: ", status)

	return status
}

// GetXplaneContainer return the default image for helm
func (m *Crossplane) GetXplaneContainer() *dagger.Container {
	return dag.Container().
		From("ghcr.io/stuttgart-things/crossplane-cli:v1.18.0")
}

// Init Crossplane Package based on custom templates and a configuration file
func (m *Crossplane) InitCustomPackage(ctx context.Context, kind string) *dagger.Directory {

	// DEFINE INTERFACE MAP FOR TEMPLATE DATA INLINE - LATER LOAD AS YAML FILE
	// DEFINE A STRUCT WITH THE NEEDED PACKAGE FOLDER STRUCTURE AND TARGET PATHS
	// RENDER THE TEMPLATES WITH THE DATA
	// COPY TO CONTAINER AND RETURN OR TRY TO RETURN FOR EXPORTING WITHOUT USING A CONTAINER

	xplane := m.XplaneContainer

	packageName := strings.ToLower(kind)
	workingDir := "/" + packageName + "/"

	// Data to be used with the template
	data := map[string]interface{}{
		"namespace":           "crossplane-system",
		"claimName":           "incluster",
		"apiGroup":            "resources.stuttgart-things.com",
		"claimApiVersion":     "v1alpha1",
		"maintainer":          "patrick.hermann@sva.de",
		"source":              "github.com/stuttgart-things/stuttgart-things",
		"license":             "Apache-2.0",
		"crossplaneVersion":   ">=v1.14.1-0",
		"kindLower":           packageName,
		"kindLowerX":          "x" + packageName,
		"kind":                "X" + ToTitle(packageName),
		"plural":              "x" + packageName + "s",
		"claimKind":           ToTitle(packageName),
		"claimPlural":         packageName + "s",
		"compositeApiVersion": "apiextensions.crossplane.io/v1",
	}

	fmt.Println("KINDS: ", data)

	for _, template := range templates.PackageFiles {
		rendered := templates.RenderTemplate(template.Template, data)
		xplane = xplane.WithNewFile(workingDir+template.Destination, rendered)
	}

	return xplane.Directory(workingDir)
}

// Init Crossplane Package
func (m *Crossplane) InitPackage(ctx context.Context, name string) *dagger.Directory {

	output := m.XplaneContainer.
		WithExec([]string{"crossplane", "xpkg", "init", name, "configuration-template", "-d", name}).
		WithExec([]string{"ls", "-lta", name}).
		WithExec([]string{"rm", "-rf", name + "/NOTES.txt"})

	return output.Directory(name)
}

func New(
	// xplane container
	// It need contain xplane
	// +optional
	xplaneContainer *dagger.Container,

) *Crossplane {
	xplane := &Crossplane{}

	if xplaneContainer != nil {
		xplane.XplaneContainer = xplaneContainer
	} else {
		xplane.XplaneContainer = xplane.GetXplaneContainer()
	}
	return xplane
}

func ToTitle(str string) string {
	letters := strings.Split(str, "")
	return strings.ToUpper(letters[0]) + strings.Join(letters[1:], "")
}
