// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGCPIntegrationResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read
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
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_gcp_integration.test", "id"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "key", prefixName("gcp-integration")),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.infrastructure.project_id", prefixName("stacklet-project")),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.infrastructure.resource_location", "us-central1"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.infrastructure.resource_prefix", "stacklet"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.organizations.#", "1"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.organizations.0.org_id", "123456789"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.organizations.0.folder_ids.0", "111111111"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.organizations.0.project_ids.0", "my-project-1"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.cost_sources.#", "1"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.cost_sources.0.billing_table", "my-billing-project.my_dataset.gcp_billing_export"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.security_contexts.#", "1"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.security_contexts.0.name", "custom-ctx"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.security_contexts.0.extra_roles.0", "roles/MyRole"),
			),
		},
		// ImportState
		{
			ResourceName:      "stacklet_gcp_integration.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_gcp_integration.test.key"),
			ImportStateVerifyIgnore: []string{
				"access_config_blob_input",
				"customer_config_input",
			},
		},
		// Update customer_config_input
		{
			Config: `
					resource "stacklet_gcp_integration" "test" {
						key = "{{.Prefix}}-gcp-integration"
						customer_config_input = {
							infrastructure = {
								project_id        = "{{.Prefix}}-stacklet-project"
								resource_location = "us-east1"
								resource_prefix   = "stacklet"
							}
							organizations = [
								{
									org_id      = "123456789"
									folder_ids  = ["111111111", "222222222"]
									project_ids = ["my-project-1", "my-project-2"]
								}
							]
							cost_sources = [
								{
									billing_table = "my-billing-project.my_dataset.gcp_billing_export"
								},
								{
									billing_table = "my-billing-project.my_dataset.gcp_billing_export_v2"
								}
							]
							security_contexts = [
								{
									name        = "custom-ctx"
									extra_roles = ["roles/MyRole", "roles/MyOtherRole"]
								}
							]
						}
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.infrastructure.resource_location", "us-east1"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.organizations.0.folder_ids.#", "2"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.organizations.0.folder_ids.1", "222222222"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.organizations.0.project_ids.#", "2"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.organizations.0.project_ids.1", "my-project-2"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.cost_sources.#", "2"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.cost_sources.1.billing_table", "my-billing-project.my_dataset.gcp_billing_export_v2"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.security_contexts.0.extra_roles.#", "2"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.security_contexts.0.extra_roles.0", "roles/MyOtherRole"),
				resource.TestCheckResourceAttr("stacklet_gcp_integration.test", "customer_config.security_contexts.0.extra_roles.1", "roles/MyRole"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccGCPIntegrationResource", steps)
}
