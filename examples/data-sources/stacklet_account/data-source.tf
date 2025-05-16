# Fetch an AWS account by its account ID
data "stacklet_account" "example" {
  cloud_provider = "AWS"
  key            = "123456789012"
}

# Fetch an Azure subscription
data "stacklet_account" "azure_dev" {
  cloud_provider = "Azure"
  key            = "00000000-0000-0000-0000-000000000000" # Azure subscription ID
}

# Fetch a GCP project
data "stacklet_account" "gcp_staging" {
  cloud_provider = "GCP"
  key            = "my-project-id" # GCP project ID
}

# Fetch a Tencent Cloud account
data "stacklet_account" "tencent_prod" {
  cloud_provider = "TencentCloud"
  key            = "1234567890" # Tencent Cloud account ID
}
