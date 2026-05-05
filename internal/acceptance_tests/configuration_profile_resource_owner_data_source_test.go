// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationProfileResourceOwnerDataSource(t *testing.T) {
	baseline := `
		resource "stacklet_configuration_profile_resource_owner" "test" {
			default = ["owner1@example.com", "owner2@example.com"]
			org_domain = "example.com"
			tags = ["owner", "team"]
		}
	`
	steps := []resource.TestStep{
		{
			Config: baseline + `
				data "stacklet_configuration_profile_resource_owner" "test" {
					depends_on = [stacklet_configuration_profile_resource_owner.test]
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_resource_owner.test", "id", "data.stacklet_configuration_profile_resource_owner.test", "id"),
				resource.TestCheckResourceAttrPair("stacklet_configuration_profile_resource_owner.test", "profile", "data.stacklet_configuration_profile_resource_owner.test", "profile"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_resource_owner.test", "default.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_resource_owner.test", "default.0", "owner1@example.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_resource_owner.test", "default.1", "owner2@example.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_resource_owner.test", "org_domain", "example.com"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_resource_owner.test", "tags.#", "2"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_resource_owner.test", "tags.0", "owner"),
				resource.TestCheckResourceAttr("data.stacklet_configuration_profile_resource_owner.test", "tags.1", "team"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccConfigurationProfileResourceOwnerDataSource", steps)
}
