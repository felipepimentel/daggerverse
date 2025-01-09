---
layout: default
title: Daggerverse Documentation
nav_order: 1
---

# Daggerverse Documentation

Welcome to the Daggerverse documentation! This repository contains a collection of reusable [Dagger](https://dagger.io) modules that can help you build better CI/CD pipelines.

## Available Modules

### Libraries

- [AWS CLI](/daggerverse/libraries/aws-cli) - AWS Command Line Interface module
- [Caddy](/daggerverse/libraries/caddy) - Caddy web server module
- [Docker](/daggerverse/libraries/docker) - Docker container and image management module
- [Docker Compose](/daggerverse/libraries/docker-compose) - Docker Compose module
- [Docusaurus](/daggerverse/libraries/docusaurus) - Docusaurus documentation site module
- [Envoy](/daggerverse/libraries/envoy) - Envoy proxy module
- [GitHub](/daggerverse/libraries/gh) - GitHub operations module
- [Helm](/daggerverse/libraries/helm) - Helm package manager module
- [JFrog CLI](/daggerverse/libraries/jfrogcli) - JFrog CLI module
- [Kafka](/daggerverse/libraries/kafka) - Apache Kafka module
- [LaunchDarkly](/daggerverse/libraries/launchdarkly) - LaunchDarkly feature flag module
- [Nginx](/daggerverse/libraries/nginx) - Nginx web server module
- [OpenAPI Changes](/daggerverse/libraries/openapi-changes) - OpenAPI diff module
- [OpenAPI Codegen](/daggerverse/libraries/openapi-codegen) - OpenAPI code generator module
- [Poetry](/daggerverse/libraries/poetry) - Python Poetry module
- [PostgreSQL](/daggerverse/libraries/postgres) - PostgreSQL database module
- [PSQL](/daggerverse/libraries/psql) - PostgreSQL client module
- [PyPI](/daggerverse/libraries/pypi) - Python Package Index module
- [Secrets Manager](/daggerverse/libraries/secretsmanager) - AWS Secrets Manager module
- [Spectral](/daggerverse/libraries/spectral) - OpenAPI linter module

### Essentials

- [Alpine](/daggerverse/essentials/alpine) - Alpine Linux base image module
- [Apko](/daggerverse/essentials/apko) - Alpine Package Keeper module
- [Checksum](/daggerverse/essentials/checksum) - File checksum module
- [DateTime](/daggerverse/essentials/datetime) - Date and time utilities module
- [Get IP](/daggerverse/essentials/get-ip) - IP address utilities module
- [Git](/daggerverse/essentials/git) - Git operations module
- [Git Changelog](/daggerverse/essentials/git-chglog) - Git changelog generator module
- [Inline Node](/daggerverse/essentials/inline-node) - Node.js inline execution module
- [Inline Python](/daggerverse/essentials/inline-python) - Python inline execution module
- [Make](/daggerverse/essentials/make) - GNU Make module
- [Registry Config](/daggerverse/essentials/registry-config) - Container registry configuration module
- [Ruff](/daggerverse/essentials/ruff) - Python linter module
- [SSH](/daggerverse/essentials/ssh) - SSH operations module
- [SSH Keygen](/daggerverse/essentials/ssh-keygen) - SSH key generation module
- [TTL.sh](/daggerverse/essentials/ttlsh) - Ephemeral container registry module
- [Versioner](/daggerverse/essentials/versioner) - Version management module
- [Wolfi](/daggerverse/essentials/wolfi) - Wolfi container base image module

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