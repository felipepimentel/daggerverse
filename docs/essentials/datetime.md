---
layout: default
title: DateTime
parent: Essentials
nav_order: 4
---

# DateTime

O módulo DateTime fornece uma interface para manipulação e formatação de datas e horários em pipelines Dagger, permitindo operações precisas com timestamps e fusos horários.

## Features

- Manipulação de datas
- Formatação de timestamps
- Suporte a fusos horários
- Cálculos de duração
- Conversão de formatos
- Validação de datas
- Operações aritméticas
- Parsing de strings

## Instalação

Para usar o módulo DateTime em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/datetime"
)
```

## Exemplos de Uso

### Data Atual

```go
func (m *MyModule) CurrentDate(ctx context.Context) (string, error) {
    datetime := dag.DateTime().
        WithFormat("2006-01-02")
    
    // Obter data atual
    return datetime.Now(ctx)
}
```

### Formatação de Timestamp

```go
func (m *MyModule) FormatTimestamp(ctx context.Context) (string, error) {
    datetime := dag.DateTime().
        WithTimestamp(1234567890).
        WithFormat("2006-01-02 15:04:05")
    
    // Formatar timestamp
    return datetime.Format(ctx)
}
```

### Operações com Datas

```go
func (m *MyModule) DateOperations(ctx context.Context) (string, error) {
    datetime := dag.DateTime().
        WithDate("2023-01-01").
        WithFormat("2006-01-02")
    
    // Adicionar dias
    return datetime.AddDays(ctx, 7)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: DateTime Operations
on: [push]

jobs:
  process:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Process Dates
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/datetime
          args: |
            do -p '
              datetime := DateTime().
                WithFormat("2006-01-02").
                Now(ctx)
            '
```

## Referência da API

### DateTime

Estrutura principal que fornece acesso à funcionalidade de data e hora.

#### Construtor

- `New() *DateTime`
  - Cria uma nova instância do DateTime

#### Métodos de Configuração

- `WithFormat(format string) *DateTime`
  - Define formato de data/hora
  - Parâmetro:
    - `format`: Formato Go de data/hora

- `WithTimestamp(ts int64) *DateTime`
  - Define timestamp Unix
  - Parâmetro:
    - `ts`: Timestamp em segundos

- `WithDate(date string) *DateTime`
  - Define data inicial
  - Parâmetro:
    - `date`: Data em string

- `WithTimezone(tz string) *DateTime`
  - Define fuso horário
  - Parâmetro:
    - `tz`: Nome do fuso horário

#### Métodos de Operação

- `Now(ctx context.Context) (string, error)`
  - Retorna data/hora atual formatada
  - Retorna string formatada

- `Format(ctx context.Context) (string, error)`
  - Formata timestamp configurado
  - Retorna string formatada

- `AddDays(ctx context.Context, days int) (string, error)`
  - Adiciona dias à data
  - Parâmetro:
    - `days`: Número de dias
  - Retorna data resultante

- `Parse(ctx context.Context, dateStr string) (int64, error)`
  - Converte string para timestamp
  - Parâmetro:
    - `dateStr`: Data em string
  - Retorna timestamp Unix

## Boas Práticas

1. **Formatação**
   - Use formatos ISO quando possível
   - Documente formatos usados
   - Mantenha consistência

2. **Fusos Horários**
   - Sempre especifique timezone
   - Use UTC para logs
   - Considere localização

3. **Validação**
   - Valide entradas
   - Trate erros de parsing
   - Verifique limites

4. **Performance**
   - Cache resultados comuns
   - Otimize operações em lote
   - Minimize conversões

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Formato**
   ```
   Erro: Invalid format
   Solução: Verifique string de formato
   ```

2. **Erro de Parsing**
   ```
   Erro: Cannot parse date
   Solução: Verifique formato da entrada
   ```

3. **Erro de Timezone**
   ```
   Erro: Unknown timezone
   Solução: Verifique nome do timezone
   ```

## Exemplo de Configuração

```yaml
# datetime.yaml
format:
  default: "2006-01-02 15:04:05"
  date: "2006-01-02"
  time: "15:04:05"
  iso: "2006-01-02T15:04:05Z07:00"

timezone:
  default: "UTC"
  allowed:
    - "America/New_York"
    - "Europe/London"
    - "Asia/Tokyo"

validation:
  min_year: 1970
  max_year: 2100
  strict_parsing: true
```

## Uso Avançado

### Pipeline de Processamento

```go
func (m *MyModule) ProcessingPipeline(ctx context.Context) error {
    // Configurar datetime
    datetime := dag.DateTime().
        WithFormat("2006-01-02").
        WithTimezone("UTC")
    
    // Obter data atual
    now, err := datetime.Now(ctx)
    if err != nil {
        return err
    }
    
    // Adicionar uma semana
    nextWeek, err := datetime.
        WithDate(now).
        AddDays(ctx, 7)
    if err != nil {
        return err
    }
    
    // Converter para timestamp
    ts, err := datetime.Parse(ctx, nextWeek)
    if err != nil {
        return err
    }
    
    // Formatar resultado
    return datetime.
        WithTimestamp(ts).
        WithFormat("2006-01-02 15:04:05").
        Format(ctx)
}
```

### Manipulação Avançada

```go
func (m *MyModule) AdvancedManipulation(ctx context.Context) error {
    // Configurar datetime
    datetime := dag.DateTime().
        WithFormat("2006-01-02T15:04:05Z07:00").
        WithTimezone("America/New_York")
    
    // Processar datas
    dates := []string{
        "2023-01-01",
        "2023-06-15",
        "2023-12-31",
    }
    
    for _, date := range dates {
        // Converter para timestamp
        ts, err := datetime.Parse(ctx, date)
        if err != nil {
            return err
        }
        
        // Formatar em diferentes timezones
        result, err := datetime.
            WithTimestamp(ts).
            WithTimezone("Europe/London").
            Format(ctx)
        if err != nil {
            return err
        }
        
        fmt.Printf("NY: %s, London: %s\n", date, result)
    }
    
    return nil
} 