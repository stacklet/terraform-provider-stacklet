// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWriteOnlyRequiresVersion(t *testing.T) {
	// Optional write-only attributes, whose version is required when they're set.
	missingVersion := func(name string) *regexp.Regexp {
		return regexp.MustCompile(`"` + name + `_version" must be specified`)
	}
	// Required write-only attributes, whose version is required too.
	requiredVersion := func(name string) *regexp.Regexp {
		return regexp.MustCompile(`"` + name + `_version"\s+is\s+required`)
	}

	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_account_discovery_gcp" "test" {
  name               = "{{.Prefix}}-discovery"
  org_id             = "123456789012"
  credential_json_wo = "{}"
}
`,
			ExpectError: requiredVersion("credential_json_wo"),
		},
		{
			Config: `
resource "stacklet_anthropic_integration" "test" {
  org_id       = "` + testAnthropicOrgID + `"
  admin_key_wo = "` + testAnthropicAdminKey + `"
}
`,
			ExpectError: requiredVersion("admin_key_wo"),
		},
		{
			Config: `
resource "stacklet_configuration_profile_jira" "test" {
  url        = "https://example.atlassian.net"
  user       = "test@example.com"
  api_key_wo = "api-key"

  project {
    name          = "foo"
    project       = "FOO"
    issue_type    = "Task"
    closed_status = "Done"
  }
}
`,
			ExpectError: requiredVersion("api_key_wo"),
		},
		{
			Config: `
resource "stacklet_configuration_profile_servicenow" "test" {
  endpoint     = "https://dev12345.service-now.com"
  username     = "test-user"
  password_wo  = "test-password"
  issue_type   = "incident"
  closed_state = "closed"
}
`,
			ExpectError: requiredVersion("password_wo"),
		},
		{
			Config: `
resource "stacklet_configuration_profile_slack" "test" {
  webhook {
    name   = "foo"
    url_wo = "https://example.com/webhooks/one"
  }
}
`,
			ExpectError: requiredVersion("url_wo"),
		},
		{
			Config: `
resource "stacklet_configuration_profile_symphony" "test" {
  agent_domain    = "example.symphony.com"
  service_account = "test-service-account"
  private_key_wo  = "key"
}
`,
			ExpectError: requiredVersion("private_key_wo"),
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
			ExpectError: requiredVersion("token_wo"),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "ssh://git@github.com/test-org/test-repo"

  ssh_auth = {
    user           = "git"
    private_key_wo = "key"
  }
}
`,
			ExpectError: requiredVersion("private_key_wo"),
		},
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "ssh://git@github.com/test-org/test-repo"

  ssh_auth = {
    user                   = "git"
    private_key_wo         = "key"
    private_key_wo_version = "1"
    passphrase_wo          = "secret"
  }
}
`,
			ExpectError: missingVersion("ssh_auth.passphrase_wo"),
		},
		{
			Config: `
resource "stacklet_account" "test" {
  name                = "{{.Prefix}}-account"
  key                 = "999999999999"
  cloud_provider      = "AWS"
  security_context_wo = "arn:aws:iam::123456789012:role/stacklet-execution"
}
`,
			ExpectError: missingVersion("security_context_wo"),
		},
		{
			Config: `
resource "stacklet_binding" "test" {
  name                   = "{{.Prefix}}-binding"
  account_group_uuid     = "00000000-0000-0000-0000-000000000001"
  policy_collection_uuid = "00000000-0000-0000-0000-000000000002"
  security_context_wo    = "arn:aws:iam::123456789012:role/stacklet-execution"
}
`,
			ExpectError: missingVersion("security_context_wo"),
		},
		{
			Config: `
resource "stacklet_configuration_profile_email" "test" {
  from = "user@example.com"

  smtp = {
    server      = "smtp.example.com"
    port        = "1234"
    password_wo = "secret"
  }
}
`,
			ExpectError: missingVersion("smtp.password_wo"),
		},
		{
			Config: `
resource "stacklet_configuration_profile_slack" "test" {
  token_wo = "xoxb-token"

  webhook {
    name           = "foo"
    url_wo         = "https://example.com/webhooks/one"
    url_wo_version = "1"
  }
}
`,
			ExpectError: missingVersion("token_wo"),
		},
		{
			Config: `
resource "stacklet_datadog_integration" "test" {
  api_key_wo         = "api-key"
  app_key_wo         = "app-key"
  app_key_wo_version = "1"
}
`,
			ExpectError: missingVersion("api_key_wo"),
		},
		{
			Config: `
resource "stacklet_datadog_integration" "test" {
  api_key_wo         = "api-key"
  api_key_wo_version = "1"
  app_key_wo         = "app-key"
}
`,
			ExpectError: missingVersion("app_key_wo"),
		},
		{
			Config: `
resource "stacklet_repository" "test" {
  url           = "https://github.com/test-org/test-repo"
  name          = "{{.Prefix}}-repo"
  auth_token_wo = "secret"
}
`,
			ExpectError: missingVersion("auth_token_wo"),
		},
		{
			Config: `
resource "stacklet_repository" "test" {
  url                = "ssh://git@github.com/test-org/test-repo"
  name               = "{{.Prefix}}-repo"
  ssh_private_key_wo = "key"
}
`,
			ExpectError: missingVersion("ssh_private_key_wo"),
		},
		{
			Config: `
resource "stacklet_repository" "test" {
  url                        = "ssh://git@github.com/test-org/test-repo"
  name                       = "{{.Prefix}}-repo"
  ssh_private_key_wo         = "key"
  ssh_private_key_wo_version = "1"
  ssh_passphrase_wo          = "secret"
}
`,
			ExpectError: missingVersion("ssh_passphrase_wo"),
		},
	}
	for i := range steps {
		steps[i].PlanOnly = true
	}
	runRecordedAccTest(t, "TestAccWriteOnlyRequiresVersion", steps)
}

func TestAccWriteOnlyVersionRequiresValue(t *testing.T) {
	missingValue := func(name string) *regexp.Regexp {
		return regexp.MustCompile(`"` + name + `" must be specified`)
	}

	steps := []resource.TestStep{
		{
			Config: `
resource "stacklet_iac_policy_source" "test" {
  name = "{{.Prefix}}-policy-source"
  url  = "ssh://git@github.com/test-org/test-repo"

  ssh_auth = {
    user                   = "git"
    private_key_wo         = "key"
    private_key_wo_version = "1"
    passphrase_wo_version  = "1"
  }
}
`,
			ExpectError: missingValue("ssh_auth.passphrase_wo"),
		},
		{
			Config: `
resource "stacklet_account" "test" {
  name                        = "{{.Prefix}}-account"
  key                         = "999999999999"
  cloud_provider              = "AWS"
  security_context_wo_version = "1"
}
`,
			ExpectError: missingValue("security_context_wo"),
		},
		{
			Config: `
resource "stacklet_binding" "test" {
  name                        = "{{.Prefix}}-binding"
  account_group_uuid          = "00000000-0000-0000-0000-000000000001"
  policy_collection_uuid      = "00000000-0000-0000-0000-000000000002"
  security_context_wo_version = "1"
}
`,
			ExpectError: missingValue("security_context_wo"),
		},
		{
			Config: `
resource "stacklet_configuration_profile_email" "test" {
  from = "user@example.com"

  smtp = {
    server              = "smtp.example.com"
    port                = "1234"
    password_wo_version = "1"
  }
}
`,
			ExpectError: missingValue("smtp.password_wo"),
		},
		{
			Config: `
resource "stacklet_configuration_profile_slack" "test" {
  token_wo_version = "1"

  webhook {
    name           = "foo"
    url_wo         = "https://example.com/webhooks/one"
    url_wo_version = "1"
  }
}
`,
			ExpectError: missingValue("token_wo"),
		},
		{
			Config: `
resource "stacklet_datadog_integration" "test" {
  api_key_wo_version = "1"
  app_key_wo         = "app-key"
  app_key_wo_version = "1"
}
`,
			ExpectError: missingValue("api_key_wo"),
		},
		{
			Config: `
resource "stacklet_datadog_integration" "test" {
  api_key_wo         = "api-key"
  api_key_wo_version = "1"
  app_key_wo_version = "1"
}
`,
			ExpectError: missingValue("app_key_wo"),
		},
		{
			Config: `
resource "stacklet_repository" "test" {
  url                   = "https://github.com/test-org/test-repo"
  name                  = "{{.Prefix}}-repo"
  auth_token_wo_version = "1"
}
`,
			ExpectError: missingValue("auth_token_wo"),
		},
		{
			Config: `
resource "stacklet_repository" "test" {
  url                        = "ssh://git@github.com/test-org/test-repo"
  name                       = "{{.Prefix}}-repo"
  ssh_private_key_wo_version = "1"
}
`,
			ExpectError: missingValue("ssh_private_key_wo"),
		},
		{
			Config: `
resource "stacklet_repository" "test" {
  url                        = "ssh://git@github.com/test-org/test-repo"
  name                       = "{{.Prefix}}-repo"
  ssh_private_key_wo         = "key"
  ssh_private_key_wo_version = "1"
  ssh_passphrase_wo_version  = "1"
}
`,
			ExpectError: missingValue("ssh_passphrase_wo"),
		},
	}
	for i := range steps {
		steps[i].PlanOnly = true
	}
	runRecordedAccTest(t, "TestAccWriteOnlyVersionRequiresValue", steps)
}
