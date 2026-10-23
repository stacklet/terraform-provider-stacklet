# Anthropic only shows an admin key when it is created, so read it from
# somewhere it can be fetched again rather than from Terraform state.
variable "anthropic_admin_key" {
  description = "Anthropic Admin API key for the organization."
  type        = string
  sensitive   = true
}

resource "stacklet_anthropic_integration" "example" {
  org_id = "00000000-0000-0000-0000-000000000000"

  # To rotate the key, change the value and bump its version.
  admin_key_wo         = var.anthropic_admin_key
  admin_key_wo_version = "1"
}
