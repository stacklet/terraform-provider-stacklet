package acceptance_tests

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepositoryResource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccRepositoryResource")
	if err != nil {
		t.Fatal(err)
	}

	// Set up the HTTP client with our recorded transport
	http.DefaultTransport = rt

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
					resource "stacklet_repository" "test" {
						name = "test-repo"
						url = "https://github.com/test-org/test-repo"
						description = "Test repository"
						policy_file_suffix = [".yaml"]
						policy_directories = ["policies"]
						branch_name = "main"
						deep_import = true
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_repository.test", "name", "test-repo"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Test repository"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_file_suffix.#", "1"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_file_suffix.0", ".yaml"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_directories.#", "1"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_directories.0", "policies"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "branch_name", "main"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "deep_import", "true"),
					resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token"),
					resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key"),
					resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase"),
					resource.TestCheckResourceAttrSet("stacklet_repository.test", "uuid"),
					resource.TestCheckResourceAttrSet("stacklet_repository.test", "id"),
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
						policy_file_suffix = [".yaml", ".yml"]
						policy_directories = ["policies", "custom-policies"]
						branch_name = "develop"
						deep_import = false
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_repository.test", "name", "test-repo-updated"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Updated test repository"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_file_suffix.#", "2"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_file_suffix.0", ".yaml"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_file_suffix.1", ".yml"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_directories.#", "2"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_directories.0", "policies"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_directories.1", "custom-policies"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "branch_name", "develop"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "deep_import", "false"),
					resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token"),
					resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key"),
					resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase"),
					resource.TestCheckResourceAttrSet("stacklet_repository.test", "uuid"),
					resource.TestCheckResourceAttrSet("stacklet_repository.test", "id"),
				),
			},
		},
	})
}
