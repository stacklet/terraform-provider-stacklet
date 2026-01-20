// Copyright (c) 2026 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAPIKeyDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					resource "stacklet_api_key" "test" {
						description = "Test API key"
					}

					data "stacklet_api_key" "test" {
						identity = stacklet_api_key.test.identity

						depends_on = [stacklet_api_key.test]
					}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "id"),
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "identity"),
				resource.TestCheckResourceAttr("stacklet_api_key.test", "description", "Test API key"),
				resource.TestCheckResourceAttrSet("stacklet_api_key.test", "secret"),
				resource.TestCheckResourceAttrSet("data.stacklet_api_key.test", "id"),
				resource.TestCheckResourceAttr("data.stacklet_api_key.test", "description", "Test API key"),
				resource.TestCheckResourceAttrPair("data.stacklet_api_key.test", "identity", "stacklet_api_key.test", "identity"),
				resource.TestCheckResourceAttrPair("data.stacklet_api_key.test", "id", "stacklet_api_key.test", "id"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccAPIKeyDataSource", steps)
}
