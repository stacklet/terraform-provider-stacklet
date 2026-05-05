// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGCPIntegrationDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_gcp_integration" "test" {
						key = "{{.Prefix}}-gcp-integration"
						customer_config_input = {
							infrastructure = {
								project_id        = "{{.Prefix}}-stacklet-project"
								resource_location = "us-central1"
								resource_prefix   = "stacklet"
							}
							organizations = [
								{
									org_id      = "123456789"
									folder_ids  = ["111111111"]
									project_ids = ["my-project-1"]
								}
							]
							cost_sources = [
								{
									billing_table = "my-billing-project.my_dataset.gcp_billing_export"
								}
							]
							security_contexts = [
								{
									name        = "custom-ctx"
									extra_roles = ["roles/MyRole"]
								}
							]
						}
					}

					data "stacklet_gcp_integration" "test" {
						key = stacklet_gcp_integration.test.key
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_gcp_integration.test", "id"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "key", prefixName("gcp-integration")),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.infrastructure.project_id", prefixName("stacklet-project")),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.infrastructure.resource_location", "us-central1"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.infrastructure.resource_prefix", "stacklet"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.organizations.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.organizations.0.org_id", "123456789"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.organizations.0.folder_ids.0", "111111111"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.organizations.0.project_ids.0", "my-project-1"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.cost_sources.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.cost_sources.0.billing_table", "my-billing-project.my_dataset.gcp_billing_export"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.security_contexts.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.security_contexts.0.name", "custom-ctx"),
				resource.TestCheckResourceAttr("data.stacklet_gcp_integration.test", "customer_config.security_contexts.0.extra_roles.0", "roles/MyRole"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccGCPIntegrationDataSource", steps)
}
