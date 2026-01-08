// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleAssignmentsDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				data "stacklet_role_assignments" "system" {
				  target = "system:all"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check that assignments list is present (may be empty)
				resource.TestCheckResourceAttrSet("data.stacklet_role_assignments.system", "assignments.#"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRoleAssignmentsDataSource", steps)
}
