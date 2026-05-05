// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGCPIntegrationSurfaceDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				data "stacklet_gcp_integration_surface" "test" {}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_gcp_integration_surface.test", "trust_aws.account_id"),
				resource.TestCheckResourceAttrSet("data.stacklet_gcp_integration_surface.test", "trust_aws.assetdb_role_name"),
				resource.TestCheckResourceAttrSet("data.stacklet_gcp_integration_surface.test", "trust_aws.cost_query_role_name"),
				resource.TestCheckResourceAttrSet("data.stacklet_gcp_integration_surface.test", "trust_aws.execution_role_name"),
				resource.TestCheckResourceAttrSet("data.stacklet_gcp_integration_surface.test", "trust_aws.platform_role_name"),
				resource.TestCheckResourceAttrSet("data.stacklet_gcp_integration_surface.test", "aws_relay.bus_arn"),
				resource.TestCheckResourceAttrSet("data.stacklet_gcp_integration_surface.test", "aws_relay.role_arn"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccGCPIntegrationSurfaceDataSource", steps)
}
