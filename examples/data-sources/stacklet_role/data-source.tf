# Fetch the "owner" system role
data "stacklet_role" "owner" {
  name = "owner"
}

# Fetch the "viewer" system role
data "stacklet_role" "viewer" {
  name = "viewer"
}

# Fetch the "editor" system role
data "stacklet_role" "editor" {
  name = "editor"
}

# Fetch the "admin" system role
data "stacklet_role" "admin" {
  name = "admin"
}

# Output the permissions for the owner role
output "owner_permissions" {
  value = data.stacklet_role.owner.permissions
}
