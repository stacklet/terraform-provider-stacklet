// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyCollectionDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_policy_collection" "test" {
						name = "test-collection-ds"
						description = "Test policy collection"
						cloud_provider = "AWS"
					}

					data "stacklet_policy_collection" "test" {
						name = stacklet_policy_collection.test.name
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_policy_collection.test", "name", "test-collection-ds"),
				resource.TestCheckResourceAttr("data.stacklet_policy_collection.test", "description", "Test policy collection"),
				resource.TestCheckResourceAttr("data.stacklet_policy_collection.test", "cloud_provider", "AWS"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccPolicyCollectionDataSource", steps)
}
