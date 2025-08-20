// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccReportGroupResource(t *testing.T) {
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
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_report_group.test", "id"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "name", prefixName("report-group")),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "enabled", "true"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "bindings.#", "1"),
				resource.TestCheckResourceAttrPair("stacklet_report_group.test", "bindings.0", "stacklet_binding.b", "uuid"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "schedule", "0 12 * * *"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "group_by.#", "2"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "group_by.0", "account"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "group_by.1", "region"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "use_message_settings", "true"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_report_group.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_report_group.test.name"),
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_report_group" "test" {
						name = "{{.Prefix}}-report-group"
	                    enabled = false
						bindings = []
						schedule = "0 6 * * *"
						group_by = ["account"]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_report_group.test", "enabled", "false"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "bindings.#", "0"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "schedule", "0 6 * * *"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "group_by.#", "1"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "group_by.0", "account"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "use_message_settings", "true"),
			),
		},
		// Test that updating name causes replacement
		{
			Config: `
					resource "stacklet_report_group" "test" {
						name = "{{.Prefix}}-report-group-renamed"
	                    enabled = true
						bindings = []
						schedule = "0 12 * * *"
						group_by = ["account"]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_report_group.test", "name", prefixName("report-group-renamed")),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "enabled", "true"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "bindings.#", "0"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "schedule", "0 12 * * *"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "group_by.#", "1"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "group_by.0", "account"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "use_message_settings", "true"),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_report_group.test", plancheck.ResourceActionReplace),
				},
			},
		},
	}
	runRecordedAccTest(t, "TestAccReportGroupResource", steps)
}

func TestAccReportGroupResource_DeliverySettings(t *testing.T) {
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

	                    email_delivery_settings {
                            template = "template.html"
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
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_report_group.test", "email_delivery_settings.#", "1"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "email_delivery_settings.0.template", "template.html"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "email_delivery_settings.0.subject", "Matched resources"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "email_delivery_settings.0.recipients.#", "2"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "email_delivery_settings.0.recipients.0.resource_owner", "true"),
				resource.TestCheckResourceAttr("stacklet_report_group.test", "email_delivery_settings.0.recipients.1.value", "user@example.com"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccReportGroupResource_DeliverySettings", steps)
}
