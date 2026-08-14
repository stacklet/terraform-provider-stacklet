// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSAMLProviderDataSource(t *testing.T) {
	baseline := `
		resource "stacklet_saml_provider" "test" {
			name           = "{{.Prefix}}-saml-provider-ds"
			display_name   = "Test SAML Provider DS"
			metadata_url   = "https://example.com/saml/metadata-ds"
			enable_signout = true
		}
	`
	steps := []resource.TestStep{
		{
			Config: baseline + `
				data "stacklet_saml_provider" "test" {
					name = stacklet_saml_provider.test.name
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_saml_provider.test", "name", prefixName("saml-provider-ds")),
				resource.TestCheckResourceAttr("data.stacklet_saml_provider.test", "display_name", "Test SAML Provider DS"),
				resource.TestCheckResourceAttr("data.stacklet_saml_provider.test", "metadata_url", "https://example.com/saml/metadata-ds"),
				resource.TestCheckResourceAttr("data.stacklet_saml_provider.test", "enable_signout", "true"),
				resource.TestCheckResourceAttrSet("data.stacklet_saml_provider.test", "id"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccSAMLProviderDataSource", steps)
}
