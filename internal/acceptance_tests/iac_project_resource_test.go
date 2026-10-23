// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIaCProjectResource(t *testing.T) {
	collections := `
resource "stacklet_iac_policy_collection" "one" {
  name = "{{.Prefix}}-collection-one"
}

resource "stacklet_iac_policy_collection" "two" {
  name = "{{.Prefix}}-collection-two"
}
`
	steps := []resource.TestStep{
		// Create and Read
		{
			Config: collections + `
resource "stacklet_iac_project" "test" {
  name        = "{{.Prefix}}-project"
  collections = [stacklet_iac_policy_collection.one.name, stacklet_iac_policy_collection.two.name]
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "name", prefixName("project")),
				resource.TestCheckNoResourceAttr("stacklet_iac_project.test", "client_id"),
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "collections.#", "2"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_project.test", "collections.*", prefixName("collection-one")),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_project.test", "collections.*", prefixName("collection-two")),
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "read_groups.#", "0"),
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
			Config: collections + `
resource "stacklet_iac_project" "test" {
  name        = "{{.Prefix}}-project"
  collections = [stacklet_iac_policy_collection.one.name]
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "name", prefixName("project")),
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "collections.#", "1"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_project.test", "collections.*", prefixName("collection-one")),
				resource.TestCheckResourceAttr("stacklet_iac_project.test", "read_groups.#", "0"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCProjectResource", steps)
}
