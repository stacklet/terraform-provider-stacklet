package acceptance_tests

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBindingDataSource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccBindingDataSource")
	if err != nil {
		t.Fatal(err)
	}

	// Set up the HTTP client with our recorded transport
	http.DefaultTransport = rt

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a binding to test the data source
			{
				Config: `
					resource "stacklet_account_group" "test" {
						name = "test-binding-ds-group"
						description = "Test account group for binding data source"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					resource "stacklet_policy_collection" "test" {
						name = "test-binding-ds-collection"
						description = "Test policy collection for binding data source"
						cloud_provider = "AWS"
					}

					resource "stacklet_binding" "test" {
						name = "test-binding-ds"
						description = "Test binding for data source"
						account_group_uuid = stacklet_account_group.test.uuid
						policy_collection_uuid = stacklet_policy_collection.test.uuid
						auto_deploy = true
						schedule = "rate(1 hour)"
						variables = jsonencode({
							environment = "test"
							region = "us-east-1"
						})
					}

					# Test lookup by name
					data "stacklet_binding" "by_name" {
						name = stacklet_binding.test.name
					}

					# Test lookup by UUID
					data "stacklet_binding" "by_uuid" {
						uuid = stacklet_binding.test.uuid
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify lookup by name
					resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "name", "test-binding-ds"),
					resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "description", "Test binding for data source"),
					resource.TestCheckResourceAttrSet("data.stacklet_binding.by_name", "account_group_uuid"),
					resource.TestCheckResourceAttrSet("data.stacklet_binding.by_name", "policy_collection_uuid"),
					resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "auto_deploy", "true"),
					resource.TestCheckResourceAttr("data.stacklet_binding.by_name", "schedule", "rate(1 hour)"),
					resource.TestCheckResourceAttrSet("data.stacklet_binding.by_name", "id"),
					resource.TestCheckResourceAttrSet("data.stacklet_binding.by_name", "uuid"),

					// Verify lookup by UUID
					resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "name", "test-binding-ds"),
					resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "description", "Test binding for data source"),
					resource.TestCheckResourceAttrSet("data.stacklet_binding.by_uuid", "account_group_uuid"),
					resource.TestCheckResourceAttrSet("data.stacklet_binding.by_uuid", "policy_collection_uuid"),
					resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "auto_deploy", "true"),
					resource.TestCheckResourceAttr("data.stacklet_binding.by_uuid", "schedule", "rate(1 hour)"),
					resource.TestCheckResourceAttrSet("data.stacklet_binding.by_uuid", "id"),
					resource.TestCheckResourceAttrSet("data.stacklet_binding.by_uuid", "uuid"),
				),
			},
		},
	})
}
