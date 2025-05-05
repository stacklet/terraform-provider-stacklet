# Build the provider
build:
    go build -o terraform-provider-stacklet

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
