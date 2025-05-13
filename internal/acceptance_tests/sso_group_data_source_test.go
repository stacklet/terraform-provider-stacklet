package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSOGroupDataSource(t *testing.T) {
	steps := []resource.TestStep{
		// Create a resource first
		{
			Config: `
					resource "stacklet_sso_group" "test" {
						name = "test-group-ds"
						roles = ["admin"]
						account_group_uuids = ["e2e040cf-6f10-4cdf-94b1-9600b2ee36ca"]
					}
				`,
		},
		// Read testing by name
		{
			Config: `
					resource "stacklet_sso_group" "test" {
						name = "test-group-ds"
						roles = ["admin"]
						account_group_uuids = ["e2e040cf-6f10-4cdf-94b1-9600b2ee36ca"]
					}

					data "stacklet_sso_group" "test" {
						name = stacklet_sso_group.test.name
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_sso_group.test", "name", "test-group-ds"),
				resource.TestCheckResourceAttr("data.stacklet_sso_group.test", "roles.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_sso_group.test", "roles.0", "admin"),
				resource.TestCheckResourceAttr("data.stacklet_sso_group.test", "account_group_uuids.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_sso_group.test", "account_group_uuids.0", "e2e040cf-6f10-4cdf-94b1-9600b2ee36ca"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccSSOGroupDataSource", steps)
}
