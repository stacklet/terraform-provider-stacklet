package acceptance_tests

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountResource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccAccountResource")
	if err != nil {
		t.Fatal(err)
	}

	// Set up the HTTP client with our recorded transport
	http.DefaultTransport = rt

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
					resource "stacklet_account" "test" {
						name = "test-account"
						key = "999999999999"
						cloud_provider = "AWS"
						description = "Test AWS account"
						short_name = "test"
						email = "test@example.com"
						variables = "{\"environment\": \"test\"}"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_account.test", "name", "test-account"),
					resource.TestCheckResourceAttr("stacklet_account.test", "key", "999999999999"),
					resource.TestCheckResourceAttr("stacklet_account.test", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr("stacklet_account.test", "description", "Test AWS account"),
					resource.TestCheckResourceAttr("stacklet_account.test", "short_name", "test"),
					resource.TestCheckResourceAttr("stacklet_account.test", "email", "test@example.com"),
					resource.TestCheckResourceAttr("stacklet_account.test", "variables", "{\"environment\": \"test\"}"),
					resource.TestCheckResourceAttrSet("stacklet_account.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "stacklet_account.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           "aws:999999999999",
				ImportStateVerifyIgnore: []string{"security_context"},
			},
			// Update and Read testing
			{
				Config: `
					resource "stacklet_account" "test" {
						name = "test-account-updated"
						key = "999999999999"
						cloud_provider = "AWS"
						description = "Updated AWS account"
						short_name = "test-updated"
						email = "updated@example.com"
						variables = "{\"environment\": \"staging\"}"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_account.test", "name", "test-account-updated"),
					resource.TestCheckResourceAttr("stacklet_account.test", "key", "999999999999"),
					resource.TestCheckResourceAttr("stacklet_account.test", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr("stacklet_account.test", "description", "Updated AWS account"),
					resource.TestCheckResourceAttr("stacklet_account.test", "short_name", "test-updated"),
					resource.TestCheckResourceAttr("stacklet_account.test", "email", "updated@example.com"),
					resource.TestCheckResourceAttr("stacklet_account.test", "variables", "{\"environment\": \"staging\"}"),
				),
			},
		},
	})
}
