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
						name = "{{.Prefix}}-collection"
						description = "Test policy collection"
						cloud_provider = "AWS"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "name", prefixName("collection")),
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
						name = "{{.Prefix}}-collection-updated"
						description = "Updated policy collection"
						cloud_provider = "AWS"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "name", prefixName("collection-updated")),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "description", "Updated policy collection"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "cloud_provider", "AWS"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccPolicyCollectionResource", steps)
}

func TestAccPolicyCollectionResource_Dynamic(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "{{.Prefix}}-repo"
					}

					resource "stacklet_policy_collection" "test" {
						name = "{{.Prefix}}-collection-dynamic"
                        cloud_provider = "AWS"
						description = "Dynamic policy collection"
                        auto_update = true
						dynamic_config = {
							repository_uuid = stacklet_repository.test.uuid
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "name", prefixName("collection-dynamic")),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "auto_update", "true"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "description", "Dynamic policy collection"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection.test", "dynamic_config.repository_uuid"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection.test", "dynamic_config.namespace"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.branch_name", ""),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.policy_directories.#", "0"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.policy_file_suffixes.#", "2"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.policy_file_suffixes.0", ".yaml"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.policy_file_suffixes.1", ".yml"),
			),
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "{{.Prefix}}-repo"
					}

					resource "stacklet_policy_collection" "test" {
						name = "{{.Prefix}}-collection-dynamic"
						cloud_provider = "AWS"
						description = "Dynamic policy collection updated"
						auto_update = true
						dynamic_config = {
							repository_uuid = stacklet_repository.test.uuid
                            policy_directories = ["dir1", "dir2"]
                            policy_file_suffixes = [".yaml"]
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "name", prefixName("collection-dynamic")),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "auto_update", "true"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "description", "Dynamic policy collection updated"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection.test", "dynamic_config.repository_uuid"),
				resource.TestCheckResourceAttrSet("stacklet_policy_collection.test", "dynamic_config.namespace"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.branch_name", ""),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.policy_directories.#", "2"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.policy_directories.0", "dir1"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.policy_directories.1", "dir2"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.policy_file_suffixes.#", "1"),
				resource.TestCheckResourceAttr("stacklet_policy_collection.test", "dynamic_config.policy_file_suffixes.0", ".yaml"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccPolicyCollectionResource_Dynamic", steps)
}
