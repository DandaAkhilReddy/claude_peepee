# Contributing to Claude PP

Thank you for your interest in contributing to Claude PP!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/claude_pp.git`
3. Create a feature branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `make test`
6. Commit and push
7. Open a Pull Request

## Development Setup

```bash
# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Run linter
make lint
```

## Code Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Add comments for exported functions
- Keep functions small and focused

### Security Requirements

**IMPORTANT**: All contributions must follow security best practices:

1. **Never hardcode secrets** - No API keys, passwords, or tokens in code
2. **Use parameterized queries** - Never interpolate user input into SQL
3. **Validate input** - Check size limits and sanitize user data
4. **Handle errors properly** - Don't leak sensitive info in error messages
5. **Add tests** - Include security-focused test cases

### Commit Messages

Use clear, descriptive commit messages:

```
Add feature X for Y purpose

- Detailed point 1
- Detailed point 2
```

### Pull Request Guidelines

1. **One feature per PR** - Keep changes focused
2. **Include tests** - All new code should have test coverage
3. **Update documentation** - Keep README and docs current
4. **No breaking changes** - Maintain backwards compatibility
5. **Security review** - Flag any security-sensitive changes

## Testing

```bash
# Run all tests
CGO_ENABLED=1 go test -tags "fts5" ./...

# Run with verbose output
CGO_ENABLED=1 go test -tags "fts5" -v ./...

# Run specific package tests
CGO_ENABLED=1 go test -tags "fts5" -v ./internal/store/...
```

## Reporting Issues

### Bug Reports

Include:
- Claude PP version (`claude_pp version`)
- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior

### Security Issues

**DO NOT** open public issues for security vulnerabilities!

See [SECURITY.md](SECURITY.md) for responsible disclosure.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
