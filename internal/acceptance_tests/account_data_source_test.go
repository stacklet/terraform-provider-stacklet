// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountDataSource(t *testing.T) {
	steps := []resource.TestStep{
		// Create a resource first
		{
			Config: `
					resource "stacklet_account" "test" {
						name = "{{.Prefix}}-account-ds"
						key = "999999999999"
						cloud_provider = "AWS"
						description = "Test AWS account"
						short_name = "{{.Prefix}}-account-ds"
						email = "test@example.com"
						variables = jsonencode({
                            environment = "test"
                        })
					}
				`,
		},
		// Read testing
		{
			Config: `
					resource "stacklet_account" "test" {
						name = "{{.Prefix}}-account-ds"
						key = "999999999999"
						cloud_provider = "AWS"
						description = "Test AWS account"
						short_name = "{{.Prefix}}-account-ds"
						email = "test@example.com"
						variables = jsonencode({
                            environment = "test"
                        })
					}

					data "stacklet_account" "test" {
						key = stacklet_account.test.key
						cloud_provider = stacklet_account.test.cloud_provider
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_account.test", "name", prefixName("account-ds")),
				resource.TestCheckResourceAttr("data.stacklet_account.test", "key", "999999999999"),
				resource.TestCheckResourceAttr("data.stacklet_account.test", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttr("data.stacklet_account.test", "description", "Test AWS account"),
				resource.TestCheckResourceAttr("data.stacklet_account.test", "short_name", prefixName("account-ds")),
				resource.TestCheckResourceAttr("data.stacklet_account.test", "email", "test@example.com"),
				resource.TestCheckResourceAttr("data.stacklet_account.test", "variables", "{\"environment\":\"test\"}"),
				resource.TestCheckResourceAttrSet("data.stacklet_account.test", "id"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAccountDataSource", steps)
}
