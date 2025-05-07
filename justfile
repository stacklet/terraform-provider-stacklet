# Build the provider
build:
    go build -o terraform-provider-stacklet

# Format code
format:
    terraform fmt -recursive
    go fmt ./...

# Run linters
lint:
    terraform fmt -recursive -check
    golangci-lint run --fix

# Run tests
test *args:
    TF_ACC=1 STACKLET_ENDPOINT=fake STACKLET_API_KEY=fake go test ./internal/... {{ args }}

# Generate provider docs
docs:
    go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

