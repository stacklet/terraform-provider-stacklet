resource "stacklet_configuration_profile_resource_owner" "example" {
  tags = ["Owner", "ResourceOwner"]
  default = [
    "user1@example.com",
    "user2@example.com",
  ]
}
