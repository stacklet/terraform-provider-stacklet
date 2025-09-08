// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileSymphonyResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
				resource "stacklet_configuration_profile_symphony" "test" {
					agent_domain = "example.symphony.com"
					service_account = "test-service-account"
					private_key_wo = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC7\n-----END PRIVATE KEY-----"
					private_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_symphony.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_symphony.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "agent_domain", "example.symphony.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "service_account", "test-service-account"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "private_key_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_symphony.test", "private_key"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_configuration_profile_symphony.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"private_key_wo", "private_key_wo_version",
			},
		},
		// Update and Read testing
		{
			Config: `
				resource "stacklet_configuration_profile_symphony" "test" {
					agent_domain = "updated.symphony.com"
					service_account = "updated-service-account"
					private_key_wo = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC7\n-----END PRIVATE KEY-----"
					private_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_symphony.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_symphony.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "agent_domain", "updated.symphony.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "service_account", "updated-service-account"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "private_key_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_symphony.test", "private_key"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileSymphonyResource", steps)
}

func TestAccConfigurationProfileSymphonyResource_PrivateKeyChange(t *testing.T) {
	var returnedPrivateKey string

	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_configuration_profile_symphony" "test" {
					agent_domain = "example.symphony.com"
					service_account = "test-service-account"
					private_key_wo = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC7\n-----END PRIVATE KEY-----"
					private_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_symphony.test", "id"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "agent_domain", "example.symphony.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "service_account", "test-service-account"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "private_key_wo_version", "1"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_symphony.test", "private_key"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_symphony.test", "private_key", func(value string) error {
					returnedPrivateKey = value
					return nil
				}),
			),
		},
		// Same version, private key should remain the same
		{
			Config: `
				resource "stacklet_configuration_profile_symphony" "test" {
					agent_domain = "example.symphony.com"
					service_account = "test-service-account"
					private_key_wo = "-----BEGIN PRIVATE KEY-----\nNEW_KEY_CONTENT_HERE\n-----END PRIVATE KEY-----"
					private_key_wo_version = "1"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "agent_domain", "example.symphony.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "service_account", "test-service-account"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "private_key_wo_version", "1"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_symphony.test", "private_key", func(value string) error {
					if value != returnedPrivateKey {
						return fmt.Errorf("private_key should not have changed, was %s, now: %s", returnedPrivateKey, value)
					}
					return nil
				}),
			),
		},
		// Version changed, private key should change
		{
			Config: `
				resource "stacklet_configuration_profile_symphony" "test" {
					agent_domain = "example.symphony.com"
					service_account = "test-service-account"
					private_key_wo = "-----BEGIN PRIVATE KEY-----\nNEW_KEY_CONTENT_HERE\n-----END PRIVATE KEY-----"
					private_key_wo_version = "2"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "agent_domain", "example.symphony.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "service_account", "test-service-account"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_symphony.test", "private_key_wo_version", "2"),
				resource.TestCheckResourceAttrWith("stacklet_configuration_profile_symphony.test", "private_key", func(value string) error {
					if value == returnedPrivateKey {
						return fmt.Errorf("private_key should have changed")
					}
					return nil
				}),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileSymphonyResource_PrivateKeyChange", steps)
}
