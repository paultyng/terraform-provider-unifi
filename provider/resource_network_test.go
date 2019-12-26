package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetwork_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: providers,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

const testAccNetworkConfig = `
variable "subnet" {
	default = "10.0.202.1/24"
}

resource "unifi_network" "test" {
	name    = "foo"
	purpose = "corporate"

	subnet       = var.subnet
	vlan_id      = 202
	dhcp_start   = cidrhost(var.subnet, 6)
	dhcp_stop    = cidrhost(var.subnet, 254)
	dhcp_enabled = true
}
`
