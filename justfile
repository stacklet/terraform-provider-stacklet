# Build the provider
build:
    go build -o terraform-provider-stacklet

# Format code
fmt:
    go fmt ./...

# Run linters
lint:
    golangci-lint run --fix

# Run tests
test:
    TF_ACC=1 STACKLET_ENDPOINT=fake STACKLET_API_KEY=fake go test ./internal/...

# Generate provider docs
docs:
    go generate
