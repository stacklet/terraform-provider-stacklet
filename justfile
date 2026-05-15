package := "./internal/..."

# Build the provider
build *args:
    go build -o terraform-provider-stacklet {{ args }}

# Run linting
lint *args:
    uvx prek run --all-files {{ args }}

# Run tests using recorded API requests/responses.
test *args:
    TF_ACC=1 STACKLET_UNRELEASED_FEATURES=1 go test {{ package }} {{ args }}

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

# Update golang dependencies
update-deps-go:
    go get -u ./...
    go mod tidy
    env -C tools go get -u
    env -C tools go mod tidy

# validate terraform example files
tf-examples-validate: build
    #!/usr/bin/env bash
    set -e

    validate() {
      local dir="$1"

      echo "Validating $dir"

      local module="$(mktemp -d)"
      trap "rm -rf $module" EXIT

      just _declare_provider "$module"
      cp -a "$dir"/* "$module"
      just _tf -chdir="$module" validate

      rm -rf "$module"
    }

    for dir in examples/provider examples/data-sources/* examples/resources/*; do
      validate "$dir"
    done

# check that all resources and datasources have the required example files
tf-examples-check: build
    #!/usr/bin/env bash
    set -e

    module="$(mktemp -d)"
    trap "rm -rf $module" EXIT
    just _declare_provider "$module"

    schema=$(just _tf -chdir="$module" providers schema -json)
    errors=0

    check() {
      if [[ ! -f "$1" ]]; then
        echo "MISSING: $1"
        errors=$((errors + 1))
      fi
    }

    while IFS= read -r name; do
      check "examples/resources/$name/resource.tf"
      check "examples/resources/$name/import-by-string-id.tf"
      check "examples/resources/$name/import.sh"
    done < <(jq -r '.provider_schemas["registry.terraform.io/stacklet/stacklet"].resource_schemas | keys[]' <<< "$schema")

    while IFS= read -r name; do
      check "examples/data-sources/$name/data-source.tf"
    done < <(jq -r '.provider_schemas["registry.terraform.io/stacklet/stacklet"].data_source_schemas | keys[]' <<< "$schema")

    if [[ $errors -gt 0 ]]; then
      echo "$errors missing file(s)"
      exit 1
    fi

# Create release tag for the specified version (without leading 'v')
tag-release version ref="HEAD":
    git tag -a v{{ version }} -m 'Version {{ version }}' {{ ref }}


tf_config := '''
provider_installation {
  dev_overrides {
    "stacklet/stacklet" = "$PWD"
  }
  direct {}
}
'''

# run terraform with the local provider configured. Ensure calling targets
# depend on "build" to use the current code
_tf *args:
    #!/usr/bin/env bash
    set -e

    terraformrc="$(mktemp)"
    trap "rm -rf $terraformrc" EXIT
    cat > "$terraformrc" <<EOF
    {{ tf_config }}
    EOF
    TF_CLI_CONFIG_FILE="$terraformrc" terraform {{ args }}
    rm -rf "$terraformrc"

tf_provider_config := '''
terraform {
  required_providers {
    stacklet = {
      source = "stacklet/stacklet"
    }
  }
}
'''

_declare_provider module_path:
    #!/usr/bin/env bash
    set -e

    cat > "{{module_path}}/provider.tf" <<EOF
    {{ tf_provider_config }}
    EOF
