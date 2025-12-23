// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRoleAssignmentResource(t *testing.T) {
	testUsername := os.Getenv("TF_ACC_TEST_USERNAME")

	if testUsername == "" {
		t.Skip("TF_ACC_TEST_USERNAME environment variable must be set to run this test")
	}

	userDataSourceConfig := fmt.Sprintf(`username = %q`, testUsername)

	steps := []resource.TestStep{
		// Create and Read testing - role assignment on account group
		{
			Config: fmt.Sprintf(`
				data "stacklet_user" "test" {
					%s
				}

				resource "stacklet_account_group" "test" {
					name           = "{{.Prefix}}-role-assignment-test"
					description    = "Test account group for role assignment"
					cloud_provider = "AWS"
					regions        = ["us-east-1"]
				}

				resource "stacklet_role_assignment" "test" {
					role_name = "viewer"
					principal = data.stacklet_user.test.role_assignment_principal
					target    = stacklet_account_group.test.role_assignment_target
				}
			`, userDataSourceConfig),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "role_name", "viewer"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "principal"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "target"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_role_assignment.test",
			ImportState:       true,
			ImportStateVerify: true,
			// Role assignments use composite key: role_name,principal,target
			ImportStateIdFunc: func(s *terraform.State) (string, error) {
				rs := s.RootModule().Resources["stacklet_role_assignment.test"]
				roleName := rs.Primary.Attributes["role_name"]
				principal := rs.Primary.Attributes["principal"]
				target := rs.Primary.Attributes["target"]
				return fmt.Sprintf("%s,%s,%s", roleName, principal, target), nil
			},
		},
		// Test replacement with different role
		{
			Config: fmt.Sprintf(`
				data "stacklet_user" "test" {
					%s
				}

				resource "stacklet_account_group" "test" {
					name           = "{{.Prefix}}-role-assignment-test"
					description    = "Test account group for role assignment"
					cloud_provider = "AWS"
					regions        = ["us-east-1"]
				}

				resource "stacklet_role_assignment" "test" {
					role_name = "editor"
					principal = data.stacklet_user.test.role_assignment_principal
					target    = stacklet_account_group.test.role_assignment_target
				}
			`, userDataSourceConfig),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "role_name", "editor"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "principal"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "target"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentResource", steps)
}

func TestAccRoleAssignmentResource_SSOGroup(t *testing.T) {
	// Require TF_ACC_TEST_SSO_GROUP environment variable
	testSSOGroup := os.Getenv("TF_ACC_TEST_SSO_GROUP")
	if testSSOGroup == "" {
		t.Skip("TF_ACC_TEST_SSO_GROUP environment variable must be set to run this test")
	}

	steps := []resource.TestStep{
		// Create role assignment for SSO group
		{
			Config: fmt.Sprintf(`
				data "stacklet_sso_group" "test" {
					name = %q
				}

				resource "stacklet_role_assignment" "test" {
					role_name = "viewer"
					principal = data.stacklet_sso_group.test.role_assignment_principal
					target    = "system:all"
				}
			`, testSSOGroup),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "role_name", "viewer"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "principal"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "target", "system:all"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentResource_SSOGroup", steps)
}

func TestAccRoleAssignmentResource_AccountGroup(t *testing.T) {
	testUsername := os.Getenv("TF_ACC_TEST_USERNAME")

	if testUsername == "" {
		t.Skip("TF_ACC_TEST_USERNAME environment variable must be set to run this test")
	}

	userDataSourceConfig := fmt.Sprintf(`username = %q`, testUsername)

	steps := []resource.TestStep{
		// Create role assignment on account group target
		{
			Config: fmt.Sprintf(`
				data "stacklet_user" "test" {
					%s
				}

				resource "stacklet_account_group" "test" {
					name           = "{{.Prefix}}-role-test-group"
					description    = "Test account group for role assignment"
					cloud_provider = "AWS"
					regions        = ["us-east-1"]
				}

				resource "stacklet_role_assignment" "test" {
					role_name = "viewer"
					principal = data.stacklet_user.test.role_assignment_principal
					target    = stacklet_account_group.test.role_assignment_target
				}
			`, userDataSourceConfig),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "role_name", "viewer"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "principal"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "target"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentResource_AccountGroup", steps)
}
