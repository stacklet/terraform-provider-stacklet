// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileServiceNowResource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_servicenow" "test" {
					endpoint = "https://dev12345.service-now.com"
					username = "test-user"
					password_wo = "test-password"
					password_wo_version = "1"
					issue_type = "incident"
					closed_state = "closed"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_servicenow.test", "id"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "profile", "servicenow"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "endpoint", "https://dev12345.service-now.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "username", "test-user"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_servicenow.test", "password"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "issue_type", "incident"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "closed_state", "closed"),
			),
		},
		{
			ResourceName:      "stacklet_configuration_profile_servicenow.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"password_wo",
				"password_wo_version",
			},
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileServiceNowResource", steps)
}

func TestAccConfigurationProfileServiceNowResource_PasswordChange(t *testing.T) {
	var returnedPassword string

	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_servicenow" "test" {
					endpoint = "https://dev12345.service-now.com"
					username = "test-user"
					password_wo = "initial-password"
					password_wo_version = "1"
					issue_type = "incident"
					closed_state = "closed"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_servicenow.test", "id"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "endpoint", "https://dev12345.service-now.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "username", "test-user"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "password_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_servicenow.test", "password"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_servicenow.test", "password", func(value string) error {
					returnedPassword = value
					return nil
				}),
			),
		},
		// Same version, password should remain the same
		{
			Config: `
				resource "stacklet_configuration_profile_servicenow" "test" {
					endpoint = "https://dev12345.service-now.com"
					username = "test-user"
					password_wo = "new-password-content"
					password_wo_version = "1"
					issue_type = "incident"
					closed_state = "closed"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "endpoint", "https://dev12345.service-now.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "username", "test-user"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "password_wo_version", "1"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_servicenow.test", "password", func(value string) error {
					if value != returnedPassword {
						return fmt.Errorf("password should not have changed, was %s, now: %s", returnedPassword, value)
					}
					return nil
				}),
			),
		},
		// Version changed, password should change
		{
			Config: `
				resource "stacklet_configuration_profile_servicenow" "test" {
					endpoint = "https://dev12345.service-now.com"
					username = "test-user"
					password_wo = "new-password-content"
					password_wo_version = "2"
					issue_type = "incident"
					closed_state = "closed"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "endpoint", "https://dev12345.service-now.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "username", "test-user"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_servicenow.test", "password_wo_version", "2"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_servicenow.test", "password", func(value string) error {
					if value == returnedPassword {
						return fmt.Errorf("password should have changed")
					}
					return nil
				}),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileServiceNowResource_PasswordChange", steps)
}
