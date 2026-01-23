package := "./internal/..."

# Build the provider
build *args:
    go build -o terraform-provider-stacklet {{ args }}

# Format code
format: format-go format-tf

# Format terraform code
format-tf:
    terraform fmt -recursive

# Format golang code
format-go:
    go fmt ./...    

# Run linters
lint: lint-go lint-tf lint-docs lint-copyright

# Run linters for terraform
lint-tf:
    terraform fmt -recursive -check
    just validate-tf

# Run linters for golang
lint-go:
    go vet {{ package }}
    golangci-lint run --fix

# Run linter for generated docs
lint-docs:
    env -C tools go generate -run=validate-docs

# Run linter for copyright headers
lint-copyright:
    env -C tools go generate -run=validate-copyright

# Run tests using recorded API requests/responses.
test *args:
    TF_ACC=1 go test {{ package }} {{ args }}

# Run tests against a live deployment. Requires real STACKLET_ENDPOINT and STACKLET_API_KEY or logged in stacklet-admin.
test-live *args:
    TF_ACC_MODE=live just test -count=1 {{ args }}
    
# Record API request/responses for an acceptance test. Requires real STACKLET_ENDPOINT and STACKLET_API_KEY or logged in stacklet-admin.
test-record testname:
    TF_ACC_MODE=record just test -count=1 -run {{ testname }}

# Run only unit tests (no acceptance)
test-unit *args:
    go test {{ package }} {{ args }}

# Generate provider documentation
docs:
    env -C tools go generate -run=generate-docs

# Add copyright headers
copyright:
    env -C tools go generate -run=generate-copyright

# Update golang dependencies
update-deps-go:
    go get -u ./...
    go mod tidy
    env -C tools go get -u
    env -C tools go mod tidy

tf_config := '''
provider_installation {
  dev_overrides {
    "stacklet/stacklet" = "$PWD"
  }
  direct {}
}
'''
tf_provider_config := '''
terraform {
  required_providers {
    stacklet = {
      source = "stacklet/stacklet"
    }
  }
}
'''

# validate terraform example files
validate-tf: build
    #!/usr/bin/env bash
    set -e

    validate() {
      local dir="$1"

      echo "Validating $dir"

      local module="$(mktemp -d)"
      trap "rm -rf $module" EXIT

      local terraformrc="$module/.terraformrc"
      cat > "$terraformrc" <<EOF
    {{ tf_config }}
    EOF
      cat > "$module/provider.tf" <<EOF
    {{ tf_provider_config }}
    EOF

      cp -a "$dir"/* "$module"
      TF_CLI_CONFIG_FILE="$terraformrc" terraform -chdir="$module" validate
      rm -rf $module
    }

    for dir in examples/provider examples/data-sources/* examples/resources/*; do
      validate $dir
    done

# Create release tag for the specified version (without leading 'v')
tag-release version ref="HEAD":
    git tag -a v{{ version }} -m 'Version {{ version }}' {{ ref }}
