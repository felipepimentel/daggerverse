---
layout: default
title: Inline-Node
parent: Essentials
nav_order: 8
---

# Inline-Node

O módulo Inline-Node fornece uma interface para executar código Node.js diretamente em pipelines Dagger, permitindo automação e scripts personalizados.

## Features

- Execução de código Node.js
- Gerenciamento de dependências
- Ambiente isolado
- Suporte a módulos
- Integração com npm
- Scripts personalizados
- Gestão de pacotes
- Execução assíncrona

## Instalação

Para usar o módulo Inline-Node em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/inline-node"
)
```

## Exemplos de Uso

### Execução Básica

```go
func (m *MyModule) RunScript(ctx context.Context) (string, error) {
    node := dag.InlineNode().
        WithCode(`
            console.log('Hello from Node.js!');
            return 'Success';
        `)
    
    // Executar código
    return node.Run(ctx)
}
```

### Com Dependências

```go
func (m *MyModule) WithDependencies(ctx context.Context) (string, error) {
    node := dag.InlineNode().
        WithPackages([]string{"axios", "lodash"}).
        WithCode(`
            const axios = require('axios');
            const _ = require('lodash');
            
            const response = await axios.get('https://api.example.com/data');
            return _.get(response, 'data.value', 'default');
        `)
    
    // Executar com dependências
    return node.Run(ctx)
}
```

### Script Assíncrono

```go
func (m *MyModule) AsyncScript(ctx context.Context) (string, error) {
    node := dag.InlineNode().
        WithCode(`
            async function processData() {
                await new Promise(resolve => setTimeout(resolve, 1000));
                return 'Processed';
            }
            
            return await processData();
        `)
    
    // Executar assincronamente
    return node.Run(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: Node.js Script
on: [push]

jobs:
  script:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Node.js Script
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/inline-node
          args: |
            do -p '
              node := InlineNode().
                WithCode("console.log(\"Hello from Node.js!\")")
              node.Run(ctx)
            '
```

## Referência da API

### InlineNode

Estrutura principal que fornece acesso à funcionalidade Node.js.

#### Construtor

- `New() *InlineNode`
  - Cria uma nova instância do InlineNode

#### Métodos de Configuração

- `WithCode(code string) *InlineNode`
  - Define código a ser executado
  - Parâmetro:
    - `code`: Código Node.js

- `WithPackages(packages []string) *InlineNode`
  - Define pacotes npm
  - Parâmetro:
    - `packages`: Lista de pacotes

- `WithNodeVersion(version string) *InlineNode`
  - Define versão do Node.js
  - Parâmetro:
    - `version`: Versão do Node.js

- `WithEnv(key string, value string) *InlineNode`
  - Define variável de ambiente
  - Parâmetros:
    - `key`: Nome da variável
    - `value`: Valor da variável

#### Métodos de Operação

- `Run(ctx context.Context) (string, error)`
  - Executa código Node.js
  - Retorna resultado da execução

- `RunAsync(ctx context.Context) error`
  - Executa código assincronamente
  - Retorna erro se falhar

- `Install(ctx context.Context) error`
  - Instala dependências
  - Retorna erro se falhar

## Boas Práticas

1. **Dependências**
   - Especifique versões
   - Use package.json
   - Minimize dependências

2. **Código**
   - Mantenha código limpo
   - Use async/await
   - Trate erros

3. **Performance**
   - Otimize execução
   - Cache dependências
   - Minimize operações

4. **Segurança**
   - Valide entradas
   - Proteja credenciais
   - Atualize pacotes

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Dependência**
   ```
   Erro: Package not found
   Solução: Verifique nome do pacote
   ```

2. **Erro de Execução**
   ```
   Erro: Runtime error
   Solução: Verifique sintaxe do código
   ```

3. **Erro de Versão**
   ```
   Erro: Node version not supported
   Solução: Use versão compatível
   ```

## Exemplo de Configuração

```yaml
# node.yaml
version: 18
packages:
  - axios@0.24.0
  - lodash@4.17.21
  - moment@2.29.4

environment:
  NODE_ENV: production
  DEBUG: false

scripts:
  - name: process
    code: |
      const data = await processData();
      return JSON.stringify(data);
  - name: validate
    code: |
      return validateInput(process.env.INPUT);
```

## Uso Avançado

### Pipeline Completo

```go
func (m *MyModule) CompletePipeline(ctx context.Context) error {
    // Configurar Node.js
    node := dag.InlineNode().
        WithNodeVersion("18").
        WithPackages([]string{
            "axios",
            "lodash",
            "moment",
        }).
        WithEnv("NODE_ENV", "production")
    
    // Instalar dependências
    if err := node.Install(ctx); err != nil {
        return err
    }
    
    // Executar processamento
    result, err := node.
        WithCode(`
            const axios = require('axios');
            const _ = require('lodash');
            const moment = require('moment');
            
            async function processData() {
                const response = await axios.get('https://api.example.com/data');
                const data = _.get(response, 'data', {});
                return {
                    ...data,
                    timestamp: moment().format(),
                };
            }
            
            return await processData();
        `).
        Run(ctx)
    
    if err != nil {
        return err
    }
    
    fmt.Println("Resultado:", result)
    return nil
}
```

### Processamento Avançado

```go
func (m *MyModule) AdvancedProcessing(ctx context.Context) error {
    // Configurar processamento
    node := dag.InlineNode().
        WithNodeVersion("18").
        WithPackages([]string{"sharp", "fs-extra"}).
        WithCode(`
            const sharp = require('sharp');
            const fs = require('fs-extra');
            
            async function processImages() {
                const files = await fs.readdir('/input');
                
                for (const file of files) {
                    if (file.match(/\.(jpg|jpeg|png)$/i)) {
                        await sharp('/input/' + file)
                            .resize(800, 600)
                            .jpeg({ quality: 80 })
                            .toFile('/output/' + file);
                    }
                }
                
                return 'Processed ' + files.length + ' images';
            }
            
            return await processImages();
        `)
    
    // Executar processamento
    return node.RunAsync(ctx)
} 