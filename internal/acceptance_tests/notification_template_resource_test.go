// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccNotificationTemplateResource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_notification_template" "test" {
						name = "{{.Prefix}}-notification-template"
						description = "Test notification template"
						transport = "email"
						content = "sample content"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_notification_template.test", "id"),
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "name", prefixName("notification-template")),
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "description", "Test notification template"),
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "transport", "email"),
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "content", "sample content"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_notification_template.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: importStateIDFuncFromAttrs("stacklet_notification_template.test.name"),
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_notification_template" "test" {
						name = "{{.Prefix}}-notification-template"
						description = "Test updated notification template"
						transport = "slack"
						content = "sample updated content"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "description", "Test updated notification template"),
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "transport", "slack"),
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "content", "sample updated content"),
			),
		},
		// Test that updating name causes replacement
		{
			Config: `
					resource "stacklet_notification_template" "test" {
						name = "{{.Prefix}}-notification-template-renamed"
						description = "Test updated notification template"
						transport = "slack"
						content = "sample updated content"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "name", prefixName("notification-template-renamed")),
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "description", "Test updated notification template"),
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "transport", "slack"),
				resource.TestCheckResourceAttr("stacklet_notification_template.test", "content", "sample updated content"),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_notification_template.test", plancheck.ResourceActionReplace),
				},
			},
		},
	}
	runRecordedAccTest(t, "TestAccNotificationTemplateResource", steps)
}
