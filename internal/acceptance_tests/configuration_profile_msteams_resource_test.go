// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccConfigurationProfileMSTeamsResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
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
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "profile"),

				// Check access_config computed values
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "access_config.client_id", "e90b9a7a-f726-44f4-af92-9e5827c465f8"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "access_config.roundtrip_digest", "724ba7cc82663bc247b5a100b3ca2ece"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "access_config.tenant_id", "408b7351-82bd-44b5-aed5-59198cd1c1c6"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "access_config.bot_application.download_url"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "access_config.bot_application.version"),

				// Check customer_config computed values
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "customer_config.prefix", "stacklet-test"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "customer_config.bot_endpoint"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "customer_config.oidc_client"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "customer_config.oidc_issuer"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "customer_config.roundtrip_digest"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "customer_config.tags.env", "test"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "customer_config.tags.team", "platform"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "customer_config.terraform_module.repository_url"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "customer_config.terraform_module.source"),

				// Check channel mappings
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.0.name", "alerts"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.0.team_id", "e22bd265-dfcb-448d-a05b-e4d110d2266e"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.0.channel_id", "19:hZZSubNbJL7A5cYRGKnK_AiL3ytC2gNl6yFh8_LVzbM1@thread.tacv2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.1.name", "notifications"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.1.team_id", "e22bd265-dfcb-448d-a05b-e4d110d2266e"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.1.channel_id", "19:deb9db569d964adf94ddf02e7c5ce4b9@thread.tacv2"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_configuration_profile_msteams.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"access_config_input",
				"customer_config_input",
			},
		},
		// Update and Read testing
		{
			Config: `
				resource "stacklet_configuration_profile_msteams" "test" {
					access_config_input = {
						client_id        = "e90b9a7a-f726-44f4-af92-9e5827c465f8"
						roundtrip_digest = "724ba7cc82663bc247b5a100b3ca2ece"
						tenant_id        = "408b7351-82bd-44b5-aed5-59198cd1c1c6"
					}

					customer_config_input = {
						prefix = "stacklet-prod"
						tags = {
							env = "production"
							team = "devops"
						}
					}

					channel_mapping {
						name       = "alerts-updated"
						team_id    = "e22bd265-dfcb-448d-a05b-e4d110d2266e"
						channel_id = "19:hZZSubNbJL7A5cYRGKnK_AiL3ytC2gNl6yFh8_LVzbM1@thread.tacv2"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "profile"),

				// Check updated customer_config
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "customer_config.prefix", "stacklet-prod"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "customer_config.tags.env", "production"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "customer_config.tags.team", "devops"),

				// Check updated channel mappings
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.#", "1"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.0.name", "alerts-updated"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.0.team_id", "e22bd265-dfcb-448d-a05b-e4d110d2266e"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_msteams.test", "channel_mapping.0.channel_id", "19:hZZSubNbJL7A5cYRGKnK_AiL3ytC2gNl6yFh8_LVzbM1@thread.tacv2"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileMSTeamsResource", steps)
}

func TestAccConfigurationProfileMSTeamsResource_RequiresReplaceIfAccessConfigUnset(t *testing.T) {
	steps := []resource.TestStep{
		// Create initial resource with both access_config_input and customer_config_input
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
						}
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_msteams.test", "profile"),
			),
		},
		// Try to remove access_config_input - should require replacement
		{
			Config: `
				resource "stacklet_configuration_profile_msteams" "test" {
					customer_config_input = {
						prefix = "stacklet-test"
						tags = {
							env = "test"
						}
					}
				}
			`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_configuration_profile_msteams.test", plancheck.ResourceActionDestroyBeforeCreate),
				},
			},
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileMSTeamsResource_RequiresReplaceIfAccessConfigUnset", steps)
}
