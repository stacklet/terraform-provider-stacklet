// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountGroupDataSource(t *testing.T) {
	baseline := `
		resource "stacklet_account_group" "test" {
			name = "{{.Prefix}}-group-ds"
			description = "Test account group"
			cloud_provider = "AWS"
			regions = ["us-east-1"]
		}
	`
	steps := []resource.TestStep{
		{
			Config: baseline + `
				data "stacklet_account_group" "test" {
					name = stacklet_account_group.test.name
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_account_group.test", "name", prefixName("group-ds")),
				resource.TestCheckResourceAttr("data.stacklet_account_group.test", "description", "Test account group"),
				resource.TestCheckResourceAttr("data.stacklet_account_group.test", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttr("data.stacklet_account_group.test", "regions.0", "us-east-1"),
				resource.TestCheckNoResourceAttr("stacklet_account_group.test", "dynamic_filter"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAccountGroupDataSource", steps)
}
