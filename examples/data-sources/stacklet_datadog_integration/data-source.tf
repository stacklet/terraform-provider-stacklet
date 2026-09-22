data "stacklet_datadog_integration" "example" {}

output "datadog_ready" {
  description = "Whether credentials are held and every execution region that has reported holds the current configuration."
  value       = data.stacklet_datadog_integration.example.configured && data.stacklet_datadog_integration.example.status.status == "SUCCESS"
}
