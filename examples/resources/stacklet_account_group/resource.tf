provider "stacklet" {
  endpoint = "http://localhost:8080/graphql"
  api_key  = "your-api-key"
}

resource "stacklet_account_group" "example" {
  name           = "example-group"
  cloud_provider = "AWS"
  description    = "Example account group"
}

output "account_group_id" {
  value = stacklet_account_group.example.id
}

output "account_group_uuid" {
  value = stacklet_account_group.example.uuid
}

output "account_group_description" {
  value = stacklet_account_group.example.description
}

output "account_group_cloud_provider" {
  value = stacklet_account_group.example.cloud_provider
} 