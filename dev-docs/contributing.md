# Contributing

## Adding a New Resource

1. Create model in `internal/models/`
2. Create API methods in `internal/api/`
3. Create resource in `internal/resources/`
4. Register in `internal/resources/resources.go` in the `RESOURCES` slice
5. Add examples in `examples/resources/`
6. Add acceptance test in `acceptance_tests/`
7. Generate docs: `just docs`

## Adding a New Data Source

1. Create or reuse model in `internal/models/`
2. Create or reuse API methods in `internal/api/`
3. Create data source in `internal/datasources/`
4. Register in `internal/datasources/datasources.go` in the `DATASOURCES` slice
5. Add examples in `examples/data-sources/`
6. Add acceptance test in `acceptance_tests/`
7. Generate docs: `just docs`

**Note**: Data sources often reuse models and API methods from corresponding resources.

## Recording Tests

When API changes or new tests are added:

1. Set up credentials:
   ```bash
   export STACKLET_ENDPOINT="https://api.instance.stacklet.io/"
   export STACKLET_API_KEY="your-api-key"
   ```

2. Record the test:
   ```bash
   just test-record TestAccAccountResource
   ```

3. Verify recording works:
   ```bash
   just test -run TestAccAccountResource
   ```

### RBAC Tests (Special Requirements)

RBAC tests require existing user and SSO group identities from your Stacklet environment:

```bash
# Required environment variables for RBAC tests
# For user tests
export TF_ACC_TEST_USERNAME="username"

# For SSO group tests
export TF_ACC_TEST_SSO_GROUP="YourSSOGroupName"

# Record RBAC tests
just test-record TestAccUserDataSource
just test-record TestAccSSOGroupDataSource
just test-record TestAccRoleAssignmentResource
```

Tests will automatically skip if these environment variables are not set.

## Key Resources

See `internal/resources/resources.go` and `internal/datasources/datasources.go` for complete lists.

**Key Features**:
- Write-only fields use `_wo` suffix with version tracking via `_wo_version`.
