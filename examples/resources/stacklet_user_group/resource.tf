resource "stacklet_user_group" "example" {
  name         = "engineering"
  display_name = "Engineering"
}

# Grant the members of the group a role on a target (group as principal).
resource "stacklet_role_assignment" "engineering_viewer" {
  role_name = "viewer"
  principal = stacklet_user_group.example.role_assignment_principal
  target    = "system:all"
}
