// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	testUsername := getenvOrSkip(t, "TF_ACC_TEST_USERNAME")
	steps := []resource.TestStep{
		{
			Config: fmt.Sprintf(`
 				data "stacklet_user" "test" {
					username = "%s"
				}
			`, testUsername),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_user.test", "username", testUsername),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "key"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "name"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "email"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "role_assignment_principal"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "active"),
				resource.TestCheckResourceAttrSet("data.stacklet_user.test", "sso_user"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccUserDataSource", steps)
}
