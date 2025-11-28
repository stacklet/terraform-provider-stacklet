# Architecture

## Core Structure

```
.
├── main.go                      # Provider entry point
├── internal/
│   ├── provider/                # Provider configuration and registration
│   ├── api/                     # GraphQL API client and wrappers
│   ├── resources/               # Terraform resource implementations
│   ├── datasources/             # Terraform data source implementations
│   ├── models/                  # Data model definitions
│   ├── acceptance_tests/        # Acceptance tests with HTTP recording
│   ├── errors/                  # Error handling utilities (legacy)
│   ├── modelupdate/             # Helper functions for model updates
│   ├── planmodifiers/           # Custom Terraform plan modifiers
│   ├── providerdata/            # Provider data management
│   ├── schemavalidate/          # Custom schema validators
│   ├── schemadefault/           # Custom schema default value handlers
│   └── typehelpers/             # Type conversion and manipulation helpers
├── examples/                    # Example Terraform configurations
├── docs/                        # Auto-generated documentation
├── templates/                   # Documentation templates
└── tools/                       # Build and documentation tools
```

## API Integration

- **GraphQL Client**: Uses `github.com/hasura/go-graphql-client` for GraphQL communication
- **API Wrapper**: `internal/api/api.go` provides a typed `API` struct with methods for each resource type
- **Modular API Files**: Each resource type has its own API file
- **HTTP Transport**: Custom transport layers for authentication (`authTransport`) and logging (`logTransport`)
- **Enums**: Strongly typed enums in `internal/api/enums.go` (CloudProvider, ReportSource, etc.)
- **Pagination**: Generic helpers in `internal/api/pagination.go` for GraphQL connection pattern pagination
  - `FindInPaginatedQuery[T, R]()` - Searches through paginated results for a specific item
  - `CollectAllPages[T]()` - Collects all items from all pages
- **Filtering**: Filter API in `internal/api/filter.go` for constructing GraphQL filter queries
  - `FilterElementInput` and `FilterValueInput` types for building filter expressions
  - `newExactMatchFilter()` helper for creating exact-match filters with "equals" operator

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

3. **Resource Registration**: All resources are registered in `internal/resources/resources.go` in the `RESOURCES` slice

4. **Model Separation**: Models are defined in `internal/models/` separate from resource logic

5. **Import Support**: Centralized import utilities in `internal/resources/import.go`:
   - `splitImportID()` - Helper for parsing composite import IDs
   - Import IDs use colon-separated format (e.g., `aws:123456789012` for accounts)

6. **Plan Modifiers**: Custom plan modifiers in `internal/planmodifiers/`:
   - `RequiresReplaceIfFieldsChanged()` - Forces recreation when specific nested object fields change
   - `RequiresReplaceIfUnset()` - Forces recreation when object transitions from non-null to null
   - `RequiresReplaceIfNullStringChange()` - Forces recreation when string toggles between null and non-null
   - `DefaultObject()` - Sets default object values to match API defaults

## Data Source Architecture

Data sources follow similar patterns to resources:

1. **Data Source Definition**: Each implements:
   - `datasource.DataSource` - Basic data source interface
   - `datasource.DataSourceWithConfigure` - For provider configuration

2. **Data Source Registration**: All registered in `internal/datasources/datasources.go` in the `DATASOURCES` slice

3. **Shared Models**: Often share model definitions with their corresponding resources

## Model Patterns

Models in `internal/models/` define Terraform schema structure with `tfsdk` struct tags. Key conventions:
- Separate types for resources and data sources when needed
- Write-only fields use `_wo` suffix with `_wo_version` for change detection
- All `Update()` methods return `diag.Diagnostics` for consistent error handling
- Use typehelpers for complex conversions (JSON, objects, lists)

## Type Helpers and Validators

**Type Helpers** (`internal/typehelpers/`): All return `(result, diag.Diagnostics)` for proper error propagation:
- `JSONString()` - Normalizes JSON strings for consistent comparison
- `ObjectValue()` / `ObjectList()` - Converts values/slices to Terraform objects/lists
- `FilteredObject()` / `UpdatedObject()` - Extracts/updates object attributes
- `ObjectStringIdentifier()` / `ListItemsIdentifiers()` - Extracts identifiers
- `ListSortedEntries()` - Sorts list entries by identifier order

**Usage Pattern**: Call helper, append diagnostics with `diags.Append(d...)`, check `diags.HasError()` for early returns when needed.

**Schema Validators** (`internal/schemavalidate/`):
- `OneOfCloudProviders()` - Validates cloud provider values
- `unique_string_attribute.go` - Ensures string uniqueness in lists

**Schema Defaults** (`internal/schemadefault/`):
- `EmptyListDefault()` / `EmptyMapDefault()` - Provide empty defaults for optional attributes

## Helper Packages

**Error Handling** (`internal/errors/`):
- **Legacy package** - Being phased out in favor of direct `diag.Diagnostics` usage
- New code should use `diag.Diagnostics` directly and the `typehelpers` package functions which all return diagnostics

**Model Update** (`internal/modelupdate/`):
- Helper functions for updating Terraform models from API results
- Centralizes common update patterns

**Provider Data** (`internal/providerdata/`):
- `ProviderData` struct holds API client for resources/data sources
- `GetResourceProviderData()` - Retrieves provider data in resources
- `GetDataSourceProviderData()` - Retrieves provider data in data sources

## Architectural Patterns

### Diagnostics Pattern

All functions that can fail return `diag.Diagnostics` (plural). Pattern: initialize with `var diags diag.Diagnostics`, accumulate with `diags.Append(d...)`, check `diags.HasError()` for early returns, add errors with `diags.AddError("Summary", "Detail")`.

### Type Conversion Pattern

Use `internal/typehelpers` for all type conversions: `JSONString()`, `ObjectValue()`, `ObjectList()`, `ObjectStringIdentifier()`, `ListItemsIdentifiers()`, `ListSortedEntries()`. All return `(result, diag.Diagnostics)`.

### Model Update Pattern

Models implement `Update(ctx, obj) diag.Diagnostics`. Resource models extend data source models by calling parent `Update()` then applying adjustments. For order preservation, save identifiers before update, restore with `ListSortedEntries()`.

## Code Conventions

### File Organization
- One resource/data source per file
- File names match resource type (e.g., `account.go` for `stacklet_account`)
- Models in separate `models/` package
- API methods in separate `api/` package

### Naming Conventions
- Resource constructors: `NewXxxResource()`
- Data source constructors: `NewXxxDataSource()`
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
  - `DefaultObject()` - Sets default objects matching API defaults

### Write-Only Fields
- Use `_wo` suffix for write-only fields (e.g., `security_context_wo`)
- Use `_wo_version` suffix for version tracking
- Mark as `Sensitive: true` and `WriteOnly: true` in schema

### Collection Initialization

Use `make()` for empty collections. Composite literals only when providing immediate values or as direct function arguments.
