# Configure an AWS account
resource "stacklet_account" "aws_prod" {
  cloud_provider   = "aws"
  key              = "123456789012" # AWS account ID
  name             = "Production AWS Account"
  description      = "Main production AWS account"
  email            = "cloud-team@example.com"
  security_context = "arn:aws:iam::123456789012:role/stacklet-execution"

  variables = jsonencode({
    environment = "production"
    team        = "platform"
    cost_center = "12345"
  })
}

# Configure an Azure subscription
resource "stacklet_account" "azure_dev" {
  cloud_provider = "azure"
  key            = "00000000-0000-0000-0000-000000000000" # Azure subscription ID
  name           = "Development Azure Subscription"
  description    = "Development environment in Azure"
  security_context = jsonencode({
    tenant_id     = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
    client_id     = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
    client_secret = "seKret"
  })
}

# Configure a GCP project
resource "stacklet_account" "gcp_staging" {
  cloud_provider   = "gcp"
  key              = "my-project-id" # GCP project ID
  name             = "Staging GCP Project"
  description      = "Staging environment in GCP"
  security_context = "arn:aws:secretsmanager:us-east-11:12345678912:secret:gcp-staging" # ARN of the secret containing the configuration
}

# Configure a Tencent Cloud account
resource "stacklet_account" "tencent_prod" {
  cloud_provider   = "tencentcloud"
  key              = "1234567890" # Tencent Cloud account ID
  name             = "Production Tencent Cloud Account"
  description      = "Production environment in Tencent Cloud"
  security_context = "arn:aws:secretsmanager:us-east-11:12345678912:secret:tencent-prod" # ARN of the secret containing the configuration
}
