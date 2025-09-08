// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileServiceNowDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_servicenow" "test" {
					endpoint = "https://dev12345.service-now.com"
					username = "test-user"
					password_wo = "test-password"
					password_wo_version = "1"
					issue_type = "incident"
					closed_state = "closed"
				}

				data "stacklet_configuration_profile_servicenow" "test" {
                    depends_on = [stacklet_configuration_profile_servicenow.test]
                }
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_servicenow.test", "id"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_servicenow.test", "profile", "servicenow"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_servicenow.test", "endpoint", "https://dev12345.service-now.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_servicenow.test", "username", "test-user"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_servicenow.test", "password"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_servicenow.test", "issue_type", "incident"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_servicenow.test", "closed_state", "closed"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileServiceNowDataSource", steps)
}
