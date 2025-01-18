---
layout: default
title: Inline-Python
parent: Essentials
nav_order: 9
---

# Inline-Python

O módulo Inline-Python fornece uma interface para executar código Python diretamente em pipelines Dagger, permitindo automação e scripts personalizados.

## Features

- Execução de código Python
- Gerenciamento de pacotes
- Ambiente virtual
- Suporte a módulos
- Integração com pip
- Scripts personalizados
- Gestão de dependências
- Execução assíncrona

## Instalação

Para usar o módulo Inline-Python em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/inline-python"
)
```

## Exemplos de Uso

### Execução Básica

```go
func (m *MyModule) RunScript(ctx context.Context) (string, error) {
    python := dag.InlinePython().
        WithCode(`
            print('Hello from Python!')
            return 'Success'
        `)
    
    // Executar código
    return python.Run(ctx)
}
```

### Com Dependências

```go
func (m *MyModule) WithDependencies(ctx context.Context) (string, error) {
    python := dag.InlinePython().
        WithPackages([]string{"requests", "pandas"}).
        WithCode(`
            import requests
            import pandas as pd
            
            response = requests.get('https://api.example.com/data')
            df = pd.DataFrame(response.json())
            return df.to_json()
        `)
    
    // Executar com dependências
    return python.Run(ctx)
}
```

### Script Assíncrono

```go
func (m *MyModule) AsyncScript(ctx context.Context) (string, error) {
    python := dag.InlinePython().
        WithCode(`
            import asyncio
            
            async def process_data():
                await asyncio.sleep(1)
                return 'Processed'
            
            return await process_data()
        `)
    
    // Executar assincronamente
    return python.Run(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: Python Script
on: [push]

jobs:
  script:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Python Script
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/inline-python
          args: |
            do -p '
              python := InlinePython().
                WithCode("print(\"Hello from Python!\")")
              python.Run(ctx)
            '
```

## Referência da API

### InlinePython

Estrutura principal que fornece acesso à funcionalidade Python.

#### Construtor

- `New() *InlinePython`
  - Cria uma nova instância do InlinePython

#### Métodos de Configuração

- `WithCode(code string) *InlinePython`
  - Define código a ser executado
  - Parâmetro:
    - `code`: Código Python

- `WithPackages(packages []string) *InlinePython`
  - Define pacotes pip
  - Parâmetro:
    - `packages`: Lista de pacotes

- `WithPythonVersion(version string) *InlinePython`
  - Define versão do Python
  - Parâmetro:
    - `version`: Versão do Python

- `WithEnv(key string, value string) *InlinePython`
  - Define variável de ambiente
  - Parâmetros:
    - `key`: Nome da variável
    - `value`: Valor da variável

#### Métodos de Operação

- `Run(ctx context.Context) (string, error)`
  - Executa código Python
  - Retorna resultado da execução

- `RunAsync(ctx context.Context) error`
  - Executa código assincronamente
  - Retorna erro se falhar

- `Install(ctx context.Context) error`
  - Instala dependências
  - Retorna erro se falhar

## Boas Práticas

1. **Dependências**
   - Use requirements.txt
   - Especifique versões
   - Minimize dependências

2. **Código**
   - Siga PEP 8
   - Use type hints
   - Documente funções

3. **Performance**
   - Otimize imports
   - Use async quando possível
   - Gerencie recursos

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
   Erro: Syntax error
   Solução: Verifique sintaxe do código
   ```

3. **Erro de Versão**
   ```
   Erro: Python version not supported
   Solução: Use versão compatível
   ```

## Exemplo de Configuração

```yaml
# python.yaml
version: 3.11
packages:
  - requests==2.31.0
  - pandas==2.1.3
  - numpy==1.24.3

environment:
  PYTHONPATH: /app
  PYTHONUNBUFFERED: "1"

scripts:
  - name: process
    code: |
      import pandas as pd
      data = pd.read_csv('data.csv')
      return data.to_json()
  - name: validate
    code: |
      def validate_input(data):
          return all(k in data for k in ['id', 'name'])
```

## Uso Avançado

### Pipeline Completo

```go
func (m *MyModule) CompletePipeline(ctx context.Context) error {
    // Configurar Python
    python := dag.InlinePython().
        WithPythonVersion("3.11").
        WithPackages([]string{
            "requests",
            "pandas",
            "numpy",
        }).
        WithEnv("PYTHONPATH", "/app")
    
    // Instalar dependências
    if err := python.Install(ctx); err != nil {
        return err
    }
    
    // Executar processamento
    result, err := python.
        WithCode(`
            import requests
            import pandas as pd
            import numpy as np
            
            def process_data():
                response = requests.get('https://api.example.com/data')
                df = pd.DataFrame(response.json())
                return {
                    'mean': np.mean(df['values']),
                    'std': np.std(df['values']),
                    'count': len(df)
                }
            
            return process_data()
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
    python := dag.InlinePython().
        WithPythonVersion("3.11").
        WithPackages([]string{"pillow", "numpy"}).
        WithCode(`
            from PIL import Image
            import numpy as np
            import os
            
            def process_images():
                files = os.listdir('/input')
                processed = 0
                
                for file in files:
                    if file.lower().endswith(('.png', '.jpg', '.jpeg')):
                        img = Image.open(f'/input/{file}')
                        arr = np.array(img)
                        
                        # Aplicar processamento
                        processed_arr = np.flip(arr, axis=1)  # Espelhar imagem
                        
                        # Salvar resultado
                        processed_img = Image.fromarray(processed_arr)
                        processed_img.save(f'/output/{file}')
                        processed += 1
                
                return f'Processed {processed} images'
            
            return process_images()
        `)
    
    // Executar processamento
    return python.RunAsync(ctx)
} 