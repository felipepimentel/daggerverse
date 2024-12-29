package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// Create Dagger client
	client, err := dagger.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Get repository root directory (2 levels up from current)
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		panic(err)
	}

	// Get source directory from repository root
	src := client.Host().Directory(repoRoot)

	// Detect changes (find modules)
	modules, err := detectModules(ctx, client, src)
	if err != nil {
		panic(err)
	}

	if len(modules) == 0 {
		fmt.Println("No modules found")
		return
	}

	// Process each module
	for _, module := range modules {
		fmt.Printf("Processing module: %s\n", module)
		if err := releaseModule(ctx, client, src, module); err != nil {
			fmt.Printf("Error processing module %s: %v\n", module, err)
			continue
		}
	}
}

func detectModules(ctx context.Context, client *dagger.Client, src *dagger.Directory) ([]string, error) {
	// Use find to locate all dagger.json files
	container := client.Container().
		From("alpine:latest").
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"find", ".", "-name", "dagger.json", "-exec", "dirname", "{}", ";"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return nil, err
	}

	// Process the output
	var modules []string
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		// Remove "./" prefix if present
		module := strings.TrimPrefix(line, "./")
		if module != "" {
			modules = append(modules, module)
		}
	}

	return modules, nil
}

func releaseModule(ctx context.Context, client *dagger.Client, src *dagger.Directory, module string) error {
	// Setup container with Node.js and required tools
	container := client.Container().
		From("node:lts-slim").
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN")).
		WithEnvVariable("MODULE_PATH", module).
		WithEnvVariable("MODULE_NAME", module)

	// Install dependencies
	container = container.WithExec([]string{
		"npm", "install", "-g",
		"semantic-release",
		"@semantic-release/changelog",
		"@semantic-release/git",
		"@semantic-release/github",
	})

	// Configure Git
	container = container.WithExec([]string{
		"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com",
	}).WithExec([]string{
		"git", "config", "--global", "user.name", "github-actions[bot]",
	})

	// Run semantic-release
	container = container.WithExec([]string{"npx", "semantic-release"})

	// Get latest tag
	tagOutput, err := container.WithExec([]string{
		"git", "describe", "--tags", "--abbrev=0", "--match", fmt.Sprintf("%s/v*", module),
	}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error getting latest tag: %v", err)
	}

	if tagOutput == "" {
		fmt.Printf("No tag found for module %s, skipping publish\n", module)
		return nil
	}

	// Publish to Daggerverse
	publishContainer := client.Container().
		From("alpine:latest").
		WithDirectory("/src", src).
		WithWorkdir(fmt.Sprintf("/src/%s", module)).
		WithExec([]string{"dagger", "publish"})

	_, err = publishContainer.Sync(ctx)
	if err != nil {
		return fmt.Errorf("error publishing to Daggerverse: %v", err)
	}

	return nil
} 