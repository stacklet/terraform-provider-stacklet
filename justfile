# Format code
fmt:
    go fmt ./...

# Run linters
lint:
    golangci-lint run

# Run tests
test:
    go test ./internal/...

# Generate provider docs
docs:
    go generate
