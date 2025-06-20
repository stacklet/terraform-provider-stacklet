// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBindingDataSource(t *testing.T) {
	steps := []resource.TestStep{
		// Create a binding to test the data source
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-binding-ds-group"
						description = "Test account group for binding data source"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_policy_collection" "test" {
						name = "{{.Prefix}}-binding-ds-collection"
						description = "Test policy collection for binding data source"
						cloud_provider = "AWS"
					}

					resource "stacklet_binding" "test" {
						name = "{{.Prefix}}-binding-ds"
						description = "Test binding for data source"
						account_group_uuid = stacklet_account_group.test.uuid
						policy_collection_uuid = stacklet_policy_collection.test.uuid
						auto_deploy = true
						schedule = "rate(1 hour)"
						dry_run = true
						resource_limits = {
							max_count = 10
							max_percentage = 20
							requires_both = true
						}
						policy_resource_limit {
							policy_name = "policy"
							max_count = 90
							max_percentage = 50.0
							requires_both = true
						}
						variables = jsonencode({
							environment = "test"
							region = "us-east-1"
						})
					}

					# Test lookup by name
					data "stacklet_binding" "by_name" {
						name = stacklet_binding.test.name
					}

					# Test lookup by UUID
					data "stacklet_binding" "by_uuid" {
						uuid = stacklet_binding.test.uuid
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				// Verify lookup by name
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "name", prefixName("binding-ds")),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "description", "Test binding for data source"),
				resource.TestCheckResourceAttrSet("data.stacklet_binding.by_name", "account_group_uuid"),
				resource.TestCheckResourceAttrSet("data.stacklet_binding.by_name", "policy_collection_uuid"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "auto_deploy", "true"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "schedule", "rate(1 hour)"),
				resource.TestCheckResourceAttrSet("data.stacklet_binding.by_name", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_binding.by_name", "uuid"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "dry_run", "true"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "resource_limits.max_count", "10"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "resource_limits.max_percentage", "20"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "resource_limits.requires_both", "true"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "policy_resource_limit.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "policy_resource_limit.0.policy_name", "policy"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "policy_resource_limit.0.max_count", "90"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "policy_resource_limit.0.max_percentage", "50"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "policy_resource_limit.0.requires_both", "true"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "variables", "{\"environment\":\"test\",\"region\":\"us-east-1\"}"),
				resource.TestCheckNoResourceAttr("data.stacklet_binding.by_name", "security_context"),

				// Verify lookup by UUID
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "name", prefixName("binding-ds")),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "description", "Test binding for data source"),
				resource.TestCheckResourceAttrSet("data.stacklet_binding.by_uuid", "account_group_uuid"),
				resource.TestCheckResourceAttrSet("data.stacklet_binding.by_uuid", "policy_collection_uuid"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "auto_deploy", "true"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "schedule", "rate(1 hour)"),
				resource.TestCheckResourceAttrSet("data.stacklet_binding.by_uuid", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_binding.by_uuid", "uuid"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "dry_run", "true"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "resource_limits.max_count", "10"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "resource_limits.max_percentage", "20"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "resource_limits.requires_both", "true"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "policy_resource_limit.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "policy_resource_limit.0.policy_name", "policy"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "policy_resource_limit.0.max_count", "90"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "policy_resource_limit.0.max_percentage", "50"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "policy_resource_limit.0.requires_both", "true"),
				resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "variables", "{\"environment\":\"test\",\"region\":\"us-east-1\"}"),
				resource.TestCheckNoResourceAttr("data.stacklet_binding.by_uuid", "security_context"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccBindingDataSource", steps)
}
