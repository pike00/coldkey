# Contributing to coldkey

First off, thank you for considering contributing to coldkey! It's people like you that make open-source software such a great community to learn, inspire, and create.

We welcome all contributions, from bug reports and feature requests to code contributions and documentation improvements.

## Setting Up Your Development Environment

To work on coldkey, you will need:

- **Go**: The project is written in Go. You'll need a recent version installed.
- **Make/Just**: We use [`just`](https://github.com/casey/just) (a modern `make` alternative) to manage build tasks.

To get started:

1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/coldkey.git
   cd coldkey
   ```
3. Run the CI pipeline locally to ensure everything works:
   ```bash
   just ci
   ```

## Building Locally

To build the `coldkey` binary locally, run:

```bash
just build
```

This will produce the `coldkey` binary in the root directory.

## Running Tests

Testing is a crucial part of the development process. To run the test suite, use:

```bash
just test
```

## Code Style Expectations

We follow standard Go formatting conventions and use linters to maintain code quality. Before submitting your code, please ensure it complies with our standards:

- **Formatting**: Run `just fmt` to format your code using `gofumpt`.
- **Linting**: Run `just lint` to catch common issues using `golangci-lint`.
- **Vetting**: Run `just vet` for standard Go static analysis.

## Pull Request Process

When you're ready to submit your changes:

1. **Fork and Branch**: Create a new branch for your feature or bug fix (`git checkout -b my-new-feature`).
2. **Commit Messages**: Write clear, concise, and descriptive commit messages.
3. **Push**: Push your branch to your fork (`git push origin my-new-feature`).
4. **Open a PR**: Open a Pull Request against the `main` branch of this repository. Provide a detailed description of your changes and reference any related issues.

## Reporting Issues and Security Vulnerabilities

- **Bug Reports & Feature Requests**: Please use the GitHub Issues tracker to report bugs or request features. Search existing issues first to avoid duplicates.
- **Security Vulnerabilities**: If you discover a security vulnerability, please do NOT open a public issue. Instead, responsibly disclose it by contacting the project maintainers privately (e.g., via a GitHub security advisory or email if provided).

Thank you again for your time and contributions!
