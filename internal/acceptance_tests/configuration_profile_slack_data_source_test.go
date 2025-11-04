// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileSlackDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_slack" "test" {
 	                user_fields = ["username", "email"]

					webhook {
						name = "bar"
						url_wo = "https://example.com/webhooks/one"
						url_wo_version = "1"
					}

					webhook {
						name = "foo"
						url_wo = "https://example.com/webhooks/two"
						url_wo_version = "1"
					}
				}

				data "stacklet_configuration_profile_slack" "test" {
					depends_on = [stacklet_configuration_profile_slack.test]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_slack.test", "id", "data.stacklet_configuration_profile_slack.test", "id"),
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_slack.test", "profile", "data.stacklet_configuration_profile_slack.test", "profile"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_slack.test", "user_fields.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_slack.test", "user_fields.0", "username"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_slack.test", "user_fields.1", "email"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_slack.test", "webhook.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_slack.test", "webhook.0.name", "bar"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_slack.test", "webhook.0.url"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_slack.test", "webhook.1.name", "foo"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_slack.test", "webhook.1.url"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileSlackDataSource", steps)
}
