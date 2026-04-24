// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRoleAssignmentResource_AccountGroup(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing - role assignment on account group
		{
			Config: `
				resource "stacklet_user" "test" {
					name     = "{{.Prefix}}-role-user"
					username = "{{.Prefix}}_role_user"
					email    = "test@stacklet.io"
				}

				resource "stacklet_account_group" "test" {
					name           = "{{.Prefix}}-role-assignment-test"
					description    = "Test account group for role assignment"
					cloud_provider = "AWS"
					regions        = ["us-east-1"]
				}

				resource "stacklet_role_assignment" "test" {
					role_name = "viewer"
					principal = stacklet_user.test.role_assignment_principal
					target    = stacklet_account_group.test.role_assignment_target
				}
			`,
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
			Config: `
				resource "stacklet_user" "test" {
					name     = "{{.Prefix}}-role-user"
					username = "{{.Prefix}}_role_user"
					email    = "test@stacklet.io"
				}

				resource "stacklet_account_group" "test" {
					name           = "{{.Prefix}}-role-assignment-test"
					description    = "Test account group for role assignment"
					cloud_provider = "AWS"
					regions        = ["us-east-1"]
				}

				resource "stacklet_role_assignment" "test" {
					role_name = "editor"
					principal = stacklet_user.test.role_assignment_principal
					target    = stacklet_account_group.test.role_assignment_target
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "role_name", "editor"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "principal"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "target"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentResource_AccountGroup", steps)
}

func TestAccRoleAssignmentResource_SSOGroup(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_sso_group" "test" {
					name = "{{.Prefix}}-sso-role-group"
				}

				resource "stacklet_role_assignment" "test" {
					role_name = "viewer"
					principal = stacklet_sso_group.test.role_assignment_principal
					target    = "system:all"
				}
			`,
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
