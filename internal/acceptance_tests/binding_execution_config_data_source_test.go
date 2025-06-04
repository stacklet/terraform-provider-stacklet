// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBindingExecutionConfigDataSource(t *testing.T) {
	steps := []resource.TestStep{
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

					data "stacklet_binding_execution_config" "test" {
						binding_uuid = stacklet_binding.test.uuid

						# ensure the resource is created first
						depends_on = [stacklet_binding_execution_config.test]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_binding_execution_config.test", "binding_uuid"),
				resource.TestCheckResourceAttr("data.stacklet_binding_execution_config.test", "dry_run", "true"),
				resource.TestCheckResourceAttr("data.stacklet_binding_execution_config.test", "variables", "{\"env\":\"staging\"}"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccBindingExecutionConfigDataSource", steps)
}
