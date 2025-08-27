// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileTeamsDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_teams" "test" {
					webhook {
						name = "bar"
						url_wo = "https://outlook.office.com/webhook/bar"
						url_wo_version = "1"
					}

					webhook {
						name = "foo"
						url_wo = "https://outlook.office.com/webhook/foo"
						url_wo_version = "1"
					}
				}

				data "stacklet_configuration_profile_teams" "test" {
					depends_on = [stacklet_configuration_profile_teams.test]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_teams.test", "id", "data.stacklet_configuration_profile_teams.test", "id"),
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_teams.test", "profile", "data.stacklet_configuration_profile_teams.test", "profile"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_teams.test", "webhook.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_teams.test", "webhook.0.name", "bar"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_teams.test", "webhook.0.url"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_teams.test", "webhook.1.name", "foo"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_teams.test", "webhook.1.url"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileTeamsDataSource", steps)
}
