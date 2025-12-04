// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleAssignmentResource_System(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing - system target
		{
			Config: `
					resource "stacklet_role_assignment" "test" {
						role_name = "admin"

						principal {
							type = "user"
							id   = 123
						}

						target {
							type = "system"
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "role_name", "admin"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "principal.type", "user"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "principal.id", "123"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "target.type", "system"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "id"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_role_assignment.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_role_assignment.test.id"),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentResource_System", steps)
}

func TestAccRoleAssignmentResource_AccountGroup(t *testing.T) {
	steps := []resource.TestStep{
		// Create account group for testing
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-role-test"
						description = "Test account group for role assignment"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_role_assignment" "test" {
						role_name = "owner"

						principal {
							type = "user"
							id   = 456
						}

						target {
							type = "account-group"
							uuid = stacklet_account_group.test.uuid
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "role_name", "owner"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "principal.type", "user"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "principal.id", "456"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "target.type", "account-group"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "target.uuid"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "id"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_role_assignment.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_role_assignment.test.id"),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentResource_AccountGroup", steps)
}

func TestAccRoleAssignmentResource_PolicyCollection(t *testing.T) {
	steps := []resource.TestStep{
		// Create policy collection for testing
		{
			Config: `
					resource "stacklet_policy_collection" "test" {
						name = "{{.Prefix}}-role-test"
						description = "Test policy collection for role assignment"
						cloud_provider = "AWS"
					}

					resource "stacklet_role_assignment" "test" {
						role_name = "editor"

						principal {
							type = "sso-group"
							id   = 789
						}

						target {
							type = "policy-collection"
							uuid = stacklet_policy_collection.test.uuid
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "role_name", "editor"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "principal.type", "sso-group"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "principal.id", "789"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "target.type", "policy-collection"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "target.uuid"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "id"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_role_assignment.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_role_assignment.test.id"),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentResource_PolicyCollection", steps)
}

func TestAccRoleAssignmentResource_Repository(t *testing.T) {
	steps := []resource.TestStep{
		// Create repository for testing
		{
			Config: `
					resource "stacklet_repository" "test" {
						name = "{{.Prefix}}-role-test"
						description = "Test repository for role assignment"
						provider_type = "github"
						url = "https://github.com/test/repo"
					}

					resource "stacklet_role_assignment" "test" {
						role_name = "viewer"

						principal {
							type = "user"
							id   = 101
						}

						target {
							type = "repository"
							uuid = stacklet_repository.test.uuid
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "role_name", "viewer"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "principal.type", "user"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "principal.id", "101"),
				resource.TestCheckResourceAttr("stacklet_role_assignment.test", "target.type", "repository"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "target.uuid"),
				resource.TestCheckResourceAttrSet("stacklet_role_assignment.test", "id"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_role_assignment.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_role_assignment.test.id"),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentResource_Repository", steps)
}
