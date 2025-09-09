// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileEmailDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_email" "test" {
					from = "user@example.com"
	                ses_region = "us-east-1"

 		            smtp = {
                        server = "smtp.example.com"
                        port = "1234"
                        ssl = true
                        username = "user"

                        password_wo = "secret"
                        password_wo_version = "1"
                    }
                }

				data "stacklet_configuration_profile_email" "test" {
                    depends_on = [stacklet_configuration_profile_email.test]
                }
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_email.test", "id"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_email.test", "profile", "email"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_email.test", "from", "user@example.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_email.test", "ses_region", "us-east-1"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_email.test", "smtp.server", "smtp.example.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_email.test", "smtp.port", "1234"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_email.test", "smtp.ssl", "true"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_email.test", "smtp.username", "user"),
				resource.TestCheckResourceAttrSet("data.stacklet_configuration_profile_email.test", "smtp.password"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileEmailDataSource", steps)
}
