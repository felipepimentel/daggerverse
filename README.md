# Daggerverse

Collection of Dagger modules for various languages and tools.

## Development Setup

1. Install dependencies:

```bash
npm install
```

2. Setup git hooks:

```bash
npm run prepare
```

## Commit Convention

This repository enforces conventional commits with scopes. Your commits must follow this pattern:

```bash
type(scope): description

# Examples:
feat(python): add new testing feature
fix(python): resolve token handling
docs(global): update main README
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style changes
- `refactor`: Code changes that neither fix bugs nor add features
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `build`: Build system or external dependencies
- `ci`: CI/CD changes
- `chore`: Other changes that don't modify src or test files
- `revert`: Reverts a previous commit

### Scopes

- `python`: Python module changes
- `nodejs`: Node.js module changes
- `ruby`: Ruby module changes
- `global`: Repository-wide changes

### Rules

- Type must be one of the allowed types
- Scope is required and must be one of the defined scopes
- Description must be in lower case
- Breaking changes must be indicated by `!` after the type/scope

### Examples

```bash
# ✅ Good commits:
feat(python): add new testing feature
fix(python): resolve token handling issue
docs(global): update main readme
feat(python)!: change api interface

# ❌ Bad commits:
feat: add feature            # Missing scope
fix(invalid): something      # Invalid scope
feat(python): Add Feature    # Upper case in description
chore: some change          # Missing scope
```

## Modules

- [Python](python/README.md) - Python module with Poetry support
- [More modules coming soon...]
