# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Terraform provider for Stacklet's cloud governance platform. It allows managing cloud infrastructure and governance resources like accounts, policy collections, bindings, and configuration profiles through Infrastructure as Code.

Built with Terraform Plugin Framework and Go 1.24, the provider uses GraphQL to communicate with the Stacklet API.

## Development Commands

All development commands use the `just` task runner. Key commands:

- `just build` - Build the provider binary
- `just format` - Format all code (Go and Terraform)
- `just format-go` - Format Go code only
- `just format-tf` - Format Terraform code only
- `just lint` - Run all linters (Go, Terraform, docs, copyright)
- `just lint-go` - Run Go linters only
- `just lint-tf` - Run Terraform linters only
- `just lint-docs` - Validate generated documentation
- `just lint-copyright` - Validate copyright headers
- `just test` - Run acceptance tests using recorded API responses
- `just test-live` - Run tests against live API (requires authentication)
- `just test-record <testname>` - Record new API responses for a specific test
- `just docs` - Generate provider documentation
- `just validate-tf` - Validate Terraform example files
- `just copyright` - Add copyright headers to files
- `just update-deps-go` - Update Go dependencies
- `just tag-release X.Y.Z` - Create release tag

For testing individual acceptance tests:
```bash
just test -run TestAccAccountResource
```

For testing with custom prefix (useful for live tests):
```bash
TF_ACC_PREFIX=mytest just test-live -run TestAccAccountResource
```

## Architecture

### Core Structure

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

### API Integration

- **GraphQL Client**: Uses `github.com/hasura/go-graphql-client` for GraphQL communication
- **API Wrapper**: `internal/api/api.go` provides a typed `API` struct with methods for each resource type
- **Modular API Files**: Each resource type has its own API file:
  - `account.go` - Account management
  - `account_discovery.go` - Account discovery for AWS/Azure/GCP
  - `account_group.go` - Account group management
  - `binding.go` - Policy binding management
  - `configuration_profile.go` - Configuration profiles
  - `policy.go`, `policy_collection.go` - Policy management
  - `platform.go` - Platform information
  - `repository.go` - Repository management
  - `notification_template.go` - Notification template management
  - `report_group.go` - Report group management
- **HTTP Transport**: Custom transport layers for authentication (`authTransport`) and logging (`logTransport`)
- **Enums**: Strongly typed enums in `internal/api/enums.go` (CloudProvider, ReportSource, etc.)
- **Pagination**: Generic helpers in `internal/api/pagination.go` for GraphQL connection pattern pagination
- **Filtering**: Filter API in `internal/api/filter.go` for constructing GraphQL filter queries
  - `FilterElementInput` and `FilterValueInput` types for building filter expressions
  - `newExactMatchFilter()` helper for creating exact-match filters with "equals" operator

#### Pagination Helpers

The API package provides generic helpers for handling GraphQL connections with pagination (`internal/api/pagination.go`):

**`FindInPaginatedQuery[T, R]()`** - Searches through paginated results for a specific item:
```go
uuid, err := FindInPaginatedQuery(
    ctx,
    // Fetcher: executes query and returns items, hasNextPage, endCursor
    func(ctx context.Context, cursor string) ([]NodeType, bool, string, error) {
        var q struct {
            Conn struct {
                Edges []struct { Node NodeType }
                PageInfo struct {
                    HasNextPage bool
                    EndCursor   string
                }
                Problems []Problem
            } `graphql:"someConnection(first: 100, after: $cursor)"`
        }
        if err := client.Query(ctx, &q, map[string]any{"cursor": graphql.String(cursor)}); err != nil {
            return nil, false, "", NewAPIError(err)
        }
        if err := FromProblems(ctx, q.Conn.Problems); err != nil {
            return nil, false, "", err
        }
        nodes := make([]NodeType, 0, len(q.Conn.Edges))
        for _, edge := range q.Conn.Edges {
            nodes = append(nodes, edge.Node)
        }
        return nodes, q.Conn.PageInfo.HasNextPage, q.Conn.PageInfo.EndCursor, nil
    },
    // Matcher: checks each item and returns result when found
    func(node NodeType) (string, bool, error) {
        if node.Field == targetValue {
            return node.UUID, true, nil
        }
        return "", false, nil
    },
    NotFound{"Item not found"},
)
```

**`CollectAllPages[T]()`** - Collects all items from all pages:
```go
allItems, err := CollectAllPages(
    ctx,
    func(ctx context.Context, cursor string) ([]NodeType, bool, string, error) {
        // Same fetcher pattern as above
        return nodes, hasNextPage, endCursor, nil
    },
)
```

See `internal/api/repository.go:124` for a real-world example.

### Resource Architecture

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
   - `RequiresReplaceIfUnset()` - Forces recreation when an object field transitions from non-null to null
   - `RequiresReplaceIfNullStringChange()` - Forces recreation when a string attribute transitions between null and non-null (in either direction)
   - `DefaultObject()` - Sets a default object value when the planned value is null (useful for reflecting API defaults)

### Data Source Architecture

Data sources follow similar patterns to resources:

1. **Data Source Definition**: Each implements:
   - `datasource.DataSource` - Basic data source interface
   - `datasource.DataSourceWithConfigure` - For provider configuration

2. **Data Source Registration**: All registered in `internal/datasources/datasources.go` in the `DATASOURCES` slice

3. **Shared Models**: Often share model definitions with their corresponding resources

### Model Patterns

Models in `internal/models/` define the Terraform schema structure:

- Use `tfsdk` struct tags for attribute mapping
- Implement `AttributeTypes()` methods for nested objects
- Support both resource and data source variants (e.g., `ConfigurationProfileTeamsResource` and `ConfigurationProfileTeamsDataSource`)
- Write-only fields use `_wo` suffix (e.g., `security_context_wo`)
- Write-only version fields use `_wo_version` suffix for change detection
- All model `Update()` methods return `diag.Diagnostics` for consistent error handling

**Model Update Pattern**:
```go
func (m *ModelDataSource) Update(ctx context.Context, apiObj *api.Object) diag.Diagnostics {
    var diags diag.Diagnostics

    // Update fields from API object
    m.ID = types.StringValue(apiObj.ID)
    m.Name = types.StringValue(apiObj.Name)

    // Use typehelpers for complex conversions
    jsonValue, d := typehelpers.JSONString(apiObj.Variables)
    diags.Append(d...)
    if diags.HasError() {
        return diags
    }
    m.Variables = jsonValue

    return diags
}
```

Common model types:
- `Tag` - Key-value tags
- `TerraformModule` - Terraform module definitions
- Configuration profile models for each integration type

### Testing Strategy

The provider uses a sophisticated testing approach with HTTP recording:

1. **Test Modes** (controlled via `TF_ACC_MODE` environment variable):
   - `replay` (default) - Uses recorded HTTP interactions from `acceptance_tests/recordings/`
   - `record` - Records new HTTP interactions to files
   - `live` - Tests against live API without recording

2. **Recorded Transport**: `acceptance_tests/recorded_transport.go` implements HTTP recording/replay
   - Matches requests by GraphQL query and variables
   - Stores recordings as JSON files per test
   - Enables fast, deterministic tests without live API

3. **Test Setup**:
   - Tests use black-box testing approach (`package acceptance_test`)
   - `TestMain` in `setup_test.go` configures test environment
   - Sets fake credentials for replay mode
   - `runRecordedAccTest()` helper manages test execution

4. **Test Helpers**:
   - `importStateIDFuncFromAttrs()` - Creates import IDs from resource attributes
   - `prefixName()` - Adds configurable prefix to resource names (via `TF_ACC_PREFIX`)
   - Template rendering for test configurations
   - `ConfigPlanChecks` - Validates expected plan behaviors (resource actions, attribute changes)
   - `plancheck.ExpectResourceAction()` - Asserts specific resource actions (create, update, replace, etc.)

5. **Plan Validation**:
   Tests can validate expected Terraform plan behaviors using `ConfigPlanChecks`:
   ```go
   {
       Config: `...`,
       ConfigPlanChecks: resource.ConfigPlanChecks{
           PreApply: []plancheck.PlanCheck{
               plancheck.ExpectResourceAction("stacklet_account_group.test", plancheck.ResourceActionDestroyBeforeCreate),
           },
       },
   }
   ```
   Common resource actions:
   - `ResourceActionDestroyBeforeCreate` - Replace (delete before create)
   - `ResourceActionCreate` - Create new resource
   - `ResourceActionUpdate` - Update in place
   - `ResourceActionNoop` - No changes

6. **Unit Test Patterns**:
   For testing internal packages, use testify test suites for shared setup:
   ```go
   type GetCredentialsTestSuite struct {
       suite.Suite
       tmpDir      string
       stackletDir string
   }

   func TestGetCredentialsSuite(t *testing.T) {
       suite.Run(t, new(GetCredentialsTestSuite))
   }

   func (s *GetCredentialsTestSuite) SetupTest() {
       s.tmpDir = s.T().TempDir()
       s.T().Setenv("HOME", s.tmpDir)
       s.stackletDir = path.Join(s.tmpDir, ".stacklet")
   }

   func (s *GetCredentialsTestSuite) TestFromConfig() {
       // Test implementation
   }
   ```
   Benefits: Shared setup/teardown, cleaner test organization, and DRY principles.

7. **Running Tests**:
   ```bash
   # Replay mode (default) - fast, no API needed
   just test

   # Live mode - requires STACKLET_ENDPOINT and STACKLET_API_KEY
   just test-live

   # Record mode - update recordings for a test
   just test-record TestAccAccountResource

   # Run unit tests for internal packages
   go test ./internal/provider/
   ```

### Configuration Profiles

The provider supports various configuration profile types for integrations:

- **Account Owners** - Define account ownership
- **Email** - Email notifications
- **Jira** - Jira ticketing integration
- **Microsoft Teams** - MS Teams notifications
- **Resource Owner** - Resource ownership tracking
- **ServiceNow** - ServiceNow ticketing integration
- **Slack** - Slack notifications
- **Symphony** - Symphony notifications
- **Teams** - Generic Teams integration

Each configuration profile type has:
- Dedicated resource implementation (in `internal/resources/`)
- Dedicated data source implementation (in `internal/datasources/`)
- Model definition (in `internal/models/`)
- Common pattern with type-specific fields (webhooks, endpoints, credentials)

### Type Helpers and Validators

**Type Helpers** (`internal/typehelpers/`):

All helper functions follow a consistent pattern of returning `(result, diag.Diagnostics)` for proper error propagation:

- `JSONString(*string)` - Normalizes JSON strings for consistent comparison, returns `(types.String, diag.Diagnostics)`
- `ObjectValue[Type, Value]()` - Converts values to Terraform objects, returns `(types.Object, diag.Diagnostics)`
- `ObjectList[Type]()` - Converts slices to Terraform lists of objects, returns `(types.List, diag.Diagnostics)`
- `FilteredObject[Type]()` - Extracts specific fields from objects, returns `(types.Object, diag.Diagnostics)`
- `UpdatedObject()` - Updates object attributes, returns `(types.Object, diag.Diagnostics)`
- `ObjectStringIdentifier()` - Extracts string identifier from object, returns `(string, diag.Diagnostics)`
- `ListItemsIdentifiers()` - Extracts identifiers from list items, returns `([]string, diag.Diagnostics)`
- `ListSortedEntries[Type]()` - Sorts list entries by identifier order, returns `(types.List, diag.Diagnostics)`

**Usage Pattern**:
```go
// Call helper function
result, d := typehelpers.SomeHelper(input)
diags.Append(d...)
if diags.HasError() {
    return diags
}
// Use result
```

**Schema Validators** (`internal/schemavalidate/`):
- `OneOfCloudProviders()` - Validates cloud provider values
- `unique_string_attribute.go` - Ensures string uniqueness in lists

**Schema Defaults** (`internal/schemadefault/`):
- `EmptyListDefault()` - Returns an empty list default for optional list attributes
- `EmptyMapDefault()` - Returns an empty map default for optional map attributes
- Used to provide consistent defaults when fields are omitted in configuration

### Helper Packages

**Error Handling** (`internal/errors/`):
- **Legacy package** - Being phased out in favor of direct `diag.Diagnostics` usage
- `DiagError` interface for Terraform diagnostics
- `AddDiagError()` - Adds errors to diagnostics
- `AddDiagAttributeError()` - Adds attribute-specific errors

**Note**: New code should use `diag.Diagnostics` directly and the `typehelpers` package functions which all return diagnostics. Avoid using the `errors` package for new code.

**Model Update** (`internal/modelupdate/`):
- Helper functions for updating Terraform models from API results
- Centralizes common update patterns

**Provider Data** (`internal/providerdata/`):
- `ProviderData` struct holds API client for resources/data sources
- `GetResourceProviderData()` - Retrieves provider data in resources
- `GetDataSourceProviderData()` - Retrieves provider data in data sources

### Authentication

Provider supports multiple authentication methods with strict precedence order (tested in `internal/provider/provider_test.go`):

1. **Direct Configuration** (highest priority):
   ```hcl
   provider "stacklet" {
     endpoint = "https://api.instance.stacklet.io/"
     api_key  = "your-api-key"
   }
   ```

2. **Environment Variables**:
   - `STACKLET_ENDPOINT` - API endpoint URL
   - `STACKLET_API_KEY` - API key for authentication

3. **stacklet-admin CLI** (lowest priority):
   - `~/.stacklet/config.json` - Contains endpoint configuration (JSON with `"api"` field)
   - `~/.stacklet/credentials` - Contains API key (plain text file)
   - Automatically used if logged in via `stacklet-admin login`

**Credential Resolution**: Each credential (endpoint, API key) is resolved independently. For example, the endpoint can come from config while the API key comes from environment variables. The provider tries each source in order for each credential, using the first non-empty value found.

**Testing**: Comprehensive unit tests in `internal/provider/provider_test.go` verify credential precedence, mixed sources, and fallback behavior using testify test suites.

## Supported Resources

### Account Management
- `stacklet_account` - Cloud accounts
- `stacklet_account_discovery_aws` - AWS account discovery
- `stacklet_account_discovery_azure` - Azure subscription discovery
- `stacklet_account_discovery_gcp` - GCP project discovery
- `stacklet_account_group` - Account groups with optional dynamic filtering
  - Supports `dynamic_filter` attribute for automatic account matching
  - Dynamic filter changes (null ↔ non-null) trigger resource replacement
  - Uses `RequiresReplaceIfNullStringChange()` plan modifier
- `stacklet_account_group_mapping` - Account to group mappings

### Policy Management
- `stacklet_policy_collection` - Policy collections
- `stacklet_policy_collection_mapping` - Policy to collection mappings
- `stacklet_binding` - Bindings (policy collection + account group)

### Configuration
- `stacklet_configuration_profile_account_owners` - Account owners profile
- `stacklet_configuration_profile_email` - Email profile
- `stacklet_configuration_profile_jira` - Jira profile
- `stacklet_configuration_profile_resource_owner` - Resource owner profile
- `stacklet_configuration_profile_servicenow` - ServiceNow profile
- `stacklet_configuration_profile_slack` - Slack profile
- `stacklet_configuration_profile_symphony` - Symphony profile
- `stacklet_configuration_profile_teams` - Teams profile
- `stacklet_notification_template` - Notification templates
- `stacklet_report_group` - Report groups

### Repository Management
- `stacklet_repository` - Source code repositories

## Supported Data Sources

All resources have corresponding data sources, plus:
- `stacklet_policy` - Query individual policies
- `stacklet_platform` - Platform information and module versions
- `stacklet_configuration_profile_msteams` - Microsoft Teams profile (read-only)

## Development Workflow

1. **Build locally**:
   ```bash
   just build
   ```

2. **Configure Terraform to use local provider**:
   Create `~/.terraformrc`:
   ```hcl
   provider_installation {
     dev_overrides {
       "stacklet/stacklet" = "/absolute/path/to/terraform-provider-stacklet"
     }
     direct {}
   }
   ```

3. **Test changes**:
   ```bash
   just test              # Recorded tests (fast)
   just test-live         # Live API tests
   just test-record Foo   # Record new test
   ```

4. **Format and lint**:
   ```bash
   just format
   just lint
   ```

5. **Update documentation**:
   ```bash
   just docs
   ```

6. **Debug the provider**:
   ```bash
   # Terminal 1: Start provider in debug mode
   ./terraform-provider-stacklet -debug

   # Terminal 2: Export TF_REATTACH_PROVIDERS and run terraform
   export TF_REATTACH_PROVIDERS='...'  # From Terminal 1 output
   terraform plan
   ```

## Release Process

1. **Update CHANGELOG.md**:
   Add entry for the new version with changes

2. **Create release tag**:
   ```bash
   just tag-release X.Y.Z
   ```

3. **Push tag to trigger release**:
   ```bash
   git push origin vX.Y.Z
   ```

4. **GitHub Actions workflow**:
   - Builds provider binaries for multiple platforms
   - Signs releases with GPG
   - Creates GitHub release
   - Published to Terraform Registry automatically

## Architectural Patterns

### Diagnostics Pattern

The codebase follows a consistent diagnostics pattern throughout:

1. **Return Type**: All functions that can fail return `diag.Diagnostics` (plural), not `error` or `diag.Diagnostic` (singular)

2. **Initialization**: Start functions with `var diags diag.Diagnostics`

3. **Error Accumulation**: Use `diags.Append(d...)` to accumulate diagnostics from helper functions

4. **Early Returns**: Check `diags.HasError()` when you need to stop processing on errors

5. **Adding Errors**: Use `diags.AddError("Summary", "Detail")` for new errors

**Example**:
```go
func (m *Model) Update(ctx context.Context, apiObj *api.Object) diag.Diagnostics {
    var diags diag.Diagnostics

    // Simple field assignment
    m.ID = types.StringValue(apiObj.ID)

    // Complex conversion with error handling
    jsonValue, d := typehelpers.JSONString(apiObj.Variables)
    diags.Append(d...)
    if diags.HasError() {
        return diags
    }
    m.Variables = jsonValue

    // Another helper call
    obj, d := typehelpers.ObjectValue(ctx, apiObj.Config, constructorFunc)
    diags.Append(d...)
    m.Config = obj

    return diags
}
```

### Type Conversion Pattern

Use `internal/typehelpers` for all type conversions and manipulations:

- **JSON normalization**: `typehelpers.JSONString()` instead of manual JSON handling
- **Object conversion**: `typehelpers.ObjectValue()` for creating Terraform objects
- **List conversion**: `typehelpers.ObjectList()` for creating lists of objects
- **Identifier extraction**: `typehelpers.ObjectStringIdentifier()` and `typehelpers.ListItemsIdentifiers()`
- **List sorting**: `typehelpers.ListSortedEntries()` to preserve user-defined order

All typehelpers functions return `(result, diag.Diagnostics)` for consistent error handling.

### Model Update Pattern

Models follow a consistent Update() pattern:

1. **Data Source Models**: `Update()` populates the model from API objects
   ```go
   func (m *XxxDataSource) Update(ctx context.Context, obj *api.Xxx) diag.Diagnostics
   ```

2. **Resource Models**: `Update()` extends the data source pattern with resource-specific logic
   ```go
   func (m *XxxResource) Update(ctx context.Context, obj *api.Xxx) diag.Diagnostics {
       // Save original values if needed for special handling
       originalValue := m.SomeField

       // Call parent Update()
       diags := m.XxxDataSource.Update(ctx, obj)

       // Resource-specific adjustments
       if shouldResetValue {
           m.SomeField = originalValue
       }

       return diags
   }
   ```

3. **Order Preservation**: For lists where user-defined order matters:
   ```go
   // Before calling parent Update(), save identifiers
   identifiers, d := typehelpers.ListItemsIdentifiers(m.Items, "name")
   diags.Append(d...)

   // Call parent Update()
   d = m.XxxDataSource.Update(ctx, obj)
   diags.Append(d...)

   // Restore user-defined order
   if identifiers != nil {
       sorted, d := typehelpers.ListSortedEntries[ItemType](m.Items, "name", identifiers)
       diags.Append(d...)
       m.Items = sorted
   }
   ```

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
- **Always** return `diag.Diagnostics` (plural) from helper functions and model methods
- Use `diags.AddError("Summary", "Detail")` to add errors
- Use `diags.Append(d...)` to accumulate diagnostics from helper functions
- Check `diags.HasError()` for early returns when needed
- Avoid using the legacy `internal/errors` package for new code
- Provide clear error summaries and details

**Pattern**:
```go
func DoSomething() diag.Diagnostics {
    var diags diag.Diagnostics

    result, d := typehelpers.SomeHelper(input)
    diags.Append(d...)
    if diags.HasError() {
        return diags
    }

    // ... use result

    return diags
}
```

### Plan Modifiers
- Use `stringplanmodifier.RequiresReplace()` for immutable fields
- Use `stringplanmodifier.UseStateForUnknown()` for computed fields
- Custom modifiers in `internal/planmodifiers/` for complex logic:
  - **`RequiresReplaceIfNullStringChange()`**: Use when string fields that toggle between static and dynamic modes
    ```go
    "dynamic_filter": schema.StringAttribute{
        Optional: true,
        PlanModifiers: []planmodifier.String{
            planmodifiers.RequiresReplaceIfNullStringChange(),
        },
    }
    ```
  - **`RequiresReplaceIfFieldsChanged(fieldNames...)`**: Use for object attributes where specific nested fields require replacement
    ```go
    PlanModifiers: []planmodifier.Object{
        planmodifiers.RequiresReplaceIfFieldsChanged("repository_uuid"),
    }
    ```
  - **`RequiresReplaceIfUnset()`**: Use for object attributes that cannot transition from configured to unconfigured
  - **`DefaultObject()`**: Use to set default objects that match API defaults, preventing plan inconsistencies

### Write-Only Fields
- Use `_wo` suffix for write-only fields (e.g., `security_context_wo`)
- Use `_wo_version` suffix for version tracking
- Mark as `Sensitive: true` and `WriteOnly: true` in schema

### Collection Initialization

Use `make()` for empty collections, composite literals only when providing values or as direct function arguments:

```go
// Correct
values := make([]attr.Value, 0)
tagsMap := make(map[string]attr.Value)

// Incorrect
values := []attr.Value{}
tagsMap := map[string]attr.Value{}

// Exception - immediate values or function arguments
items := []string{"foo", "bar"}
types.ListValueMust(attrType, []attr.Value{})
```

## Documentation

- **Auto-generated**: Documentation in `docs/` is generated from code
- **Templates**: Use `templates/` for customizing documentation
- **Examples**: All examples in `examples/` must be valid Terraform
- **Import examples**: Each resource has `import.sh` example

## CI/CD

### Continuous Integration (`.github/workflows/ci.yaml`)
- Runs on push to `main` and all PRs
- **Security**: All GitHub Actions use SHA pinning instead of tags for supply chain security
- Uses `wistia/parse-tool-versions` to extract tool versions from `.tool-versions`
- Separate jobs for:
  - Go linting (`golangci-lint`)
  - Terraform linting
  - Documentation validation
  - Copyright header validation
  - Acceptance tests (replay mode)

### Release Workflow (`.github/workflows/release.yaml`)
- **Public Registry Release**: Triggered by version tags (`v*.*.*`)
  - Uses GoReleaser to build multi-platform binaries
  - Signs releases with GPG
  - Publishes to public Terraform Registry via GitHub releases
- **Private Registry Release** (`release-main` job): Currently disabled
  - Would trigger on commits to `main` branch
  - Uses GoReleaser with `--snapshot` flag
  - Integrates with `thechrisjohnson/terraform-cloud-provider-publish` action
  - Publishes to Terraform Cloud private registry

## Development Tools

Required tools (defined in `.tool-versions`):
- Go 1.24
- golangci-lint v2.1.5
- just 1.40.0

Additional tools (in `go.mod`):
- `github.com/hashicorp/terraform-plugin-docs` - Documentation generation
- `github.com/hashicorp/copywrite` - Copyright header management

## Common Development Tasks

### Adding a New Resource

1. Create model in `internal/models/`
2. Create API methods in `internal/api/`
3. Create resource in `internal/resources/`
4. Register in `internal/resources/resources.go`
5. Add examples in `examples/resources/`
6. Add acceptance test in `acceptance_tests/`
7. Generate docs: `just docs`

### Adding a New Data Source

1. Create or reuse model in `internal/models/`
2. Create or reuse API methods in `internal/api/`
3. Create data source in `internal/datasources/`
4. Register in `internal/datasources/datasources.go`
5. Add examples in `examples/data-sources/`
6. Add acceptance test in `acceptance_tests/`
7. Generate docs: `just docs`

### Recording New Test Interactions

When API changes or new tests are added:
```bash
# Set up credentials
export STACKLET_ENDPOINT="https://api.instance.stacklet.io/"
export STACKLET_API_KEY="your-api-key"

# Record the test
just test-record TestAccAccountResource

# Verify recording works
just test -run TestAccAccountResource
```

### Updating Dependencies

```bash
just update-deps-go
```

This updates both main dependencies and tools dependencies.
