package acceptance_tests

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyCollectionItemResource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccPolicyCollectionItemResource")
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
						name = "test-collection-items"
						description = "Test policy collection"
						cloud_provider = "AWS"
					}

					resource "stacklet_policy_collection_item" "test" {
						collection_uuid = stacklet_policy_collection.test.uuid
						policy_uuid = "d741e99b-5d4d-44c9-83e8-884c614cfc8f"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("stacklet_policy_collection_item.test", "collection_uuid"),
					resource.TestCheckResourceAttr("stacklet_policy_collection_item.test", "policy_uuid", "d741e99b-5d4d-44c9-83e8-884c614cfc8f"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "stacklet_policy_collection_item.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "659ad69f-27a4-46ea-815d-8b87b15b2df1:d741e99b-5d4d-44c9-83e8-884c614cfc8f",
			},
			// Update and Read testing
			{
				Config: `
					resource "stacklet_policy_collection" "test" {
						name = "test-collection-items"
						description = "Test policy collection"
						cloud_provider = "AWS"
					}

					resource "stacklet_policy_collection_item" "test" {
						collection_uuid = stacklet_policy_collection.test.uuid
						policy_uuid = "d741e99b-5d4d-44c9-83e8-884c614cfc8f"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("stacklet_policy_collection_item.test", "collection_uuid"),
					resource.TestCheckResourceAttr("stacklet_policy_collection_item.test", "policy_uuid", "d741e99b-5d4d-44c9-83e8-884c614cfc8f"),
				),
			},
		},
	})
}
