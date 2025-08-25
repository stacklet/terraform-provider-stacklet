// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileAccountOwnersDataSource(t *testing.T) {
	steps := []resource.TestStep{
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
					tags = ["owner", "team"]
				}

				data "stacklet_configuration_profile_account_owners" "test" {
                    depends_on = [stacklet_configuration_profile_account_owners.test]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_account_owners.test", "id", "data.stacklet_configuration_profile_account_owners.test", "id"),
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_account_owners.test", "profile", "data.stacklet_configuration_profile_account_owners.test", "profile"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_account_owners.test", "default.#", "1"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_account_owners.test", "default.0.account", "123456789012"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_account_owners.test", "default.0.owners.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_account_owners.test", "default.0.owners.0", "owner1@example.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_account_owners.test", "default.0.owners.1", "owner2@example.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_account_owners.test", "org_domain", "example.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_account_owners.test", "tags.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_account_owners.test", "tags.0", "owner"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_account_owners.test", "tags.1", "team"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileAccountOwnersDataSource", steps)
}
