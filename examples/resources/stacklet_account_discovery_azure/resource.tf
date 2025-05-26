resource "stacklet_account_discovery_azure" "example" {
  name             = "test-azure"
  description      = "Azure tenant discovery"
  tenant_id        = "00000000-0000-0000-0000-000000000000"
  client_id        = "11111111-1111-1111-1111-111111111111"
  client_secret_wo = "your-client-secret"
}
