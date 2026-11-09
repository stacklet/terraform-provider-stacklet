// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIaCProjectDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_collection" "one" {
  name = "{{.Prefix}}-collection-one"
}

resource "stacklet_iac_policy_collection" "two" {
  name = "{{.Prefix}}-collection-two"
}

resource "stacklet_iac_project" "test" {
  name        = "{{.Prefix}}-project"
  collections = [stacklet_iac_policy_collection.one.name, stacklet_iac_policy_collection.two.name]
}

data "stacklet_iac_project" "test" {
  name = stacklet_iac_project.test.name
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_iac_project.test", "name", prefixName("project")),
				resource.TestCheckNoResourceAttr("data.stacklet_iac_project.test", "client_id"),
				resource.TestCheckResourceAttr("data.stacklet_iac_project.test", "collections.#", "2"),
				resource.TestCheckTypeSetElemAttr("data.stacklet_iac_project.test", "collections.*", prefixName("collection-one")),
				resource.TestCheckTypeSetElemAttr("data.stacklet_iac_project.test", "collections.*", prefixName("collection-two")),
				resource.TestCheckResourceAttr("data.stacklet_iac_project.test", "read_groups.#", "0"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCProjectDataSource", steps)
}
