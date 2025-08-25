// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileAccountOwnersResource(t *testing.T) {
	steps := []resource.TestStep{
		// Minimal configuration
		{
			Config: `
				resource "stacklet_configuration_profile_account_owners" "test" {
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_account_owners.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_account_owners.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.#", "0"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "tags.#", "0"),
				resource.TestCheckNoResourceAttr("stacklet_configuration_profile_account_owners.test", "org_domain"),
				resource.TestCheckNoResourceAttr("stacklet_configuration_profile_account_owners.test", "org_domain_tag"),
			),
		},
		// Full configuration
		{
			Config: `
				resource "stacklet_configuration_profile_account_owners" "test" {
					default = [
						{
							account = "123456789012"
							owners = ["owner1@example.com", "owner2@example.com"]
						}
					]
					org_domain = "example.com"
                    org_domain_tag = "domain"
					tags = ["owner", "team"]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_account_owners.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_account_owners.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.#", "1"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.0.account", "123456789012"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.0.owners.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.0.owners.0", "owner1@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.0.owners.1", "owner2@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "org_domain", "example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "org_domain_tag", "domain"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "tags.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "tags.0", "owner"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "tags.1", "team"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_configuration_profile_account_owners.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
		// Update and Read testing
		{
			Config: `
				resource "stacklet_configuration_profile_account_owners" "test" {
					default = [
						{
							account = "123456789012"
							owners = ["owner1@example.com"]
						},
						{
							account = "210987654321"
							owners = ["owner3@example.com", "owner4@example.com"]
						}
					]
					org_domain = "updated.com"
					org_domain_tag = "domain_tag"
					tags = ["updated_owner", "updated_team", "department"]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.0.account", "123456789012"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.0.owners.#", "1"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.0.owners.0", "owner1@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.1.account", "210987654321"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.1.owners.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.1.owners.0", "owner3@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "default.1.owners.1", "owner4@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "org_domain", "updated.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "org_domain_tag", "domain_tag"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "tags.#", "3"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "tags.0", "updated_owner"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "tags.1", "updated_team"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_account_owners.test", "tags.2", "department"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileAccountOwnersResource", steps)
}
