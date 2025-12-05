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
go build ./cmd/flux
./flux -v
```

### Running Tests

```bash
go test ./... -v
```

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
