# Fetch an account group by UUID
data "stacklet_account_group" "example" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Fetch an account group by name
data "stacklet_account_group" "example" {
  name = "production-accounts"
}

# Output account group details
output "group_cloud_provider" {
  value = data.stacklet_account_group.example.cloud_provider
}

output "group_description" {
  value = data.stacklet_account_group.example.description
}
