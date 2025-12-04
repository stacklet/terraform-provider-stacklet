# System-level admin role assignment for a user
resource "stacklet_role_assignment" "platform_admin" {
  role_name = "admin"

  principal {
    type = "user"
    id   = 123
  }

  target {
    type = "system"
  }
}

# Account group owner role assignment for an SSO group
resource "stacklet_role_assignment" "prod_owners" {
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

# Policy collection editor role assignment for a user
resource "stacklet_role_assignment" "policy_editor" {
  role_name = "editor"

  principal {
    type = "user"
    id   = 789
  }

  target {
    type = "policy-collection"
    uuid = stacklet_policy_collection.security.uuid
  }
}

# Repository viewer role assignment for an SSO group
resource "stacklet_role_assignment" "repo_viewers" {
  role_name = "viewer"

  principal {
    type = "sso-group"
    id   = 101
  }

  target {
    type = "repository"
    uuid = stacklet_repository.policies.uuid
  }
}

# Example: Query role assignments for an account group
data "stacklet_role_assignments" "prod_acl" {
  target {
    type = "account-group"
    uuid = stacklet_account_group.production.uuid
  }
}

output "production_access" {
  description = "All role assignments for production account group"
  value       = data.stacklet_role_assignments.prod_acl.assignments
}
