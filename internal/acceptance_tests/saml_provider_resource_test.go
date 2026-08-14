// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSAMLProviderResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
				resource "stacklet_saml_provider" "test" {
					name         = "{{.Prefix}}-saml-provider"
					display_name = "Test SAML Provider"
					metadata_url = "https://example.com/saml/metadata"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_saml_provider.test", "name", prefixName("saml-provider")),
				resource.TestCheckResourceAttr("stacklet_saml_provider.test", "display_name", "Test SAML Provider"),
				resource.TestCheckResourceAttr("stacklet_saml_provider.test", "metadata_url", "https://example.com/saml/metadata"),
				resource.TestCheckResourceAttr("stacklet_saml_provider.test", "enable_signout", "false"),
				resource.TestCheckNoResourceAttr("stacklet_saml_provider.test", "metadata_xml"),
				resource.TestCheckNoResourceAttr("stacklet_saml_provider.test", "idp_alias"),
				resource.TestCheckResourceAttrSet("stacklet_saml_provider.test", "id"),
			),
		},
		// ImportState testing using the name
		{
			ResourceName:      "stacklet_saml_provider.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_saml_provider.test.name"),
		},
		// Update and Read testing
		{
			Config: `
				resource "stacklet_saml_provider" "test" {
					name           = "{{.Prefix}}-saml-provider"
					display_name   = "Updated SAML Provider"
					metadata_url   = "https://example.com/saml/metadata-updated"
					enable_signout = true
					idp_alias      = "{{.Prefix}}-alias"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_saml_provider.test", "name", prefixName("saml-provider")),
				resource.TestCheckResourceAttr("stacklet_saml_provider.test", "display_name", "Updated SAML Provider"),
				resource.TestCheckResourceAttr("stacklet_saml_provider.test", "metadata_url", "https://example.com/saml/metadata-updated"),
				resource.TestCheckResourceAttr("stacklet_saml_provider.test", "enable_signout", "true"),
				resource.TestCheckResourceAttr("stacklet_saml_provider.test", "idp_alias", prefixName("alias")),
			),
		},
	}
	runRecordedAccTest(t, "TestAccSAMLProviderResource", steps)
}
