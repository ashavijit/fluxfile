# Contributing to Flux

Thank you for your interest in contributing to Flux!

## How to Contribute

### Reporting Issues

- Check existing issues before creating a new one
- Include FluxFile version (`flux -v`)
- Provide a minimal reproducible example

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit with a descriptive message
6. Push and create a Pull Request

### Code Style

- Follow standard Go conventions
- Run `go fmt ./...` before committing
- Add tests for new functionality

### Development Setup

```bash
git clone https://github.com/ashavijit/fluxfile
cd fluxfile
flux build
./bin/fluxfile -v
```

### Running Tests

### Running Tests

```bash
flux test
```


### Git Hooks

To ensure code quality, we use git hooks to run linting and tests before pushing checks.
Please run the following command to set up the hooks:

```bash
flux setup-hooks
```

This will install a pre-push hook that runs `go fmt`, `golangci-lint`, and `go test` before you push your changes.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
