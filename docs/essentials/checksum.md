---
layout: default
title: Checksum
parent: Essentials
nav_order: 3
---

# Checksum

O módulo Checksum fornece uma interface para calcular e verificar checksums de arquivos e diretórios, garantindo integridade e segurança dos dados.

## Features

- Cálculo de checksums
- Suporte a múltiplos algoritmos
- Verificação de integridade
- Processamento de diretórios
- Geração de relatórios
- Comparação de checksums
- Cache de resultados
- Validação de arquivos

## Instalação

Para usar o módulo Checksum em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/checksum"
)
```

## Exemplos de Uso

### Cálculo Básico

```go
func (m *MyModule) CalculateChecksum(ctx context.Context) (string, error) {
    checksum := dag.Checksum().
        WithFile(dag.Directory(".")).
        WithAlgorithm("sha256")
    
    // Calcular checksum
    return checksum.Calculate(ctx)
}
```

### Verificação de Arquivo

```go
func (m *MyModule) VerifyFile(ctx context.Context) error {
    checksum := dag.Checksum().
        WithFile(dag.File("./package.tar.gz")).
        WithExpected("abc123...")
    
    // Verificar checksum
    return checksum.Verify(ctx)
}
```

### Processamento de Diretório

```go
func (m *MyModule) ProcessDirectory(ctx context.Context) (string, error) {
    checksum := dag.Checksum().
        WithDirectory(dag.Directory("./src")).
        WithAlgorithm("sha512").
        WithRecursive(true)
    
    // Calcular checksum do diretório
    return checksum.Calculate(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: Checksum Verification
on: [push]

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Verify Checksums
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/checksum
          args: |
            do -p '
              checksum := Checksum().
                WithFile(Directory(".")).
                WithAlgorithm("sha256")
              checksum.Calculate(ctx)
            '
```

## Referência da API

### Checksum

Estrutura principal que fornece acesso à funcionalidade de checksum.

#### Construtor

- `New() *Checksum`
  - Cria uma nova instância do Checksum

#### Métodos de Configuração

- `WithFile(file *File) *Checksum`
  - Define arquivo para processamento
  - Parâmetro:
    - `file`: Arquivo a ser processado

- `WithDirectory(dir *Directory) *Checksum`
  - Define diretório para processamento
  - Parâmetro:
    - `dir`: Diretório a ser processado

- `WithAlgorithm(algo string) *Checksum`
  - Define algoritmo de hash
  - Parâmetro:
    - `algo`: "md5", "sha1", "sha256", "sha512"

- `WithRecursive(recursive bool) *Checksum`
  - Define processamento recursivo
  - Parâmetro:
    - `recursive`: Se true, processa subdiretórios

- `WithExpected(expected string) *Checksum`
  - Define checksum esperado
  - Parâmetro:
    - `expected`: Hash esperado

#### Métodos de Operação

- `Calculate(ctx context.Context) (string, error)`
  - Calcula checksum
  - Retorna hash calculado

- `Verify(ctx context.Context) error`
  - Verifica se checksum corresponde ao esperado
  - Retorna erro se não corresponder

- `GenerateReport(ctx context.Context) (*File, error)`
  - Gera relatório de checksums
  - Retorna arquivo com relatório

## Boas Práticas

1. **Algoritmos**
   - Use SHA-256 ou superior
   - Evite MD5 para segurança
   - Documente algoritmos usados

2. **Verificação**
   - Sempre verifique downloads
   - Compare com fonte confiável
   - Mantenha registros

3. **Performance**
   - Use cache quando possível
   - Processe arquivos grandes em partes
   - Otimize para diretórios grandes

4. **Segurança**
   - Valide fontes de hash
   - Proteja relatórios
   - Monitore alterações

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Algoritmo**
   ```
   Erro: Invalid algorithm
   Solução: Verifique nome do algoritmo
   ```

2. **Erro de Verificação**
   ```
   Erro: Checksum mismatch
   Solução: Verifique integridade do arquivo
   ```

3. **Erro de Acesso**
   ```
   Erro: Permission denied
   Solução: Verifique permissões do arquivo
   ```

## Exemplo de Configuração

```yaml
# checksum.yaml
algorithm: sha256
recursive: true
ignore:
  - "*.tmp"
  - "*.log"
  - ".git/"

report:
  format: json
  output: checksums.json

verify:
  strict: true
  fail_fast: true
```

## Uso Avançado

### Pipeline de Verificação

```go
func (m *MyModule) VerificationPipeline(ctx context.Context) error {
    // Configurar checksum
    checksum := dag.Checksum().
        WithDirectory(dag.Directory("./dist")).
        WithAlgorithm("sha512").
        WithRecursive(true)
    
    // Calcular checksums
    result, err := checksum.Calculate(ctx)
    if err != nil {
        return err
    }
    
    // Gerar relatório
    report, err := checksum.GenerateReport(ctx)
    if err != nil {
        return err
    }
    
    // Verificar checksums anteriores
    return checksum.
        WithExpected(result).
        Verify(ctx)
}
```

### Processamento Customizado

```go
func (m *MyModule) CustomProcessing(ctx context.Context) error {
    // Criar configuração
    config := `
    algorithm: sha512
    recursive: true
    ignore:
      - "*.tmp"
      - "*.log"
    report:
      format: json
      output: checksums.json
    `
    
    // Configurar checksum
    checksum := dag.Checksum().
        WithFile(dag.File("checksum.yaml", config)).
        WithDirectory(dag.Directory("./src"))
    
    // Processar arquivos
    result, err := checksum.Calculate(ctx)
    if err != nil {
        return err
    }
    
    // Salvar resultado
    return checksum.WithFile(
        dag.File("checksum.txt", result),
    ).Sync(ctx)
} 