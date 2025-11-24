// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccAccountGroupResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-group"
						description = "Test account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group.test", "name", prefixName("group")),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "description", "Test account group"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "regions.0", "us-east-1"),
				resource.TestCheckResourceAttrSet("stacklet_account_group.test", "uuid"),
				resource.TestCheckNoResourceAttr("stacklet_account_group.test", "dynamic_filter"),
			),
		},
		// ImportState testing using UUID
		{
			ResourceName:      "stacklet_account_group.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_account_group.test.uuid"),
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-group-updated"
						description = "Updated account group"
						cloud_provider = "AWS"
						regions = ["us-east-1", "us-east-2"]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group.test", "name", prefixName("group-updated")),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "description", "Updated account group"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "regions.0", "us-east-1"),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "regions.1", "us-east-2"),
				resource.TestCheckResourceAttrSet("stacklet_account_group.test", "uuid"),
			),
		},
		//
		{
			Config: `
					resource "stacklet_account_group" "different" {
						name = "{{.Prefix}}-another"
						description = "Different account group"
						cloud_provider = "Azure"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group.different", "name", prefixName("another")),
				resource.TestCheckResourceAttr("stacklet_account_group.different", "description", "Different account group"),
				resource.TestCheckResourceAttr("stacklet_account_group.different", "cloud_provider", "Azure"),
				resource.TestCheckResourceAttr("stacklet_account_group.different", "regions.#", "0"),
				resource.TestCheckResourceAttrSet("stacklet_account_group.different", "uuid"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAccountGroupResource", steps)
}

func TestAccAccountGroupResource_Dynamic(t *testing.T) {
	steps := []resource.TestStep{
		// Create account group without dynamic filter
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-dynamic-group"
						description = "Test dynamic account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group.test", "name", prefixName("dynamic-group")),
				resource.TestCheckNoResourceAttr("stacklet_account_group.test", "dynamic_filter"),
			),
		},
		// Add dynamic filter - should trigger recreation
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-dynamic-group"
						description = "Test dynamic account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
						dynamic_filter = "tag:Environment=prod"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group.test", "name", prefixName("dynamic-group")),
				resource.TestCheckResourceAttr("stacklet_account_group.test", "dynamic_filter", "tag:Environment=prod"),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_account_group.test", plancheck.ResourceActionDestroyBeforeCreate),
				},
			},
		},
		// Remove dynamic filter - should trigger recreation
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-dynamic-group"
						description = "Test dynamic account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group.test", "name", prefixName("dynamic-group")),
				resource.TestCheckNoResourceAttr("stacklet_account_group.test", "dynamic_filter"),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_account_group.test", plancheck.ResourceActionDestroyBeforeCreate),
				},
			},
		},
	}
	runRecordedAccTest(t, "TestAccAccountGroupResource_Dynamic", steps)
}
