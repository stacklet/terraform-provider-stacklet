package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountGroupResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "test-group"
						description = "Test account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group.test", "name", "test-group"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "description", "Test account group"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "regions.0", "us-east-1"),
				resource.TestCheckResourceAttrSet("stacklet_account_group.test", "uuid"),
			),
		},
		// ImportState testing using UUID
		{
			ResourceName:      "stacklet_account_group.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateId:     "d9784826-dba3-4bb1-8df3-3dd60c8983e1",
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "test-group-updated"
						description = "Updated account group"
						cloud_provider = "AWS"
						regions = ["us-east-1", "us-east-2"]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group.test", "name", "test-group-updated"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "description", "Updated account group"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "regions.0", "us-east-1"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "regions.1", "us-east-2"),
				resource.TestCheckResourceAttrSet("stacklet_account_group.test", "uuid"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAccountGroupResource", steps)
}
