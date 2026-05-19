resource "stacklet_gcp_integration" "example" {
  key = "my-gcp-integration"

  customer_config_input = {
    infrastructure = {
      project_id        = "my-stacklet-project"
      resource_location = "us-central1"
      resource_prefix   = "stacklet"
    }
    organizations = [
      {
        org_id     = "123456789012"
        folder_ids = ["folders/111111111", "folders/222222222"]
      }
    ]
    cost_sources = [
      {
        billing_table = "my-billing-project.billing_dataset.gcp_billing_export_v1_ABCDEF_123456_789012"
      }
    ]
  }
}
