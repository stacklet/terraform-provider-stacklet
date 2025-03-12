package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountGroupItemsDataSource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccAccountGroupItemsDataSource")
	if err != nil {
		t.Fatal(err)
	}

	// Set up the HTTP client with our recorded transport
	http.DefaultTransport = rt

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "stacklet_account" "test1" {
						name = "test-account-1-ds"
						key = "111111111111"
						cloud_provider = "aws"
						description = "Test AWS account 1"
					}

					resource "stacklet_account_group" "test" {
						name = "test-group-items-ds"
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

					data "stacklet_account_group_items" "test" {
						group_name = stacklet_account_group.test.name
						depends_on = [stacklet_account_group_items.test]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.stacklet_account_group_items.test", "group_name", "test-group-items-ds"),
					resource.TestCheckResourceAttr("data.stacklet_account_group_items.test", "items.#", "1"),
					resource.TestCheckResourceAttr("data.stacklet_account_group_items.test", "items.0.provider", "aws"),
					resource.TestCheckResourceAttr("data.stacklet_account_group_items.test", "items.0.key", "111111111111"),
					resource.TestCheckResourceAttr("data.stacklet_account_group_items.test", "items.0.variables", "{\"environment\": \"test\"}"),
				),
			},
		},
	})
}
