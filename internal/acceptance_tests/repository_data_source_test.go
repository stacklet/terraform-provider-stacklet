package acceptance_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRepositoryDataSource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccRepositoryDataSource")
	if err != nil {
		t.Fatal(err)
	}

	// Set up the HTTP client with our recorded transport
	http.DefaultTransport = rt

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a resource first
			{
				Config: `
					resource "stacklet_repository" "test" {
						name = "test-repo-ds"
						url = "https://github.com/test-org/test-repo"
						description = "Test repository"
						policy_file_suffix = [".yaml"]
						policy_directories = ["policies"]
						branch_name = "main"
						deep_import = false
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check every single attribute to see what might be different
					resource.TestCheckResourceAttr("stacklet_repository.test", "name", "test-repo-ds"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Test repository"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_file_suffix.#", "1"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_file_suffix.0", ".yaml"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_directories.#", "1"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_directories.0", "policies"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "branch_name", "main"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "deep_import", "false"),
					resource.TestCheckNoResourceAttr("stacklet_repository.test", "auth_token"),
					resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_private_key"),
					resource.TestCheckNoResourceAttr("stacklet_repository.test", "ssh_passphrase"),
					resource.TestCheckResourceAttrSet("stacklet_repository.test", "uuid"),
					// Add a custom check that logs all attributes
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["stacklet_repository.test"]
						if !ok {
							return fmt.Errorf("resource not found")
						}
						t.Logf("Resource Attributes: %#v", rs.Primary.Attributes)
						return nil
					},
				),
				ExpectNonEmptyPlan: false,
			},
			// Read testing by name
			{
				Config: `
					resource "stacklet_repository" "test" {
						name = "test-repo-ds"
						url = "https://github.com/test-org/test-repo"
						description = "Test repository"
						policy_file_suffix = [".yaml"]
						policy_directories = ["policies"]
						branch_name = "main"
						deep_import = false
					}

					data "stacklet_repository" "test" {
						name = stacklet_repository.test.name
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check resource attributes
					resource.TestCheckResourceAttr("stacklet_repository.test", "name", "test-repo-ds"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Test repository"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_file_suffix.#", "1"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_file_suffix.0", ".yaml"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_directories.#", "1"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "policy_directories.0", "policies"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "branch_name", "main"),
					resource.TestCheckResourceAttr("stacklet_repository.test", "deep_import", "false"),
					// Check data source attributes match exactly
					resource.TestCheckResourceAttr("data.stacklet_repository.test", "name", "test-repo-ds"),
					resource.TestCheckResourceAttr("data.stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
					resource.TestCheckResourceAttr("data.stacklet_repository.test", "description", "Test repository"),
					resource.TestCheckResourceAttr("data.stacklet_repository.test", "policy_file_suffix.#", "1"),
					resource.TestCheckResourceAttr("data.stacklet_repository.test", "policy_file_suffix.0", ".yaml"),
					resource.TestCheckResourceAttr("data.stacklet_repository.test", "policy_directories.#", "1"),
					resource.TestCheckResourceAttr("data.stacklet_repository.test", "policy_directories.0", "policies"),
					resource.TestCheckResourceAttr("data.stacklet_repository.test", "branch_name", "main"),
					// Add debug logging for both resource and data source
					func(s *terraform.State) error {
						rs := s.RootModule().Resources["stacklet_repository.test"]
						ds := s.RootModule().Resources["data.stacklet_repository.test"]
						t.Logf("Resource Attributes: %#v", rs.Primary.Attributes)
						t.Logf("Data Source Attributes: %#v", ds.Primary.Attributes)
						return nil
					},
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
