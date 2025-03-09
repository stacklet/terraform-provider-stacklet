provider "stacklet" {
  endpoint = "http://localhost:8080/graphql"
  api_key  = "your-api-key"
}

data "stacklet_policy_collection" "example" {
  name = "example-collection"
}

output "policy_collection_id" {
  value = data.stacklet_policy_collection.example.id
}

output "policy_collection_uuid" {
  value = data.stacklet_policy_collection.example.uuid
}

output "policy_collection_description" {
  value = data.stacklet_policy_collection.example.description
}

output "policy_collection_cloud_provider" {
  value = data.stacklet_policy_collection.example.cloud_provider
}

output "policy_collection_auto_update" {
  value = data.stacklet_policy_collection.example.auto_update
}

output "policy_collection_system" {
  value = data.stacklet_policy_collection.example.system
}

output "policy_collection_repository" {
  value = data.stacklet_policy_collection.example.repository
} 