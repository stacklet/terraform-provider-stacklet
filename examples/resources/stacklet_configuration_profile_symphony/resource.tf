resource "stacklet_configuration_profile_symphony" "example" {
  agent_domain    = "mydomain"
  service_account = "acc1"

  private_key_wo         = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC7\n-----END PRIVATE KEY-----"
  private_key_wo_version = "1"
}
