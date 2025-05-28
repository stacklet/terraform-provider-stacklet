// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountGroupDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_account_group" "test" {
						name = "{{.Prefix}}-group-ds"
						description = "Test account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					data "stacklet_account_group" "test" {
						name = stacklet_account_group.test.name
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_account_group.test", "name", prefixName("group-ds")),
				resource.TestCheckResourceAttr("data.stacklet_account_group.test", "description", "Test account group"),
				resource.TestCheckResourceAttr("data.stacklet_account_group.test", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttr("data.stacklet_account_group.test", "regions.0", "us-east-1"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAccountGroupDataSource", steps)
}
