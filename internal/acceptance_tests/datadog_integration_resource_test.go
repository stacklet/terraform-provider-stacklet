// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatadogIntegrationResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing, leaving the site to its default.
		{
			Config: `
				resource "stacklet_datadog_integration" "test" {
					api_key_wo         = "api-key-one"
					api_key_wo_version = "1"
					app_key_wo         = "app-key-one"
					app_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "id", "datadog"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "enabled", "true"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "site", "datadoghq.com"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "configured", "true"),
				// Two details per region while the integration is on, both
				// outstanding until the propagation reports, which rolls up to
				// SKIPPED rather than to a failure.
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "status.status", "SKIPPED"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "status.details.#", "4"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "status.details.0.component", "us-east-1-secret"),
				resource.TestCheckNoResourceAttr("stacklet_datadog_integration.test", "status.details.0.status"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "status.details.1.component", "us-east-1-execution"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "status.details.2.component", "us-west-2-secret"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "status.details.3.component", "us-west-2-execution"),
				// The credentials are write-only, so nothing about them is in state.
				resource.TestCheckNoResourceAttr("stacklet_datadog_integration.test", "api_key_wo"),
				resource.TestCheckNoResourceAttr("stacklet_datadog_integration.test", "app_key_wo"),
			),
		},
		// ImportState testing; the import ID is ignored since the integration is global.
		{
			ResourceName:      "stacklet_datadog_integration.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				// The credentials are never returned, so an import can't know
				// which versions the stored ones correspond to.
				"api_key_wo_version", "app_key_wo_version",
			},
		},
		// Switching the integration off keeps the stored credentials.
		{
			Config: `
				resource "stacklet_datadog_integration" "test" {
					enabled = false

					api_key_wo         = "api-key-one"
					api_key_wo_version = "1"
					app_key_wo         = "app-key-one"
					app_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "enabled", "false"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "configured", "true"),
				// Switched off, a region keeps only its -execution detail, until
				// it confirms it has let go of the credentials.
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "status.details.#", "2"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "status.details.0.component", "us-east-1-execution"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "status.details.1.component", "us-west-2-execution"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccDatadogIntegrationResource", steps)
}

// TestAccDatadogIntegrationResource_KeyRotation checks which credentials are
// sent to the API. The recorded requests are the assertion: a key is sent only
// when its version changes, so an unrelated update leaves the stored one alone.
func TestAccDatadogIntegrationResource_KeyRotation(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_datadog_integration" "test" {
					site = "datadoghq.com"

					api_key_wo         = "api-key-one"
					api_key_wo_version = "1"
					app_key_wo         = "app-key-one"
					app_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "site", "datadoghq.com"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "configured", "true"),
			),
		},
		// Moving to another site with the key values changed but their versions
		// untouched: neither key is sent, and the API rechecks the stored ones
		// against the new site.
		{
			Config: `
				resource "stacklet_datadog_integration" "test" {
					site = "datadoghq.eu"

					api_key_wo         = "api-key-two"
					api_key_wo_version = "1"
					app_key_wo         = "app-key-two"
					app_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "site", "datadoghq.eu"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "configured", "true"),
			),
		},
		// Bumping only the API key version rotates that key on its own.
		{
			Config: `
				resource "stacklet_datadog_integration" "test" {
					site = "datadoghq.eu"

					api_key_wo         = "api-key-two"
					api_key_wo_version = "2"
					app_key_wo         = "app-key-two"
					app_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "api_key_wo_version", "2"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "app_key_wo_version", "1"),
				resource.TestCheckResourceAttr("stacklet_datadog_integration.test", "configured", "true"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccDatadogIntegrationResource_KeyRotation", steps)
}
