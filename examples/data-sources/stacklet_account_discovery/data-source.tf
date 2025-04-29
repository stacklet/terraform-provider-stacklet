# Fetch an AWS account discovery configuration
data "stacklet_account_discovery" "aws_example" {
  name = "aws-production"
}

# Output the configuration details
output "aws_discovery_config" {
  value     = jsondecode(data.stacklet_account_discovery.aws_example.config)
  sensitive = true
}

# Output the validity status
output "aws_discovery_validity" {
  value = jsondecode(data.stacklet_account_discovery.aws_example.validity)
}

# Fetch an Azure account discovery configuration
data "stacklet_account_discovery" "azure_example" {
  name = "azure-production"
}

# Fetch a GCP account discovery configuration
data "stacklet_account_discovery" "gcp_example" {
  name = "gcp-production"
} 