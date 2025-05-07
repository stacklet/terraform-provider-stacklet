package := "./internal/..."

# Build the provider
build:
    go mod download
    go build -o terraform-provider-stacklet

# Format code
format:
    terraform fmt -recursive
    go fmt ./...

# Run linters
lint:
    terraform fmt -recursive -check
    go vet {{ package }}
    golangci-lint run --fix

# Run tests
test *args:
    TF_ACC=1 go test {{ package }} {{ args }}


# Record API request/responses for an acceptance test.
# Requires real STACKLET_ENDPOINT and STACKLET_API_KEY or logged in stacklet-admin.
test-record testname:
    TF_ACC_RECORD=1 just test -run {{ testname }}

# Generate provider docs
docs:
    go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

