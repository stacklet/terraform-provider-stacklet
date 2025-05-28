// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyCollectionResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
					resource "stacklet_policy_collection" "test" {
						name = "test-collection"
						description = "Test policy collection"
						cloud_provider = "AWS"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "name", "test-collection"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "description", "Test policy collection"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection.test", "id"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_policy_collection.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_policy_collection.test.uuid"),
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_policy_collection" "test" {
						name = "test-collection-updated"
						description = "Updated policy collection"
						cloud_provider = "AWS"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "name", "test-collection-updated"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "description", "Updated policy collection"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "cloud_provider", "AWS"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccPolicyCollectionResource", steps)
}
