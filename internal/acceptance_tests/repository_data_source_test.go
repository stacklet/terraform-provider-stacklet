// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRepositoryDataSource(t *testing.T) {
	baseline := `
		resource "stacklet_repository" "test" {
			name = "{{.Prefix}}-repo-ds"
			url = "https://github.com/test-org/test-repo"
			description = "Test repository"
		}
	`
	steps := []resource.TestStep{
		{
			Config: baseline + `
				data "stacklet_repository" "test_uuid" {
					uuid = stacklet_repository.test.uuid
				}

				data "stacklet_repository" "test_url" {
					url = stacklet_repository.test.url
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check resource attributes
				resource.TestCheckResourceAttr("stacklet_repository.test", "name", prefixName("repo-ds")),
				resource.TestCheckResourceAttr("stacklet_repository.test", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("stacklet_repository.test", "description", "Test repository"),
				// Check data source attributes match exactly
				resource.TestCheckResourceAttr("data.stacklet_repository.test_uuid", "name", "test-repo-ds"),
				resource.TestCheckResourceAttr("data.stacklet_repository.test_uuid", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("data.stacklet_repository.test_uuid", "description", "Test repository"),
				resource.TestCheckResourceAttr("data.stacklet_repository.test_url", "name", "test-repo-ds"),
				resource.TestCheckResourceAttr("data.stacklet_repository.test_url", "url", "https://github.com/test-org/test-repo"),
				resource.TestCheckResourceAttr("data.stacklet_repository.test_url", "description", "Test repository"),
				// Add debug logging for both resource and data source
				func(s *terraform.State) error {
					rs := s.RootModule().Resources["stacklet_repository.test"]
					ds_url := s.RootModule().Resources["data.stacklet_repository.test_url"]
					ds_uuid := s.RootModule().Resources["data.stacklet_repository.test_uuid"]
					t.Logf("Resource Attributes: %#v", rs.Primary.Attributes)
					t.Logf("Data Source (URL) Attributes: %#v", ds_url.Primary.Attributes)
					t.Logf("Data Source (UUID) Attributes: %#v", ds_uuid.Primary.Attributes)
					return nil
				},
			),
		},
	}
	runRecordedAccTest(t, "TestAccRepositoryDataSource", steps)
}
