set script-interpreter := ['bash', '-euo', 'pipefail']

package := "./internal/..."

# Build the provider
build *args:
    go build -o terraform-provider-stacklet {{ args }}

# Run linting
lint *args:
    prek run --all-files {{ args }}

# Run tests using recorded API requests/responses.
test *args:
    TF_ACC=1 STACKLET_UNRELEASED_FEATURES=1 go test {{ package }} {{ args }}

# Run tests against a live deployment. Requires real STACKLET_ENDPOINT and STACKLET_API_KEY or logged in stacklet-admin.
test-live *args:
    TF_ACC_MODE=live just test -count=1 {{ args }}

# Record API request/responses for an acceptance test. Requires real STACKLET_ENDPOINT and STACKLET_API_KEY or logged in stacklet-admin.
test-record testname:
    TF_ACC_MODE=record just test -count=1 -run {{ testname }}

# Check that all acceptance tests have recordings and there are no leftover ones
[script]
test-recordings-check:
    dir=internal/acceptance_tests
    used=$(go test "./$dir" -list '^TestAcc' | grep '^TestAcc' | sort)
    existing=$(find "$dir/recordings" -name '*.json' -exec basename {} .json \; | sort)

    missing=$(comm -23 <(echo "$used") <(echo "$existing"))
    leftover=$(comm -13 <(echo "$used") <(echo "$existing"))

    if [[ -n "$missing" ]]; then
      echo "Missing recordings:"
      sed 's/^/  /' <<< "$missing"
    fi
    if [[ -n "$leftover" ]]; then
      echo "Leftover recordings:"
      sed 's/^/  /' <<< "$leftover"
    fi
    [[ -z "$missing" && -z "$leftover" ]]

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
[script]
tf-examples-validate: build
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
[script]
tf-examples-check: build
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
[script]
_tf *args:
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

[script]
_declare_provider module_path:
    cat > "{{module_path}}/provider.tf" <<EOF
    {{ tf_provider_config }}
    EOF
