package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountGroupItemsResource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccAccountGroupItemsResource")
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
					resource "stacklet_account" "test1" {
						name = "test-account-1"
						key = "111111111111"
						cloud_provider = "aws"
						description = "Test AWS account 1"
					}

					resource "stacklet_account" "test2" {
						name = "test-account-2"
						key = "222222222222"
						cloud_provider = "aws"
						description = "Test AWS account 2"
					}

					resource "stacklet_account_group" "test" {
						name = "test-group-items"
						description = "Test account group"
						policy_collection_name = "default"
						variables = "{\"environment\": \"test\"}"
					}

					resource "stacklet_account_group_items" "test" {
						group_name = stacklet_account_group.test.name
						items = [
							{
								provider = "aws"
								key = stacklet_account.test1.key
								variables = "{\"environment\": \"test\"}"
							}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "group_name", "test-group-items"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.#", "1"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.0.provider", "aws"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.0.key", "111111111111"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.0.variables", "{\"environment\": \"test\"}"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "stacklet_account_group_items.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "test-group-items",
			},
			// Update and Read testing
			{
				Config: `
					resource "stacklet_account" "test1" {
						name = "test-account-1"
						key = "111111111111"
						cloud_provider = "aws"
						description = "Test AWS account 1"
					}

					resource "stacklet_account" "test2" {
						name = "test-account-2"
						key = "222222222222"
						cloud_provider = "aws"
						description = "Test AWS account 2"
					}

					resource "stacklet_account_group" "test" {
						name = "test-group-items"
						description = "Test account group"
						policy_collection_name = "default"
						variables = "{\"environment\": \"test\"}"
					}

					resource "stacklet_account_group_items" "test" {
						group_name = stacklet_account_group.test.name
						items = [
							{
								provider = "aws"
								key = stacklet_account.test1.key
								variables = "{\"environment\": \"staging\"}"
							},
							{
								provider = "aws"
								key = stacklet_account.test2.key
								variables = "{\"environment\": \"prod\"}"
							}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "group_name", "test-group-items"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.#", "2"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.0.provider", "aws"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.0.key", "111111111111"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.0.variables", "{\"environment\": \"staging\"}"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.1.provider", "aws"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.1.key", "222222222222"),
					resource.TestCheckResourceAttr("stacklet_account_group_items.test", "items.1.variables", "{\"environment\": \"prod\"}"),
				),
			},
		},
	})
}
