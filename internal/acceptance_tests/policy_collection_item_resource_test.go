// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyCollectionMappingResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
					data "stacklet_policy" "policy1" {
						name = "cost-aws:aws-rds-instance-unused-inform"
					}

					data "stacklet_policy" "policy2" {
						name = "cost-aws:aws-redshift-unused-inform"
					}

					resource "stacklet_policy_collection" "test" {
						name = "test-collection-mapping"
						description = "Test policy collection"
						cloud_provider = "AWS"
					}

					resource "stacklet_policy_collection_mapping" "test" {
						collection_uuid = stacklet_policy_collection.test.uuid
						policy_uuid = data.stacklet_policy.policy1.uuid
						policy_version = data.stacklet_policy.policy1.version
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_policy_collection_mapping.test", "collection_uuid"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection_mapping.test", "policy_uuid"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection_mapping.test", "policy_version"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_policy_collection_mapping.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs(
				"stacklet_policy_collection.test.uuid",
				"data.stacklet_policy.policy1.uuid",
			),
		},
		// Update and Read testing
		{
			Config: `
					data "stacklet_policy" "policy1" {
						name = "cost-aws:aws-rds-instance-unused-inform"
					}

					data "stacklet_policy" "policy2" {
						name = "cost-aws:aws-redshift-unused-inform"
					}

					resource "stacklet_policy_collection" "test" {
						name = "test-collection-mapping"
						description = "Test policy collection"
						cloud_provider = "AWS"
					}

					resource "stacklet_policy_collection_mapping" "test" {
						collection_uuid = stacklet_policy_collection.test.uuid
						policy_uuid = data.stacklet_policy.policy2.uuid
						policy_version = data.stacklet_policy.policy2.version
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_policy_collection_mapping.test", "collection_uuid"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection_mapping.test", "policy_uuid"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection_mapping.test", "policy_version"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccPolicyCollectionMappingResource", steps)
}
