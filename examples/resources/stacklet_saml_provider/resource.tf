# A provider configured from a metadata endpoint.
resource "stacklet_saml_provider" "example" {
  name           = "okta"
  display_name   = "Okta"
  metadata_url   = "https://example.okta.com/app/exampleapp/sso/saml/metadata"
  enable_signout = true
  idp_alias      = "corp"
}

# A provider configured from an inline metadata document. Exactly one of
# metadata_url or metadata_xml must be set.
variable "entra_metadata_xml" {
  description = "SAML metadata XML document for the Entra ID application."
  type        = string
}

resource "stacklet_saml_provider" "inline" {
  name         = "entra"
  display_name = "Microsoft Entra ID"
  metadata_xml = var.entra_metadata_xml
}
