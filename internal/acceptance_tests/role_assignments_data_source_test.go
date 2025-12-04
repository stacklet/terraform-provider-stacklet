// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleAssignmentsDataSource_System(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					# Create a role assignment to query
					resource "stacklet_role_assignment" "test" {
						role_name = "admin"

						principal {
							type = "user"
							id   = 789
						}

						target {
							type = "system"
						}
					}

					# Query all system-level role assignments
					data "stacklet_role_assignments" "system" {
						target {
							type = "system"
						}

						depends_on = [stacklet_role_assignment.test]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				// Should have at least one assignment
				resource.TestCheckResourceAttrSet("data.stacklet_role_assignments.system", "assignments.#"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentsDataSource_System", steps)
}

func TestAccRoleAssignmentsDataSource_AccountGroup(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					# Create account group for testing
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-role-assignments-test"
						description = "Test account group for role assignments data source"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					# Create multiple role assignments for the account group
					resource "stacklet_role_assignment" "owner" {
						role_name = "owner"

						principal {
							type = "user"
							id   = 100
						}

						target {
							type = "account-group"
							uuid = stacklet_account_group.test.uuid
						}
					}

					resource "stacklet_role_assignment" "editor" {
						role_name = "editor"

						principal {
							type = "user"
							id   = 101
						}

						target {
							type = "account-group"
							uuid = stacklet_account_group.test.uuid
						}
					}

					# Query all role assignments for the account group
					data "stacklet_role_assignments" "test" {
						target {
							type = "account-group"
							uuid = stacklet_account_group.test.uuid
						}

						depends_on = [
							stacklet_role_assignment.owner,
							stacklet_role_assignment.editor,
						]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				// Should have assignments
				resource.TestCheckResourceAttrSet("data.stacklet_role_assignments.test", "assignments.#"),
				// Target should match our account group
				resource.TestCheckResourceAttr("data.stacklet_role_assignments.test", "target.type", "account-group"),
				resource.TestCheckResourceAttrSet("data.stacklet_role_assignments.test", "target.uuid"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentsDataSource_AccountGroup", steps)
}

func TestAccRoleAssignmentsDataSource_PolicyCollection(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					# Create policy collection for testing
					resource "stacklet_policy_collection" "test" {
						name = "{{.Prefix}}-role-assignments-test"
						description = "Test policy collection for role assignments data source"
					}

					# Create role assignment for the policy collection
					resource "stacklet_role_assignment" "test" {
						role_name = "editor"

						principal {
							type = "sso-group"
							id   = 200
						}

						target {
							type = "policy-collection"
							uuid = stacklet_policy_collection.test.uuid
						}
					}

					# Query role assignments for the policy collection
					data "stacklet_role_assignments" "test" {
						target {
							type = "policy-collection"
							uuid = stacklet_policy_collection.test.uuid
						}

						depends_on = [stacklet_role_assignment.test]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_role_assignments.test", "assignments.#"),
				resource.TestCheckResourceAttr("data.stacklet_role_assignments.test", "target.type", "policy-collection"),
				resource.TestCheckResourceAttrSet("data.stacklet_role_assignments.test", "target.uuid"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentsDataSource_PolicyCollection", steps)
}
