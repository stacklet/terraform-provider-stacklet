// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMSTeamsIntegrationSurfaceDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					data "stacklet_msteams_integration_surface" "test" {}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_msteams_integration_surface.test", "bot_endpoint"),
				resource.TestCheckResourceAttrSet("data.stacklet_msteams_integration_surface.test", "wif_issuer_url"),
				resource.TestCheckResourceAttrSet("data.stacklet_msteams_integration_surface.test", "trust_role_arn"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccMSTeamsIntegrationSurfaceDataSource", steps)
}
