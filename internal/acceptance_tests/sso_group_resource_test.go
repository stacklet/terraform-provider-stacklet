// Copyright Stacklet, Inc. 2025, 2026



package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSOGroupResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read
		{
			Config: `
				resource "stacklet_sso_group" "test" {
					name = "{{.Prefix}}-sso-group"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_sso_group.test", "name", prefixName("sso-group")),
				resource.TestCheckResourceAttrSet("stacklet_sso_group.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_sso_group.test", "role_assignment_principal"),
			),
		},
		// ImportState
		{
			ResourceName:      "stacklet_sso_group.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_sso_group.test.name"),
		},
		// Update
		{
			Config: `
				resource "stacklet_sso_group" "test" {
					name         = "{{.Prefix}}-sso-group"
					display_name = "Test SSO Group"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_sso_group.test", "name", prefixName("sso-group")),
				resource.TestCheckResourceAttr("stacklet_sso_group.test", "display_name", "Test SSO Group"),
				resource.TestCheckResourceAttrSet("stacklet_sso_group.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_sso_group.test", "role_assignment_principal"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccSSOGroupResource", steps)
}
