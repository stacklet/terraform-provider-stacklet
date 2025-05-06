package acceptance_tests

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccountGroupDataSource(t *testing.T) {
	rt, err := setupRecordedTest(t, "TestAccAccountGroupDataSource")
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
					resource "stacklet_account_group" "test" {
						name = "test-group-ds"
						description = "Test account group"
						cloud_provider = "AWS"
						regions = ["us-east-1"]
					}

					data "stacklet_account_group" "test" {
						name = stacklet_account_group.test.name
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.stacklet_account_group.test", "name", "test-group-ds"),
					resource.TestCheckResourceAttr("data.stacklet_account_group.test", "description", "Test account group"),
					resource.TestCheckResourceAttr("data.stacklet_account_group.test", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr("data.stacklet_account_group.test", "regions.0", "us-east-1"),
				),
			},
		},
	})
}
