# Fetch user group information by name
data "stacklet_user_group" "example" {
  name = "engineering"
}

# Use the role_assignment_principal to grant roles to the group's members
output "group_principal" {
  value = data.stacklet_user_group.example.role_assignment_principal
}

# Use the role_assignment_target to grant roles on the group itself
output "group_target" {
  value = data.stacklet_user_group.example.role_assignment_target
}
