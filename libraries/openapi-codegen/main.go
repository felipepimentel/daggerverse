// A Dagger module to generate code for OpenAPI specifications
package main

import (
	"github.com/felipepimentel/daggerverse/libraries/openapi-codegen/internal/dagger"
)

type OpenApiCodegen struct{}

// Returns a directory with the generated code from the OpenAPI spec
func (m *OpenApiCodegen) OpenApiCodegen(
	// The OpenAPI spec file to generate the code from
	Spec *dagger.File,
	// The Generator to use (e.g. "go" | "rust" | ...)
	Generator string,
) *dagger.Directory {
	return dag.Container().
		From("openapitools/openapi-generator-cli").
		WithFile("/codegen/spec", Spec).
		WithExec([]string{"/usr/local/bin/docker-entrypoint.sh", "generate", "-i", "/codegen/spec", "-g", Generator, "-o", "/codegen/out"}).
		Directory("codegen/out")
}
