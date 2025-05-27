// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountGroupMappingResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
					resource "stacklet_account" "test1" {
						name = "test-account-1"
						key = "111111111111"
						cloud_provider = "AWS"
						description = "Test AWS account 1"
					}

					resource "stacklet_account" "test2" {
						name = "test-account-2"
						key = "222222222222"
						cloud_provider = "AWS"
						description = "Test AWS account 2"
					}

					resource "stacklet_account_group" "test" {
						name = "test-group-mappings"
						description = "Test account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_account_group_mapping" "test" {
						group_uuid = stacklet_account_group.test.uuid
						account_key = stacklet_account.test1.key
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group_mapping.test", "account_key", "111111111111"),
				resource.TestCheckResourceAttrSet("stacklet_account_group_mapping.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_account_group_mapping.test", "group_uuid"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_account_group_mapping.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs(
				"stacklet_account_group_mapping.test.group_uuid",
				"stacklet_account_group_mapping.test.account_key",
			),
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_account" "test1" {
						name = "test-account-1"
						key = "111111111111"
						cloud_provider = "AWS"
						description = "Test AWS account 1"
					}

					resource "stacklet_account" "test2" {
						name = "test-account-2"
						key = "222222222222"
						cloud_provider = "AWS"
						description = "Test AWS account 2"
					}

					resource "stacklet_account_group" "test" {
						name = "test-group-mappings"
						description = "Test account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_account_group_mapping" "test" {
						group_uuid = stacklet_account_group.test.uuid
						account_key = stacklet_account.test2.key
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account_group_mapping.test", "account_key", "222222222222"),
				resource.TestCheckResourceAttrSet("stacklet_account_group_mapping.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_account_group_mapping.test", "group_uuid"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAccountGroupMappingResource", steps)
}
