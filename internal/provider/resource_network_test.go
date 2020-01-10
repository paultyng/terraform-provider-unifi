package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetwork_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { preCheck(t) },
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig("10.0.202.1/24", 202),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "domain_name", "foo.local"),
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.202.1/24"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "202"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig("10.0.203.1/24", 203),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.203.1/24"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "203"),
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

func testAccNetworkConfig(subnet string, vlan int) string {
	return fmt.Sprintf(`
variable "subnet" {
	default = "%s"
}

resource "unifi_network" "test" {
	name    = "tfacc"
	purpose = "corporate"

	subnet       = var.subnet
	vlan_id      = %d
	dhcp_start   = cidrhost(var.subnet, 6)
	dhcp_stop    = cidrhost(var.subnet, 254)
	dhcp_enabled = true
	domain_name  = "foo.local"
}
`, subnet, vlan)
}
