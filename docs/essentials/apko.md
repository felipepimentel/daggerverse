---
layout: default
title: APKO
parent: Essentials
nav_order: 2
---

# APKO

O módulo APKO fornece uma interface para trabalhar com a ferramenta APKO (Alpine Package Keeper Orchestrator), permitindo a criação de imagens de container otimizadas e seguras baseadas em Alpine Linux.

## Features

- Criação de imagens Alpine
- Gerenciamento de pacotes
- Configuração de imagem
- Suporte a multi-arquitetura
- Otimização de tamanho
- Segurança aprimorada
- Reprodutibilidade
- Integração com registries

## Instalação

Para usar o módulo APKO em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/apko"
)
```

## Exemplos de Uso

### Imagem Básica

```go
func (m *MyModule) BasicImage(ctx context.Context) (*Container, error) {
    apko := dag.APKO().
        WithPackages([]string{"python3", "git"})
    
    // Criar imagem
    return apko.Build(ctx)
}
```

### Configuração Personalizada

```go
func (m *MyModule) CustomConfig(ctx context.Context) (*Container, error) {
    apko := dag.APKO().
        WithPackages([]string{"python3", "nodejs"}).
        WithEnv(map[string]string{
            "PYTHON_PATH": "/usr/lib/python3.9",
            "NODE_ENV": "production",
        }).
        WithUser("nonroot")
    
    // Criar imagem com configuração personalizada
    return apko.Build(ctx)
}
```

### Multi-arquitetura

```go
func (m *MyModule) MultiArch(ctx context.Context) error {
    apko := dag.APKO().
        WithPackages([]string{"python3"}).
        WithPlatforms([]string{
            "linux/amd64",
            "linux/arm64",
        })
    
    // Construir para múltiplas arquiteturas
    return apko.BuildAll(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: APKO Build
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build with APKO
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/apko
          args: |
            do -p '
              apko := APKO().
                WithPackages(["python3", "git"]).
                WithUser("nonroot")
              apko.Build(ctx)
            '
```

## Referência da API

### APKO

Estrutura principal que fornece acesso à funcionalidade do APKO.

#### Construtor

- `New() *APKO`
  - Cria uma nova instância do APKO

#### Métodos de Configuração

- `WithPackages(packages []string) *APKO`
  - Define pacotes a serem instalados
  - Parâmetro:
    - `packages`: Lista de pacotes

- `WithEnv(env map[string]string) *APKO`
  - Define variáveis de ambiente
  - Parâmetro:
    - `env`: Mapa de variáveis

- `WithUser(user string) *APKO`
  - Define usuário da imagem
  - Parâmetro:
    - `user`: Nome do usuário

- `WithPlatforms(platforms []string) *APKO`
  - Define plataformas alvo
  - Parâmetro:
    - `platforms`: Lista de plataformas

- `WithConfig(config *File) *APKO`
  - Define arquivo de configuração
  - Parâmetro:
    - `config`: Arquivo de configuração APKO

#### Métodos de Operação

- `Build(ctx context.Context) (*Container, error)`
  - Constrói imagem para plataforma padrão
  - Retorna container e erro

- `BuildAll(ctx context.Context) error`
  - Constrói imagens para todas plataformas
  - Retorna erro se falhar

- `Push(ctx context.Context, ref string) error`
  - Publica imagem em registry
  - Parâmetros:
    - `ref`: Referência da imagem
  - Retorna erro se falhar

## Boas Práticas

1. **Otimização**
   - Minimize número de pacotes
   - Use camadas eficientemente
   - Limpe caches

2. **Segurança**
   - Use usuário não-root
   - Mantenha pacotes atualizados
   - Siga princípio do menor privilégio

3. **Multi-arquitetura**
   - Teste em todas plataformas
   - Valide compatibilidade
   - Otimize para cada arquitetura

4. **Reprodutibilidade**
   - Use versões específicas
   - Documente configurações
   - Mantenha builds determinísticos

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Build**
   ```
   Erro: Failed to build image
   Solução: Verifique configuração e dependências
   ```

2. **Erro de Pacote**
   ```
   Erro: Package not found
   Solução: Verifique repositórios e nome do pacote
   ```

3. **Erro de Plataforma**
   ```
   Erro: Unsupported platform
   Solução: Verifique suporte à arquitetura
   ```

## Exemplo de Configuração

```yaml
# apko.yaml
contents:
  repositories:
    - https://dl-cdn.alpinelinux.org/alpine/edge/main
    - https://dl-cdn.alpinelinux.org/alpine/edge/community
  
  packages:
    - python3
    - nodejs
    - git
    - curl

environment:
  PYTHON_PATH: /usr/lib/python3.9
  NODE_ENV: production

accounts:
  groups:
    - nonroot
  users:
    - name: nonroot
      uid: 65532

entrypoint:
  command: /bin/sh

cmd: -c

archs:
  - x86_64
  - aarch64
```

## Uso Avançado

### Pipeline Completo

```go
func (m *MyModule) CompletePipeline(ctx context.Context) error {
    // Configurar APKO
    apko := dag.APKO().
        WithPackages([]string{
            "python3",
            "nodejs",
            "git",
        }).
        WithEnv(map[string]string{
            "PYTHON_PATH": "/usr/lib/python3.9",
            "NODE_ENV": "production",
        }).
        WithUser("nonroot").
        WithPlatforms([]string{
            "linux/amd64",
            "linux/arm64",
        })
    
    // Construir imagens
    if err := apko.BuildAll(ctx); err != nil {
        return err
    }
    
    // Publicar imagens
    return apko.Push(ctx, "registry.example.com/app:latest")
}
```

### Configuração Avançada

```go
func (m *MyModule) AdvancedConfig(ctx context.Context) (*Container, error) {
    // Criar arquivo de configuração
    config := `
    contents:
      repositories:
        - https://dl-cdn.alpinelinux.org/alpine/edge/main
        - https://dl-cdn.alpinelinux.org/alpine/edge/community
      
      packages:
        - python3
        - python3-dev
        - build-base
        - git
    
    environment:
      PYTHON_PATH: /usr/lib/python3.9
      PYTHONUNBUFFERED: "1"
    
    accounts:
      groups:
        - nonroot
      users:
        - name: nonroot
          uid: 65532
    
    entrypoint:
      command: /usr/bin/python3
    `
    
    // Configurar APKO com arquivo
    apko := dag.APKO().
        WithConfig(dag.File("apko.yaml", config))
    
    // Construir imagem
    return apko.Build(ctx)
} 