# N8N Deployment Troubleshooting

## Issue #1: Connection and Service Verification

### Problem
- Deployment process is taking too long with multiple retry attempts
- Service verification is failing
- DNS propagation checks are timing out
- Múltiplas chaves SSH sendo criadas sem limpeza

### Investigation Steps

1. Verificação dos droplets ativos:
```bash
$ doctl compute droplet list --format "ID,Name,PublicIPv4,Status"
ID           Name    Public IPv4        Status
473946812    n8n     134.122.113.92     active
473973038    n8n     198.199.64.245     active
473976101    n8n     134.209.221.159    active
```

2. Verificação das chaves SSH:
```bash
$ doctl compute ssh-key list
# Resultado: Mais de 70 chaves SSH registradas para o mesmo projeto
```

### Problemas Identificados
1. Múltiplos droplets ativos causando conflitos de DNS
2. Falha na autenticação SSH - chave pública não está corretamente configurada
3. Processo de limpeza de recursos antigos não está funcionando como esperado
4. Acúmulo excessivo de chaves SSH (mais de 70 chaves)
5. Cada execução do deploy está gerando uma nova chave SSH

### Análise do Problema
1. Processo atual:
   - Gera nova chave SSH a cada deploy
   - Registra a chave no DigitalOcean
   - Cria novo droplet
   - Não limpa recursos antigos

2. Problemas no fluxo:
   - Chaves SSH não são reutilizadas
   - Recursos antigos não são limpos
   - Falta persistência de estado entre execuções

3. Impacto:
   - Acúmulo de recursos
   - Falhas de autenticação
   - Processo instável

### Ações Corretivas

1. Limpar recursos antigos:
```bash
# Remover droplets antigos
doctl compute droplet delete 473946812 473973038 --force

# Listar todas as chaves SSH antigas
doctl compute ssh-key list --format ID,Name --no-header | grep "n8n-deploy-" | awk '{print $1}' > /tmp/old-keys.txt

# Remover chaves antigas
while read -r key_id; do
    doctl compute ssh-key delete "$key_id" --force
done < /tmp/old-keys.txt
```

2. Implementar nova estratégia de chaves SSH:
```bash
# Criar chave SSH persistente
ssh-keygen -t ed25519 -f ~/.ssh/n8n-deploy -N ""

# Registrar chave no DigitalOcean com nome fixo
doctl compute ssh-key create n8n-deploy-key --public-key-file ~/.ssh/n8n-deploy.pub
```

3. Ajustar código do deploy:
- Usar chave SSH persistente
- Implementar limpeza de recursos
- Melhorar verificação de status

### Modificações no Código

1. Alterações necessárias em `main.go`:
```go
// Usar chave SSH persistente
func (n *N8N) getSSHKey() (*SSHKeys, error) {
    // Tentar ler chave existente
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    keyPath := filepath.Join(homeDir, ".ssh", "n8n-deploy")
    if _, err := os.Stat(keyPath); err == nil {
        // Ler chave existente
        privateKey, err := os.ReadFile(keyPath)
        if err != nil {
            return nil, err
        }
        publicKey, err := os.ReadFile(keyPath + ".pub")
        if err != nil {
            return nil, err
        }
        return &SSHKeys{
            name:       "n8n-deploy-key",
            privateKey: string(privateKey),
            publicKey:  string(publicKey),
        }, nil
    }

    // Se não existir, criar nova
    return n.generateSSHKeys()
}

// Implementar limpeza de recursos
func (n *N8N) cleanup() error {
    // Implementar limpeza de droplets e chaves SSH antigas
}
```

### Status Atual
- Identificados problemas com gerenciamento de recursos
- Definida nova estratégia para chaves SSH
- Em processo de implementação das correções

### Próximos Passos
1. Implementar persistência da chave SSH
2. Ajustar processo de limpeza de recursos
3. Melhorar verificação de status dos serviços
4. Implementar retry com backoff exponencial
5. Adicionar logs detalhados para troubleshooting

### Resolução
(Em andamento - Implementando nova estratégia de gerenciamento de recursos) 