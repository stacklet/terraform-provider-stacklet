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
						execution_config = {
							variables = jsonencode({
								environment = "test"
								region = "us-east-1"
							})
						}
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
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.dry_run", "false"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.variables", "{\"environment\":\"test\",\"region\":\"us-east-1\"}"),
				resource.TestCheckNoResourceAttr("stacklet_binding.test", "execution_config.security_context_wo_version"),
				resource.TestCheckNoResourceAttr("stacklet_binding.test", "execution_config.security_context"),
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
						execution_config = {
							dry_run = true
							variables = jsonencode({
								environment = "staging"
								region = "us-west-2"
							})
						}
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
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.dry_run", "true"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.variables", "{\"environment\":\"staging\",\"region\":\"us-west-2\"}"),
				resource.TestCheckNoResourceAttr("stacklet_binding.test", "execution_config.security_context_wo_version"),
				resource.TestCheckNoResourceAttr("stacklet_binding.test", "execution_config.security_context"),
			),
		},
		// Update unsetting execution config
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
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckNoResourceAttr("stacklet_binding.test", "execution_config"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccBindingResource", steps)
}

func TestAccBindingResourceExecutionConfigSecurityContext(t *testing.T) {
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
						execution_config = {
							security_context_wo = "arn:aws:iam::123456789012:role/test-role"
							security_context_wo_version = "1"
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.dry_run", "false"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.security_context", "arn:aws:iam::123456789012:role/test-role"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.security_context_wo_version", "1"),
			),
		},
		// No change in security context version, so value is not updated
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
						execution_config = {
							dry_run = true
							security_context_wo = "arn:aws:iam::123456789012:role/new-role"
							security_context_wo_version = "1"
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.dry_run", "true"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.security_context", "arn:aws:iam::123456789012:role/test-role"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.security_context_wo_version", "1"),
			),
		},
		// Version is updated, value is updated too
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
						execution_config = {
							dry_run = true
							security_context_wo = "arn:aws:iam::123456789012:role/new-role"
							security_context_wo_version = "2"
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.dry_run", "true"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.security_context", "arn:aws:iam::123456789012:role/new-role"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.security_context_wo_version", "2"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccBindingResourceExecutionConfigSecurityContext", steps)
}

func TestAccBindingResourceExecutionConfigResourceLimitsDefault(t *testing.T) {
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
						execution_config = {
							resource_limits = {}
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckNoResourceAttr("stacklet_binding.test", "execution_config.resource_limits.default"),
			),
		},
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
						execution_config = {
							resource_limits = {
								default = {
									max_count = 100
									max_percentage = 20.1
									requires_both = true
								}
							}
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.resource_limits.default.max_count", "100"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.resource_limits.default.max_percentage", "20.1"),
				resource.TestCheckResourceAttr("stacklet_binding.test", "execution_config.resource_limits.default.requires_both", "true"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccBindingResourceExecutionConfigResourceLimitsDefault", steps)
}
