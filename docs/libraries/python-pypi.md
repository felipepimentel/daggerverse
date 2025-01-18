---
layout: default
title: Python PyPI
parent: Libraries
nav_order: 2
---

# Python PyPI

O módulo Python PyPI fornece uma interface para interagir com o Python Package Index (PyPI), permitindo publicar, gerenciar e baixar pacotes Python.

## Features

- Publicação de pacotes no PyPI
- Download de pacotes
- Autenticação segura
- Suporte a TestPyPI
- Gerenciamento de versões
- Validação de pacotes
- Configuração de repositório
- Integração com pip

## Instalação

Para usar o módulo Python PyPI em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/libraries/python-pypi"
)
```

## Exemplos de Uso

### Publicação de Pacote

```go
func (m *MyModule) PublishPackage(ctx context.Context) error {
    pypi := dag.PythonPyPI().
        WithToken(dag.SetSecret("PYPI_TOKEN", "seu-token")).
        WithPackageDir(dag.Directory("."))
    
    // Publicar pacote
    return pypi.Publish(ctx)
}
```

### Download de Pacote

```go
func (m *MyModule) DownloadPackage(ctx context.Context) error {
    pypi := dag.PythonPyPI().
        WithPackageName("requests").
        WithVersion("2.31.0")
    
    // Baixar pacote
    return pypi.Download(ctx)
}
```

### Uso do TestPyPI

```go
func (m *MyModule) TestPublish(ctx context.Context) error {
    pypi := dag.PythonPyPI().
        WithToken(dag.SetSecret("TEST_PYPI_TOKEN", "seu-token")).
        WithPackageDir(dag.Directory(".")).
        WithTestPyPI(true)
    
    // Publicar no TestPyPI
    return pypi.Publish(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: PyPI Publish
on: [push]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Publish to PyPI
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/libraries/python-pypi
          args: |
            do -p '
              pypi := PythonPyPI().
                WithToken(SetSecret("PYPI_TOKEN", "${{ secrets.PYPI_TOKEN }}")).
                WithPackageDir(Directory("."))
              pypi.Publish(ctx)
            '
```

## Referência da API

### PythonPyPI

Estrutura principal que fornece acesso à funcionalidade do PyPI.

#### Construtor

- `New() *PythonPyPI`
  - Cria uma nova instância do PythonPyPI

#### Métodos de Configuração

- `WithToken(token *Secret) *PythonPyPI`
  - Define o token de autenticação PyPI
  - Parâmetro:
    - `token`: Token PyPI ou TestPyPI

- `WithPackageDir(dir *Directory) *PythonPyPI`
  - Define o diretório do pacote
  - Parâmetro:
    - `dir`: Diretório contendo setup.py ou pyproject.toml

- `WithPackageName(name string) *PythonPyPI`
  - Define o nome do pacote para download
  - Parâmetro:
    - `name`: Nome do pacote no PyPI

- `WithVersion(version string) *PythonPyPI`
  - Define a versão do pacote
  - Parâmetro:
    - `version`: Versão específica do pacote

- `WithTestPyPI(useTest bool) *PythonPyPI`
  - Configura para usar TestPyPI
  - Parâmetro:
    - `useTest`: Se true, usa TestPyPI em vez de PyPI

#### Métodos de Operação

- `Publish(ctx context.Context) error`
  - Publica o pacote no PyPI/TestPyPI
  - Retorna erro se a publicação falhar

- `Download(ctx context.Context) error`
  - Baixa o pacote do PyPI
  - Retorna erro se o download falhar

- `Validate(ctx context.Context) error`
  - Valida o pacote antes da publicação
  - Retorna erro se a validação falhar

## Boas Práticas

1. **Autenticação**
   - Use tokens API seguros
   - Nunca exponha credenciais
   - Rotacione tokens regularmente

2. **Publicação**
   - Teste no TestPyPI primeiro
   - Valide pacotes antes de publicar
   - Mantenha versionamento semântico

3. **Download**
   - Especifique versões exatas
   - Verifique hashes
   - Use mirrors confiáveis

4. **Segurança**
   - Use HTTPS sempre
   - Verifique integridade
   - Monitore dependências

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Autenticação**
   ```
   Erro: Invalid PyPI token
   Solução: Verifique token e permissões
   ```

2. **Erro de Publicação**
   ```
   Erro: Version already exists
   Solução: Atualize versão no setup.py/pyproject.toml
   ```

3. **Erro de Download**
   ```
   Erro: Package not found
   Solução: Verifique nome e versão do pacote
   ```

## Exemplo de Configuração

```toml
# pyproject.toml para publicação
[build-system]
requires = ["setuptools>=61.0"]
build-backend = "setuptools.build_meta"

[project]
name = "meu-pacote"
version = "1.0.0"
authors = [
    { name="Seu Nome", email="seu@email.com" },
]
description = "Descrição do pacote"
readme = "README.md"
requires-python = ">=3.8"
classifiers = [
    "Programming Language :: Python :: 3",
    "License :: OSI Approved :: MIT License",
    "Operating System :: OS Independent",
]

[project.urls]
"Homepage" = "https://github.com/username/projeto"
"Bug Tracker" = "https://github.com/username/projeto/issues"
```

## Uso Avançado

### Pipeline Completo

```go
func (m *MyModule) CompletePipeline(ctx context.Context) error {
    pypi := dag.PythonPyPI().
        WithToken(dag.SetSecret("PYPI_TOKEN", "token")).
        WithPackageDir(dag.Directory("."))
    
    // Validar pacote
    if err := pypi.Validate(ctx); err != nil {
        return err
    }
    
    // Testar no TestPyPI primeiro
    if err := pypi.
        WithTestPyPI(true).
        WithToken(dag.SetSecret("TEST_PYPI_TOKEN", "test-token")).
        Publish(ctx); err != nil {
        return err
    }
    
    // Se sucesso, publicar no PyPI
    return pypi.
        WithTestPyPI(false).
        WithToken(dag.SetSecret("PYPI_TOKEN", "prod-token")).
        Publish(ctx)
}
```

### Download com Verificação

```go
func (m *MyModule) VerifiedDownload(ctx context.Context) error {
    pypi := dag.PythonPyPI().
        WithPackageName("requests").
        WithVersion("2.31.0")
    
    // Baixar e verificar pacote
    if err := pypi.Download(ctx); err != nil {
        return err
    }
    
    // Verificar hash (exemplo)
    return pypi.Run(ctx, []string{
        "pip", "hash", "--algorithm", "sha256",
        "requests-2.31.0.tar.gz",
    })
} 