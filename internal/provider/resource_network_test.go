package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/paultyng/terraform-provider-unifi/unifi"
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
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_network.test", "domain_name", "foo.local"),
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.202.1/24"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "202"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig("10.0.203.1/24", 203),
				Check: resource.ComposeTestCheckFunc(
					testCheckNetworkExists(t, "tfacc", nil),
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.203.1/24"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "203"),
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

func testCheckNetworkExists(t *testing.T, name string, network *unifi.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		networks, err := testClient.ListNetwork("default")
		if err != nil {
			return err
		}

		for _, net := range networks {
			if net.Name == name {
				if network != nil {
					*network = net
				}
				return nil
			}
		}

		return fmt.Errorf("unable to find network %q", name)
	}
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
