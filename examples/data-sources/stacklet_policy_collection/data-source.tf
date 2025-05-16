# Fetch a policy collection by UUID
data "stacklet_policy_collection" "by_uuid" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Fetch a policy collection by name
data "stacklet_policy_collection" "by_name" {
  name = "aws-security-policies"
}
