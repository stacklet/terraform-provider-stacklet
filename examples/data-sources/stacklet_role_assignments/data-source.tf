# Query all system-level role assignments
data "stacklet_role_assignments" "system_access" {
  target {
    type = "system"
  }
}

# Query role assignments for a specific account group
data "stacklet_role_assignments" "production_acl" {
  target {
    type = "account-group"
    uuid = stacklet_account_group.production.uuid
  }
}

# Query role assignments for a policy collection
data "stacklet_role_assignments" "security_policies_access" {
  target {
    type = "policy-collection"
    uuid = stacklet_policy_collection.security.uuid
  }
}

# Query role assignments for a repository
data "stacklet_role_assignments" "repo_access" {
  target {
    type = "repository"
    uuid = stacklet_repository.custom_policies.uuid
  }
}

# Output all system administrators
output "system_admins" {
  description = "All principals with system-level access"
  value = [
    for assignment in data.stacklet_role_assignments.system_access.assignments :
    {
      role      = assignment.role_name
      principal = assignment.principal
    }
  ]
}

# Output production account group access control list
output "production_access_summary" {
  description = "Summary of who has access to the production account group"
  value = {
    total_assignments = length(data.stacklet_role_assignments.production_acl.assignments)
    assignments       = data.stacklet_role_assignments.production_acl.assignments
  }
}

# Check if specific user has access
locals {
  user_id_to_check = 123

  user_has_production_access = anytrue([
    for assignment in data.stacklet_role_assignments.production_acl.assignments :
    assignment.principal.type == "user" && assignment.principal.id == local.user_id_to_check
  ])
}

output "user_has_access" {
  description = "Whether user 123 has any role on the production account group"
  value       = local.user_has_production_access
}
