# Fetch a policy by UUID
data "stacklet_policy" "example" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Fetch a policy by name
data "stacklet_policy" "s3_encryption" {
  name = "s3-bucket-encryption"
}
