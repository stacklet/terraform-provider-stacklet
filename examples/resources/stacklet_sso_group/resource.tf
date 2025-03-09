# Create an SSO group for administrators
resource "stacklet_sso_group" "admins" {
  name = "Administrators"
  roles = [
    "admin",
    "viewer"
  ]
  account_group_uuids = [
    "00000000-0000-0000-0000-000000000001",  # Production accounts
    "00000000-0000-0000-0000-000000000002"   # Development accounts
  ]
}

# Create an SSO group for developers
resource "stacklet_sso_group" "developers" {
  name = "Developers"
  roles = [
    "viewer"
  ]
  account_group_uuids = [
    "00000000-0000-0000-0000-000000000002"   # Development accounts only
  ]
}

# Output the group IDs
output "admin_group_id" {
  value = stacklet_sso_group.admins.id
}

output "developer_group_id" {
  value = stacklet_sso_group.developers.id
} 