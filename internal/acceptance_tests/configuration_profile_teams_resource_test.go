// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileTeamsResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing - webhooks intentionally not in alphabetical order
		{
			Config: `
				resource "stacklet_configuration_profile_teams" "test" {
					webhook {
						name = "foo"
						url_wo = "https://outlook.office.com/webhook/foo"
						url_wo_version = "1"
					}

					webhook {
						name = "bar"
						url_wo = "https://outlook.office.com/webhook/bar"
						url_wo_version = "1"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.name", "foo"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.url_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "webhook.0.url"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.name", "bar"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.url_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "webhook.1.url"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_configuration_profile_teams.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				// Import returns webhooks in API order (alphabetical), not config order
				"webhook.0.url_wo", "webhook.1.url_wo",
				"webhook.0.url_wo_version", "webhook.1.url_wo_version",
				"webhook.0.name", "webhook.1.name", "webhook.0.url", "webhook.1.url",
			},
		},
		// Update and Read testing - different order and updated values
		{
			Config: `
				resource "stacklet_configuration_profile_teams" "test" {
					webhook {
						name = "foo-new"
						url_wo = "https://outlook.office.com/webhook/foo-new"
						url_wo_version = "1"
					}

					webhook {
						name = "bar-new"
						url_wo = "https://outlook.office.com/webhook/bar-new"
						url_wo_version = "1"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.name", "foo-new"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.url_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "webhook.0.url"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.name", "bar-new"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.url_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "webhook.1.url"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileTeamsResource", steps)
}

func TestAccConfigurationProfileTeamsResource_URLChange(t *testing.T) {
	var returnedURL1, returnedURL2 string

	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_teams" "test" {
					webhook {
						name = "foo"
						url_wo = "https://outlook.office.com/webhook/foo"
						url_wo_version = "1"
					}

					webhook {
						name = "bar"
						url_wo = "https://outlook.office.com/webhook/bar"
						url_wo_version = "1"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "id"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.name", "foo"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.url_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "webhook.0.url"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.name", "bar"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.url_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_teams.test", "webhook.1.url"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_teams.test", "webhook.0.url", func(value string) error {
					returnedURL1 = value
					return nil
				}),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_teams.test", "webhook.1.url", func(value string) error {
					returnedURL2 = value
					return nil
				}),
			),
		},
		// Same version, URLs should remain the same
		{
			Config: `
				resource "stacklet_configuration_profile_teams" "test" {
					webhook {
						name = "foo"
						url_wo = "https://outlook.office.com/webhook/foo-new"
						url_wo_version = "1"
					}

					webhook {
						name = "bar"
						url_wo = "https://outlook.office.com/webhook/bar-new"
						url_wo_version = "1"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.name", "foo"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.url_wo_version", "1"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.name", "bar"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.url_wo_version", "1"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_teams.test", "webhook.0.url", func(value string) error {
					if value != returnedURL1 {
						return fmt.Errorf("webhook.0.url should not have changed, was %s, now: %s", returnedURL1, value)
					}
					return nil
				}),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_teams.test", "webhook.1.url", func(value string) error {
					if value != returnedURL2 {
						return fmt.Errorf("webhook.1.url should not have changed, was %s, now: %s", returnedURL2, value)
					}
					return nil
				}),
			),
		},
		// Version changed for first webhook, URL should change
		{
			Config: `
				resource "stacklet_configuration_profile_teams" "test" {
					webhook {
						name = "foo"
						url_wo = "https://outlook.office.com/webhook/foo-new"
						url_wo_version = "2"
					}

					webhook {
						name = "bar"
						url_wo = "https://outlook.office.com/webhook/bar-new"
						url_wo_version = "1"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.name", "foo"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.0.url_wo_version", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.name", "bar"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_teams.test", "webhook.1.url_wo_version", "1"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_teams.test", "webhook.0.url", func(value string) error {
					if value == returnedURL1 {
						return fmt.Errorf("webhook.0.url should have changed")
					}
					return nil
				}),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_teams.test", "webhook.1.url", func(value string) error {
					if value != returnedURL2 {
						return fmt.Errorf("webhook.1.url should not have changed, was %s, now: %s", returnedURL2, value)
					}
					return nil
				}),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileTeamsResource_URLChange", steps)
}
