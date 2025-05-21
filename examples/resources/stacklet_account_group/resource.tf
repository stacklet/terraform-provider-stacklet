# Create an AWS account group
resource "stacklet_account_group" "production" {
  name           = "production-accounts"
  cloud_provider = "AWS"
  description    = "Production AWS accounts"
  regions        = ["us-east-1", "us-west-2"]
}

# Create an Azure account group
resource "stacklet_account_group" "development" {
  name           = "development-accounts"
  cloud_provider = "Azure"
  description    = "Development Azure accounts"
}
