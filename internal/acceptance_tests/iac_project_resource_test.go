// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Collections and groups must exist when recording the test.
func TestAccIaCProjectResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read
		{
			Config: `
resource "stacklet_iac_project" "test" {
  name        = "{{.Prefix}}-project"
  collections = ["collection-one", "collection-two"]
  read_groups = ["group-one"]
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "name", prefixName("project")),
				resource.TestCheckNoResourceAttr("stacklet_iac_project.test", "client_id"),
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "collections.#", "2"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_project.test", "collections.*", "collection-one"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_project.test", "collections.*", "collection-two"),
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "read_groups.#", "1"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_project.test", "read_groups.*", "group-one"),
			),
		},
		// ImportState
		{
			ResourceName:                         "stacklet_iac_project.test",
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateId:                        prefixName("project"),
			ImportStateVerifyIdentifierAttribute: "name",
		},
		// Update
		{
			Config: `
resource "stacklet_iac_project" "test" {
  name        = "{{.Prefix}}-project"
  collections = ["collection-one"]
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "name", prefixName("project")),
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "collections.#", "1"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_project.test", "collections.*", "collection-one"),
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "read_groups.#", "0"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCProjectResource", steps)
}
