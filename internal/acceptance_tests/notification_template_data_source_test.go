// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationTemplateDataSource(t *testing.T) {
	baseline := `
		resource "stacklet_notification_template" "test" {
			name = "{{.Prefix}}-notification-template"
			description = "Test notification template"
			transport = "email"
			content = "sample content"
		}
	`
	steps := []resource.TestStep{
		{
			Config: baseline + `
				data "stacklet_notification_template" "test" {
					name = stacklet_notification_template.test.name
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_notification_template.test", "id"),
				resource.TestCheckResourceAttr("data.stacklet_notification_template.test", "name", prefixName("notification-template")),
				resource.TestCheckResourceAttr("data.stacklet_notification_template.test", "description", "Test notification template"),
				resource.TestCheckResourceAttr("data.stacklet_notification_template.test", "transport", "email"),
				resource.TestCheckResourceAttr("data.stacklet_notification_template.test", "content", "sample content"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccNotificationTemplateDataSource", steps)
}
