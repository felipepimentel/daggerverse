---
layout: default
title: Get-IP
parent: Essentials
nav_order: 5
---

# Get-IP

O módulo Get-IP fornece uma interface para obter e gerenciar endereços IP em pipelines Dagger, permitindo consultas a serviços de IP públicos e privados.

## Features

- Obtenção de IP público
- Consulta de IP privado
- Suporte a IPv4 e IPv6
- Validação de endereços
- Cache de resultados
- Múltiplos provedores
- Formatação de saída
- Verificação de conectividade

## Instalação

Para usar o módulo Get-IP em seu pipeline Dagger:

```go
import (
    "dagger.io/dagger"
    "github.com/felipepimentel/daggerverse/essentials/get-ip"
)
```

## Exemplos de Uso

### IP Público

```go
func (m *MyModule) GetPublicIP(ctx context.Context) (string, error) {
    getip := dag.GetIP().
        WithProvider("ipify")
    
    // Obter IP público
    return getip.Public(ctx)
}
```

### IP Privado

```go
func (m *MyModule) GetPrivateIP(ctx context.Context) (string, error) {
    getip := dag.GetIP().
        WithInterface("eth0")
    
    // Obter IP privado
    return getip.Private(ctx)
}
```

### Verificação de IP

```go
func (m *MyModule) CheckIP(ctx context.Context) error {
    getip := dag.GetIP().
        WithProvider("ipify").
        WithValidation(true)
    
    // Verificar IP
    return getip.Validate(ctx)
}
```

## Integração com GitHub Actions

Exemplo de workflow usando o módulo:

```yaml
name: IP Check
on: [push]

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Get Public IP
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/essentials/get-ip
          args: |
            do -p '
              getip := GetIP().
                WithProvider("ipify")
              getip.Public(ctx)
            '
```

## Referência da API

### GetIP

Estrutura principal que fornece acesso à funcionalidade de IP.

#### Construtor

- `New() *GetIP`
  - Cria uma nova instância do GetIP

#### Métodos de Configuração

- `WithProvider(provider string) *GetIP`
  - Define provedor de IP público
  - Parâmetro:
    - `provider`: Nome do provedor

- `WithInterface(iface string) *GetIP`
  - Define interface de rede
  - Parâmetro:
    - `iface`: Nome da interface

- `WithValidation(validate bool) *GetIP`
  - Ativa validação de IP
  - Parâmetro:
    - `validate`: Se true, valida IPs

- `WithTimeout(seconds int) *GetIP`
  - Define timeout para requisições
  - Parâmetro:
    - `seconds`: Timeout em segundos

#### Métodos de Operação

- `Public(ctx context.Context) (string, error)`
  - Obtém IP público
  - Retorna endereço IP

- `Private(ctx context.Context) (string, error)`
  - Obtém IP privado
  - Retorna endereço IP

- `Validate(ctx context.Context) error`
  - Valida endereço IP
  - Retorna erro se inválido

- `Info(ctx context.Context) (map[string]string, error)`
  - Obtém informações detalhadas do IP
  - Retorna mapa com informações

## Boas Práticas

1. **Provedores**
   - Use provedores confiáveis
   - Implemente fallback
   - Monitore limites

2. **Cache**
   - Cache resultados frequentes
   - Defina TTL apropriado
   - Limpe cache regularmente

3. **Validação**
   - Valide formatos
   - Verifique ranges
   - Trate exceções

4. **Performance**
   - Use timeouts adequados
   - Minimize requisições
   - Otimize cache

## Solução de Problemas

Problemas comuns e soluções:

1. **Erro de Provedor**
   ```
   Erro: Provider unavailable
   Solução: Tente provedor alternativo
   ```

2. **Erro de Timeout**
   ```
   Erro: Request timeout
   Solução: Aumente timeout ou verifique conexão
   ```

3. **Erro de Validação**
   ```
   Erro: Invalid IP format
   Solução: Verifique formato do endereço
   ```

## Exemplo de Configuração

```yaml
# getip.yaml
providers:
  - name: ipify
    url: https://api.ipify.org
    timeout: 5
  - name: icanhazip
    url: https://icanhazip.com
    timeout: 5

interfaces:
  - eth0
  - wlan0

validation:
  enabled: true
  formats:
    - ipv4
    - ipv6

cache:
  enabled: true
  ttl: 300
```

## Uso Avançado

### Pipeline de Verificação

```go
func (m *MyModule) VerificationPipeline(ctx context.Context) error {
    // Configurar GetIP
    getip := dag.GetIP().
        WithProvider("ipify").
        WithValidation(true).
        WithTimeout(10)
    
    // Obter IP público
    publicIP, err := getip.Public(ctx)
    if err != nil {
        return err
    }
    
    // Obter informações detalhadas
    info, err := getip.Info(ctx)
    if err != nil {
        return err
    }
    
    // Validar e processar
    if err := getip.Validate(ctx); err != nil {
        return err
    }
    
    fmt.Printf("IP: %s\nInfo: %v\n", publicIP, info)
    return nil
}
```

### Múltiplos Provedores

```go
func (m *MyModule) MultiProvider(ctx context.Context) (string, error) {
    // Lista de provedores
    providers := []string{
        "ipify",
        "icanhazip",
        "ipapi",
    }
    
    // Tentar cada provedor
    for _, provider := range providers {
        ip, err := dag.GetIP().
            WithProvider(provider).
            WithTimeout(5).
            Public(ctx)
        
        if err == nil {
            return ip, nil
        }
        
        fmt.Printf("Provedor %s falhou: %v\n", provider, err)
    }
    
    return "", fmt.Errorf("todos os provedores falharam")
} 