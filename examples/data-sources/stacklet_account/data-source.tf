# Fetch an AWS account by its account ID
data "stacklet_account" "example" {
  cloud_provider = "aws"
  key           = "123456789012"
}

# Output the account name
output "account_name" {
  value = data.stacklet_account.example.name
}

# Output the account path
output "account_path" {
  value = data.stacklet_account.example.path
} 