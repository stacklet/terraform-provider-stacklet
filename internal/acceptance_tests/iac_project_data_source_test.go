// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Collections and groups must exist when recording the test.
func TestAccIaCProjectDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_project" "test" {
  name        = "{{.Prefix}}-project"
  collections = ["collection-one", "collection-two"]
  read_groups = ["group-one"]
}

data "stacklet_iac_project" "test" {
  name = stacklet_iac_project.test.name
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_iac_project.test", "name", prefixName("project")),
				resource.TestCheckNoResourceAttr("data.stacklet_iac_project.test", "client_id"),
				resource.TestCheckResourceAttr("data.stacklet_iac_project.test", "collections.#", "2"),
				resource.TestCheckTypeSetElemAttr("data.stacklet_iac_project.test", "collections.*", "collection-one"),
				resource.TestCheckTypeSetElemAttr("data.stacklet_iac_project.test", "collections.*", "collection-two"),
				resource.TestCheckResourceAttr("data.stacklet_iac_project.test", "read_groups.#", "1"),
				resource.TestCheckTypeSetElemAttr("data.stacklet_iac_project.test", "read_groups.*", "group-one"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCProjectDataSource", steps)
}
