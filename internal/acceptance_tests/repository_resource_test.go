package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepositoryResource(t *testing.T) {
	steps := []resource.TestStep{
		// Create and Read testing
		{
			Config: `
					resource "stacklet_repository" "test" {
						name = "test-repo"
						url = "https://github.com/test-org/test-repo"
						description = "Test repository"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_repository.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_repository.test", "uuid"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "name", "test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Test repository"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase"),
			),
		},
		// ImportState testing
		{
			ResourceName:            "stacklet_repository.test",
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateId:           "https://github.com/test-org/test-repo",
			ImportStateVerifyIgnore: []string{"auth_token", "ssh_private_key", "ssh_passphrase", "deep_import"},
		},
		// Update and Read testing
		{
			Config: `
					resource "stacklet_repository" "test" {
						name = "test-repo-updated"
						url = "https://github.com/test-org/test-repo"
						description = "Updated test repository"
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_repository.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_repository.test", "uuid"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "name", "test-repo-updated"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Updated test repository"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key"),
				resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccRepositoryResource", steps)
}
