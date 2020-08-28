package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetwork_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig("10.0.202.0/24", 202, true, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "domain_name", "foo.local"),
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.202.0/24"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "202"),
					resource.TestCheckResourceAttr("unifi_network.test", "igmp_snooping", "true"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig("10.0.203.0/24", 203, false, nil),
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
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig("10.0.204.3/24", 204, true, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.204.0/24"),
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

func TestAccNetwork_dhcp_dns(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig("10.0.205.0/24", 205, true, []string{"192.168.1.101"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.0", "192.168.1.101"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig("10.0.205.0/24", 205, true, []string{"192.168.1.101", "192.168.1.102"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.0", "192.168.1.101"),
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.1", "192.168.1.102"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig("10.0.205.0/24", 205, true, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("unifi_network.test", "dhcp_dns"),
				),
			},
			{
				Config: testAccNetworkConfig("10.0.205.0/24", 205, true, []string{"192.168.1.101"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.0", "192.168.1.101"),
				),
			},
		},
	})
}

func TestAccNetwork_v6(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfigV6("10.0.206.0/24", 206, "static", "fd6a:37be:e362::1/64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "domain_name", "foo.local"),
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.206.0/24"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "206"),
					resource.TestCheckResourceAttr("unifi_network.test", "ipv6_static_subnet", "fd6a:37be:e362::1/64"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfigV6("10.0.207.0/24", 207, "static", "fd6a:37be:e363::1/64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "10.0.207.0/24"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "207"),
					resource.TestCheckResourceAttr("unifi_network.test", "ipv6_static_subnet", "fd6a:37be:e363::1/64"),
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

// TODO: ipv6 prefix delegation test

func quoteStrings(src []string) []string {
	dst := make([]string, 0, len(src))
	for _, s := range src {
		dst = append(dst, fmt.Sprintf("%q", s))
	}
	return dst
}

func testAccNetworkConfig(subnet string, vlan int, igmpSnoop bool, dhcpDNS []string) string {
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

	dhcp_dns = [%s]
}
`, subnet, vlan, igmpSnoop, strings.Join(quoteStrings(dhcpDNS), ","))
}

func testAccNetworkConfigV6(subnet string, vlan int, ipv6Type string, ipv6Subnet string) string {
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

	ipv6_interface_type = "%s"
	ipv6_static_subnet = "%s"
	ipv6_ra_enable = true
}
`, subnet, vlan, ipv6Type, ipv6Subnet)
}
