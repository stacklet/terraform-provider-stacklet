// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlatformDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
					data "stacklet_platform" "test" {}
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.stacklet_platform.test", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_platform.test", "external_id"),
				// at least one region is enabled
				resource.TestCheckResourceAttrSet("data.stacklet_platform.test", "execution_regions.0"),
				resource.TestCheckResourceAttrSet("data.stacklet_platform.test", "default_role"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccPlatformDataSource", steps)
}
