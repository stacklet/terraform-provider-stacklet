package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountDataSource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccAccountDataSource")
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
					resource "stacklet_account" "test" {
						name = "test-account-ds"
						key = "999999999999"
						cloud_provider = "AWS"
						description = "Test AWS account"
						short_name = "test"
						path = "/test"
						email = "test@example.com"
						active = true
						variables = "{\"environment\": \"test\"}"
					}
				`,
			},
			// Read testing
			{
				Config: `
					resource "stacklet_account" "test" {
						name = "test-account-ds"
						key = "999999999999"
						cloud_provider = "AWS"
						description = "Test AWS account"
						short_name = "test"
						path = "/test"
						email = "test@example.com"
						active = true
						variables = "{\"environment\": \"test\"}"
					}

					data "stacklet_account" "test" {
						key = stacklet_account.test.key
						cloud_provider = stacklet_account.test.cloud_provider
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.stacklet_account.test", "name", "test-account-ds"),
					resource.TestCheckResourceAttr("data.stacklet_account.test", "key", "999999999999"),
					resource.TestCheckResourceAttr("data.stacklet_account.test", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr("data.stacklet_account.test", "description", "Test AWS account"),
					resource.TestCheckResourceAttr("data.stacklet_account.test", "short_name", "test"),
					resource.TestCheckResourceAttr("data.stacklet_account.test", "path", "/test"),
					resource.TestCheckResourceAttr("data.stacklet_account.test", "email", "test@example.com"),
					resource.TestCheckResourceAttr("data.stacklet_account.test", "active", "true"),
					testAccCheckMapValues("data.stacklet_account.test", "variables", map[string]string{"environment": "test"}),
					resource.TestCheckResourceAttrSet("data.stacklet_account.test", "id"),
				),
			},
		},
	})
}
