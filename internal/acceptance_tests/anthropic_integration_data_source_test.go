// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAnthropicIntegrationDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_anthropic_integration" "test" {
					org_id = "` + testAnthropicOrgID + `"

					admin_key_wo         = "` + testAnthropicAdminKey + `"
					admin_key_wo_version = "1"
				}

				data "stacklet_anthropic_integration" "test" {
					org_id = stacklet_anthropic_integration.test.org_id
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPair("data.stacklet_anthropic_integration.test", "id", "stacklet_anthropic_integration.test", "id"),
				resource.TestCheckResourceAttr("data.stacklet_anthropic_integration.test", "org_id", testAnthropicOrgID),
				resource.TestCheckResourceAttr("data.stacklet_anthropic_integration.test", "configured", "true"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAnthropicIntegrationDataSource", steps)
}
