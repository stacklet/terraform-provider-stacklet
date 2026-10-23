# Architecture

## Core Structure

```
.
├── main.go                      # Provider entry point
├── internal/
│   ├── provider/                # Provider configuration and registration
│   ├── api/                     # Shared HTTP transports and errors for API clients
│   │   ├── iac/                 # IaC (Sinistral) REST API client
│   │   └── platform/            # Platform GraphQL API client and wrappers
│   ├── resources/               # Terraform resource implementations
│   ├── datasources/             # Terraform data source implementations
│   ├── models/                  # Data model definitions
│   ├── acceptance_tests/        # Acceptance tests with HTTP recording
│   ├── errors/                  # Error handling utilities (legacy)
│   ├── planmodifiers/           # Custom Terraform plan modifiers
│   ├── providerdata/            # Provider data management
│   ├── schemavalidate/          # Custom schema validators
│   ├── schemadefault/           # Custom schema default value handlers
│   └── typehelpers/             # Type conversion and manipulation helpers
├── examples/                    # Example Terraform configurations
├── docs/                        # Auto-generated documentation
└── tools/                       # Build and documentation tools
```

## API Integration

- **Shared API code** (`internal/api/`): Functionality common to all API clients
  - HTTP transport layers for authentication (`AuthTransport`) and logging (`LogTransport`), combined by `NewTransport()`
  - Error types: `APIError`, `NotFound` and `Conflict` (resources remove objects from state on `api.NotFound`)
  - `MaskSecrets()` returns a context which masks the given secrets in request/response logs. Call it before API requests that send secrets
- **GraphQL Client**: The platform client (`internal/api/platform/`) uses `github.com/hasura/go-graphql-client` for GraphQL communication
- **API Wrapper**: `internal/api/platform/api.go` provides a typed `API` struct with methods for each resource type
- **Modular API Files**: Each resource type has its own API file
- **Enums**: Strongly typed enums in `internal/api/platform/enums.go`.
- **Filtering**: Filter API in `internal/api/platform/filter.go` for constructing GraphQL filter queries
  - `filterElementInput` and `filterValueInput` types for building filter expressions
  - `newExactMatchFilter()` helper for creating exact-match filters with "equals" operator
- **REST Client**: The IaC client (`internal/api/iac/`) talks to the IaC (Sinistral) REST API, using the same API key as the platform one
  - `Get()`, `Post()`, `Put()`, `Patch()`, `Delete()` methods for JSON requests
  - `resourcePath()` builds request paths, escaping each component
  - Error responses are converted to `api.NotFound` (404) or `api.APIError`, reporting FastAPI `detail` messages (including validation errors)

## Resource Architecture

Resources follow the Terraform Plugin Framework patterns:

1. **Resource Definition**: Each resource implements:
   - `resource.Resource` - Basic resource interface
   - `resource.ResourceWithConfigure` - For provider configuration
   - `resource.ResourceWithImportState` - For import support

2. **CRUD Operations**: Resources implement:
   - `Create()` - Create new resource
   - `Read()` - Read existing resource
   - `Update()` - Update existing resource
   - `Delete()` - Delete resource
   - `ImportState()` - Import existing resource

3. **Resource Registration**: All resources are registered in the `Resources` variable in `internal/resources/resources.go`, via `newFactory()`. Resources in the `Unreleased` list are only enabled when `STACKLET_UNRELEASED_FEATURES` is set

4. **Model Separation**: Models are defined in `internal/models/` separate from resource logic

5. **Import Support**: Centralized import utilities in `internal/resources/import.go`:
   - `splitImportID()` - Helper for parsing composite import IDs
   - Import IDs use colon-separated format (e.g., `aws:123456789012` for accounts)

6. **Plan Modifiers**: Custom plan modifiers in `internal/planmodifiers/`:
   - `RequiresReplaceIfFieldsChanged()` - Forces recreation when specific nested object fields change
   - `RequiresReplaceIfUnset()` - Forces recreation when object transitions from non-null to null
   - `RequiresReplaceIfNullStringChange()` - Forces recreation when string toggles between null and non-null
   - `TrimWhiteSpace()` - Strips leading/trailing whitespace from planned strings, for APIs that normalize it
   - `DefaultObject()` - Sets default object values to match API defaults

## Data Source Architecture

Data sources follow similar patterns to resources:

1. **Data Source Definition**: Each implements:
   - `datasource.DataSource` - Basic data source interface
   - `datasource.DataSourceWithConfigure` - For provider configuration

2. **Data Source Registration**: All registered in the `DataSources` variable in `internal/datasources/datasources.go`, with the same `Released`/`Unreleased` split as resources

3. **Shared Models**: Often share model definitions with their corresponding resources

## Model Patterns

Models in `internal/models/` define Terraform schema structure with `tfsdk` struct tags. Key conventions:
- Separate types for resources and data sources when needed
- Write-only fields use `_wo` suffix with `_wo_version` for change detection
- All `Update()` methods return `diag.Diagnostics` for consistent error handling
- Use typehelpers for complex conversions (JSON, objects, lists)

## Type Helpers and Validators

**Type Helpers** (`internal/typehelpers/`): Helpers that can fail return `(result, diag.Diagnostics)` for proper error propagation:
- `JSONString()` - Normalizes JSON strings for consistent comparison
- `ObjectValue()` / `ObjectList()` - Converts values/slices to Terraform objects/lists
- `FilteredObject()` / `UpdatedObject()` - Extracts/updates object attributes
- `ObjectStringIdentifier()` / `ListItemsIdentifiers()` - Extracts identifiers
- `ListSortedEntries()` - Sorts list entries by identifier order
- `StringsList()` / `GraphQLIDValue()` - Simple conversions (no diagnostics)
- `PreserveIfNull()` / `PreserveIfEmptyJSON()` - Keep the original value when the API returns null or `{}` for it, avoiding perpetual diffs

**Usage Pattern**: Call helper, append diagnostics with `diags.Append(d...)`, check `diags.HasError()` for early returns when needed.

**Schema Validators** (`internal/schemavalidate/`):
- `OneOfCloudProviders()` - Validates cloud provider values
- `UniqueStringAttribute()` - Ensures string uniqueness in lists
- `UUID()` - Validates UUID strings

**Schema Defaults** (`internal/schemadefault/`):
- `EmptyListDefault()` / `EmptyMapDefault()` - Provide empty defaults for optional attributes

## Helper Packages

**Error Handling** (`internal/errors/`):
- **Legacy package** - Being phased out in favor of direct `diag.Diagnostics` usage
- New code should use `diag.Diagnostics` directly and the `typehelpers` package functions which all return diagnostics

**Provider Data** (`internal/providerdata/`):
- `providerData` struct holds API clients for resources/data sources, built by `New()`
- `GetForResource()` - Retrieves provider data in resources
- `GetForDataSource()` - Retrieves provider data in data sources

## Architectural Patterns

### Diagnostics Pattern

All functions that can fail return `diag.Diagnostics` (plural). Pattern: initialize with `var diags diag.Diagnostics`, accumulate with `diags.Append(d...)`, check `diags.HasError()` for early returns, add errors with `diags.AddError("Summary", "Detail")`.

### Type Conversion Pattern

Use `internal/typehelpers` for type conversions: `JSONString()`, `ObjectValue()`, `ObjectList()`, `ObjectStringIdentifier()`, `ListItemsIdentifiers()`, `ListSortedEntries()`. These return `(result, diag.Diagnostics)`.

### Model Update Pattern

Models implement `Update(ctx, obj) diag.Diagnostics`. Resource models extend data source models by calling parent `Update()` then applying adjustments. For order preservation, save identifiers before update, restore with `ListSortedEntries()`.

## Code Conventions

### File Organization
- One resource/data source per file
- File names match resource type (e.g., `account.go` for `stacklet_account`)
- Models in separate `models/` package
- API methods in the `api/platform/` package

### Naming Conventions
- Resource types: `xxxResource`, registered with `newFactory(&xxxResource{})`
- Data source types: `xxxDataSource`, registered with `newFactory(&xxxDataSource{})`
- API types: Input types use `XxxInput` suffix
- Models: Separate types for resources vs data sources when needed

### Error Handling
- Always return `diag.Diagnostics` (plural) from helper functions and model methods
- Use `diags.AddError("Summary", "Detail")` and `diags.Append(d...)` to accumulate diagnostics
- Check `diags.HasError()` for early returns when needed
- Avoid legacy `internal/errors` package for new code

### Plan Modifiers
- Use `stringplanmodifier.RequiresReplace()` for immutable fields
- Use `stringplanmodifier.UseStateForUnknown()` for computed fields
- Custom modifiers in `internal/planmodifiers/`:
  - `RequiresReplaceIfNullStringChange()` - For string fields that toggle between null and non-null
  - `RequiresReplaceIfFieldsChanged()` - For objects where specific nested fields require replacement
  - `RequiresReplaceIfUnset()` - For objects that cannot transition from configured to null
  - `TrimWhiteSpace()` - For strings where surrounding whitespace is not significant
  - `DefaultObject()` - Sets default objects matching API defaults

### Write-Only Fields
- Use `_wo` suffix for write-only fields (e.g., `security_context_wo`)
- Use `_wo_version` suffix for version tracking
- Mark as `Sensitive: true` and `WriteOnly: true` in schema

### Collection Initialization

Use `make()` for empty collections. Composite literals only when providing immediate values or as direct function arguments.
