# Fetch an SSO group configuration
data "stacklet_sso_group" "example" {
  name = "Administrators"
}

# Output the roles assigned to the group
output "admin_roles" {
  value = data.stacklet_sso_group.example.roles
}

# Output the account groups accessible to the group
output "admin_account_groups" {
  value = data.stacklet_sso_group.example.account_group_uuids
} 