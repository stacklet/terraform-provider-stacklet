# Create an AWS account group
resource "stacklet_account_group" "production" {
  name           = "production-accounts"
  cloud_provider = "aws"
  description    = "Production AWS accounts"
  regions        = ["us-east-1", "us-west-2"]
}

# Create an Azure account group
resource "stacklet_account_group" "development" {
  name           = "development-accounts"
  cloud_provider = "azure"
  description    = "Development Azure accounts"
}
