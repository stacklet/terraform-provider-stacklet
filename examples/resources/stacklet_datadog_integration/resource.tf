variable "datadog_api_key" {
  description = "Datadog API key, which identifies the organization."
  type        = string
  sensitive   = true
}

# Datadog only shows an application key when it is created, so read it from
# somewhere it can be fetched again rather than from Terraform state.
variable "datadog_app_key" {
  description = "Datadog application key, scoped to query timeseries metrics."
  type        = string
  sensitive   = true
}

resource "stacklet_datadog_integration" "example" {
  # Defaults to true. Set it to false to stop policies using the
  # datadog-metrics filter while keeping the stored credentials, so the
  # integration can be switched back on without re-entering them.
  enabled = true

  site = "datadoghq.com"

  # To rotate a key, change the value and bump its version. Leaving the other
  # version alone keeps the other stored key untouched.
  api_key_wo         = var.datadog_api_key
  api_key_wo_version = "1"
  app_key_wo         = var.datadog_app_key
  app_key_wo_version = "1"
}
