// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBindingExecutionConfigResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-binding-group"
						description = "Test account group for binding"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_policy_collection" "test" {
						name = "{{.Prefix}}-binding-collection"
						description = "Test policy collection for binding"
						cloud_provider = "AWS"
					}

					resource "stacklet_binding" "test" {
						name = "{{.Prefix}}-binding"
						description = "Test binding"
						account_group_uuid = stacklet_account_group.test.uuid
						policy_collection_uuid = stacklet_policy_collection.test.uuid
					}

					resource "stacklet_binding_execution_config" "test" {
						binding_uuid = stacklet_binding.test.uuid
						dry_run = true
						variables = jsonencode({env="staging"})
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_binding_execution_config.test", "binding_uuid"),
				resource.TestCheckResourceAttr("stacklet_binding_execution_config.test", "dry_run", "true"),
				resource.TestCheckResourceAttr("stacklet_binding_execution_config.test", "variables", "{\"env\":\"staging\"}"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_binding_execution_config.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_binding.test.uuid"),
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-binding-group"
						description = "Test account group for binding"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_policy_collection" "test" {
						name = "{{.Prefix}}-binding-collection"
						description = "Test policy collection for binding"
						cloud_provider = "AWS"
					}

					resource "stacklet_binding" "test" {
						name = "{{.Prefix}}-binding"
						description = "Test binding"
						account_group_uuid = stacklet_account_group.test.uuid
						policy_collection_uuid = stacklet_policy_collection.test.uuid
					}

					resource "stacklet_binding_execution_config" "test" {
						binding_uuid = stacklet_binding.test.uuid
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_binding_execution_config.test", "binding_uuid"),
				resource.TestCheckResourceAttr("stacklet_binding_execution_config.test", "dry_run", "false"),
				resource.TestCheckNoResourceAttr("stacklet_binding_execution_config.test", "variables"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccBindingExecutionConfigResource", steps)
}
