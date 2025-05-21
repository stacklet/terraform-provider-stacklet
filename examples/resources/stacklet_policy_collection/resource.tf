# Create an AWS policy collection
resource "stacklet_policy_collection" "aws_security" {
  name           = "aws-security-policies"
  cloud_provider = "aws"
  description    = "Security policies for AWS resources"
  auto_update    = true
}

# Create an Azure policy collection
resource "stacklet_policy_collection" "azure_compliance" {
  name           = "azure-compliance-policies"
  cloud_provider = "azure"
  description    = "Compliance policies for Azure resources"
  auto_update    = false
}

# Create a GCP policy collection
resource "stacklet_policy_collection" "gcp_cost" {
  name           = "gcp-cost-policies"
  cloud_provider = "gcp"
  description    = "Cost optimization policies for GCP resources"
  auto_update    = true
}
