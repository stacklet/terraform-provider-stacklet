# Development

## Setup

### Local Development

Use `~/.terraformrc` with `dev_overrides` pointing to build directory:

```hcl
provider_installation {
  dev_overrides {
    "stacklet/stacklet" = "/absolute/path/to/terraform-provider-stacklet"
  }
  direct {}
}
```

### Debug Mode

Run `./terraform-provider-stacklet -debug` and use `TF_REATTACH_PROVIDERS` environment variable from output.

## Commands

All development commands use the `just` task runner:

- `just build` - Build the provider binary
- `just lint` - Format and lint code (Go, Terraform, docs, copyright)
- `just test` - Run acceptance tests (replay mode, no API needed)
- `just test-unit` - Run unit tests only
- `just test-live` - Run tests against live API
- `just test-record <testname>` - Record new API responses for a test
- `just docs` - Generate provider documentation
- `just update-deps-go` - Update Go dependencies

Testing individual tests: `just test -run TestAccAccountResource`

## Testing Strategy

**HTTP Recording**: Tests use recorded HTTP interactions (`internal/acceptance_tests/recordings/`) for fast, deterministic runs without live API. Three modes via `TF_ACC_MODE`: `replay` (default), `record`, `live`. See `internal/acceptance_tests/recorded_transport.go`.

**Test Helpers**:
- `importStateIDFuncFromAttrs()` - Creates import IDs from resource attributes
- `prefixName()` - Adds test prefixes via `TF_ACC_PREFIX`
- `ConfigPlanChecks` with `plancheck.ExpectResourceAction()` - Validates plan actions
- Template rendering for test configurations

**Unit Tests**: Table-driven tests with testify assertions (e.g., `internal/provider/provider_test.go`). Run with: `just test-unit`

## Authentication

Provider supports multiple authentication methods (in priority order):
1. Direct configuration (endpoint/api_key in provider block)
2. Environment variables (`STACKLET_ENDPOINT`, `STACKLET_API_KEY`)
3. stacklet-admin CLI (`~/.stacklet/config.json`, `~/.stacklet/credentials`)

Each credential is resolved independently using the first non-empty value found.

The IaC API endpoint (`iac_endpoint`) is resolved from the provider block, then the `STACKLET_IAC_ENDPOINT` environment variable. If neither is set, it's derived from the API endpoint by replacing the `api.` hostname prefix with `api.sinistral.` (e.g. `https://api.example.stacklet.io/` becomes `https://api.sinistral.example.stacklet.io/`). If it can't be derived, IaC resources and data sources fail when used. The API key is shared by both APIs. Unit tests in `internal/provider/provider_test.go` verify precedence and mixed sources.

## Development Tools

Required tools are defined in `.tool-versions`.
Additional tools are defined in `tools/go.mod`.

## Documentation

- **Auto-generated**: Documentation in `docs/` is generated from code
- **Examples**: All examples in `examples/` must be valid Terraform
- **Import examples**: Each resource has `import.sh` example

## CI/CD

### Continuous Integration (`.github/workflows/ci.yaml`)
- Runs on push to `main` and PRs targeting `main`
- All GitHub Actions and reusable workflows use SHA pinning for security
- `Lint` job runs `just lint` (all pre-commit hooks: Go linting, Terraform formatting and examples validation, documentation validation, copyright headers, spelling)
- `Tests` job runs `just test` (acceptance tests in replay mode)

### Release Workflow (`.github/workflows/release.yaml`)
- **Public Registry**: Triggered by version tags (`v*.*.*`), builds with GoReleaser, signs with GPG, publishes to Terraform Registry
- **Private Registry** (`release-main` job): Currently disabled, would publish to Terraform Cloud private registry on `main` commits

## Release Process

Update CHANGELOG.md, create tag with `just tag-release X.Y.Z`, push tag. GitHub Actions builds, signs with GPG, and publishes to Terraform Registry.
