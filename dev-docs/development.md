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
- `just format` / `just lint` - Format and lint code (Go, Terraform, docs, copyright)
- `just test` - Run acceptance tests (replay mode, no API needed)
- `just test-live` - Run tests against live API
- `just test-record <testname>` - Record new API responses for a test
- `just docs` - Generate provider documentation
- `just update-deps-go` - Update Go dependencies

Testing individual tests: `just test -run TestAccAccountResource`

## Testing Strategy

**HTTP Recording**: Tests use recorded HTTP interactions (`acceptance_tests/recordings/`) for fast, deterministic runs without live API. Three modes via `TF_ACC_MODE`: `replay` (default), `record`, `live`. See `acceptance_tests/recorded_transport.go`.

**Test Helpers**:
- `importStateIDFuncFromAttrs()` - Creates import IDs from resource attributes
- `prefixName()` - Adds test prefixes via `TF_ACC_PREFIX`
- `ConfigPlanChecks` with `plancheck.ExpectResourceAction()` - Validates plan actions
- Template rendering for test configurations

**Unit Tests**: Use testify test suites for shared setup (e.g., `internal/provider/provider_test.go`). Run with: `go test ./internal/provider/`

## Authentication

Provider supports multiple authentication methods (in priority order):
1. Direct configuration (endpoint/api_key in provider block)
2. Environment variables (`STACKLET_ENDPOINT`, `STACKLET_API_KEY`)
3. stacklet-admin CLI (`~/.stacklet/config.json`, `~/.stacklet/credentials`)

Each credential is resolved independently using the first non-empty value found. Unit tests in `internal/provider/provider_test.go` verify precedence and mixed sources.

## Workflow

1. **Build**: `just build`
2. **Test**: `just test` (replay), `just test-live` (live API), `just test-record <name>` (record new)
3. **Format/Lint**: `just format` and `just lint` before committing
4. **Documentation**: `just docs` to regenerate provider docs

## Development Tools

Required tools (defined in `.tool-versions`):
- Go 1.24
- golangci-lint v2.1.5
- just 1.40.0

Additional tools (in `go.mod`):
- `github.com/hashicorp/terraform-plugin-docs` - Documentation generation
- `github.com/hashicorp/copywrite` - Copyright header management

## Documentation

- **Auto-generated**: Documentation in `docs/` is generated from code
- **Templates**: Use `templates/` for customizing documentation
- **Examples**: All examples in `examples/` must be valid Terraform
- **Import examples**: Each resource has `import.sh` example

## CI/CD

### Continuous Integration (`.github/workflows/ci.yaml`)
- Runs on push to `main` and all PRs
- All GitHub Actions use SHA pinning for security
- Separate jobs for: Go linting, Terraform linting, documentation validation, copyright headers, and acceptance tests (replay mode)

### Release Workflow (`.github/workflows/release.yaml`)
- **Public Registry**: Triggered by version tags (`v*.*.*`), builds with GoReleaser, signs with GPG, publishes to Terraform Registry
- **Private Registry** (`release-main` job): Currently disabled, would publish to Terraform Cloud private registry on `main` commits

## Release Process

Update CHANGELOG.md, create tag with `just tag-release X.Y.Z`, push tag. GitHub Actions builds, signs with GPG, and publishes to Terraform Registry.
