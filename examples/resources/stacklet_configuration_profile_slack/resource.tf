resource "stacklet_configuration_profile_slack" "example" {
  user_fields = ["username", "email"]

  webhook {
    name           = "webhook1"
    url_wo         = "https://example.com/webhooks/secret1"
    url_wo_version = "1"
  }

  webhook {
    name           = "webhook2"
    url_wo         = "https://example.com/webhooks/secret2"
    url_wo_version = "1"
  }
}
