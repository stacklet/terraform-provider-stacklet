// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
                    data "stacklet_role" "owner" {
                      name = "owner"
                    }

                    data "stacklet_role" "viewer" {
                      name = "viewer"
                    }
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check owner role
				resource.TestCheckResourceAttr("data.stacklet_role.owner", "name", "owner"),
				resource.TestCheckResourceAttr("data.stacklet_role.owner", "system", "true"),
				resource.TestCheckResourceAttrSet("data.stacklet_role.owner", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_role.owner", "permissions.#"),
				// Check viewer role
				resource.TestCheckResourceAttr("data.stacklet_role.viewer", "name", "viewer"),
				resource.TestCheckResourceAttr("data.stacklet_role.viewer", "system", "true"),
				resource.TestCheckResourceAttrSet("data.stacklet_role.viewer", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_role.viewer", "permissions.#"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRoleDataSource", steps)
}
