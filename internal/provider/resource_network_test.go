package provider

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetwork_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet1, vlan1 := getTestVLAN(t)
	subnet2, vlan2 := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig(name, subnet1, vlan1, true, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "domain_name", "foo.local"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlan1)),
					resource.TestCheckResourceAttr("unifi_network.test", "igmp_snooping", "true"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig(name, subnet2, vlan2, false, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlan2)),
					resource.TestCheckResourceAttr("unifi_network.test", "igmp_snooping", "false"),
				),
			},
			importStep("unifi_network.test"),
			// re-test import here with default site, but full ID string
			{
				ResourceName:      "unifi_network.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_network.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetwork_weird_cidr(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig(name, subnet, vlan, true, nil),
				Check:  resource.ComposeTestCheckFunc(
				// TODO: ...
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

func TestAccNetwork_dhcp_dns(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig(name, subnet, vlan, true, []string{"192.168.1.101"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.0", "192.168.1.101"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig(name, subnet, vlan, true, []string{"192.168.1.101", "192.168.1.102"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.0", "192.168.1.101"),
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.1", "192.168.1.102"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig(name, subnet, vlan, true, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.#", "0"),
				),
			},
			{
				Config: testAccNetworkConfig(name, subnet, vlan, true, []string{"192.168.1.101"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.0", "192.168.1.101"),
				),
			},
		},
	})
}

func TestAccNetwork_dhcp_boot(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfigDHCPBoot(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// TODO: ...
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

func TestAccNetwork_v6(t *testing.T) {
	t.Skip()
	name := acctest.RandomWithPrefix("tfacc")
	subnet1, vlan1 := getTestVLAN(t)
	subnet2, vlan2 := getTestVLAN(t)
	subnet3, vlan3 := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfigV6(name, subnet1, vlan1, "static", "fd6a:37be:e362::1/64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "domain_name", "foo.local"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlan1)),
					resource.TestCheckResourceAttr("unifi_network.test", "ipv6_static_subnet", "fd6a:37be:e362::1/64"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfigV6(name, subnet2, vlan2, "static", "fd6a:37be:e363::1/64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlan2)),
					resource.TestCheckResourceAttr("unifi_network.test", "ipv6_static_subnet", "fd6a:37be:e363::1/64"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfigDhcpV6(
					name,
					subnet3,
					vlan3,
					"fd6a:37be:e364::1/64",
					"fd6a:37be:e364::2",
					"fd6a:37be:e364::7d1",
					[]string{"2001:4860:4860::8888", "2001:4860:4860::8844"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlan3)),
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_v6_start", "fd6a:37be:e364::2"),
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_v6_stop", "fd6a:37be:e364::7d1"),
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_v6_lease", strconv.Itoa(12*60*60)),
				),
			},
			{
				Config: testAccNetworkConfigDhcpV6(
					name,
					subnet3,
					vlan3,
					"fd6a:37be:e365::1/64",
					"fd6a:37be:e364::2",
					"fd6a:37be:e364::7d1",
					[]string{"2001:4860:4860::8888", "2001:4860:4860::8844"}),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("api.err.InvalidDHCPv6Range")),
			},
		},
	})
}

func TestAccNetwork_wan(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testWanNetworkConfig(name, "WAN", "pppoe", "192.168.1.1", 1, "username", "password", "8.8.8.8", "4.4.4.4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_networkgroup", "WAN"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_type", "pppoe"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_ip", "192.168.1.1"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_egress_qos", "1"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_username", "username"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "x_wan_password", "password"),

					resource.TestCheckOutput("wan_dns1", "8.8.8.8"),
					resource.TestCheckOutput("wan_dns2", "4.4.4.4"),
				),
			},
			importStep("unifi_network.wan_test"),
			// remove qos
			{
				Config: testWanNetworkConfig(name, "WAN", "pppoe", "192.168.1.1", 0, "username", "password", "8.8.8.8", "4.4.4.4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_networkgroup", "WAN"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_type", "pppoe"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_ip", "192.168.1.1"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_egress_qos", "0"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_username", "username"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "x_wan_password", "password"),

					resource.TestCheckOutput("wan_dns1", "8.8.8.8"),
					resource.TestCheckOutput("wan_dns2", "4.4.4.4"),
				),
			},
			importStep("unifi_network.wan_test"),
			{
				Config: testWanNetworkConfig(name, "WAN", "pppoe", "192.168.1.1", 1, "username", "password", "8.8.8.8", "4.4.4.4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_networkgroup", "WAN"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_type", "pppoe"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_ip", "192.168.1.1"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_egress_qos", "1"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_username", "username"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "x_wan_password", "password"),

					resource.TestCheckOutput("wan_dns1", "8.8.8.8"),
					resource.TestCheckOutput("wan_dns2", "4.4.4.4"),
				),
			},
			importStep("unifi_network.wan_test"),
			{
				Config:      testWanV6NetworkConfig(name, "dhcpv6", 47),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("expected wan_dhcp_v6_pd_size to be in the range (48 - 64)")),
			},
			{
				Config:      testWanV6NetworkConfig(name, "invalid", 48),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("invalid value for wan_type_v6")),
			},
			{
				Config: testWanV6NetworkConfig(name, "dhcpv6", 48),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_type_v6", "dhcpv6"),
					resource.TestCheckResourceAttr("unifi_network.wan_test", "wan_dhcp_v6_pd_size", "48"),
				),
			},
		},
	})
}

func TestAccNetwork_differentSite(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet1, vlan1 := getTestVLAN(t)
	subnet2, vlan2 := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkWithSiteConfig(name, subnet1, vlan1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("unifi_network.test", "site", "unifi_site.test", "name"),
				),
			},
			{
				ResourceName:      "unifi_network.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_network.test"),
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkWithSiteConfig(name, subnet2, vlan2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("unifi_network.test", "site", "unifi_site.test", "name"),
				),
			},
			{
				ResourceName:      "unifi_network.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_network.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetwork_importByName(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet1, vlan1 := getTestVLAN(t)
	subnet2, vlan2 := getTestVLAN(t)
	subnet3, vlan3 := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Apply and import network by name.
			{
				Config: testAccNetworkConfig(name, subnet1, vlan1, true, nil),
			},
			{
				Config:            testAccNetworkConfig(name, subnet1, vlan1, true, nil),
				ResourceName:      "unifi_network.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("name=%s", name),
			},
			// Apply and test errors.
			{
				Config: testAccNetworkWithDuplicateNames(subnet2, vlan2, subnet3, vlan3, "DUPLICATE_NAME"),
			},
			// Test error on name that doesn't exist.
			{
				Config:            testAccNetworkWithDuplicateNames(subnet2, vlan2, subnet3, vlan3, "DUPLICATE_NAME"),
				ResourceName:      "unifi_network.test1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "name=BAD_NAME",
				ExpectError:       regexp.MustCompile("BAD_NAME"),
			},
			// Test error on multiple matches.
			{
				Config:            testAccNetworkWithDuplicateNames(subnet2, vlan2, subnet3, vlan3, "DUPLICATE_NAME"),
				ResourceName:      "unifi_network.test1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "name=DUPLICATE_NAME",
				ExpectError:       regexp.MustCompile("DUPLICATE_NAME"),
			},
		},
	})
}

func TestAccNetwork_dhcpRelay(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfigDHCPRelay(name, subnet, vlan, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_relay_enabled", "true"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfigDHCPRelay(name, subnet, vlan, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_relay_enabled", "false"),
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

func TestAccNetwork_vlanOnly(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	_, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkVlanOnly(name, vlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlan)),
				),
			},
			{
				ResourceName:      "unifi_network.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_network.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetwork_mdns(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckMinVersion(t, controllerV7)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfigMDNS(name, subnet, vlan, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "multicast_dns", "true"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfigMDNS(name, subnet, vlan, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "multicast_dns", "false"),
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

func testAccNetworkConfigDHCPBoot(name string, subnet *net.IPNet, vlan int) string {
	return fmt.Sprintf(`
locals {
	subnet  = "%[2]s"
	vlan_id = %[3]d
}

resource "unifi_network" "test" {
	name     = "%[1]s"
	purpose = "corporate"

	subnet        = local.subnet
	vlan_id       = local.vlan_id
	dhcp_start    = cidrhost(local.subnet, 6)
	dhcp_stop     = cidrhost(local.subnet, 254)
	dhcp_enabled  = true
	domain_name   = "foo.local"

	dhcpd_boot_enabled  = true
	dhcpd_boot_server   = "192.168.1.180"
	dhcpd_boot_filename = "test.boot"

	dhcp_dns = ["192.168.1.101", "192.168.1.102"]
}
`, name, subnet, vlan)
}

func testAccNetworkConfig(name string, subnet *net.IPNet, vlan int, igmpSnoop bool, dhcpDNS []string) string {
	return fmt.Sprintf(`
locals {
	subnet  = "%[2]s"
	vlan_id = %[3]d
}

resource "unifi_network" "test" {
	name    = "%[1]s"
	purpose = "corporate"

	subnet        = local.subnet
	vlan_id       = local.vlan_id
	dhcp_start    = cidrhost(local.subnet, 6)
	dhcp_stop     = cidrhost(local.subnet, 254)
	dhcp_enabled  = true
	domain_name   = "foo.local"
	igmp_snooping = %[4]t

	dhcp_dns = [%[5]s]
}
`, name, subnet, vlan, igmpSnoop, strings.Join(quoteStrings(dhcpDNS), ","))
}

func testAccNetworkConfigV6(name string, subnet *net.IPNet, vlan int, ipv6Type string, ipv6Subnet string) string {
	return fmt.Sprintf(`
locals {
	subnet  = "%[2]s"
	vlan_id = %[3]d
}
	
resource "unifi_network" "test" {
	name    = "%[1]s"
	purpose = "corporate"

	subnet        = local.subnet
	vlan_id       = local.vlan_id
	dhcp_start    = cidrhost(local.subnet, 6)
	dhcp_stop     = cidrhost(local.subnet, 254)
	dhcp_enabled  = true
	domain_name   = "foo.local"

	ipv6_interface_type = "%[4]s"
	ipv6_static_subnet  = "%[5]s"
	ipv6_ra_enable      = true
}
`, name, subnet, vlan, ipv6Type, ipv6Subnet)
}

func testWanNetworkConfig(name string, networkGroup string, wanType string, wanIP string, wanEgressQOS int, wanUsername string, wanPassword string, wanDNS1 string, wanDNS2 string) string {
	return fmt.Sprintf(`
resource "unifi_network" "wan_test" {
	name             = "%s"
	purpose          = "wan"
	wan_networkgroup = "%s"
	wan_type         = "%s"
	wan_ip           = "%s"
	wan_egress_qos   = %d
	wan_username     = "%s"
	x_wan_password   = "%s"

	wan_dns = ["%s", "%s"]
}

output "wan_dns1" {
	value = unifi_network.wan_test.wan_dns[0]
}

output "wan_dns2" {
	value = unifi_network.wan_test.wan_dns[1]
}
`, name, networkGroup, wanType, wanIP, wanEgressQOS, wanUsername, wanPassword, wanDNS1, wanDNS2)
}

func testWanV6NetworkConfig(name string, wanTypeV6 string, wanDhcpV6PdSize int) string {
	return fmt.Sprintf(`
resource "unifi_network" "wan_test" {
	name             = "%s"
	purpose          = "wan"
	wan_networkgroup = "WAN"
	wan_type         = "pppoe"
	wan_ip           = "192.168.1.1"
	wan_egress_qos   = 1
	wan_username     = "username"
	x_wan_password   = "password"

	wan_dns = ["8.8.8.8", "4.4.4.4"]

	wan_type_v6 = "%s"
	wan_dhcp_v6_pd_size = %d
}
`, name, wanTypeV6, wanDhcpV6PdSize)
}

func testAccNetworkWithSiteConfig(name string, subnet *net.IPNet, vlan int) string {
	return fmt.Sprintf(`
locals {
	subnet  = "%[2]s"
	vlan_id = %[3]d
}

resource "unifi_site" "test" {
  description = "%[1]s"
}

resource "unifi_network" "test" {
	site    = unifi_site.test.name
	name    = "%[1]s"
	purpose = "corporate"

	subnet        = local.subnet
	vlan_id       = local.vlan_id
	dhcp_start    = cidrhost(local.subnet, 6)
	dhcp_stop     = cidrhost(local.subnet, 254)
	dhcp_enabled  = true
	domain_name   = "foo.local"
	igmp_snooping = true
}
`, name, subnet, vlan)
}

func testAccNetworkWithDuplicateNames(subnet1 *net.IPNet, vlan1 int, subnet2 *net.IPNet, vlan2 int, networkName string) string {
	return fmt.Sprintf(`
locals {
	subnet1  = "%[1]s"
	vlan_id1 = %[2]d
	subnet2  = "%[3]s"
	vlan_id2 = %[4]d
}

resource "unifi_network" "test1" {
	name    = "%[5]s"
	purpose = "corporate"

	subnet  = local.subnet1
	vlan_id = local.vlan_id1
}

resource "unifi_network" "test2" {
	name    = "%[5]s"
	purpose = "corporate"

	subnet  = local.subnet2
	vlan_id = local.vlan_id2
}
`, subnet1, vlan1, subnet2, vlan2, networkName)
}

func testAccNetworkConfigDHCPRelay(name string, subnet *net.IPNet, vlan int, dhcpRelay bool) string {
	return fmt.Sprintf(`
locals {
	subnet  = "%[2]s"
	vlan_id = %[3]d
}

resource "unifi_network" "test" {
	name    = "%[1]s"
	purpose = "corporate"

	subnet      = local.subnet
	vlan_id     = local.vlan_id
	domain_name = "foo.local"
	
	dhcp_relay_enabled = %[4]t
}
`, name, subnet, vlan, dhcpRelay)
}

func testAccNetworkVlanOnly(name string, vlan int) string {
	return fmt.Sprintf(`
resource "unifi_site" "test" {
  description = "%[1]s"
}

resource "unifi_network" "test" {
  site    = unifi_site.test.name
  name    = "test"
  purpose = "vlan-only"
  vlan_id = %[2]d
}
`, name, vlan)
}

func testAccNetworkConfigDhcpV6(name string, subnet *net.IPNet, vlan int, gatewayIP string, dhcpdV6Start string, dhcpdV6Stop string, dhcpV6DNS []string) string {
	return fmt.Sprintf(`
locals {
  subnet  = "%[2]s"
	vlan_id = %[3]d
}

resource "unifi_network" "test" {
	name    = "%[1]s"
	purpose = "corporate"

	subnet        = local.subnet
	vlan_id       = local.vlan_id

	ipv6_static_subnet  = "%[4]s"

	dhcp_v6_dns_auto = false
	dhcp_v6_dns = [%[7]s]
	dhcp_v6_enabled = true
	dhcp_v6_start = "%[5]s"
	dhcp_v6_stop = "%[6]s"
	dhcp_v6_lease = 12 * 60 * 60
}
`, name, subnet, vlan, gatewayIP, dhcpdV6Start, dhcpdV6Stop, strings.Join(quoteStrings(dhcpV6DNS), ","))
}

func testAccNetworkConfigMDNS(name string, subnet *net.IPNet, vlan int, mdns bool) string {
	return fmt.Sprintf(`
resource "unifi_network" "test" {
	name    = "%[1]s"
	purpose = "corporate"
	subnet  = "%[2]s"
	vlan_id = %[3]d

	multicast_dns = %[4]t
}
`, name, subnet, vlan, mdns)
}
