// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIaCPoliciesDataSource_InvalidConfig(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
data "stacklet_iac_policies" "test" {
  source     = "source"
  collection = "collection"
}
`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`2 attributes specified when one \(and only one\) of`),
		},
		{
			Config: `
data "stacklet_iac_policies" "test" {}
`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of`),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPoliciesDataSource_InvalidConfig", steps)
}
