package provider

import (
	"testing"
)

func TestAccExampleDataSource(t *testing.T) {
	//resource.Test(t, resource.TestCase{
	//	PreCheck:                 func() { testAccPreCheck(t) },
	//	ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
	//	Steps: []resource.TestStep{
	//		// Read testing
	//		{
	//			Config: testAccExampleDataSourceConfig,
	//			Check: resource.ComposeAggregateTestCheckFunc(
	//				resource.TestCheckNoResourceAttr("data.pmwrp_brokers.default", "foo"),
	//			),
	//		},
	//	},
	//})
}

const testAccExampleDataSourceConfig = `data "pmwrp_brokers" "default" {}`
