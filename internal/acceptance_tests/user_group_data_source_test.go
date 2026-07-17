// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupDataSource(t *testing.T) {
	baseline := `
		resource "stacklet_user_group" "test" {
			name         = "{{.Prefix}}-user-group-ds"
			display_name = "Test User Group DS"
		}
	`
	steps := []resource.TestStep{
		{
			Config: baseline + `
				data "stacklet_user_group" "test" {
					name = stacklet_user_group.test.name
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_user_group.test", "name", prefixName("user-group-ds")),
				resource.TestCheckResourceAttr("data.stacklet_user_group.test", "display_name", "Test User Group DS"),
				resource.TestCheckResourceAttrSet("data.stacklet_user_group.test", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_user_group.test", "uuid"),
				resource.TestCheckResourceAttrSet("data.stacklet_user_group.test", "role_assignment_principal"),
				resource.TestCheckResourceAttrSet("data.stacklet_user_group.test", "role_assignment_target"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccUserGroupDataSource", steps)
}
