# Create a policy collection
resource "stacklet_policy_collection" "aws_security" {
  name           = "aws-security-policies"
  cloud_provider = "AWS"
  description    = "Security policies for AWS resources"
  auto_update    = true
}

data "stacklet_repository" "policies" {
  url = "ssh://git@example.com/my-policies.git"
}

# Create a dynamic policy collection
resource "stacklet_policy_collection" "policies" {
  name           = "my-aws-policies"
  cloud_provider = "AWS"
  auto_update    = true
  dynamic_config = {
    repository_uuid = data.stacklet_repository.policies.uuid
    branch_name     = "aws"
  }
}
