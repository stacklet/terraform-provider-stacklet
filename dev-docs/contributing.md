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

## Key Resources

See `internal/resources/resources.go` and `internal/datasources/datasources.go` for complete lists.

**Key Features**:
- Write-only fields use `_wo` suffix with version tracking via `_wo_version`.
