# Contributing to Vodka

First off, thank you for considering contributing to Vodka.

Vodka is a modern Go web framework focused on developer experience, full-stack workflow, and fast iteration. Contributions of all kinds are welcome — whether it's fixing bugs, improving documentation, optimizing performance, or proposing new ideas.

---

# Before Contributing

Please:

- Read the README carefully
- Search existing issues before creating a new one
- Keep pull requests focused and minimal
- Avoid large unrelated refactors in a single PR

---

# Development Setup

## Fork the Repository

Fork the repository in github

## Clone the Repository

```bash
git clone https://github.com/<your-github-username>/vodka.git
cd vodka
```

---

## Install Dependencies

```bash
go mod tidy
```

---

## Run Tests

```bash
go test ./...
```

---

## Run Race Detection

```bash
go test -race ./...
```

Concurrency safety is extremely important for framework development.

---

# Project Structure

```text
cmd/           -> CLI tooling
examples/      -> Example applications
mixers/        -> Built-in middlewares
```

---

# Coding Guidelines

## Keep APIs Clean

Vodka prioritizes developer experience and readability.

Avoid:
- unnecessary abstractions
- overengineering
- excessive magic behavior

Prefer:
- explicit APIs
- clean naming
- predictable behavior

---

## Maintain Consistency

Please follow the existing naming conventions and project structure.

Consistency is more important than personal style preferences.

---

## Write Minimal and Focused PRs

Smaller pull requests are easier to review and merge.

Good:
- fixing one bug
- improving one middleware
- adding one feature

Bad:
- rewriting half the framework in one PR titled "minor improvements"

---

# Commit Message Style

Recommended format:

```text
feat: add validation middleware
fix: resolve route param parsing issue
docs: improve README examples
refactor: simplify middleware chain logic
```

---

# Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Open a pull request with a clear description

---

# Reporting Bugs

When reporting bugs, include:

- Go version
- Operating system
- Steps to reproduce
- Expected behavior
- Actual behavior

If possible, include a minimal reproducible example.

---

# Feature Requests

Feature requests are welcome, but please ensure they align with Vodka's philosophy:

- fast development workflow
- strong developer experience
- lightweight architecture
- minimal boilerplate
- practical APIs

---

# Philosophy

Vodka is intentionally designed to stay lightweight and developer-focused.

The goal is not to become an everything-framework.

The goal is to provide:
- clean APIs
- modern tooling
- fast iteration
- practical abstractions

---

# Questions

If you have questions, feel free to open an issue or discussion.

Thanks again for contributing to Vodka 🚀
