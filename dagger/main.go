package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/felipepimentel/daggerverse/dagger/internal/dagger"
)

type Daggerverse struct {
	Source   *dagger.Directory
	RepoName string
}

func New(
	// +optional
	// +defaultPath="/"
	source *dagger.Directory,
	// +optional
	// +default="vbehar/daggerverse"
	repoName string,
) *Daggerverse {
	return &Daggerverse{
		Source:   source,
		RepoName: repoName,
	}
}

func (d *Daggerverse) Release(
	ctx context.Context,
	gitToken *dagger.Secret,
) (string, error) {
	nextVersion, err := dagger.JxReleaseVersion().NextVersion(ctx,
		d.Source.Directory(".git"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to determine next version: %w", err)
	}

	tagName := "v" + nextVersion

	err = dagger.Gh(dagger.GhOpts{
		Token:  gitToken,
		Source: d.Source,
		Repo:   d.RepoName,
	}).Release().Create(ctx, tagName, tagName, dagger.GhReleaseCreateOpts{
		GenerateNotes: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create release %q: %w", tagName, err)
	}

	return tagName, nil
}

func (d *Daggerverse) Publish(
	ctx context.Context,
	// +optional
	// +default=false
	dryRun bool,
) ([]string, error) {
	return dagger.DaggerverseCockpit().Publish(ctx, d.Source, dagger.DaggerverseCockpitPublishOpts{
		DryRun: dryRun,
		Exclude: []string{
			"dagger.json",
		},
	})
}

func main() {
	// Verifica se os argumentos foram passados
	if len(os.Args) < 2 {
		log.Fatalf("No command provided. Usage: <command> [options]")
	}

	// Comando principal
	command := os.Args[1]

	daggerverse := New(nil, "vbehar/daggerverse") // Placeholder para inicialização

	// Define os handlers disponíveis inline
	handlers := map[string]func(){
		"release": func() {
			ctx := context.Background()
			gitToken := &dagger.Secret{} // Placeholder para o token

			tag, err := daggerverse.Release(ctx, gitToken)
			if err != nil {
				log.Fatalf("Release failed: %v", err)
			}
			log.Printf("Release created with tag: %s", tag)
		},
		"publish": func() {
			ctx := context.Background()

			results, err := daggerverse.Publish(ctx, false)
			if err != nil {
				log.Fatalf("Publish failed: %v", err)
			}
			log.Printf("Publish completed: %v", results)
		},
	}

	// Verifica se o comando é válido
	handlerFunc, exists := handlers[command]
	if !exists {
		log.Fatalf("Invalid command: %s", command)
	}

	// Executa o handler correspondente
	handlerFunc()
}
