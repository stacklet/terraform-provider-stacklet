resource "stacklet_iac_policy_collection" "example" {
  name       = "example"
  is_default = true
  policies = [
    "s3-bucket-encryption",
    "ec2-no-public-ip",
  ]
}
