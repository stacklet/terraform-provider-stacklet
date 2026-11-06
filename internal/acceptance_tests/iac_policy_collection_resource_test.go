// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccIaCPolicyCollectionResource_Attrs(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_collection" "test" {
  name       = "{{.Prefix}}-policy-collection"
  is_default = true
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_collection.test", "name", prefixName("policy-collection")),
				resource.TestCheckResourceAttr("stacklet_iac_policy_collection.test", "policies.#", "0"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_collection.test", "is_default", "true"),
			),
		},
		{
			ResourceName:                         "stacklet_iac_policy_collection.test",
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "name",
			ImportStateIdFunc:                    importStateIDFuncFromAttrs("stacklet_iac_policy_collection.test.name"),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicyCollectionResource_Attrs", steps)
}

func TestAccIaCPolicyCollectionResource_Update(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_collection" "test" {
  name = "{{.Prefix}}-policy-collection"
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_collection.test", "policies.#", "0"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_collection.test", "is_default", "false"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_collection" "test" {
  name       = "{{.Prefix}}-policy-collection"
  is_default = true
}
`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_iac_policy_collection.test", plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_collection.test", "is_default", "true"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_collection" "test" {
  name = "{{.Prefix}}-policy-collection"
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_collection.test", "is_default", "false"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_collection" "test" {
  name = "{{.Prefix}}-other-policy-collection"
}
`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_iac_policy_collection.test", plancheck.ResourceActionReplace),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_collection.test", "name", prefixName("other-policy-collection")),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicyCollectionResource_Update", steps)
}

func TestAccIaCPolicyCollectionResource_InvalidConfig(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_collection" "test" {
  name = "{{.Prefix}}-policy-collection-mirror"
}
`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`must not end with "-mirror"`),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicyCollectionResource_InvalidConfig", steps)
}
