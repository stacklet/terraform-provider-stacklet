// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileSymphonyDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_symphony" "test" {
					agent_domain = "example.symphony.com"
					service_account = "test-service-account"
					private_key_wo = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC7\n-----END PRIVATE KEY-----"
					private_key_wo_version = "1"
				}

				data "stacklet_configuration_profile_symphony" "test" {
                    depends_on = [stacklet_configuration_profile_symphony.test]
                }
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_symphony.test", "id"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_symphony.test", "profile", "symphony"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_symphony.test", "agent_domain", "example.symphony.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_symphony.test", "service_account", "test-service-account"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_symphony.test", "private_key"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileSymphonyDataSource", steps)
}
