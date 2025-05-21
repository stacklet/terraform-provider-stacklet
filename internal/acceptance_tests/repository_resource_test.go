package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccRepositoryResourceAttrs(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "test-repo"
						description = "Test repository"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_repository.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_repository.test", "uuid"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "name", "test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Test repository"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "system", "false"),
				resource.TestCheckResourceAttrSet("stacklet_repository.test", "webhook_url"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_auth_token", "false"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_ssh_private_key", "false"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_ssh_passphrase", "false"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_user"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token_wo"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token_wo_version"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_public_key"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key_wo"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key_wo_version"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase_wo"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase_wo_version"),
			),
		},
		{
			ResourceName:      "stacklet_repository.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateId:     "https://github.com/test-org/test-repo",
		},
	}
	runRecordedAccTest(t, "TestAccRepositoryResourceAttrs", steps)
}

func TestAccRepositoryResourceUpdate(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "test-repo"
						description = "Test repository"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "name", "test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Test repository"),
			),
		},
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "real-test-repo"
						description = "Real test repository"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "name", "real-test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Real test repository"),
			),
		},
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/real-test-repo"
						name = "real-test-repo"
						description = "Real test repository"
					}
				`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("stacklet_repository.test", plancheck.ResourceActionReplace),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/real-test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "name", "real-test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Real test repository"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRepositoryResourceUpdate", steps)
}

func TestAccRepositoryResourceHTTPAuthSet(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "test-repo"

					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_user"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_auth_token", "false"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token_wo_version"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token_wo"),
			),
		},
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "test-repo"

						auth_token_wo = "secret"
					}
				`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// auth_token_wo alone doesn't trigger a change
					plancheck.ExpectResourceAction("stacklet_repository.test", plancheck.ResourceActionNoop),
				},
			},
		},
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "test-repo"

						auth_user = "bill"
						auth_token_wo = "secret"
						auth_token_wo_version = 1
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_repository.test", "auth_user", "bill"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_auth_token", "true"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "auth_token_wo_version", "1"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token_wo"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRepositoryResourceHTTPAuthSet", steps)
}

func TestAccRepositoryResourceHTTPAuthClear(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "test-repo"

						auth_user = "bill"
						auth_token_wo = "secret"
						auth_token_wo_version = 1
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_repository.test", "auth_user", "bill"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_auth_token", "true"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "auth_token_wo_version", "1"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token_wo"),
			),
		},
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "test-repo"

						auth_user = "bill"
						auth_token_wo = "different-secret"
						auth_token_wo_version = 1
					}
				`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// auth_token_wo alone doesn't trigger a change
					plancheck.ExpectResourceAction("stacklet_repository.test", plancheck.ResourceActionNoop),
				},
			},
		},
		{
			Config: `
					resource "stacklet_repository" "test" {
						url = "https://github.com/test-org/test-repo"
						name = "test-repo"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_user"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_auth_token", "false"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token_wo_version"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token_wo"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRepositoryResourceHTTPAuthClear", steps)
}

func TestAccRepositoryResourceSSHUpdate(t *testing.T) {
	steps := []resource.TestStep{
		{
			// Create with simple SSH config.
			Config: `
resource "stacklet_repository" "test" {
	url = "ssh://git@github.com/stacklet/test-repo"
	name = "test-repo"

	ssh_private_key_wo = <<EOT
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDQSLdUuvJls2k28bLyalLCrZLaBG/ObQXBkEJl8vB1GwAAAJjX7M8q1+zP
KgAAAAtzc2gtZWQyNTUxOQAAACDQSLdUuvJls2k28bLyalLCrZLaBG/ObQXBkEJl8vB1Gw
AAAEC/dsi/DISHYy8HxIrX5JWLWhYKv2XFBlL15NLRzIlA5tBIt1S68mWzaTbxsvJqUsKt
ktoEb85tBcGQQmXy8HUbAAAADnJlcG8tdGVzdC1yYXcKAQIDBAUGBw==
-----END OPENSSH PRIVATE KEY-----
EOT
	ssh_private_key_wo_version = 1
}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_repository.test", "ssh_public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINBIt1S68mWzaTbxsvJqUsKtktoEb85tBcGQQmXy8HUb"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_ssh_private_key", "true"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "ssh_private_key_wo_version", "1"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key_wo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_ssh_passphrase", "false"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase_wo_version"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase_wo"),
			),
		},
		{
			// Quick check for _wo/_wo_version behaviour.
			Config: `
resource "stacklet_repository" "test" {
	url = "ssh://git@github.com/stacklet/test-repo"
	name = "test-repo"

	ssh_private_key_wo = "ignore-because-version-unchanged"
	ssh_private_key_wo_version = 1
	ssh_passphrase_wo = "ignore-because-version-unchanged"
}
				`,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// _wo_version unchanged, nothing to do
					plancheck.ExpectResourceAction("stacklet_repository.test", plancheck.ResourceActionNoop),
				},
			},
		},
		{
			// Set a new key with a passphrase.
			Config: `
resource "stacklet_repository" "test" {
	url = "ssh://git@github.com/stacklet/test-repo"
	name = "test-repo"

	ssh_private_key_wo = <<EOT
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABBo6Jwjw7
wQYNZrr9iiO8JWAAAAGAAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIKqKHdGkr5bDRMDX
XbfWxmbCKduClJweRKUxdqHnmUHXAAAAoDnrX1ai+rLAVjCJUW1nrcEVqb+JRZ0K5dIOnz
ysKDPFz6LdY4S6uzgZrE/WOcHX7/MgeXpjne8CQuIqej8KDM9XkLGHf010/cg7Fo60YMoG
UTEPMNoh4wYqZ030I7a5iOjPSRMD2tN+xhb8NSm5gDnYdn9SDkhArS7WGQrWLDf0Eh5qB/
LhPc79SMQGPhtv5Cwpb6686bmwIJU2/0l4g3M=
-----END OPENSSH PRIVATE KEY-----
EOT
	ssh_private_key_wo_version = 2
	ssh_passphrase_wo = "secret"
	ssh_passphrase_wo_version = 2
}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_repository.test", "ssh_public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKqKHdGkr5bDRMDXXbfWxmbCKduClJweRKUxdqHnmUHX"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_ssh_private_key", "true"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "ssh_private_key_wo_version", "2"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key_wo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_ssh_passphrase", "true"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "ssh_passphrase_wo_version", "2"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase_wo"),
			),
		},
		{
			// Restore more-or-less the original config.
			Config: `
resource "stacklet_repository" "test" {
	url = "ssh://git@github.com/stacklet/test-repo"
	name = "test-repo"

	ssh_private_key_wo = <<EOT
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDQSLdUuvJls2k28bLyalLCrZLaBG/ObQXBkEJl8vB1GwAAAJjX7M8q1+zP
KgAAAAtzc2gtZWQyNTUxOQAAACDQSLdUuvJls2k28bLyalLCrZLaBG/ObQXBkEJl8vB1Gw
AAAEC/dsi/DISHYy8HxIrX5JWLWhYKv2XFBlL15NLRzIlA5tBIt1S68mWzaTbxsvJqUsKt
ktoEb85tBcGQQmXy8HUbAAAADnJlcG8tdGVzdC1yYXcKAQIDBAUGBw==
-----END OPENSSH PRIVATE KEY-----
EOT
	ssh_private_key_wo_version = "something else"
}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("stacklet_repository.test", "ssh_public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINBIt1S68mWzaTbxsvJqUsKtktoEb85tBcGQQmXy8HUb"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_ssh_private_key", "true"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "ssh_private_key_wo_version", "something else"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key_wo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "has_ssh_passphrase", "false"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase_wo_version"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase_wo"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRepositoryResourceSSHAuthUpdate", steps)
}
