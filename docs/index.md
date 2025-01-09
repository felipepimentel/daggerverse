---
layout: default
title: Daggerverse Documentation
nav_order: 1
---

# Daggerverse Documentation

Welcome to the Daggerverse documentation! This repository contains a collection of reusable [Dagger](https://dagger.io) modules that can help you build better CI/CD pipelines.

## Available Modules

### Libraries

- [AWS CLI](./libraries/aws-cli.md) - AWS Command Line Interface module
- [Caddy](./libraries/caddy.md) - Caddy web server module
- [Docker](./libraries/docker.md) - Docker container and image management module
- [Docker Compose](./libraries/docker-compose.md) - Docker Compose module
- [Docusaurus](./libraries/docusaurus.md) - Docusaurus documentation site module
- [Envoy](./libraries/envoy.md) - Envoy proxy module
- [GitHub](./libraries/gh.md) - GitHub operations module
- [Helm](./libraries/helm.md) - Helm package manager module
- [JFrog CLI](./libraries/jfrogcli.md) - JFrog CLI module
- [Kafka](./libraries/kafka.md) - Apache Kafka module
- [LaunchDarkly](./libraries/launchdarkly.md) - LaunchDarkly feature flag module
- [Nginx](./libraries/nginx.md) - Nginx web server module
- [OpenAPI Changes](./libraries/openapi-changes.md) - OpenAPI diff module
- [OpenAPI Codegen](./libraries/openapi-codegen.md) - OpenAPI code generator module
- [Poetry](./libraries/poetry.md) - Python Poetry module
- [PostgreSQL](./libraries/postgres.md) - PostgreSQL database module
- [PSQL](./libraries/psql.md) - PostgreSQL client module
- [PyPI](./libraries/pypi.md) - Python Package Index module
- [Secrets Manager](./libraries/secretsmanager.md) - AWS Secrets Manager module
- [Spectral](./libraries/spectral.md) - OpenAPI linter module

### Essentials

- [Alpine](./essentials/alpine.md) - Alpine Linux base image module
- [Apko](./essentials/apko.md) - Alpine Package Keeper module
- [Checksum](./essentials/checksum.md) - File checksum module
- [DateTime](./essentials/datetime.md) - Date and time utilities module
- [Get IP](./essentials/get-ip.md) - IP address utilities module
- [Git](./essentials/git.md) - Git operations module
- [Git Changelog](./essentials/git-chglog.md) - Git changelog generator module
- [Inline Node](./essentials/inline-node.md) - Node.js inline execution module
- [Inline Python](./essentials/inline-python.md) - Python inline execution module
- [Make](./essentials/make.md) - GNU Make module
- [Registry Config](./essentials/registry-config.md) - Container registry configuration module
- [Ruff](./essentials/ruff.md) - Python linter module
- [SSH](./essentials/ssh.md) - SSH operations module
- [SSH Keygen](./essentials/ssh-keygen.md) - SSH key generation module
- [TTL.sh](./essentials/ttlsh.md) - Ephemeral container registry module
- [Versioner](./essentials/versioner.md) - Version management module
- [Wolfi](./essentials/wolfi.md) - Wolfi container base image module

## Quick Start

To use any module in this repository, you can reference it in your Dagger pipeline using:

```go
import (
    "dagger.io/dagger"
    "<module_name>" "github.com/felipepimentel/daggerverse/<path_to_module>"
)
```

Or in your GitHub Actions workflow:

```yaml
name: CI
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Use Dagger module
        uses: dagger/dagger-action@v1
        with:
          module: github.com/felipepimentel/daggerverse/<path_to_module>
```

## Contributing

Please read our [Contributing Guidelines](../CONTRIBUTING.md) before submitting any changes.

## License

This project is licensed under the [Apache License 2.0](../LICENSE). 