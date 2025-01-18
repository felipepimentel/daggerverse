# Daggerverse Documentation

This directory contains the documentation for the Daggerverse modules. The documentation is built using [Jekyll](https://jekyllrb.com/) and the [Just the Docs](https://just-the-docs.github.io/just-the-docs/) theme.

## Local Development

To run the documentation site locally:

1. Install Ruby and Bundler
2. Install dependencies:
   ```bash
   cd docs
   bundle install
   ```
3. Start the local server:
   ```bash
   bundle exec jekyll serve
   ```
4. Open http://localhost:4000 in your browser

## Documentation Structure

- `index.md`: Main landing page
- `_config.yml`: Jekyll configuration
- `Gemfile`: Ruby dependencies
- `libraries/`: Documentation for library modules
  - `docker.md`: Docker module documentation
  - `aws-cli.md`: AWS CLI module documentation
  - `gh.md`: GitHub module documentation
  - `postgres.md`: PostgreSQL module documentation
- `essentials/`: Documentation for essential modules
  - `wolfi.md`: Wolfi module documentation

## Contributing

1. Fork the repository
2. Create a new branch for your changes
3. Make your changes
4. Submit a pull request

## Building

The documentation is automatically built and deployed to GitHub Pages when changes are pushed to the `main` branch in the `docs/` directory. 