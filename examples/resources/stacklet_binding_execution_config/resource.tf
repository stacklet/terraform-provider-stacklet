data "stacklet_binding" "example" {
  name = "apply-soc-policies"
}

resource "stacklet_binding_execution_config" "example" {
  binding_uuid = data.stacklet_binding.example.uuid

  dry_run = true
  variables = jsonencode({
    environment = "development"
    severity    = "medium"
  })
}
