provider "stacklet" {
  endpoint = "http://localhost:8080/graphql"
  api_key  = "your-api-key"
}

data "stacklet_account_group" "example" {
  name = "example-group"
}

output "account_group_id" {
  value = data.stacklet_account_group.example.id
}

output "account_group_uuid" {
  value = data.stacklet_account_group.example.uuid
}

output "account_group_description" {
  value = data.stacklet_account_group.example.description
}

output "account_group_cloud_provider" {
  value = data.stacklet_account_group.example.cloud_provider
} 