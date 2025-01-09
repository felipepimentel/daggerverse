---
layout: default
title: Versioner
parent: Libraries
nav_order: 3
---

# Versioner

O módulo Versioner fornece uma interface para gerenciar versionamento semântico em projetos, automatizando o processo de incremento de versão e geração de changelogs.

## Features

- Versionamento semântico
- Geração de changelog
- Detecção automática de versão
- Incremento de versão
- Suporte a diferentes formatos
- Integração com Git
- Validação de versão
- Customização de regras

## Instalação

Para usar o módulo Versioner em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/versioner"
)
```

## Exemplos de Uso

### Incremento de Versão

```go
func (m *MyModule) BumpVersion(ctx context.Context) (string, error) {
    versioner := dag.Versioner().
        WithSource(dag.Directory(".")).
        WithType("minor")  // major, minor, ou patch
    
    // Incrementar versão
    return versioner.Bump(ctx)
}
```

### Geração de Changelog

```go
func (m *MyModule) GenerateChangelog(ctx context.Context) error {
    versioner := dag.Versioner().
        WithSource(dag.Directory(".")).
        WithFormat("markdown")
    
    // Gerar changelog
    return versioner.GenerateChangelog(ctx)
}
```

### Detecção de Versão

```go
func (m *MyModule) GetCurrentVersion(ctx context.Context) (string, error) {
    versioner := dag.Versioner().
        WithSource(dag.Directory("."))
    
    // Obter versão atual
    return versioner.GetVersion(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: Version Bump
on:
  push:
    branches: [main]

jobs:
  version:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Bump Version
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/versioner
          args: |
            do -p '
              versioner := Versioner().
                WithSource(Directory(".")).
                WithType("minor")
              versioner.Bump(ctx)
            '
```

## Referência da API

### Versioner

Estrutura principal que fornece acesso à funcionalidade de versionamento.

#### Construtor

- `New() *Versioner`
  - Cria uma nova instância do Versioner

#### Métodos de Configuração

- `WithSource(source *Directory) *Versioner`
  - Define o diretório fonte do projeto
  - Parâmetro:
    - `source`: Diretório contendo arquivos de versão

- `WithType(type string) *Versioner`
  - Define o tipo de incremento de versão
  - Parâmetro:
    - `type`: "major", "minor", ou "patch"

- `WithFormat(format string) *Versioner`
  - Define o formato do changelog
  - Parâmetro:
    - `format`: "markdown", "text", etc.

- `WithConfig(config *File) *Versioner`
  - Define arquivo de configuração personalizado
  - Parâmetro:
    - `config`: Arquivo de configuração

#### Métodos de Operação

- `Bump(ctx context.Context) (string, error)`
  - Incrementa a versão do projeto
  - Retorna a nova versão

- `GenerateChangelog(ctx context.Context) error`
  - Gera changelog baseado em commits
  - Retorna erro se falhar

- `GetVersion(ctx context.Context) (string, error)`
  - Obtém a versão atual do projeto
  - Retorna a versão atual

- `Validate(ctx context.Context) error`
  - Valida a versão atual
  - Retorna erro se inválida

## Boas Práticas

1. **Versionamento**
   - Use versionamento semântico
   - Documente mudanças
   - Mantenha histórico

2. **Commits**
   - Use mensagens descritivas
   - Siga convenções
   - Agrupe mudanças relacionadas

3. **Changelog**
   - Mantenha atualizado
   - Seja descritivo
   - Categorize mudanças

4. **Automação**
   - Automatize incrementos
   - Valide antes de publicar
   - Mantenha consistência

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Versão**
   ```
   Erro: Invalid version format
   Solução: Verifique formato semântico
   ```

2. **Erro de Changelog**
   ```
   Erro: Failed to generate changelog
   Solução: Verifique histórico de commits
   ```

3. **Erro de Configuração**
   ```
   Erro: Config file not found
   Solução: Verifique arquivo de configuração
   ```

## Exemplo de Configuração

```yaml
# .versionrc.yml
version:
  files:
    - package.json
    - pyproject.toml
  
changelog:
  format: markdown
  sections:
    - features
    - fixes
    - breaking
  
commit:
  types:
    - feat
    - fix
    - docs
    - style
    - refactor
    - test
    - chore
  
rules:
  major:
    - type: feat
      scope: "!"
  minor:
    - type: feat
  patch:
    - type: fix
```

## Uso Avançado

### Pipeline Completo

```go
func (m *MyModule) ReleasePipeline(ctx context.Context) error {
    versioner := dag.Versioner().
        WithSource(dag.Directory(".")).
        WithType("minor")
    
    // Validar versão atual
    if err := versioner.Validate(ctx); err != nil {
        return err
    }
    
    // Gerar changelog
    if err := versioner.GenerateChangelog(ctx); err != nil {
        return err
    }
    
    // Incrementar versão
    newVersion, err := versioner.Bump(ctx)
    if err != nil {
        return err
    }
    
    // Criar tag Git
    return versioner.Run(ctx, []string{
        "git", "tag", "-a", newVersion, "-m", "Release " + newVersion,
    })
}
```

### Configuração Personalizada

```go
func (m *MyModule) CustomVersioning(ctx context.Context) error {
    versioner := dag.Versioner().
        WithSource(dag.Directory(".")).
        WithConfig(dag.File("custom-version.yml"))
    
    // Configurar regras personalizadas
    if err := versioner.Configure(ctx, map[string]interface{}{
        "rules": map[string]interface{}{
            "major": []map[string]string{
                {"type": "breaking"},
            },
            "minor": []map[string]string{
                {"type": "feature"},
            },
            "patch": []map[string]string{
                {"type": "bugfix"},
            },
        },
    }); err != nil {
        return err
    }
    
    // Aplicar versionamento
    return versioner.Apply(ctx)
} 