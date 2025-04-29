provider "stacklet" {
  endpoint = "http://localhost:8080/graphql"
  api_key  = "your-api-key"
}

resource "stacklet_policy_collection" "example" {
  name           = "example-collection"
  cloud_provider = "aws"
  description    = "Example policy collection"
  auto_update    = true
}

output "policy_collection_id" {
  value = stacklet_policy_collection.example.id
}

output "policy_collection_uuid" {
  value = stacklet_policy_collection.example.uuid
}

output "policy_collection_system" {
  value = stacklet_policy_collection.example.system
}

output "policy_collection_repository" {
  value = stacklet_policy_collection.example.repository
} 