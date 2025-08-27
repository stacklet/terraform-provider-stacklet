// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileJiraResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
				resource "stacklet_configuration_profile_jira" "test" {
					url = "https://example.atlassian.net"
					user = "test@example.com"
					api_key_wo = "initial-api-key"
					api_key_wo_version = "1"

					project {
						name = "foo"
						project = "FOO"
						issue_type = "Task"
						closed_status = "Done"
					}

					project {
						name = "bar"
						project = "BAR"
						issue_type = "Bug"
						closed_status = "Fixed"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_jira.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_jira.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "url", "https://example.atlassian.net"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "user", "test@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "api_key_wo_version", "1"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.0.name", "foo"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.0.project", "FOO"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.0.issue_type", "Task"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.0.closed_status", "Done"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.1.name", "bar"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.1.project", "BAR"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.1.issue_type", "Bug"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.1.closed_status", "Fixed"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_configuration_profile_jira.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"api_key_wo", "api_key_wo_version",
				// Import returns projects in API order (alphabetical), not config order
				"project.0.name", "project.1.name",
				"project.0.project", "project.1.project",
				"project.0.issue_type", "project.1.issue_type",
				"project.0.closed_status", "project.1.closed_status",
			},
		},
		// Update and Read testing
		{
			Config: `
				resource "stacklet_configuration_profile_jira" "test" {
					url = "https://updated.atlassian.net"
					user = "updated@example.com"
					api_key_wo = "updated-api-key"
					api_key_wo_version = "1"

					project {
						name = "Updated Project"
						project = "UPD"
						issue_type = "Bug"
						closed_status = "Closed"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "url", "https://updated.atlassian.net"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "user", "updated@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "api_key_wo_version", "1"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.#", "1"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.0.name", "Updated Project"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.0.project", "UPD"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.0.issue_type", "Bug"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "project.0.closed_status", "Closed"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileJiraResource", steps)
}

func TestAccConfigurationProfileJiraResource_APIKeyChange(t *testing.T) {
	var returnedAPIKey string

	steps := []resource.TestStep{
		// Initial API key setup
		{
			Config: `
				resource "stacklet_configuration_profile_jira" "test" {
					url = "https://example.atlassian.net"
					user = "test@example.com"
					api_key_wo = "initial-api-key"
					api_key_wo_version = "1"

					project {
						name = "Test Project"
						project = "TEST"
						issue_type = "Task"
						closed_status = "Done"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_jira.test", "id"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "api_key_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_jira.test", "api_key"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_jira.test", "api_key", func(value string) error {
					returnedAPIKey = value
					return nil
				}),
			),
		},
		// Same version, API key should remain the same
		{
			Config: `
				resource "stacklet_configuration_profile_jira" "test" {
					url = "https://example.atlassian.net"
					user = "test@example.com"
					api_key_wo = "changed-api-key-but-same-version"
					api_key_wo_version = "1"

					project {
						name = "Test Project"
						project = "TEST"
						issue_type = "Task"
						closed_status = "Done"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "api_key_wo_version", "1"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_jira.test", "api_key", func(value string) error {
					if value != returnedAPIKey {
						return fmt.Errorf("api_key should not have changed, was %s, now: %s", returnedAPIKey, value)
					}
					returnedAPIKey = value
					return nil
				}),
			),
		},
		// Version changed, API key should change
		{
			Config: `
				resource "stacklet_configuration_profile_jira" "test" {
					url = "https://example.atlassian.net"
					user = "test@example.com"
					api_key_wo = "updated-api-key"
					api_key_wo_version = "2"

					project {
						name = "Test Project"
						project = "TEST"
						issue_type = "Task"
						closed_status = "Done"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_jira.test", "api_key_wo_version", "2"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_jira.test", "api_key", func(value string) error {
					if value == returnedAPIKey {
						return fmt.Errorf("api_key should have changed")
					}
					return nil
				}),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileJiraResource_APIKeyChange", steps)
}
