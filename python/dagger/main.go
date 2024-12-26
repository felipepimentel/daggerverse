package main

import (
	"context"
	"log"

	"dagger.io/dagger"
)

// main é o ponto de entrada do módulo Dagger
func main() {
	ctx := context.Background()
	
	// Inicializa o cliente Dagger
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Writer()))
	if err != nil {
		log.Fatalf("Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	// Registre as funções disponíveis no módulo
	runner := NewRunner(client)

	if err := runner.RegisterActions(ctx); err != nil {
		log.Fatalf("Failed to register actions: %v", err)
	}
}

// Runner é responsável por agrupar todas as ações disponíveis no módulo
type Runner struct {
	client *dagger.Client
}

// NewRunner cria um novo Runner
func NewRunner(client *dagger.Client) *Runner {
	return &Runner{client: client}
}

// RegisterActions registra as funções disponíveis no módulo
func (r *Runner) RegisterActions(ctx context.Context) error {
	return dagger.Serve(
		ctx,
		dagger.ActionMap{
			"build":        r.Build,
			"test":         r.Test,
			"deploy":       r.Deploy,
			"clean":        r.Clean,
		},
	)
}

// Build é uma função exemplo para compilar código
func (r *Runner) Build(ctx context.Context, req dagger.ActionArgs) error {
	sourcePath := req.Get("source", "./src")
	outputPath := req.Get("output", "./build")

	log.Printf("Building project from %s to %s", sourcePath, outputPath)

	// Exemplo de execução com um container
	container := r.client.Container().From("golang:1.19").WithDirectory("/src", r.client.Host().Directory(sourcePath))

	if _, err := container.WithExec([]string{"go", "build", "-o", outputPath}).Sync(ctx); err != nil {
		return err
	}

	log.Println("Build completed successfully")
	return nil
}

// Test é uma função exemplo para rodar testes
func (r *Runner) Test(ctx context.Context, req dagger.ActionArgs) error {
	sourcePath := req.Get("source", "./src")
	log.Printf("Running tests for project in %s", sourcePath)

	container := r.client.Container().From("golang:1.19").WithDirectory("/src", r.client.Host().Directory(sourcePath))

	if _, err := container.WithExec([]string{"go", "test", "./..."}).Sync(ctx); err != nil {
		return err
	}

	log.Println("Tests ran successfully")
	return nil
}

// Deploy é uma função exemplo para deployment
func (r *Runner) Deploy(ctx context.Context, req dagger.ActionArgs) error {
	log.Println("Deploying project... (placeholder)")
	// Lógica de deploy modular aqui
	return nil
}

// Clean é uma função exemplo para limpeza de artefatos
func (r *Runner) Clean(ctx context.Context, req dagger.ActionArgs) error {
	log.Println("Cleaning build artifacts... (placeholder)")
	// Lógica de limpeza aqui
	return nil
}
