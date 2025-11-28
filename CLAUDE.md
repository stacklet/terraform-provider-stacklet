# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Terraform provider for Stacklet's cloud governance platform. It allows managing cloud infrastructure and governance resources like accounts, policy collections, bindings, and configuration profiles through Infrastructure as Code.

Built with Terraform Plugin Framework the provider uses GraphQL to communicate with the Stacklet API.

## Documentation Structure

- **[dev-docs/architecture.md](dev-docs/architecture.md)** - Code structure, patterns, helpers, and conventions
- **[dev-docs/development.md](dev-docs/development.md)** - Setup, commands, workflow, testing, and tools
- **[dev-docs/contributing.md](dev-docs/contributing.md)** - How to add resources/data sources and record tests

## Quick Reference

### Development Commands

All commands use the `just` task runner:
- `just build` - Build the provider binary
- `just format` / `just lint` - Format and lint code
- `just test` - Run acceptance tests (replay mode, no API needed)
- `just test-live` - Run tests against live API
- `just test-record <testname>` - Record new API responses
- `just docs` - Generate provider documentation

### Core Patterns

**Diagnostics**: All functions return `diag.Diagnostics`. Pattern: `var diags diag.Diagnostics`, `diags.Append(d...)`, `diags.HasError()`, `diags.AddError()`.

**Type Helpers**: Use `internal/typehelpers` for conversions (all return `(result, diag.Diagnostics)`). Call helper, append diagnostics, check for errors.

**Models**: Implement `Update(ctx, obj) diag.Diagnostics`. Resource models extend data source models.

### Adding Resources

1. Create model (`internal/models/`)
2. Create API methods (`internal/api/`)
3. Create resource (`internal/resources/`), register in `resources.go`
4. Add examples and tests
5. Run `just docs`

Data sources are similar, register in `internal/datasources/datasources.go`, often reuse models/API methods.

## Resources and Data Sources

See `internal/resources/resources.go` and `internal/datasources/datasources.go` for complete lists. The provider supports account management, policy management, configuration profiles (multiple integration types), and repository management.
