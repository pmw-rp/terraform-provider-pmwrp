package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.pmwrp_brokers.default", "foo"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
provider "pmwrp" {
  seed     = "seed-for-your.cluster.address.here:9092"
  username = "your-username-here"
  password = "your-password-here"
}

data "pmwrp_brokers" "default" {
  id = "foo"
}
`
