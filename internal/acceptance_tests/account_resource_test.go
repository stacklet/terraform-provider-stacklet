// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
					resource "stacklet_account" "test" {
						name = "test-account"
						key = "999999999999"
						cloud_provider = "AWS"
						description = "Test AWS account"
						short_name = "test"
						email = "test@example.com"
						variables = jsonencode({
                            environment = "test"
                        })
                        security_context_wo = "arn:aws:iam::123456789012:role/stacklet-execution"
                        security_context_wo_version = "1"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account.test", "name", "test-account"),
				resource.TestCheckResourceAttr("stacklet_account.test", "key", "999999999999"),
				resource.TestCheckResourceAttr("stacklet_account.test", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttr("stacklet_account.test", "description", "Test AWS account"),
				resource.TestCheckResourceAttr("stacklet_account.test", "short_name", "test"),
				resource.TestCheckResourceAttr("stacklet_account.test", "email", "test@example.com"),
				resource.TestCheckResourceAttr("stacklet_account.test", "variables", "{\"environment\":\"test\"}"),
				// For AWS, passing a role ARN as security_context_wo passes it through as a security_context
				resource.TestCheckResourceAttr("stacklet_account.test", "security_context", "arn:aws:iam::123456789012:role/stacklet-execution"),
				resource.TestCheckResourceAttrSet("stacklet_account.test", "id"),
			),
		},
		// ImportState testing
		{
			ResourceName:            "stacklet_account.test",
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateId:           "AWS:999999999999",
			ImportStateVerifyIgnore: []string{"security_context"},
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_account" "test" {
						name = "test-account-updated"
						key = "999999999999"
						cloud_provider = "AWS"
						description = "Updated AWS account"
						short_name = "test-updated"
						email = "updated@example.com"
						variables = jsonencode({
                            environment = "staging"
                        })
                        security_context_wo = "arn:aws:iam::123456789012:role/stacklet-execution-new"
                        security_context_wo_version = "2"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_account.test", "name", "test-account-updated"),
				resource.TestCheckResourceAttr("stacklet_account.test", "key", "999999999999"),
				resource.TestCheckResourceAttr("stacklet_account.test", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttr("stacklet_account.test", "description", "Updated AWS account"),
				resource.TestCheckResourceAttr("stacklet_account.test", "short_name", "test-updated"),
				resource.TestCheckResourceAttr("stacklet_account.test", "email", "updated@example.com"),
				resource.TestCheckResourceAttr("stacklet_account.test", "variables", "{\"environment\":\"staging\"}"),
				resource.TestCheckResourceAttr("stacklet_account.test", "security_context", "arn:aws:iam::123456789012:role/stacklet-execution-new"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAccountResource", steps)
}
