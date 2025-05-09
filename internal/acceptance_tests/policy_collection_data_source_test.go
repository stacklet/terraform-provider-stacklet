package acceptance_tests

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyCollectionDataSource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccPolicyCollectionDataSource")
	if err != nil {
		t.Fatal(err)
	}

	// Set up the HTTP client with our recorded transport
	http.DefaultTransport = rt

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "stacklet_policy_collection" "test" {
						name = "test-collection-ds"
						description = "Test policy collection"
						cloud_provider = "AWS"
					}

					data "stacklet_policy_collection" "test" {
						name = stacklet_policy_collection.test.name
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.stacklet_policy_collection.test", "name", "test-collection-ds"),
					resource.TestCheckResourceAttr("data.stacklet_policy_collection.test", "description", "Test policy collection"),
					resource.TestCheckResourceAttr("data.stacklet_policy_collection.test", "cloud_provider", "AWS"),
				),
			},
		},
	})
}
