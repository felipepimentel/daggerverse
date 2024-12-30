package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	python "github.com/felipepimentel/daggerverse/python"
)

func main() {
	ctx := context.Background()

	// Initialize Dagger client
	client, err := dagger.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Get source code directory
	source := client.Host().Directory(".")

	// Create Python module instance
	python := &python.Python{
		PackagePath: ".", // Path to your Python package
	}

	// Get PyPI token from environment
	token := client.SetSecret("PYPI_TOKEN", os.Getenv("PYPI_TOKEN"))

	// Run complete CI/CD pipeline
	version, err := python.CICD(ctx, source, token)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully released version %s\n", version)
} 