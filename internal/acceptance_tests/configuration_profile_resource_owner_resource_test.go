// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileResourceOwnerResource(t *testing.T) {
	steps := []resource.TestStep{
		// Minimal configuration
		{
			Config: `
				resource "stacklet_configuration_profile_resource_owner" "test" {}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_resource_owner.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_resource_owner.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "default.#", "0"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "tags.#", "0"),
				resource.TestCheckNoResourceAttr("stacklet_configuration_profile_resource_owner.test", "org_domain"),
				resource.TestCheckNoResourceAttr("stacklet_configuration_profile_resource_owner.test", "org_domain_tag"),
			),
		},
		// Full configuration
		{
			Config: `
				resource "stacklet_configuration_profile_resource_owner" "test" {
					default = ["owner1@example.com", "owner2@example.com"]
					org_domain = "example.com"
                    org_domain_tag = "domain"
					tags = ["owner", "team"]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_resource_owner.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_configuration_profile_resource_owner.test", "profile"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "default.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "default.0", "owner1@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "default.1", "owner2@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "org_domain", "example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "org_domain_tag", "domain"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "tags.#", "2"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "tags.0", "owner"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "tags.1", "team"),
			),
		},
		// ImportState testing
		{
			ResourceName:      "stacklet_configuration_profile_resource_owner.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
		// Update and Read testing
		{
			Config: `
				resource "stacklet_configuration_profile_resource_owner" "test" {
					default = ["owner3@example.com", "owner4@example.com", "owner5@example.com"]
					org_domain = "updated.com"
					org_domain_tag = "domain_tag"
					tags = ["updated_owner", "updated_team", "department"]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "default.#", "3"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "default.0", "owner3@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "default.1", "owner4@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "default.2", "owner5@example.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "org_domain", "updated.com"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "org_domain_tag", "domain_tag"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "tags.#", "3"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "tags.0", "updated_owner"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "tags.1", "updated_team"),
				resource.TestCheckResourceAttr("stacklet_configuration_profile_resource_owner.test", "tags.2", "department"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileResourceOwnerResource", steps)
}
