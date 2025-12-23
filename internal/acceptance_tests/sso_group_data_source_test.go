// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSOGroupDataSource(t *testing.T) {
	// Require TF_ACC_TEST_SSO_GROUP environment variable
	testSSOGroup := os.Getenv("TF_ACC_TEST_SSO_GROUP")
	if testSSOGroup == "" {
		t.Skip("TF_ACC_TEST_SSO_GROUP environment variable must be set to run this test")
	}

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
