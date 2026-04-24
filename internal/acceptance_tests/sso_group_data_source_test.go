// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSOGroupDataSource(t *testing.T) {
	testSSOGroup := getenvOrSkip(t, "TF_ACC_TEST_SSO_GROUP")
	steps := []resource.TestStep{
		{
			Config: fmt.Sprintf(`
				data "stacklet_sso_group" "test" {
				  name = %q
				}
			`, testSSOGroup),
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check SSO group attributes
				resource.TestCheckResourceAttr("data.stacklet_sso_group.test", "name", testSSOGroup),
				resource.TestCheckResourceAttrSet("data.stacklet_sso_group.test", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_sso_group.test", "role_assignment_principal"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccSSOGroupDataSource", steps)
}
