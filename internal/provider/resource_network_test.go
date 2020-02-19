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
				Config: testAccNetworkConfig("10.0.202.0/24", 202, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "domain_name", "foo.local"),
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.202.0/24"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "202"),
					resource.TestCheckResourceAttr("unifi_network.test", "igmp_snooping", "true"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig("10.0.203.0/24", 203, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.203.0/24"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "203"),
					resource.TestCheckResourceAttr("unifi_network.test", "igmp_snooping", "false"),
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

func TestAccNetwork_weird_cidr(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { preCheck(t) },
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig("10.0.202.3/24", 202, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.202.0/24"),
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

func testAccNetworkConfig(subnet string, vlan int, igmpSnoop bool) string {
	return fmt.Sprintf(`
variable "subnet" {
	default = "%s"
}

resource "unifi_network" "test" {
	name    = "tfacc"
	purpose = "corporate"

	subnet        = var.subnet
	vlan_id       = %d
	dhcp_start    = cidrhost(var.subnet, 6)
	dhcp_stop     = cidrhost(var.subnet, 254)
	dhcp_enabled  = true
	domain_name   = "foo.local"
	igmp_snooping = %t
}
`, subnet, vlan, igmpSnoop)
}
