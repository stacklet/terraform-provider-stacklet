// Copyright Stacklet, Inc. 2025, 2026



package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read
		{
			Config: `
				resource "stacklet_user" "test" {
					name     = "{{.Prefix}}-user"
					username = "test_user"
					email    = "test@stacklet.io"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_user.test", "name", prefixName("user")),
				resource.TestCheckResourceAttr("stacklet_user.test", "username", "test_user"),
				resource.TestCheckResourceAttr("stacklet_user.test", "email", "test@stacklet.io"),
				resource.TestCheckResourceAttrSet("stacklet_user.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_user.test", "key"),
				resource.TestCheckResourceAttrSet("stacklet_user.test", "role_assignment_principal"),
				resource.TestCheckResourceAttrSet("stacklet_user.test", "active"),
				resource.TestCheckResourceAttrSet("stacklet_user.test", "sso_user"),
			),
		},
		// ImportState
		{
			ResourceName:      "stacklet_user.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_user.test.username"),
		},
		// Update
		{
			Config: `
				resource "stacklet_user" "test" {
					name         = "{{.Prefix}}-user-updated"
					username     = "test_user"
					email        = "test+updated@stacklet.io"
					display_name = "Test User"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_user.test", "name", prefixName("user-updated")),
				resource.TestCheckResourceAttr("stacklet_user.test", "username", "test_user"),
				resource.TestCheckResourceAttr("stacklet_user.test", "email", "test+updated@stacklet.io"),
				resource.TestCheckResourceAttr("stacklet_user.test", "display_name", "Test User"),
				resource.TestCheckResourceAttrSet("stacklet_user.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_user.test", "key"),
				resource.TestCheckResourceAttrSet("stacklet_user.test", "role_assignment_principal"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccUserResource", steps)
}
