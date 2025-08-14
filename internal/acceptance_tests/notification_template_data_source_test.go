// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationTemplateDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_notification_template" "test" {
						name = "{{.Prefix}}-notification-template"
						description = "Test notification template"
						transport = "email"
						content = "sample content"
					}

					data "stacklet_notification_template" "test" {
						name = "{{.Prefix}}-notification-template"

		                depends_on = [stacklet_notification_template.test]
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
