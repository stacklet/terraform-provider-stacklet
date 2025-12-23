# Fetch user information by username
data "stacklet_user" "example" {
  username = "username"
}

# Use the role_assignment_principal for role assignments
output "user_principal" {
  value = data.stacklet_user.example.role_assignment_principal
}

# Output user details
output "user_name" {
  value = data.stacklet_user.example.name
}

output "user_email" {
  value = data.stacklet_user.example.email
}

output "user_roles" {
  value = data.stacklet_user.example.roles
}
