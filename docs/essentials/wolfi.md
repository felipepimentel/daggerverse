---
layout: default
title: Wolfi
parent: Essentials
nav_order: 1
---

# Wolfi

O módulo Wolfi fornece uma interface para trabalhar com o sistema operacional Wolfi, permitindo a criação e gerenciamento de containers baseados em Wolfi.

## Features

- Criação de containers Wolfi
- Gerenciamento de pacotes
- Configuração de ambiente
- Execução de comandos
- Customização de imagens
- Integração com apk
- Suporte a multi-arquitetura
- Otimização de camadas

## Instalação

Para usar o módulo Wolfi em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/wolfi"
)
```

## Exemplos de Uso

### Container Básico

```go
func (m *MyModule) BasicContainer(ctx context.Context) *Container {
    wolfi := dag.Wolfi().
        WithPackages([]string{"python-3.12", "git"})
    
    // Criar container
    return wolfi.Container()
}
```

### Instalação de Pacotes

```go
func (m *MyModule) WithPackages(ctx context.Context) *Container {
    wolfi := dag.Wolfi().
        WithPackages([]string{
            "python-3.12",
            "nodejs",
            "git",
            "curl",
            "build-base",
        })
    
    // Criar container com pacotes
    return wolfi.Container()
}
```

### Execução de Comandos

```go
func (m *MyModule) RunCommands(ctx context.Context) error {
    container := dag.Wolfi().
        WithPackages([]string{"python-3.12"}).
        Container()
    
    // Executar comando
    return container.
        WithExec([]string{
            "python",
            "-c",
            "print('Hello from Wolfi!')",
        }).
        Sync(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: Wolfi Build
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build with Wolfi
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/wolfi
          args: |
            do -p '
              wolfi := Wolfi().
                WithPackages(["python-3.12", "git"]).
                Container()
              wolfi.WithExec(["python", "--version"]).
                Sync(ctx)
            '
```

## Referência da API

### Wolfi

Estrutura principal que fornece acesso à funcionalidade do Wolfi.

#### Construtor

- `New() *Wolfi`
  - Cria uma nova instância do Wolfi

#### Métodos de Configuração

- `WithPackages(packages []string) *Wolfi`
  - Define pacotes a serem instalados
  - Parâmetro:
    - `packages`: Lista de pacotes

- `WithEnv(key string, value string) *Wolfi`
  - Define variável de ambiente
  - Parâmetros:
    - `key`: Nome da variável
    - `value`: Valor da variável

- `WithWorkdir(path string) *Wolfi`
  - Define diretório de trabalho
  - Parâmetro:
    - `path`: Caminho do diretório

- `WithUser(user string) *Wolfi`
  - Define usuário do container
  - Parâmetro:
    - `user`: Nome do usuário

#### Métodos de Operação

- `Container() *Container`
  - Cria um novo container Wolfi
  - Retorna o container configurado

- `WithExec(args []string) *Container`
  - Executa comando no container
  - Parâmetro:
    - `args`: Argumentos do comando

- `WithFile(path string, contents string) *Container`
  - Adiciona arquivo ao container
  - Parâmetros:
    - `path`: Caminho do arquivo
    - `contents`: Conteúdo do arquivo

## Boas Práticas

1. **Gerenciamento de Pacotes**
   - Instale apenas o necessário
   - Mantenha pacotes atualizados
   - Use versões específicas

2. **Otimização**
   - Minimize camadas
   - Limpe caches
   - Remova arquivos temporários

3. **Segurança**
   - Use usuário não-root
   - Atualize regularmente
   - Verifique vulnerabilidades

4. **Performance**
   - Otimize ordem de comandos
   - Combine operações relacionadas
   - Use multi-stage builds

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Pacote**
   ```
   Erro: Package not found
   Solução: Verifique nome do pacote
   ```

2. **Erro de Permissão**
   ```
   Erro: Permission denied
   Solução: Verifique permissões/usuário
   ```

3. **Erro de Execução**
   ```
   Erro: Command not found
   Solução: Verifique instalação do pacote
   ```

## Exemplo de Configuração

```yaml
# wolfi.yaml
packages:
  - python-3.12
  - nodejs
  - git
  - curl
  - build-base

environment:
  PYTHONUNBUFFERED: "1"
  NODE_ENV: "production"

workdir: /app

user: nonroot
```

## Uso Avançado

### Multi-stage Build

```go
func (m *MyModule) MultiStageBuild(ctx context.Context) *Container {
    // Build stage
    builder := dag.Wolfi().
        WithPackages([]string{
            "python-3.12",
            "build-base",
            "poetry",
        }).
        Container()
    
    // Copiar e construir aplicação
    builder = builder.
        WithDirectory("/app", dag.Directory(".")).
        WithWorkdir("/app").
        WithExec([]string{"poetry", "install"}).
        WithExec([]string{"poetry", "build"})
    
    // Runtime stage
    runtime := dag.Wolfi().
        WithPackages([]string{"python-3.12"}).
        Container()
    
    // Copiar artefatos do builder
    return runtime.
        WithDirectory(
            "/app",
            builder.Directory("/app/dist"),
        )
}
```

### Container Customizado

```go
func (m *MyModule) CustomContainer(ctx context.Context) *Container {
    wolfi := dag.Wolfi().
        WithPackages([]string{
            "python-3.12",
            "git",
        }).
        WithEnv("PYTHONUNBUFFERED", "1").
        WithWorkdir("/app").
        WithUser("nonroot")
    
    // Configurar container
    container := wolfi.Container()
    
    // Adicionar arquivo de configuração
    container = container.WithFile(
        "/app/config.py",
        `
        DEBUG = False
        HOST = "0.0.0.0"
        PORT = 8000
        `,
    )
    
    // Executar setup
    return container.WithExec([]string{
        "python",
        "-c",
        "import config; print(f'Server running on {config.HOST}:{config.PORT}')",
    })
} 