// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: `
                    data "stacklet_policy" "p1" {
                      name = "cost-aws:aws-elb-unattached-inform"
                    }

                    data "stacklet_policy" "p2" {
                      uuid = data.stacklet_policy.p1.uuid
                    }
				`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.stacklet_policy.p1", "name", "cost-aws:aws-elb-unattached-inform"),
				resource.TestCheckResourceAttr("data.stacklet_policy.p1", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttrSet("data.stacklet_policy.p1", "version"),
				resource.TestCheckResourceAttrSet("data.stacklet_policy.p1", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_policy.p1", "description"),
				// The same policy, but obtained via UUID.
				resource.TestCheckResourceAttr("data.stacklet_policy.p2", "name", "cost-aws:aws-elb-unattached-inform"),
				resource.TestCheckResourceAttr("data.stacklet_policy.p2", "cloud_provider", "AWS"),
				resource.TestCheckResourceAttrSet("data.stacklet_policy.p2", "version"),
				resource.TestCheckResourceAttrSet("data.stacklet_policy.p2", "id"),
				resource.TestCheckResourceAttrSet("data.stacklet_policy.p2", "description"),
			),
		},
	}
	runRecordedAccTest(t, "TestAccPolicyDataSource", steps)
}
