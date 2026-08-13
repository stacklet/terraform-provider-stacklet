# Fetch SAML provider information by name
data "stacklet_saml_provider" "example" {
  name = "okta"
}

output "saml_provider_display_name" {
  value = data.stacklet_saml_provider.example.display_name
}

output "saml_provider_metadata_url" {
  value = data.stacklet_saml_provider.example.metadata_url
}
