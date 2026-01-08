# Fetch SSO group information by name
data "stacklet_sso_group" "example" {
  name = "Engineering"
}

# Use the role_assignment_principal for role assignments
output "group_principal" {
  value = data.stacklet_sso_group.example.role_assignment_principal
}

# Output group details
output "group_display_name" {
  value = data.stacklet_sso_group.example.display_name
}
