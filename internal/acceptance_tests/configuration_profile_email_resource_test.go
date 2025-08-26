// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileEmailResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
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
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_email.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_email.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "from", "user@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "ses_region", "us-east-1"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.server", "smtp.example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.port", "1234"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.ssl", "true"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.username", "user"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.password_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_email.test", "smtp.password"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_configuration_profile_email.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"smtp.password_wo", "smtp.password_wo_version",
			},
		},
		// Update and Read testing
		{
			Config: `
				resource "stacklet_configuration_profile_email" "test" {
					from = "updated-user@example.com"
	                ses_region = "us-east-2"

 		            smtp = {
                        server = "new-smtp.example.com"
                        port = "5678"
                        ssl = false
                        username = "new-user"

                        password_wo = "secret"
                        password_wo_version = "1"
                    }
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_email.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_email.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "from", "updated-user@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "ses_region", "us-east-2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.server", "new-smtp.example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.port", "5678"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.ssl", "false"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.username", "new-user"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileEmailResource", steps)
}

func TestAccConfigurationProfileEmailResource_PasswordChange(t *testing.T) {
	var returnedPassword string

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
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.password_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_email.test", "smtp.password"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_email.test", "smtp.password", func(value string) error {
					returnedPassword = value
					return nil
				}),
			),
		},
		// Same version, password should remain the same
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

                        password_wo = "new-secret"
                        password_wo_version = "1"
                    }
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.password_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_email.test", "smtp.password"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_email.test", "smtp.password", func(value string) error {
					if value != returnedPassword {
						return fmt.Errorf("smtp.password should not have changed, was %s, now: %s", returnedPassword, value)
					}
					return nil
				}),
			),
		},
		// Version changed, password should change
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

                        password_wo = "new-secret"
                        password_wo_version = "2"
                    }
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_email.test", "smtp.password_wo_version", "2"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_email.test", "smtp.password"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_email.test", "smtp.password", func(value string) error {
					if value == returnedPassword {
						return fmt.Errorf("smtp.password should have changed")
					}
					return nil
				}),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileEmailResource_PasswordChange", steps)
}
