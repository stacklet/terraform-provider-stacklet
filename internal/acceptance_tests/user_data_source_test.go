// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	testUsername := os.Getenv("TF_ACC_TEST_USERNAME")

	if testUsername == "" {
		t.Skip("TF_ACC_TEST_USERNAME environment variable must be set to run this test")
	}

	config := fmt.Sprintf(`
		data "stacklet_user" "test" {
		  username = %q
		}
	`, testUsername)

	steps := []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check user attributes
				resource.TestCheckResourceAttr("data.stacklet_user.test", "username", testUsername),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "name"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "role_assignment_principal"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "active"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "sso_user"),
				// Check that role lists are present (may be empty)
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "all_roles.#"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "assigned_roles.#"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "implicit_roles.#"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "inherited_roles.#"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "roles.#"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "groups.#"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccUserDataSource", steps)
}
