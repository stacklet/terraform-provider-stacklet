// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReportGroupDataSource(t *testing.T) {
	// Create a report group to test the data source
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_account_group" "ag" {
						name = "{{.Prefix}}-rg-ag"
						description = "Test account group for report group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_policy_collection" "pc" {
						name = "{{.Prefix}}-rg-pc"
						description = "Test policy collection for report group"
						cloud_provider = "AWS"
					}

					resource "stacklet_binding" "b" {
						name = "{{.Prefix}}-rg-binding"
						description = "Test binding for report group"
						account_group_uuid = stacklet_account_group.ag.uuid
						policy_collection_uuid = stacklet_policy_collection.pc.uuid
					}

					resource "stacklet_report_group" "test" {
						name = "{{.Prefix}}-report-group"
						bindings = [stacklet_binding.b.uuid]
						schedule = "0 12 * * *"
						group_by = ["account", "region"]

	                    email_delivery_settings {
                            template = "email"
                            subject = "Matched resources"

                            recipients = [
                        	    {
                                  resource_owner = true
                                },
                        	    {
                                  value = "user@example.com"
                                },
                             ]
                        }
					}

                    data "stacklet_report_group" "test" {
						name = "{{.Prefix}}-report-group"

                        depends_on = [stacklet_report_group.test]
                    }
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_report_group.test", "id"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "name", prefixName("report-group")),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "enabled", "true"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "bindings.#", "1"),
				resource.TestCheckResourceAttrPair("data.stacklet_report_group.test", "bindings.0", "stacklet_binding.b", "uuid"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "source", "BINDING"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "schedule", "0 12 * * *"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "group_by.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "group_by.0", "account"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "group_by.1", "region"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "use_message_settings", "true"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "email_delivery_settings.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "email_delivery_settings.0.template", "email"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "email_delivery_settings.0.subject", "Matched resources"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "email_delivery_settings.0.recipients.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "email_delivery_settings.0.recipients.0.resource_owner", "true"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "email_delivery_settings.0.recipients.1.value", "user@example.com"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "slack_delivery_settings.#", "0"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "msteams_delivery_settings.#", "0"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "servicenow_delivery_settings.#", "0"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "jira_delivery_settings.#", "0"),
				resource.TestCheckResourceAttr("data.stacklet_report_group.test", "symphony_delivery_settings.#", "0"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccReportGroupDataSource", steps)
}
