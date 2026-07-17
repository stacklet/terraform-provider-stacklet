// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
				resource "stacklet_user_group" "test" {
					name         = "{{.Prefix}}-user-group"
					display_name = "Test User Group"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_user_group.test", "name", prefixName("user-group")),
				resource.TestCheckResourceAttr("stacklet_user_group.test", "display_name", "Test User Group"),
				resource.TestCheckResourceAttrSet("stacklet_user_group.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_user_group.test", "uuid"),
				resource.TestCheckResourceAttrSet("stacklet_user_group.test", "role_assignment_principal"),
				resource.TestCheckResourceAttrSet("stacklet_user_group.test", "role_assignment_target"),
			),
		},
		// ImportState testing using UUID
		{
			ResourceName:      "stacklet_user_group.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_user_group.test.uuid"),
		},
		// Update and Read testing
		{
			Config: `
				resource "stacklet_user_group" "test" {
					name         = "{{.Prefix}}-user-group-updated"
					display_name = "Updated User Group"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_user_group.test", "name", prefixName("user-group-updated")),
				resource.TestCheckResourceAttr("stacklet_user_group.test", "display_name", "Updated User Group"),
				resource.TestCheckResourceAttrSet("stacklet_user_group.test", "uuid"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccUserGroupResource", steps)
}
