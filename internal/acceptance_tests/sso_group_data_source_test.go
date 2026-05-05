// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSOGroupDataSource(t *testing.T) {
	baseline := `
		resource "stacklet_sso_group" "test" {
			name = "{{.Prefix}}-sso-group-ds"
		}
	`
	steps := []resource.TestStep{
		{
			Config: baseline + `
				data "stacklet_sso_group" "test" {
					name = stacklet_sso_group.test.name
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_sso_group.test", "name", prefixName("sso-group-ds")),
				resource.TestCheckResourceAttrSet("data.stacklet_sso_group.test", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_sso_group.test", "role_assignment_principal"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccSSOGroupDataSource", steps)
}
