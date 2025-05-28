// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBindingResource(t *testing.T) {
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
						auto_deploy = true
						schedule = "rate(1 hour)"
						variables = jsonencode({
							environment = "test"
							region = "us-east-1"
						})
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_binding.test", "name", prefixName("binding")),
				resource.TestCheckResourceAttr("stacklet_binding.test", "description", "Test binding"),
				resource.TestCheckResourceAttrSet("stacklet_binding.test", "account_group_uuid"),
				resource.TestCheckResourceAttrSet("stacklet_binding.test", "policy_collection_uuid"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "auto_deploy", "true"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "schedule", "rate(1 hour)"),
				resource.TestCheckResourceAttrSet("stacklet_binding.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_binding.test", "uuid"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_binding.test",
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
						regions = ["us-east-1", "us-east-2"]
					}

					resource "stacklet_policy_collection" "test" {
						name = "{{.Prefix}}-binding-collection"
						description = "Test policy collection for binding"
						cloud_provider = "AWS"
					}

					resource "stacklet_binding" "test" {
						name = "{{.Prefix}}-binding-updated"
						description = "Updated test binding"
						account_group_uuid = stacklet_account_group.test.uuid
						policy_collection_uuid = stacklet_policy_collection.test.uuid
						auto_deploy = false
						schedule = "rate(2 hours)"
						variables = jsonencode({
							environment = "staging"
							region = "us-west-2"
						})
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_binding.test", "name", prefixName("binding-updated")),
				resource.TestCheckResourceAttr("stacklet_binding.test", "description", "Updated test binding"),
				resource.TestCheckResourceAttrSet("stacklet_binding.test", "account_group_uuid"),
				resource.TestCheckResourceAttrSet("stacklet_binding.test", "policy_collection_uuid"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "auto_deploy", "false"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "schedule", "rate(2 hours)"),
				resource.TestCheckResourceAttrSet("stacklet_binding.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_binding.test", "uuid"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccBindingResource", steps)
}
