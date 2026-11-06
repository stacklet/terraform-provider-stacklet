// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIaCPolicyCollectionDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_collection" "test" {
  name       = "{{.Prefix}}-policy-collection-ds"
  is_default = true
}

data "stacklet_iac_policy_collection" "test" {
  name = stacklet_iac_policy_collection.test.name
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_collection.test", "name", prefixName("policy-collection-ds")),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_collection.test", "policies.#", "0"),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_collection.test", "is_default", "true"),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_collection.test", "is_mirror", "false"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicyCollectionDataSource", steps)
}

func TestAccIaCPolicyCollectionDataSource_Mirror(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source-mirror-ds"
  url  = "https://github.com/test-org/test-repo"
}

data "stacklet_iac_policy_collection" "test" {
  name = "${stacklet_iac_policy_source.test.name}-mirror"
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_collection.test", "name", prefixName("policy-source-mirror-ds-mirror")),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_collection.test", "is_default", "false"),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_collection.test", "is_mirror", "true"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicyCollectionDataSource_Mirror", steps)
}
