# Daggerverse

Coleção de módulos Dagger e workflows reutilizáveis para CI/CD.

## Módulos Disponíveis

### Python

Módulo para CI/CD de projetos Python. Inclui:

- Testes automatizados
- Linting e formatação
- Build e publicação no PyPI
- Gerenciamento de versões

#### Uso via GitHub Actions (Recomendado)

Use nosso workflow reutilizável para integração simplificada:

```yaml
jobs:
  python-pipeline:
    uses: felipepimentel/daggerverse/.github/workflows/python-ci.yml@v1
    secrets:
      PYPI_TOKEN: ${{ secrets.PYPI_TOKEN }}
```

[Documentação completa do workflow](.github/workflows/README.md)

#### Uso Direto via Dagger

Para maior flexibilidade, use o módulo Dagger diretamente:

```bash
dagger call -m github.com/felipepimentel/daggerverse/python cicd --source . --token env:PYPI_TOKEN
```

[Documentação completa do módulo Python](python/README.md)

## Contribuindo

Para contribuir com este projeto:

1. Fork o repositório
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'feat: adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.
