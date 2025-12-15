# Assign a system-wide role
resource "stacklet_role_assignment" "system_admin" {
  role_name = "admin"
  # Principal is an opaque identifier - in practice you would reference:
  # stacklet_user.example.role_assignment_principal or
  # stacklet_sso_group.example.role_assignment_principal
  principal = "user:123"
  target    = "system:all"
}

# Assign a role to a principal on a specific target resource
# resource "stacklet_role_assignment" "example" {
#   role_name = "viewer"
#   principal = stacklet_user.example.role_assignment_principal
#   target    = stacklet_account_group.example.role_assignment_target
# }
