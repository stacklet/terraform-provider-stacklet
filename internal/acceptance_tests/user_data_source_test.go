// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_user" "test" {
					name     = "{{.Prefix}}-user-ds"
					username = "test_user_ds"
					email    = "test@stacklet.io"
				}

				data "stacklet_user" "test" {
					username = stacklet_user.test.username
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_user.test", "username", "test_user_ds"),
				resource.TestCheckResourceAttr("data.stacklet_user.test", "name", prefixName("user-ds")),
				resource.TestCheckResourceAttr("data.stacklet_user.test", "email", "test@stacklet.io"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "key"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "role_assignment_principal"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "active"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "sso_user"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccUserDataSource", steps)
}
