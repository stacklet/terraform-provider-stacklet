package acceptance_tests

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSOGroupResource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccSSOGroupResource")
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
					resource "stacklet_sso_group" "test" {
						name = "test-group"
						roles = ["admin"]
						account_group_uuids = ["e2e040cf-6f10-4cdf-94b1-9600b2ee36ca"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "name", "test-group"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "roles.0", "admin"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "account_group_uuids.#", "1"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "account_group_uuids.0", "e2e040cf-6f10-4cdf-94b1-9600b2ee36ca"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "stacklet_sso_group.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           "test-group",
				ImportStateVerifyIgnore: []string{"id"}, // Ignore id as it's generated
			},
			// Update and Read testing
			{
				Config: `
					resource "stacklet_sso_group" "test" {
						name = "test-group-updated"
						roles = ["admin", "viewer"]
						account_group_uuids = [
							"d76ccc28-5bad-49b4-9caf-ec4ea71f1cdc",
							"e2e040cf-6f10-4cdf-94b1-9600b2ee36ca"
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "name", "test-group-updated"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "roles.#", "2"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "roles.0", "admin"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "roles.1", "viewer"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "account_group_uuids.#", "2"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "account_group_uuids.0", "d76ccc28-5bad-49b4-9caf-ec4ea71f1cdc"),
					resource.TestCheckResourceAttr("stacklet_sso_group.test", "account_group_uuids.1", "e2e040cf-6f10-4cdf-94b1-9600b2ee36ca"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
