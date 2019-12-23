package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetwork_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			// importStep("unifi_network.test"),
		},
	})
}

const testAccNetworkConfig = `
resource "unifi_network" "test" {
	name    = "foo"
	purpose = "corporate"

	subnet       = "10.0.202.1/24"
	vlan_id      = 202
	dhcp_start   = "10.0.202.6"
	dhcp_stop    = "10.0.202.254"
	dhcp_enabled = true
}
`
