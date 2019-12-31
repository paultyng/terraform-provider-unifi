package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataWLANGroup_default(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: providers,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDataWLANGroupConfig_default,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
		},
	})
}

const testAccDataWLANGroupConfig_default = `
data "unifi_wlan_group" "default" {
}
`
