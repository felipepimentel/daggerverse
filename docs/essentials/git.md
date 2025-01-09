---
layout: default
title: Git
parent: Essentials
nav_order: 6
---

# Git

O módulo Git fornece uma interface para operações Git em pipelines Dagger, permitindo clonagem, commit, push e outras operações de controle de versão.

## Features

- Clonagem de repositórios
- Gerenciamento de branches
- Operações de commit
- Push e pull
- Configuração de remotes
- Gerenciamento de tags
- Autenticação segura
- Integração com SSH

## Instalação

Para usar o módulo Git em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/git"
)
```

## Exemplos de Uso

### Clonagem Básica

```go
func (m *MyModule) CloneRepo(ctx context.Context) (*Directory, error) {
    git := dag.Git().
        WithURL("https://github.com/user/repo.git")
    
    // Clonar repositório
    return git.Clone(ctx)
}
```

### Commit e Push

```go
func (m *MyModule) CommitAndPush(ctx context.Context) error {
    git := dag.Git().
        WithDirectory(dag.Directory(".")).
        WithCredentials(
            dag.SetSecret("GIT_USERNAME", "user"),
            dag.SetSecret("GIT_PASSWORD", "token"),
        )
    
    // Commit e push
    return git.
        Add(ctx, ".").
        Commit(ctx, "feat: novo recurso").
        Push(ctx)
}
```

### Gerenciamento de Branches

```go
func (m *MyModule) BranchOps(ctx context.Context) error {
    git := dag.Git().
        WithDirectory(dag.Directory("."))
    
    // Criar e mudar de branch
    return git.
        Checkout(ctx, "-b", "feature/nova").
        Pull(ctx, "origin", "main")
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: Git Operations
on: [push]

jobs:
  clone:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Clone Repository
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/git
          args: |
            do -p '
              git := Git().
                WithURL("https://github.com/user/repo.git")
              git.Clone(ctx)
            '
```

## Referência da API

### Git

Estrutura principal que fornece acesso à funcionalidade Git.

#### Construtor

- `New() *Git`
  - Cria uma nova instância do Git

#### Métodos de Configuração

- `WithURL(url string) *Git`
  - Define URL do repositório
  - Parâmetro:
    - `url`: URL do repositório

- `WithDirectory(dir *Directory) *Git`
  - Define diretório de trabalho
  - Parâmetro:
    - `dir`: Diretório Git

- `WithCredentials(username *Secret, password *Secret) *Git`
  - Define credenciais
  - Parâmetros:
    - `username`: Nome de usuário
    - `password`: Senha ou token

- `WithBranch(branch string) *Git`
  - Define branch
  - Parâmetro:
    - `branch`: Nome da branch

#### Métodos de Operação

- `Clone(ctx context.Context) (*Directory, error)`
  - Clona repositório
  - Retorna diretório clonado

- `Add(ctx context.Context, path string) *Git`
  - Adiciona arquivos ao stage
  - Parâmetro:
    - `path`: Caminho dos arquivos

- `Commit(ctx context.Context, message string) *Git`
  - Cria commit
  - Parâmetro:
    - `message`: Mensagem do commit

- `Push(ctx context.Context) error`
  - Envia commits para remote
  - Retorna erro se falhar

- `Pull(ctx context.Context, remote string, branch string) error`
  - Atualiza branch local
  - Parâmetros:
    - `remote`: Nome do remote
    - `branch`: Nome da branch

## Boas Práticas

1. **Credenciais**
   - Use tokens de acesso
   - Proteja credenciais
   - Rotacione tokens

2. **Commits**
   - Use mensagens descritivas
   - Siga convenções
   - Faça commits atômicos

3. **Branches**
   - Mantenha branches atualizadas
   - Use nomes descritivos
   - Limpe branches obsoletas

4. **Segurança**
   - Verifique permissões
   - Valide certificados
   - Monitore acessos

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Autenticação**
   ```
   Erro: Authentication failed
   Solução: Verifique credenciais
   ```

2. **Erro de Clone**
   ```
   Erro: Repository not found
   Solução: Verifique URL e permissões
   ```

3. **Erro de Push**
   ```
   Erro: Remote rejected
   Solução: Atualize branch local
   ```

## Exemplo de Configuração

```yaml
# git.yaml
repository:
  url: https://github.com/user/repo.git
  branch: main
  depth: 1

credentials:
  username_env: GIT_USERNAME
  password_env: GIT_PASSWORD

ssh:
  enabled: true
  key_path: ~/.ssh/id_rsa

commit:
  name: CI Bot
  email: ci@example.com
  gpg_sign: true
```

## Uso Avançado

### Pipeline Completo

```go
func (m *MyModule) CompletePipeline(ctx context.Context) error {
    // Configurar Git
    git := dag.Git().
        WithURL("https://github.com/user/repo.git").
        WithCredentials(
            dag.SetSecret("GIT_USERNAME", "user"),
            dag.SetSecret("GIT_PASSWORD", "token"),
        ).
        WithBranch("feature/nova")
    
    // Clonar repositório
    dir, err := git.Clone(ctx)
    if err != nil {
        return err
    }
    
    // Configurar diretório
    git = git.WithDirectory(dir)
    
    // Realizar alterações
    if err := git.
        Add(ctx, ".").
        Commit(ctx, "feat: implementação completa").
        Push(ctx); err != nil {
        return err
    }
    
    // Criar tag
    return git.
        Tag(ctx, "v1.0.0", "Release v1.0.0").
        PushTags(ctx)
}
```

### Operações Avançadas

```go
func (m *MyModule) AdvancedOps(ctx context.Context) error {
    // Configurar Git com SSH
    git := dag.Git().
        WithURL("git@github.com:user/repo.git").
        WithSSHKey(dag.SetSecret("SSH_KEY", "private-key"))
    
    // Clonar e configurar
    dir, err := git.Clone(ctx)
    if err != nil {
        return err
    }
    
    git = git.WithDirectory(dir)
    
    // Configurar remotes
    if err := git.
        AddRemote(ctx, "upstream", "git@github.com:upstream/repo.git").
        FetchAll(ctx); err != nil {
        return err
    }
    
    // Rebase com upstream
    return git.
        Checkout(ctx, "main").
        Pull(ctx, "upstream", "main").
        Push(ctx, "--force-with-lease")
} 