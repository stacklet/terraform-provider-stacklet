// Copyright (c) 2026 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccAPIKeyResource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_api_key" "test" {
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "identity"),
				resource.TestCheckResourceAttr("stacklet_api_key.test", "description", ""),
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "secret"),
				resource.TestCheckNoResourceAttr("stacklet_api_key.test", "expires_at"),
				resource.TestCheckNoResourceAttr("stacklet_api_key.test", "revoked_at"),
			),
		},
		// Update description
		{
			Config: `
					resource "stacklet_api_key" "test" {
						description = "Updated API key"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "identity"),
				resource.TestCheckResourceAttr("stacklet_api_key.test", "description", "Updated API key"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_api_key.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_api_key.test.identity"),
			ImportStateVerifyIgnore: []string{
				"secret",
			},
		},
	}
	runRecordedAccTest(t, "TestAccAPIKeyResource", steps)
}

func TestAccAPIKeyResource_ExpirationForceReplace(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_api_key" "test" {
						description = "API key without expiration"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "identity"),
				resource.TestCheckResourceAttr("stacklet_api_key.test", "description", "API key without expiration"),
				resource.TestCheckNoResourceAttr("stacklet_api_key.test", "expires_at"),
			),
		},
		// Adding expiration should force replacement
		{
			Config: `
					resource "stacklet_api_key" "test" {
						description = "API key with expiration"
						expires_at = "2026-12-31T23:59:59Z"
					}
				`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_api_key.test", plancheck.ResourceActionReplace),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "identity"),
				resource.TestCheckResourceAttr("stacklet_api_key.test", "description", "API key with expiration"),
				resource.TestCheckResourceAttr("stacklet_api_key.test", "expires_at", "2026-12-31T23:59:59Z"),
			),
		},
		// Changing expiration should also force replacement
		{
			Config: `
					resource "stacklet_api_key" "test" {
						description = "API key with different expiration"
						expires_at = "2027-06-30T23:59:59+00:00"
					}
				`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_api_key.test", plancheck.ResourceActionReplace),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "identity"),
				resource.TestCheckResourceAttr("stacklet_api_key.test", "description", "API key with different expiration"),
				resource.TestCheckResourceAttr("stacklet_api_key.test", "expires_at", "2027-06-30T23:59:59+00:00"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAPIKeyResource_ExpirationForceReplace", steps)
}
