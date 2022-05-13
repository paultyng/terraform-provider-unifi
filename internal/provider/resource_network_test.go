package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetwork_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	vlanID1 := getTestVLAN(t)
	vlanID2 := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig(name, vlanID1, true, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "domain_name", "foo.local"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlanID1)),
					resource.TestCheckResourceAttr("unifi_network.test", "igmp_snooping", "true"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig(name, vlanID2, false, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlanID2)),
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
	vlanID := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig(name, vlanID, true, nil),
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
	vlanID := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig(name, vlanID, true, []string{"192.168.1.101"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.0", "192.168.1.101"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig(name, vlanID, true, []string{"192.168.1.101", "192.168.1.102"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.0", "192.168.1.101"),
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.1", "192.168.1.102"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfig(name, vlanID, true, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.#", "0"),
				),
			},
			{
				Config: testAccNetworkConfig(name, vlanID, true, []string{"192.168.1.101"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_dns.0", "192.168.1.101"),
				),
			},
		},
	})
}

func TestAccNetwork_dhcp_boot(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	vlanID := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfigDHCPBoot(name, vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// TODO: ...
				),
			},
			importStep("unifi_network.test"),
		},
	})
}

func TestAccNetwork_v6(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	vlanID1 := getTestVLAN(t)
	vlanID2 := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfigV6(name, vlanID1, "static", "fd6a:37be:e362::1/64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "domain_name", "foo.local"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlanID1)),
					resource.TestCheckResourceAttr("unifi_network.test", "ipv6_static_subnet", "fd6a:37be:e362::1/64"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfigV6(name, vlanID2, "static", "fd6a:37be:e363::1/64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", strconv.Itoa(vlanID2)),
					resource.TestCheckResourceAttr("unifi_network.test", "ipv6_static_subnet", "fd6a:37be:e363::1/64"),
				),
			},
			importStep("unifi_network.test"),
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
		},
	})
}

func TestAccNetwork_differentSite(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	vlanID1 := getTestVLAN(t)
	vlanID2 := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkWithSiteConfig(name, vlanID1),
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
				Config: testAccNetworkWithSiteConfig(name, vlanID2),
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
	vlanID1 := getTestVLAN(t)
	vlanID2 := getTestVLAN(t)
	vlanID3 := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Apply and import network by name.
			{
				Config: testAccNetworkConfig(name, vlanID1, true, nil),
			},
			{
				Config:            testAccNetworkConfig(name, vlanID1, true, nil),
				ResourceName:      "unifi_network.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("name=%s", name),
			},
			// Apply and test errors.
			{
				Config: testAccNetworkWithDuplicateNames(vlanID2, vlanID3, "DUPLICATE_NAME"),
			},
			// Test error on name that doesn't exist.
			{
				Config:            testAccNetworkWithDuplicateNames(vlanID2, vlanID3, "DUPLICATE_NAME"),
				ResourceName:      "unifi_network.test1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "name=BAD_NAME",
				ExpectError:       regexp.MustCompile("BAD_NAME"),
			},
			// Test error on multiple matches.
			{
				Config:            testAccNetworkWithDuplicateNames(vlanID2, vlanID3, "DUPLICATE_NAME"),
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
	vlanID := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfigDHCPRelay(name, vlanID, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "dhcp_relay_enabled", "true"),
				),
			},
			importStep("unifi_network.test"),
			{
				Config: testAccNetworkConfigDHCPRelay(name, vlanID, false),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkVlanOnly(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "101"),
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

func TestAccNetwork_wanGateway(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkWanGateway(name, "null"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "wan_gateway", ""),
				),
			},
			{
				ResourceName:      "unifi_network.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_network.test"),
				ImportStateVerify: true,
			},
			{
				Config:      testAccNetworkWanGateway(name, `""`),
				ExpectError: regexp.MustCompile("expected wan_gateway to contain a valid IPv4 address"),
			},
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

func testAccNetworkConfigDHCPBoot(name string, vlan int) string {
	return fmt.Sprintf(`
locals {
	subnet  = cidrsubnet("10.0.0.0/8", 6, %[2]d)
	vlan_id = %[2]d
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
`, name, vlan)
}

func testAccNetworkConfig(name string, vlan int, igmpSnoop bool, dhcpDNS []string) string {
	return fmt.Sprintf(`
locals {
	subnet  = cidrsubnet("10.0.0.0/8", 6, %[2]d)
	vlan_id = %[2]d
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
	igmp_snooping = %[3]t

	dhcp_dns = [%[4]s]
}
`, name, vlan, igmpSnoop, strings.Join(quoteStrings(dhcpDNS), ","))
}

func testAccNetworkConfigV6(name string, vlan int, ipv6Type string, ipv6Subnet string) string {
	return fmt.Sprintf(`
locals {
	subnet  = cidrsubnet("10.0.0.0/8", 6, %[2]d)
	vlan_id = %[2]d
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

	ipv6_interface_type = "%[3]s"
	ipv6_static_subnet  = "%[4]s"
	ipv6_ra_enable      = true
}
`, name, vlan, ipv6Type, ipv6Subnet)
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

func testAccNetworkWithSiteConfig(name string, vlan int) string {
	return fmt.Sprintf(`
locals {
	subnet  = cidrsubnet("10.0.0.0/8", 6, %[2]d)
	vlan_id = %[2]d
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
`, name, vlan)
}

func testAccNetworkWithDuplicateNames(vlan1, vlan2 int, networkName string) string {
	return fmt.Sprintf(`
locals {
	subnet1  = cidrsubnet("10.0.0.0/8", 6, %[1]d)
	vlan_id1 = %[1]d
	subnet2  = cidrsubnet("10.0.0.0/8", 6, %[2]d)
	vlan_id2 = %[2]d
}

resource "unifi_network" "test1" {
	name    = "%[3]s"
	purpose = "corporate"

	subnet  = local.subnet1
	vlan_id = local.vlan_id1
}

resource "unifi_network" "test2" {
	name    = "%[3]s"
	purpose = "corporate"

	subnet  = local.subnet2
	vlan_id = local.vlan_id2
}
`, vlan1, vlan2, networkName)
}

func testAccNetworkConfigDHCPRelay(name string, vlan int, dhcpRelay bool) string {
	return fmt.Sprintf(`
locals {
	subnet  = cidrsubnet("10.0.0.0/8", 6, %[2]d)
	vlan_id = %[2]d
}

resource "unifi_network" "test" {
	name    = "%[1]s"
	purpose = "corporate"

	subnet      = local.subnet
	vlan_id     = local.vlan_id
	domain_name = "foo.local"
	
	dhcp_relay_enabled = %[3]t
}
`, name, vlan, dhcpRelay)
}

func testAccNetworkVlanOnly(name string) string {
	return fmt.Sprintf(`
resource "unifi_site" "test" {
  description = "%[1]s"
}

resource "unifi_network" "test" {
  site    = unifi_site.test.name
  name    = "test"
  purpose = "vlan-only"
  vlan_id = 101
}
`, name)
}

func testAccNetworkWanGateway(site string, wanGatewayHCL string) string {
	return fmt.Sprintf(`
resource "unifi_site" "test" {
  description = "%[1]s"
}

resource "unifi_network" "test" {
  site    = unifi_site.test.name
  name       = "test"
  purpose = "vlan-only"
  vlan_id    = 107
  wan_gateway = %[2]s
}
`, site, wanGatewayHCL)
}
