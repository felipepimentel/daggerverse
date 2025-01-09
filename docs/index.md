---
layout: default
title: Daggerverse Documentation
nav_order: 1
---

# Daggerverse Documentation

Welcome to the Daggerverse documentation! This repository contains a collection of reusable [Dagger](https://dagger.io) modules that can help you build better CI/CD pipelines.

## Available Modules

### Libraries

- [AWS CLI](/libraries/aws-cli) - AWS Command Line Interface module
- [Caddy](/libraries/caddy) - Caddy web server module
- [Docker](/libraries/docker) - Docker container and image management module
- [Docker Compose](/libraries/docker-compose) - Docker Compose module
- [Docusaurus](/libraries/docusaurus) - Docusaurus documentation site module
- [Envoy](/libraries/envoy) - Envoy proxy module
- [GitHub](/libraries/gh) - GitHub operations module
- [Helm](/libraries/helm) - Helm package manager module
- [JFrog CLI](/libraries/jfrogcli) - JFrog CLI module
- [Kafka](/libraries/kafka) - Apache Kafka module
- [LaunchDarkly](/libraries/launchdarkly) - LaunchDarkly feature flag module
- [Nginx](/libraries/nginx) - Nginx web server module
- [OpenAPI Changes](/libraries/openapi-changes) - OpenAPI diff module
- [OpenAPI Codegen](/libraries/openapi-codegen) - OpenAPI code generator module
- [Poetry](/libraries/poetry) - Python Poetry module
- [PostgreSQL](/libraries/postgres) - PostgreSQL database module
- [PSQL](/libraries/psql) - PostgreSQL client module
- [PyPI](/libraries/pypi) - Python Package Index module
- [Secrets Manager](/libraries/secretsmanager) - AWS Secrets Manager module
- [Spectral](/libraries/spectral) - OpenAPI linter module

### Essentials

- [Alpine](/essentials/alpine) - Alpine Linux base image module
- [Apko](/essentials/apko) - Alpine Package Keeper module
- [Checksum](/essentials/checksum) - File checksum module
- [DateTime](/essentials/datetime) - Date and time utilities module
- [Get IP](/essentials/get-ip) - IP address utilities module
- [Git](/essentials/git) - Git operations module
- [Git Changelog](/essentials/git-chglog) - Git changelog generator module
- [Inline Node](/essentials/inline-node) - Node.js inline execution module
- [Inline Python](/essentials/inline-python) - Python inline execution module
- [Make](/essentials/make) - GNU Make module
- [Registry Config](/essentials/registry-config) - Container registry configuration module
- [Ruff](/essentials/ruff) - Python linter module
- [SSH](/essentials/ssh) - SSH operations module
- [SSH Keygen](/essentials/ssh-keygen) - SSH key generation module
- [TTL.sh](/essentials/ttlsh) - Ephemeral container registry module
- [Versioner](/essentials/versioner) - Version management module
- [Wolfi](/essentials/wolfi) - Wolfi container base image module

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

Please read our [Contributing Guidelines](https://github.com/felipepimentel/daggerverse/blob/main/CONTRIBUTING.md) before submitting any changes.

## License

This project is licensed under the [Apache License 2.0](https://github.com/felipepimentel/daggerverse/blob/main/LICENSE). 