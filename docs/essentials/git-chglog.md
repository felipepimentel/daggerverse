---
layout: default
title: Git-Chglog
parent: Essentials
nav_order: 7
---

# Git-Chglog

O módulo Git-Chglog fornece uma interface para gerar changelogs automaticamente a partir do histórico de commits Git, seguindo convenções e formatos personalizáveis.

## Features

- Geração de changelog
- Formatação personalizada
- Suporte a templates
- Filtro de commits
- Agrupamento por tipo
- Detecção de versões
- Markdown e outros formatos
- Integração com Git

## Instalação

Para usar o módulo Git-Chglog em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/git-chglog"
)
```

## Exemplos de Uso

### Geração Básica

```go
func (m *MyModule) GenerateChangelog(ctx context.Context) (*File, error) {
    chglog := dag.GitChglog().
        WithDirectory(dag.Directory("."))
    
    // Gerar changelog
    return chglog.Generate(ctx)
}
```

### Configuração Personalizada

```go
func (m *MyModule) CustomChangelog(ctx context.Context) (*File, error) {
    chglog := dag.GitChglog().
        WithDirectory(dag.Directory(".")).
        WithConfig(dag.File("config.yml")).
        WithTemplate(dag.File("template.md"))
    
    // Gerar changelog personalizado
    return chglog.Generate(ctx)
}
```

### Changelog por Tag

```go
func (m *MyModule) TagChangelog(ctx context.Context) (*File, error) {
    chglog := dag.GitChglog().
        WithDirectory(dag.Directory(".")).
        WithTag("v1.0.0")
    
    // Gerar changelog para tag
    return chglog.Generate(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: Generate Changelog
on: [push]

jobs:
  changelog:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Generate Changelog
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/git-chglog
          args: |
            do -p '
              chglog := GitChglog().
                WithDirectory(Directory("."))
              chglog.Generate(ctx)
            '
```

## Referência da API

### GitChglog

Estrutura principal que fornece acesso à funcionalidade de changelog.

#### Construtor

- `New() *GitChglog`
  - Cria uma nova instância do GitChglog

#### Métodos de Configuração

- `WithDirectory(dir *Directory) *GitChglog`
  - Define diretório do repositório
  - Parâmetro:
    - `dir`: Diretório Git

- `WithConfig(config *File) *GitChglog`
  - Define arquivo de configuração
  - Parâmetro:
    - `config`: Arquivo de configuração

- `WithTemplate(template *File) *GitChglog`
  - Define template personalizado
  - Parâmetro:
    - `template`: Arquivo de template

- `WithTag(tag string) *GitChglog`
  - Define tag específica
  - Parâmetro:
    - `tag`: Nome da tag

#### Métodos de Operação

- `Generate(ctx context.Context) (*File, error)`
  - Gera arquivo de changelog
  - Retorna arquivo gerado

- `Init(ctx context.Context) error`
  - Inicializa configuração padrão
  - Retorna erro se falhar

- `Preview(ctx context.Context) (string, error)`
  - Visualiza changelog sem salvar
  - Retorna conteúdo do changelog

## Boas Práticas

1. **Configuração**
   - Use arquivo de config
   - Personalize templates
   - Defina convenções

2. **Commits**
   - Siga padrões
   - Use tipos consistentes
   - Escreva mensagens claras

3. **Organização**
   - Agrupe por tipo
   - Ordene por data
   - Destaque breaking changes

4. **Manutenção**
   - Atualize regularmente
   - Revise mudanças
   - Mantenha histórico

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Config**
   ```
   Erro: Invalid configuration
   Solução: Verifique formato do config.yml
   ```

2. **Erro de Template**
   ```
   Erro: Template not found
   Solução: Verifique caminho do template
   ```

3. **Erro de Git**
   ```
   Erro: Git history not found
   Solução: Verifique repositório Git
   ```

## Exemplo de Configuração

```yaml
# .chglog/config.yml
style: github
template: CHANGELOG.tpl.md
info:
  title: CHANGELOG
  repository_url: https://github.com/user/repo

options:
  commits:
    filters:
      Type:
        - feat
        - fix
        - perf
        - refactor
  commit_groups:
    title_maps:
      feat: Features
      fix: Bug Fixes
      perf: Performance Improvements
      refactor: Code Refactoring
  header:
    pattern: "^(\\w*)(?:\\(([\\w\\$\\.\\-\\*\\s]*)\\))?\\:\\s(.*)$"
    pattern_maps:
      - Type
      - Scope
      - Subject
  notes:
    keywords:
      - BREAKING CHANGE
```

## Uso Avançado

### Pipeline Completo

```go
func (m *MyModule) CompletePipeline(ctx context.Context) error {
    // Configurar git-chglog
    chglog := dag.GitChglog().
        WithDirectory(dag.Directory(".")).
        WithConfig(dag.File(".chglog/config.yml")).
        WithTemplate(dag.File(".chglog/CHANGELOG.tpl.md"))
    
    // Inicializar se necessário
    if err := chglog.Init(ctx); err != nil {
        return err
    }
    
    // Gerar preview
    preview, err := chglog.Preview(ctx)
    if err != nil {
        return err
    }
    
    fmt.Println("Preview:", preview)
    
    // Gerar changelog final
    changelog, err := chglog.Generate(ctx)
    if err != nil {
        return err
    }
    
    // Usar changelog gerado
    return dag.Container().
        WithFile("/changelog.md", changelog).
        WithExec([]string{
            "cat", "/changelog.md",
        }).
        Sync(ctx)
}
```

### Configuração Avançada

```go
func (m *MyModule) AdvancedConfig(ctx context.Context) error {
    // Criar configuração personalizada
    config := `
    style: gitlab
    template: |
      {{ range .Versions }}
      ## {{ .Tag.Name }} - {{ datetime "2006-01-02" .Tag.Date }}
      {{ range .CommitGroups }}
      ### {{ .Title }}
      {{ range .Commits }}
      * {{ .Subject }}
      {{- if .Refs }}
        {{- range .Refs }}
        - {{ .Ref }}
        {{- end }}
      {{- end }}
      {{ end }}
      {{ end }}
      {{ end }}
    
    options:
      commits:
        sort_by: Scope
      commit_groups:
        group_by: Type
        sort_by: Title
        title_order:
          - feat
          - fix
          - docs
    `
    
    // Configurar git-chglog
    chglog := dag.GitChglog().
        WithDirectory(dag.Directory(".")).
        WithConfig(dag.File("config.yml", config))
    
    // Gerar changelog
    changelog, err := chglog.Generate(ctx)
    if err != nil {
        return err
    }
    
    // Salvar resultado
    return dag.Container().
        WithFile("/CHANGELOG.md", changelog).
        Sync(ctx)
} 