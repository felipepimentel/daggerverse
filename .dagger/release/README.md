# Dagger Release Tool

Esta ferramenta implementa o processo de release do repositório usando Dagger. Ela substitui o workflow do GitHub Actions (`release.yml`) por uma implementação em Go que pode ser executada tanto localmente quanto no CI.

## Funcionalidades

- Detecta módulos Dagger no repositório (procura por arquivos `dagger.json`)
- Para cada módulo:
  - Executa semantic-release para criar releases e tags
  - Publica o módulo no Daggerverse quando uma nova tag é criada
- Mantém a mesma configuração do semantic-release usada no workflow original

## Uso

1. Instale as dependências:

```bash
cd .dagger/release
go mod tidy
```

2. Execute a ferramenta:

```bash
GITHUB_TOKEN=your_token go run release.go
```

## Configuração

A configuração do semantic-release está no arquivo `.releaserc.json` e inclui:

- Análise de commits usando convenções Angular
- Geração de CHANGELOG
- Criação de tags específicas por módulo
- Integração com GitHub (comentários em PRs, labels)
