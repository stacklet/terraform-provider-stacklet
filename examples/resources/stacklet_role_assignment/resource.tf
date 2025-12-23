# Fetch user information
data "stacklet_user" "example_user" {
  email = "user@example.com"
}

# Fetch SSO group information
data "stacklet_sso_group" "example_group" {
  name = "Engineering"
}

# Fetch account group
data "stacklet_account_group" "example" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Assign a system-wide role to a user
resource "stacklet_role_assignment" "system_admin" {
  role_name = "admin"
  principal = data.stacklet_user.example_user.role_assignment_principal
  target    = "system:all"
}

# Assign a role to an SSO group on a specific target resource
resource "stacklet_role_assignment" "group_viewer" {
  role_name = "viewer"
  principal = data.stacklet_sso_group.example_group.role_assignment_principal
  target    = data.stacklet_account_group.example.role_assignment_target
}
