package acceptance_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAccountGroupItemsResource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccAccountGroupItemsResource")
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
					resource "stacklet_account" "test1" {
						name = "test-account-1"
						key = "111111111111"
						cloud_provider = "AWS"
						description = "Test AWS account 1"
					}

					resource "stacklet_account" "test2" {
						name = "test-account-2"
						key = "222222222222"
						cloud_provider = "AWS"
						description = "Test AWS account 2"
					}

					resource "stacklet_account_group" "test" {
						name = "test-group-items"
						description = "Test account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_account_group_item" "test" {
						group_uuid = stacklet_account_group.test.uuid
						account_key = stacklet_account.test1.key
						cloud_provider = stacklet_account.test1.cloud_provider
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_account_group_item.test", "account_key", "111111111111"),
					resource.TestCheckResourceAttr("stacklet_account_group_item.test", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttrSet("stacklet_account_group_item.test", "id"),
					resource.TestCheckResourceAttrSet("stacklet_account_group_item.test", "group_uuid"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["stacklet_account_group_item.test"]
						if !ok {
							return fmt.Errorf("resource not found in state")
						}
						id := rs.Primary.Attributes["id"]
						groupUUID := rs.Primary.Attributes["group_uuid"]
						accountKey := rs.Primary.Attributes["account_key"]
						expectedID := fmt.Sprintf("%s:%s", groupUUID, accountKey)
						if id != expectedID {
							return fmt.Errorf("expected ID to be %s, got %s", expectedID, id)
						}
						return nil
					},
				),
			},
			// ImportState testing
			{
				ResourceName:      "stacklet_account_group_item.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["stacklet_account_group_item.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return fmt.Sprintf("%s:%s:%s", rs.Primary.Attributes["group_uuid"], rs.Primary.Attributes["cloud_provider"], rs.Primary.Attributes["account_key"]), nil
				},
			},
			// Update and Read testing
			{
				Config: `
					resource "stacklet_account" "test1" {
						name = "test-account-1"
						key = "111111111111"
						cloud_provider = "AWS"
						description = "Test AWS account 1"
					}

					resource "stacklet_account" "test2" {
						name = "test-account-2"
						key = "222222222222"
						cloud_provider = "AWS"
						description = "Test AWS account 2"
					}

					resource "stacklet_account_group" "test" {
						name = "test-group-items"
						description = "Test account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_account_group_item" "test" {
						group_uuid = stacklet_account_group.test.uuid
						account_key = stacklet_account.test2.key
						cloud_provider = stacklet_account.test2.cloud_provider
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("stacklet_account_group_item.test", "account_key", "222222222222"),
					resource.TestCheckResourceAttr("stacklet_account_group_item.test", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttrSet("stacklet_account_group_item.test", "id"),
					resource.TestCheckResourceAttrSet("stacklet_account_group_item.test", "group_uuid"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["stacklet_account_group_item.test"]
						if !ok {
							return fmt.Errorf("resource not found in state")
						}
						id := rs.Primary.Attributes["id"]
						groupUUID := rs.Primary.Attributes["group_uuid"]
						accountKey := rs.Primary.Attributes["account_key"]
						expectedID := fmt.Sprintf("%s:%s", groupUUID, accountKey)
						if id != expectedID {
							return fmt.Errorf("expected ID to be %s, got %s", expectedID, id)
						}
						return nil
					},
				),
			},
		},
	})
}
