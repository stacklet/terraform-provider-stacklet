# Fetch a binding by UUID
data "stacklet_binding" "by_uuid" {
  uuid = "00000000-0000-0000-0000-000000000000"
}

# Fetch a binding by name
data "stacklet_binding" "by_name" {
  name = "production-security-policies"
}

# Output binding details
output "binding_name" {
  value = data.stacklet_binding.by_uuid.name
}

output "binding_uuid" {
  value = data.stacklet_binding.by_name.uuid
}
