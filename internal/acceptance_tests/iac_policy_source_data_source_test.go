// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIaCPolicySourceDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
				resource "stacklet_iac_policy_source" "test" {
					name = "{{.Prefix}}-policy-source-ds"
					url = "https://github.com/test-org/test-repo"
					branch = "main"
					paths = ["policies"]
					auto_scan = true

					token_auth = {
						user = "bill"
						token_wo = "secret"
						token_wo_version = 1
					}
				}

				data "stacklet_iac_policy_source" "test" {
					name = stacklet_iac_policy_source.test.name
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_source.test", "name", prefixName("policy-source-ds")),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_source.test", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_source.test", "branch", "main"),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_source.test", "paths.#", "1"),
				resource.TestCheckTypeSetElemAttr("data.stacklet_iac_policy_source.test", "paths.*", "policies"),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_source.test", "extensions.#", "2"),
				resource.TestCheckTypeSetElemAttr("data.stacklet_iac_policy_source.test", "extensions.*", ".yml"),
				resource.TestCheckTypeSetElemAttr("data.stacklet_iac_policy_source.test", "extensions.*", ".yaml"),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_source.test", "auto_scan", "true"),
				resource.TestCheckResourceAttrPair("data.stacklet_iac_policy_source.test", "webhook_url", "stacklet_iac_policy_source.test", "webhook_url"),
				resource.TestCheckResourceAttr("data.stacklet_iac_policy_source.test", "token_auth.user", "bill"),
				resource.TestCheckNoResourceAttr("data.stacklet_iac_policy_source.test", "ssh_auth.user"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicySourceDataSource", steps)
}
