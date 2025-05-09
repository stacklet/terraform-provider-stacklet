package acceptance_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPolicyCollectionResource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccPolicyCollectionResource")
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
					resource "stacklet_policy_collection" "test" {
						name = "test-collection"
						description = "Test policy collection"
						cloud_provider = "AWS"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_policy_collection.test", "name", "test-collection"),
					resource.TestCheckResourceAttr("stacklet_policy_collection.test", "description", "Test policy collection"),
					resource.TestCheckResourceAttrSet("stacklet_policy_collection.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "stacklet_policy_collection.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["stacklet_policy_collection.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return rs.Primary.Attributes["uuid"], nil
				},
			},
			// Update and Read testing
			{
				Config: `
					resource "stacklet_policy_collection" "test" {
						name = "test-collection-updated"
						description = "Updated policy collection"
						cloud_provider = "AWS"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_policy_collection.test", "name", "test-collection-updated"),
					resource.TestCheckResourceAttr("stacklet_policy_collection.test", "description", "Updated policy collection"),
					resource.TestCheckResourceAttr("stacklet_policy_collection.test", "cloud_provider", "AWS"),
				),
			},
		},
	})
}
