// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The recordings come from the platform's local development stack, whose fake
// Anthropic Admin API accepts these keys. testAnthropicAdminKey belongs to
// testAnthropicOrgID, and otherAnthropicAdminKey to another organization.
const (
	testAnthropicOrgID     = "org-functest"
	testAnthropicAdminKey  = "sk-ant-admin01-fake-test"
	otherAnthropicAdminKey = "sk-ant-admin01-fake"
)

func TestAccAnthropicIntegrationResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
				resource "stacklet_anthropic_integration" "test" {
					org_id = "` + testAnthropicOrgID + `"

					admin_key_wo         = "` + testAnthropicAdminKey + `"
					admin_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_anthropic_integration.test", "id"),
				resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "org_id", testAnthropicOrgID),
				resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "configured", "true"),
				resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "admin_key_wo_version", "1"),
				// The key is write-only, so it isn't in state.
				resource.TestCheckNoResourceAttr("stacklet_anthropic_integration.test", "admin_key_wo"),
			),
		},
		// ImportState testing
		{
			ResourceName:                         "stacklet_anthropic_integration.test",
			ImportState:                          true,
			ImportStateId:                        testAnthropicOrgID,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "org_id",
			ImportStateVerifyIgnore: []string{
				// The key is never returned, so an import can't know which
				// version the stored one corresponds to.
				"admin_key_wo_version",
			},
		},
	}
	runRecordedAccTest(t, "TestAccAnthropicIntegrationResource", steps)
}

// TestAccAnthropicIntegrationResource_KeyRotation checks that bumping the
// version sends the key again, which the API checks with Anthropic before it
// replaces the stored one.
func TestAccAnthropicIntegrationResource_KeyRotation(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_anthropic_integration" "test" {
					org_id = "` + testAnthropicOrgID + `"

					admin_key_wo         = "` + testAnthropicAdminKey + `"
					admin_key_wo_version = "1"
				}
			`,
			Check: resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "configured", "true"),
		},
		{
			Config: `
				resource "stacklet_anthropic_integration" "test" {
					org_id = "` + testAnthropicOrgID + `"

					admin_key_wo         = "` + testAnthropicAdminKey + `"
					admin_key_wo_version = "2"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "admin_key_wo_version", "2"),
				resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "configured", "true"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAnthropicIntegrationResource_KeyRotation", steps)
}

// TestAccAnthropicIntegrationResource_RejectedKey checks that a key the API
// rejects when it checks it with Anthropic fails the apply with the reason.
func TestAccAnthropicIntegrationResource_RejectedKey(t *testing.T) {
	steps := []resource.TestStep{
		// A key that doesn't authenticate.
		{
			Config: `
				resource "stacklet_anthropic_integration" "test" {
					org_id = "` + testAnthropicOrgID + `"

					admin_key_wo         = "sk-ant-admin01-invalid"
					admin_key_wo_version = "1"
				}
			`,
			ExpectError: regexp.MustCompile(`Failed to connect to Anthropic`),
		},
		// Another organization's key.
		{
			Config: `
				resource "stacklet_anthropic_integration" "test" {
					org_id = "` + testAnthropicOrgID + `"

					admin_key_wo         = "` + otherAnthropicAdminKey + `"
					admin_key_wo_version = "1"
				}
			`,
			ExpectError: regexp.MustCompile(`The admin key belongs to Anthropic organization`),
		},
		{
			Config: `
				resource "stacklet_anthropic_integration" "test" {
					org_id = "` + testAnthropicOrgID + `"

					admin_key_wo         = "` + testAnthropicAdminKey + `"
					admin_key_wo_version = "1"
				}
			`,
			Check: resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "configured", "true"),
		},
		// A rotation to a bad key fails too, and leaves the stored key alone.
		{
			Config: `
				resource "stacklet_anthropic_integration" "test" {
					org_id = "` + testAnthropicOrgID + `"

					admin_key_wo         = "sk-ant-admin01-invalid"
					admin_key_wo_version = "2"
				}
			`,
			ExpectError: regexp.MustCompile(`Failed to connect to Anthropic`),
		},
		// The failed rotation left state at the first version, and the API
		// still holds the first key.
		{
			Config: `
				resource "stacklet_anthropic_integration" "test" {
					org_id = "` + testAnthropicOrgID + `"

					admin_key_wo         = "` + testAnthropicAdminKey + `"
					admin_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "org_id", testAnthropicOrgID),
				resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "configured", "true"),
				resource.TestCheckResourceAttr("stacklet_anthropic_integration.test", "admin_key_wo_version", "1"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAnthropicIntegrationResource_RejectedKey", steps)
}
