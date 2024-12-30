package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

type PythonPipeline struct {
	PackagePath string // Caminho para o pacote Python dentro do source
}

// CICD executa o pipeline completo de CI/CD
func (m *PythonPipeline) CICD(ctx context.Context, source *dagger.Directory, token *dagger.Secret) (string, error) {
	client, err := dagger.Connect(ctx)
	if err != nil {
		return "", fmt.Errorf("error connecting to dagger: %v", err)
	}
	defer client.Close()

	// Build e teste
	builder := client.Container().
		From("python:3.11-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src/" + m.PackagePath).
		WithExec([]string{"pip", "install", "poetry"}).
		WithExec([]string{"poetry", "install", "--no-interaction"}).
		WithExec([]string{"poetry", "run", "pytest"}).
		WithExec([]string{"poetry", "build"})

	if _, err := builder.Sync(ctx); err != nil {
		return "", fmt.Errorf("error building package: %v", err)
	}

	// Bump version
	versioner := client.Container().
		From("node:lts-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithEnvVariable("GIT_AUTHOR_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_AUTHOR_EMAIL", "github-actions[bot]@users.noreply.github.com").
		WithEnvVariable("GIT_COMMITTER_NAME", "github-actions[bot]").
		WithEnvVariable("GIT_COMMITTER_EMAIL", "github-actions[bot]@users.noreply.github.com").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "openssh-client"}).
		WithExec([]string{
			"npm", "install", "-g",
			"semantic-release",
			"@semantic-release/commit-analyzer",
			"@semantic-release/release-notes-generator",
			"@semantic-release/changelog",
			"@semantic-release/git",
			"@semantic-release/github",
		}).
		WithExec([]string{"git", "config", "--global", "user.email", "github-actions[bot]@users.noreply.github.com"}).
		WithExec([]string{"git", "config", "--global", "user.name", "github-actions[bot]"})

	// Criar package.json
	packageJSON := `{
		"name": "@daggerverse/python",
		"version": "0.0.0-development",
		"private": true,
		"release": {
			"branches": ["main"],
			"plugins": [
				"@semantic-release/commit-analyzer",
				"@semantic-release/release-notes-generator",
				"@semantic-release/changelog",
				["@semantic-release/git", {
					"assets": ["CHANGELOG.md"],
					"message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}"
				}],
				["@semantic-release/github", {
					"assets": []
				}]
			]
		}
	}`

	versioner = versioner.
		WithExec([]string{"bash", "-c", fmt.Sprintf("echo '%s' > /src/package.json", packageJSON)})

	versionOutput, err := versioner.
		WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN")).
		WithEnvVariable("GH_TOKEN", os.Getenv("GITHUB_TOKEN")).
		WithExec([]string{
			"npx", "semantic-release",
			"--branches", "main",
			"--ci", "false",
			"--debug",
		}).Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("error bumping version: %v", err)
	}

	// Publicar no PyPI
	tokenValue, err := token.Plaintext(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting token value: %v", err)
	}

	publisher := client.Container().
		From("python:3.11-slim").
		WithDirectory("/src", source).
		WithWorkdir("/src/" + m.PackagePath).
		WithExec([]string{"pip", "install", "poetry"}).
		WithExec([]string{"poetry", "config", "pypi-token.pypi", tokenValue}).
		WithExec([]string{"poetry", "publish", "--build"})

	if _, err := publisher.Sync(ctx); err != nil {
		return "", fmt.Errorf("error publishing package: %v", err)
	}

	return versionOutput, nil
} 