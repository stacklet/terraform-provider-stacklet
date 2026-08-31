// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatadogIntegrationDataSource(t *testing.T) {
	baseline := `
		resource "stacklet_datadog_integration" "test" {
			site = "datadoghq.eu"

			api_key_wo         = "api-key-one"
			api_key_wo_version = "1"
			app_key_wo         = "app-key-one"
			app_key_wo_version = "1"
		}
	`
	steps := []resource.TestStep{
		{
			Config: baseline + `
				data "stacklet_datadog_integration" "test" {
					depends_on = [stacklet_datadog_integration.test]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_datadog_integration.test", "id", "datadog"),
				resource.TestCheckResourceAttr("data.stacklet_datadog_integration.test", "enabled", "true"),
				resource.TestCheckResourceAttr("data.stacklet_datadog_integration.test", "site", "datadoghq.eu"),
				resource.TestCheckResourceAttr("data.stacklet_datadog_integration.test", "configured", "true"),
				resource.TestCheckResourceAttrSet("data.stacklet_datadog_integration.test", "status.status"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccDatadogIntegrationDataSource", steps)
}
