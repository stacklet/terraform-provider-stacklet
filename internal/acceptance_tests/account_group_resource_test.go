package acceptance_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAccountGroupResource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccAccountGroupResource")
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
					resource "stacklet_account_group" "test" {
						name = "test-group"
						description = "Test account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_account_group.test", "name", "test-group"),
					resource.TestCheckResourceAttr("stacklet_account_group.test", "description", "Test account group"),
					resource.TestCheckResourceAttr("stacklet_account_group.test", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr("stacklet_account_group.test", "regions.0", "us-east-1"),
					resource.TestCheckResourceAttrSet("stacklet_account_group.test", "uuid"),
				),
			},
			// ImportState testing using UUID
			{
				ResourceName:      "stacklet_account_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["stacklet_account_group.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return rs.Primary.Attributes["uuid"], nil
				},
			},
			// Update and Read testing
			{
				Config: `
					resource "stacklet_account_group" "test" {
						name = "test-group-updated"
						description = "Updated account group"
						cloud_provider = "AWS"
						regions = ["us-east-1", "us-east-2"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_account_group.test", "name", "test-group-updated"),
					resource.TestCheckResourceAttr("stacklet_account_group.test", "description", "Updated account group"),
					resource.TestCheckResourceAttr("stacklet_account_group.test", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr("stacklet_account_group.test", "regions.0", "us-east-1"),
					resource.TestCheckResourceAttr("stacklet_account_group.test", "regions.1", "us-east-2"),
					resource.TestCheckResourceAttrSet("stacklet_account_group.test", "uuid"),
				),
			},
		},
	})
}
