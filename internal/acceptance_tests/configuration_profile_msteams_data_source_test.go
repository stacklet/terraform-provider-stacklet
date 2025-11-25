// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileMSTeamsDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_msteams" "test" {
					access_config_input = {
						client_id        = "e90b9a7a-f726-44f4-af92-9e5827c465f8"
						roundtrip_digest = "724ba7cc82663bc247b5a100b3ca2ece"
						tenant_id        = "408b7351-82bd-44b5-aed5-59198cd1c1c6"
					}

					customer_config_input = {
						prefix = "stacklet-test"
						tags = {
							env = "test"
							team = "platform"
						}
					}

					channel_mapping {
						name       = "alerts"
						team_id    = "e22bd265-dfcb-448d-a05b-e4d110d2266e"
						channel_id = "19:hZZSubNbJL7A5cYRGKnK_AiL3ytC2gNl6yFh8_LVzbM1@thread.tacv2"
					}

					channel_mapping {
						name       = "notifications"
						team_id    = "e22bd265-dfcb-448d-a05b-e4d110d2266e"
						channel_id = "19:deb9db569d964adf94ddf02e7c5ce4b9@thread.tacv2"
					}
				}

				data "stacklet_configuration_profile_msteams" "test" {
					depends_on = [stacklet_configuration_profile_msteams.test]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_msteams.test", "id", "data.stacklet_configuration_profile_msteams.test", "id"),
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_msteams.test", "profile", "data.stacklet_configuration_profile_msteams.test", "profile"),

				// Check access_config computed values in data source
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "access_config.client_id", "e90b9a7a-f726-44f4-af92-9e5827c465f8"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "access_config.roundtrip_digest", "724ba7cc82663bc247b5a100b3ca2ece"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "access_config.tenant_id", "408b7351-82bd-44b5-aed5-59198cd1c1c6"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_msteams.test", "access_config.bot_application.download_url"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_msteams.test", "access_config.bot_application.version"),

				// Check customer_config computed values in data source
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "customer_config.prefix", "stacklet-test"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_msteams.test", "customer_config.roundtrip_digest"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "customer_config.tags.env", "test"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "customer_config.tags.team", "platform"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_msteams.test", "customer_config.terraform_module.repository_url"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_msteams.test", "customer_config.terraform_module.source"),

				// Check channel mappings in data source
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "channel_mapping.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "channel_mapping.0.name", "alerts"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "channel_mapping.0.team_id", "e22bd265-dfcb-448d-a05b-e4d110d2266e"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "channel_mapping.0.channel_id", "19:hZZSubNbJL7A5cYRGKnK_AiL3ytC2gNl6yFh8_LVzbM1@thread.tacv2"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "channel_mapping.1.name", "notifications"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "channel_mapping.1.team_id", "e22bd265-dfcb-448d-a05b-e4d110d2266e"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_msteams.test", "channel_mapping.1.channel_id", "19:deb9db569d964adf94ddf02e7c5ce4b9@thread.tacv2"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileMSTeamsDataSource", steps)
}
