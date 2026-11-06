// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccIaCPolicySourceResource_Attrs(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name       = "{{.Prefix}}-policy-source"
  url        = "https://github.com/test-org/test-repo"
  branch     = "main"
  paths      = ["policies", "more-policies"]
  extensions = [".yaml"]
  auto_scan  = true
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "name", prefixName("policy-source")),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "branch", "main"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "paths.#", "2"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_policy_source.test", "paths.*", "policies"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_policy_source.test", "paths.*", "more-policies"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "extensions.#", "1"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_policy_source.test", "extensions.*", ".yaml"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "auto_scan", "true"),
				resource.TestCheckResourceAttrSet("stacklet_iac_policy_source.test", "webhook_url"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "token_auth.user"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.user"),
			),
		},
		{
			ResourceName:                         "stacklet_iac_policy_source.test",
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "name",
			ImportStateIdFunc:                    importStateIDFuncFromAttrs("stacklet_iac_policy_source.test.name"),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicySourceResource_Attrs", steps)
}

func TestAccIaCPolicySourceResource_Update(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "branch"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "paths"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "extensions.#", "2"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_policy_source.test", "extensions.*", ".yml"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_policy_source.test", "extensions.*", ".yaml"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "auto_scan", "false"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name       = "{{.Prefix}}-policy-source"
  url        = "https://github.com/test-org/test-other-repo"
  branch     = "main"
  paths      = ["policies"]
  extensions = [".yaml"]
  auto_scan  = true
}
`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_iac_policy_source.test", plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "url", "https://github.com/test-org/test-other-repo"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "branch", "main"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "paths.#", "1"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_policy_source.test", "paths.*", "policies"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "extensions.#", "1"),
				resource.TestCheckTypeSetElemAttr("stacklet_iac_policy_source.test", "extensions.*", ".yaml"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "auto_scan", "true"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-other-repo"
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "branch"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "paths"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "extensions.#", "2"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "auto_scan", "false"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-other-policy-source"
  url  = "https://github.com/test-org/test-other-repo"
}
`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_iac_policy_source.test", plancheck.ResourceActionReplace),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "name", prefixName("other-policy-source")),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicySourceResource_Update", steps)
}

func TestAccIaCPolicySourceResource_TokenAuthSet(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "token_auth.user"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"

  token_auth = {
    user             = "bill"
    token_wo         = "secret"
    token_wo_version = 1
  }
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.user", "bill"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.token_wo_version", "1"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "token_auth.token_wo"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.user"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"

  token_auth = {
    user             = "bill"
    token_wo         = "different-secret"
    token_wo_version = 1
  }
}
`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// token_wo alone doesn't trigger a change
					plancheck.ExpectResourceAction("stacklet_iac_policy_source.test", plancheck.ResourceActionNoop),
				},
			},
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"

  token_auth = {
    user             = "ben"
    token_wo         = "different-secret"
    token_wo_version = 1
  }
}
`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_iac_policy_source.test", plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.user", "ben"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.token_wo_version", "1"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"

  token_auth = {
    user             = "ben"
    token_wo         = "new-secret"
    token_wo_version = 2
  }
}
`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_iac_policy_source.test", plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.user", "ben"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.token_wo_version", "2"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicySourceResource_TokenAuthSet", steps)
}

func TestAccIaCPolicySourceResource_TokenAuthClear(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"

  token_auth = {
    user             = "bill"
    token_wo         = "secret"
    token_wo_version = 1
  }
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.user", "bill"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.token_wo_version", "1"),
			),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "token_auth.user"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "token_auth.token_wo_version"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicySourceResource_TokenAuthClear", steps)
}

func TestAccIaCPolicySourceResource_SSHAuthUpdate(t *testing.T) {
	steps := []resource.TestStep{
		{
			// Create with simple SSH config.
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "ssh://git@github.com/stacklet/test-repo"

  ssh_auth = {
    user                   = "git"
    private_key_wo         = file("{{ fixture_path "id_ed25519" }}")
    private_key_wo_version = 1
  }
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.user", "git"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.private_key_wo_version", "1"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.private_key_wo"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.passphrase_wo_version"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.passphrase_wo"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "token_auth.user"),
			),
		},
		{
			// Set a new key with a passphrase.
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "ssh://git@github.com/stacklet/test-repo"

  ssh_auth = {
    user                   = "git"
    private_key_wo         = file("{{ fixture_path "id_ed25519_passphrase" }}")
    private_key_wo_version = 2
    passphrase_wo          = "secret"
    passphrase_wo_version  = 1
  }
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.private_key_wo_version", "2"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.passphrase_wo_version", "1"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.passphrase_wo"),
			),
		},
		{
			// Switch to token authentication.
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/stacklet/test-repo"

  token_auth = {
    user             = "bill"
    token_wo         = "secret"
    token_wo_version = 1
  }
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.user", "bill"),
				resource.TestCheckResourceAttr("stacklet_iac_policy_source.test", "token_auth.token_wo_version", "1"),
				resource.TestCheckNoResourceAttr("stacklet_iac_policy_source.test", "ssh_auth.user"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicySourceResource_SSHAuthUpdate", steps)
}

func TestAccIaCPolicySourceResource_InvalidConfig(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"

  token_auth = {
    user             = "bill"
    token_wo         = "secret"
    token_wo_version = 1
  }
  ssh_auth = {
    user                   = "git"
    private_key_wo         = "key"
    private_key_wo_version = 1
  }
}
`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`cannot be specified when`),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "https://github.com/test-org/test-repo"

  token_auth = {
    user     = "bill"
    token_wo = "secret"
  }
}
`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`attribute "token_wo_version"\s+is required`),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "ssh://git@github.com/stacklet/test-repo"

  ssh_auth = {
    user                   = "git"
    private_key_wo         = "key"
    private_key_wo_version = 1
    passphrase_wo          = "secret"
  }
}
`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`passphrase_wo_version" must be specified`),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name       = "{{.Prefix}}-policy-source"
  url        = "https://github.com/test-org/test-repo"
  extensions = ["yaml"]
}
`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`must start with a dot`),
		},
	}
	runRecordedAccTest(t, "TestAccIaCPolicySourceResource_InvalidConfig", steps)
}
