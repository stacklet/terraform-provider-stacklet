# Add a policy to a collection
resource "stacklet_policy_collection_item" "example" {
  collection_uuid = "00000000-0000-0000-0000-000000000000"
  policy_uuid     = "11111111-1111-1111-1111-111111111111"
  policy_version  = 1
}
