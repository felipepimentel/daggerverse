---
layout: default
title: Python Pipeline
parent: Pipelines
nav_order: 1
---

# Python Pipeline

O módulo Python Pipeline fornece um pipeline completo para projetos Python usando Poetry e PyPI. Ele automatiza o processo de teste, construção e publicação de pacotes Python.

## Features

- Integração com Poetry
- Publicação no PyPI
- Testes automatizados
- Verificação de código com Ruff
- Gerenciamento de versão
- Construção de containers
- Autenticação Docker Hub
- Ambiente de desenvolvimento
- Configuração Git
- Gestão de dependências

## Instalação

Para usar o módulo Python Pipeline em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/pipelines/python"
)
```

## Exemplos de Uso

### Pipeline Básico

```go
func (m *MyModule) Example(ctx context.Context) (string, error) {
    python := dag.Python().New(
        "",           // versão Python (default: 3.12-alpine)
        "",           // email Git (default: github-actions[bot]@users.noreply.github.com)
        "",           // nome Git (default: github-actions[bot])
        "",           // usuário Docker Hub (opcional)
        nil,          // senha Docker Hub (opcional)
    )
    
    // Publicar pacote
    return python.Publish(
        ctx,
        dag.Directory("."),  // diretório fonte
        dag.SetSecret("PYPI_TOKEN", "seu-token"),  // token PyPI
    )
}
```

### Construção de Container

```go
func (m *MyModule) BuildContainer(ctx context.Context) *Container {
    python := dag.Python().New(
        "3.11-alpine",  // versão Python específica
        "",
        "",
        "username",     // usuário Docker Hub
        dag.SetSecret("DOCKER_PASSWORD", "senha"),
    )
    
    // Construir container
    return python.Build(
        ctx,
        dag.Directory("."),
    )
}
```

### Execução de Testes

```go
func (m *MyModule) RunTests(ctx context.Context) (string, error) {
    python := dag.Python().New("", "", "", "", nil)
    
    // Executar testes
    return python.Test(
        ctx,
        dag.Directory("."),
    )
}
```

## Integração com GitHub Actions

Você pode usar este módulo em seus workflows do GitHub Actions:

```yaml
name: Python Pipeline
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build and Publish
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/pipelines/python
          args: |
            do -p '
              python := Python().New("", "", "", "", nil)
              python.Publish(
                ctx,
                Directory("."),
                SetSecret("PYPI_TOKEN", "${{ secrets.PYPI_TOKEN }}"),
              )
            '
```

## Referência da API

### Python

Estrutura principal que fornece acesso à funcionalidade do pipeline Python.

#### Construtor

- `New(pythonVersion string, gitEmail string, gitName string, dockerUsername string, dockerPassword *Secret) *Python`
  - Cria uma nova instância do pipeline Python
  - Parâmetros:
    - `pythonVersion`: Versão do Python (opcional, padrão: "3.12-alpine")
    - `gitEmail`: Email para commits Git (opcional)
    - `gitName`: Nome para commits Git (opcional)
    - `dockerUsername`: Usuário Docker Hub (opcional)
    - `dockerPassword`: Senha Docker Hub (opcional)

#### Métodos

- `Publish(ctx context.Context, source *Directory, token *Secret) (string, error)`
  - Publica o pacote no PyPI
  - Parâmetros:
    - `source`: Diretório fonte
    - `token`: Token de autenticação PyPI
  - Retorna a versão publicada

- `Build(ctx context.Context, source *Directory) *Container`
  - Constrói um container com o projeto
  - Parâmetros:
    - `source`: Diretório fonte
  - Retorna o container configurado

- `Test(ctx context.Context, source *Directory) (string, error)`
  - Executa testes e verificações de qualidade
  - Parâmetros:
    - `source`: Diretório fonte
  - Retorna a saída dos testes

- `Lint(ctx context.Context, source *Directory) error`
  - Executa verificações de código com Ruff
  - Parâmetros:
    - `source`: Diretório fonte

- `BuildEnv(ctx context.Context, source *Directory) *Container`
  - Cria ambiente de desenvolvimento
  - Parâmetros:
    - `source`: Diretório fonte
  - Retorna o container configurado

## Boas Práticas

1. **Configuração do Projeto**
   - Use Poetry para gerenciamento de dependências
   - Mantenha pyproject.toml atualizado
   - Documente dependências

2. **Testes**
   - Escreva testes unitários
   - Configure cobertura de código
   - Automatize testes

3. **Qualidade de Código**
   - Use Ruff para linting
   - Mantenha padrões consistentes
   - Documente código

4. **Versionamento**
   - Use versionamento semântico
   - Atualize CHANGELOG
   - Faça tags de versão

## Solução de Problemas

Problemas comuns e soluções:

1. **Erros de Build**
   ```
   Erro: failed to build container
   Solução: Verifique dependências
   ```

2. **Falhas nos Testes**
   ```
   Erro: poetry test failed
   Solução: Verifique logs de teste
   ```

3. **Erros de Publicação**
   ```
   Erro: failed to publish to PyPI
   Solução: Verifique token PyPI
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

[tool.poetry.dev-dependencies]
pytest = "^7.4.3"
ruff = "^0.1.8"

[build-system]
requires = ["poetry-core>=1.0.0"]
build-backend = "poetry.core.masonry.api"
```

## Uso Avançado

### Pipeline Personalizado

```go
func (m *MyModule) CustomPipeline(ctx context.Context) error {
    python := dag.Python().New(
        "3.12-alpine",
        "dev@example.com",
        "Developer",
        "",
        nil,
    )
    
    // Executar testes primeiro
    testOutput, err := python.Test(ctx, dag.Directory("."))
    if err != nil {
        return err
    }
    fmt.Println("Testes:", testOutput)
    
    // Construir container
    container := python.Build(ctx, dag.Directory("."))
    
    // Executar comando personalizado
    return container.
        WithExec([]string{
            "python",
            "-c",
            "print('Pipeline personalizado executado com sucesso!')",
        }).
        Sync(ctx)
}
```

### Ambiente de Desenvolvimento

```go
func (m *MyModule) DevEnvironment(ctx context.Context) error {
    python := dag.Python().New("", "", "", "", nil)
    
    // Criar ambiente de desenvolvimento
    devEnv := python.BuildEnv(ctx, dag.Directory("."))
    
    // Configurar ambiente
    return devEnv.
        WithExec([]string{
            "sh", "-c",
            `
            # Instalar ferramentas de desenvolvimento
            pip install ipython debugpy
            
            # Configurar ambiente
            poetry install --with dev
            
            # Iniciar shell interativo
            poetry shell
            `,
        }).
        Sync(ctx)
} 