package acceptance_tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyDataSource(t *testing.T) {
	runRecordedAccTest(
		t,
		"TestAccPolicyDataSource",
		[]resource.TestStep{
			// Read resource
			{
				Config: `
                    data "stacklet_policy" "p1" {
                      name = "cost-aws:aws-elb-unattached-inform"
                    }

                    data "stacklet_policy" "p2" {
                      uuid = "bd27ffb4-2d7e-42bf-a7b3-58edcb7284a7"
                    }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.stacklet_policy.p1", "name", "cost-aws:aws-elb-unattached-inform"),
					resource.TestCheckResourceAttr("data.stacklet_policy.p1", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr("data.stacklet_policy.p1", "version", "1"),
					resource.TestCheckResourceAttrSet("data.stacklet_policy.p1", "id"),
					resource.TestCheckResourceAttrSet("data.stacklet_policy.p1", "description"),

					resource.TestCheckResourceAttr("data.stacklet_policy.p2", "name", "cost-aws:aws-rds-rightsize-disk-inform"),
					resource.TestCheckResourceAttr("data.stacklet_policy.p2", "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr("data.stacklet_policy.p2", "version", "1"),
					resource.TestCheckResourceAttrSet("data.stacklet_policy.p2", "id"),
					resource.TestCheckResourceAttrSet("data.stacklet_policy.p2", "description"),
				),
			},
		},
	)
}
