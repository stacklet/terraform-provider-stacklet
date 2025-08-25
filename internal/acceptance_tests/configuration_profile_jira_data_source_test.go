// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileJiraDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_jira" "test" {
					url = "https://example.atlassian.net"
					user = "test@example.com"
					api_key_wo = "test-api-key"
					api_key_wo_version = "1"

					project {
						name = "Test Project"
						project = "TEST"
						issue_type = "Task"
						closed_status = "Done"
					}
				}

				data "stacklet_configuration_profile_jira" "test" {
					depends_on = [stacklet_configuration_profile_jira.test]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_jira.test", "id", "data.stacklet_configuration_profile_jira.test", "id"),
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_jira.test", "profile", "data.stacklet_configuration_profile_jira.test", "profile"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_jira.test", "url", "https://example.atlassian.net"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_jira.test", "user", "test@example.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_jira.test", "project.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_jira.test", "project.0.name", "Test Project"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_jira.test", "project.0.project", "TEST"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_jira.test", "project.0.issue_type", "Task"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_jira.test", "project.0.closed_status", "Done"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileJiraDataSource", steps)
}
