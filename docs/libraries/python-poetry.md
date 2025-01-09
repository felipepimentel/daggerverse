---
layout: default
title: Python Poetry
parent: Libraries
nav_order: 1
---

# Python Poetry

O módulo Python Poetry fornece uma interface para gerenciar projetos Python usando Poetry, incluindo instalação de dependências, construção de pacotes e publicação no PyPI.

## Features

- Gerenciamento de dependências com Poetry
- Construção de pacotes Python
- Publicação no PyPI
- Ambiente virtual isolado
- Integração com pip
- Suporte a diferentes versões Python
- Configuração de projeto
- Gestão de dependências de desenvolvimento

## Instalação

Para usar o módulo Python Poetry em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/python-poetry"
)
```

## Exemplos de Uso

### Instalação de Dependências

```go
func (m *MyModule) InstallDeps(ctx context.Context) error {
    poetry := dag.PythonPoetry().
        WithPythonVersion("3.12").
        WithSource(dag.Directory("."))
    
    // Instalar dependências
    return poetry.Install(ctx)
}
```

### Construção de Pacote

```go
func (m *MyModule) BuildPackage(ctx context.Context) error {
    poetry := dag.PythonPoetry().
        WithPythonVersion("3.12").
        WithSource(dag.Directory("."))
    
    // Construir pacote
    return poetry.Build(ctx)
}
```

### Publicação no PyPI

```go
func (m *MyModule) PublishPackage(ctx context.Context) error {
    poetry := dag.PythonPoetry().
        WithPythonVersion("3.12").
        WithSource(dag.Directory(".")).
        WithPyPIToken(dag.SetSecret("PYPI_TOKEN", "seu-token"))
    
    // Publicar no PyPI
    return poetry.Publish(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: Poetry Build & Publish
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build and Publish
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/python-poetry
          args: |
            do -p '
              poetry := PythonPoetry().
                WithPythonVersion("3.12").
                WithSource(Directory(".")).
                WithPyPIToken(SetSecret("PYPI_TOKEN", "${{ secrets.PYPI_TOKEN }}"))
              poetry.Publish(ctx)
            '
```

## Referência da API

### PythonPoetry

Estrutura principal que fornece acesso à funcionalidade do Poetry.

#### Construtor

- `New() *PythonPoetry`
  - Cria uma nova instância do PythonPoetry

#### Métodos de Configuração

- `WithPythonVersion(version string) *PythonPoetry`
  - Define a versão do Python
  - Parâmetro:
    - `version`: Versão do Python (ex: "3.12")

- `WithSource(source *Directory) *PythonPoetry`
  - Define o diretório fonte do projeto
  - Parâmetro:
    - `source`: Diretório contendo pyproject.toml

- `WithPyPIToken(token *Secret) *PythonPoetry`
  - Define o token de autenticação PyPI
  - Parâmetro:
    - `token`: Token PyPI

#### Métodos de Operação

- `Install(ctx context.Context) error`
  - Instala dependências do projeto
  - Retorna erro se a instalação falhar

- `Build(ctx context.Context) error`
  - Constrói o pacote Python
  - Retorna erro se a construção falhar

- `Publish(ctx context.Context) error`
  - Publica o pacote no PyPI
  - Retorna erro se a publicação falhar

- `Run(ctx context.Context, args []string) error`
  - Executa comando Poetry personalizado
  - Parâmetros:
    - `args`: Lista de argumentos para o comando

## Boas Práticas

1. **Gerenciamento de Dependências**
   - Mantenha pyproject.toml atualizado
   - Use grupos de dependências
   - Especifique versões precisas

2. **Construção de Pacotes**
   - Inclua todos os arquivos necessários
   - Configure metadados corretamente
   - Teste o pacote antes de publicar

3. **Publicação**
   - Use tokens PyPI seguros
   - Verifique credenciais
   - Teste em TestPyPI primeiro

4. **Ambiente Virtual**
   - Use ambientes isolados
   - Ative env quando necessário
   - Mantenha env limpo

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Instalação**
   ```
   Erro: Failed to install dependencies
   Solução: Verifique pyproject.toml e versões
   ```

2. **Erro de Build**
   ```
   Erro: Failed to build package
   Solução: Verifique estrutura do projeto
   ```

3. **Erro de Publicação**
   ```
   Erro: Failed to publish to PyPI
   Solução: Verifique token e configurações
   ```

## Exemplo de Configuração

```toml
# pyproject.toml
[tool.poetry]
name = "meu-projeto"
version = "0.1.0"
description = "Descrição do projeto"
authors = ["Seu Nome <seu@email.com>"]

[tool.poetry.dependencies]
python = "^3.12"
requests = "^2.31.0"

[tool.poetry.group.dev.dependencies]
pytest = "^7.4.3"
black = "^23.11.0"

[build-system]
requires = ["poetry-core>=1.0.0"]
build-backend = "poetry.core.masonry.api"
```

## Uso Avançado

### Pipeline Personalizado

```go
func (m *MyModule) CustomPoetryPipeline(ctx context.Context) error {
    poetry := dag.PythonPoetry().
        WithPythonVersion("3.12").
        WithSource(dag.Directory("."))
    
    // Instalar dependências
    if err := poetry.Install(ctx); err != nil {
        return err
    }
    
    // Executar testes
    if err := poetry.Run(ctx, []string{"run", "pytest"}); err != nil {
        return err
    }
    
    // Construir pacote
    if err := poetry.Build(ctx); err != nil {
        return err
    }
    
    // Publicar se na branch main
    if onMain {
        return poetry.
            WithPyPIToken(dag.SetSecret("PYPI_TOKEN", "token")).
            Publish(ctx)
    }
    
    return nil
}
```

### Ambiente de Desenvolvimento

```go
func (m *MyModule) DevEnvironment(ctx context.Context) error {
    poetry := dag.PythonPoetry().
        WithPythonVersion("3.12").
        WithSource(dag.Directory("."))
    
    // Instalar dependências de desenvolvimento
    if err := poetry.Run(ctx, []string{"install", "--with", "dev"}); err != nil {
        return err
    }
    
    // Configurar ambiente
    return poetry.Run(ctx, []string{
        "run", "python", "-m", "ipython",
    })
} 