# Stage 1: RBAC Management in Terraform Provider

## Overview

Stage 1 focuses on enabling role assignment management using existing system roles. This provides immediate value while establishing patterns for Stage 2 (custom roles and principal management).

**Timeline**: 2-2.5 weeks
**Status**: 4/5 PRs Complete ✅✅✅✅

---

## PR Breakdown

### ✅ PR 1: API Client Wrappers - COMPLETED

**Goal**: Create Go API wrappers for existing GraphQL operations

**Files Created**:
- `internal/api/role.go` (~60 lines) - Read operations for roles
- `internal/api/role_assignment.go` (~145 lines) - Full CRUD for role assignments
- `internal/api/api.go` (modified) - Registered new APIs

**Key Features**:
- Role API: `Read(name)`, `List()`
- Role Assignment API: `Create()`, `Read(id)`, `List(target, principal)`, `Delete()`
- Follows existing provider patterns
- Error handling with NotFound support

**Build Status**: ✅ Passing

---

### ✅ PR 2: Data Source - stacklet_role - COMPLETED

**Goal**: Enable querying existing system roles from Terraform

**Files Created**:
- `internal/models/role.go` (~38 lines)
- `internal/datasources/role.go` (~70 lines)
- `internal/acceptance_tests/role_data_source_test.go` (~38 lines)
- `examples/data-sources/stacklet_role/data-source.tf` (~26 lines)
- `internal/datasources/datasources.go` (modified)

**Usage Example**:
```terraform
data "stacklet_role" "owner" {
  name = "owner"
}

output "owner_permissions" {
  value = data.stacklet_role.owner.permissions
}
```

**Build Status**: ✅ Passing

**Testing**: Acceptance test created, needs recording with real API

---

### ✅ PR 3: Resource - stacklet_role_assignment (Core Implementation) - COMPLETED

**Goal**: Enable managing role assignments via Terraform

**Actual Size**: Large (~565 lines with tests)

**Files Created**:
- `internal/models/role_assignment.go` (~110 lines) - Model with nested principal/target
- `internal/resources/role_assignment.go` (~180 lines) - Full CRUD resource
- `internal/acceptance_tests/role_assignment_resource_test.go` (~200 lines) - 4 comprehensive tests
- `examples/resources/stacklet_role_assignment/resource.tf` (~75 lines) - All target type examples
- `internal/resources/resources.go` (modified) - Registered resource

**Usage Example**:
```terraform
# System-level role assignment
resource "stacklet_role_assignment" "admin_system" {
  role_name = "admin"

  principal {
    type = "user"
    id   = 123
  }

  target {
    type = "system"
  }
}

# Account group role assignment
resource "stacklet_role_assignment" "owner_prod" {
  role_name = "owner"

  principal {
    type = "sso-group"
    id   = 456
  }

  target {
    type = "account-group"
    uuid = stacklet_account_group.production.uuid
  }
}
```

**Key Features**:
- Full CRUD operations (Create, Read, Delete)
- No Update support (all fields require replacement)
- ImportState by ID
- Validators for principal.type and target.type
- 4 comprehensive tests covering all target types
- Both principal types tested (user, sso-group)

**Build Status**: ✅ Passing

**Testing**: 4 acceptance tests created:
- `TestAccRoleAssignmentResource_System`
- `TestAccRoleAssignmentResource_AccountGroup`
- `TestAccRoleAssignmentResource_PolicyCollection`
- `TestAccRoleAssignmentResource_Repository`

---

### ✅ PR 4: Data Source - stacklet_role_assignments - COMPLETED

**Goal**: Enable querying existing role assignments for a target

**Actual Size**: Small (~290 lines with tests)

**Status**: Complete

**Files Created**:
- `internal/models/role_assignment.go` (extended with data source model - ~172 lines total)
- `internal/datasources/role_assignments.go` (~137 lines)
- `internal/acceptance_tests/role_assignments_data_source_test.go` (~153 lines)
- `examples/data-sources/stacklet_role_assignments/data-source.tf` (~65 lines)
- `internal/datasources/datasources.go` (modified - registered data source)

**Usage Example**:
```terraform
# Query all role assignments for a target
data "stacklet_role_assignments" "production_acl" {
  target {
    type = "account-group"
    uuid = stacklet_account_group.production.uuid
  }
}

output "production_access" {
  value = data.stacklet_role_assignments.production_acl.assignments
}
```

**Key Features**:
- Query role assignments by target (system, account-group, policy-collection, repository)
- Returns complete assignment details including role, principal, and target
- 3 comprehensive acceptance tests covering different target types
- Example usage with output formatting

**Build Status**: ✅ Passing

**Testing**: 3 acceptance tests created:
- `TestAccRoleAssignmentsDataSource_System`
- `TestAccRoleAssignmentsDataSource_AccountGroup`
- `TestAccRoleAssignmentsDataSource_PolicyCollection`

**Acceptance Criteria**:
- ✅ Can query role assignments by target
- ✅ Returns complete assignment details (role, principal, target)
- ✅ Acceptance tests created
- ✅ Example documentation created

---

### PR 5: Provider Registration, Integration & Documentation

**Goal**: Wire everything together and document usage

**Estimated Size**: Small (~200-250 lines plus documentation)

**Files to Create/Modify**:
- `internal/provider/provider.go` (register resources/data sources)
- `CHANGELOG.md` (document v0.6.0 changes)
- `docs/guides/managing-rbac.md` (user guide)
- Integration test examples

**Documentation Coverage**:
- Overview of RBAC concepts in Stacklet
- Common use cases:
  - Assigning system administrators
  - Granting account group access
  - Managing policy collection permissions
  - Repository access control
- Importing existing role assignments
- Best practices for organizing RBAC as code

**Acceptance Criteria**:
- ✅ New resources/data sources registered in provider
- ✅ Provider builds and installs successfully
- ✅ End-to-end smoke test passes
- ✅ CHANGELOG documents new features
- ✅ User guide with working examples
- ✅ Migration guide explains import workflow

---

## Dependencies

**Sequential Dependencies**:
- PR 2, 3, 4 all depend on PR 1 (API client)
- PR 5 depends on PR 2, 3, 4 (integration point)

**Critical Path**: PR 1 → PR 3 → PR 5

**Parallel Work Opportunities**:
- After PR 1: PR 2, 3, 4 can be developed in parallel (PR 2 already complete)
- PR 3 is the most complex and critical

---

## Stage 1 Deliverables

**What Stage 1 Delivers**:
- ✅ Assign system roles (owner, viewer, editor, admin) to existing principals
- ✅ Query existing roles and role assignments
- ✅ Import existing role assignments into Terraform state
- ✅ Full documentation and examples
- ✅ Foundation for Stage 2 (custom roles and principal management)

**What's Deferred to Stage 2**:
- Custom role creation (`stacklet_role` resource - write operations)
- User management (`stacklet_user` resource)
- SSO group management (`stacklet_sso_group` resource)
- Principal lookup data sources

---

## Testing Strategy

### Acceptance Tests
- Uses recorded HTTP responses pattern
- Three modes: record, replay (default), live
- Recordings stored in `internal/acceptance_tests/recordings/`

### Running Tests
```bash
# Record mode (requires real API access)
TF_ACC=1 TF_ACC_MODE=record go test -v ./internal/acceptance_tests

# Replay mode (uses recorded responses - default)
TF_ACC=1 go test -v ./internal/acceptance_tests

# Live mode (always hits real API)
TF_ACC=1 TF_ACC_MODE=live go test -v ./internal/acceptance_tests
```

### Test Coverage Requirements
- Each resource needs acceptance tests for all CRUD operations
- Import functionality must be tested
- Data sources need at least one query test
- All target types must be covered

---

## Next Steps

**Immediate**: Start PR 3 - `stacklet_role_assignment` resource
- This is the critical path and largest PR
- Most valuable feature for users
- Establishes resource patterns for Stage 2

**After PR 3**: PRs 4 and 5 are smaller and can move quickly

**Estimated Time to Complete Stage 1**:
- PR 3: 4-5 days
- PR 4: 2-3 days
- PR 5: 2-3 days
- **Total remaining**: ~1.5-2 weeks

---

## Reference Links

- Confluence Proposal: https://stacklet.atlassian.net/wiki/spaces/ENG/pages/2354577409/S117+RBAC+Management+in+Terraform+Provider+WIP
- Provider Repository: `/home/chrism/Projects/terraform-provider-stacklet`
- Terraform Plugin Framework: https://developer.hashicorp.com/terraform/plugin/framework

---

**Last Updated**: 2025-12-03
**Status**: PR 1 & PR 2 Complete, Ready for PR 3
