package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFirewallRule_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfig,
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckFirewallGroupExists(t, "name"),
				// ),
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_port(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigWithPort,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "src_port", "123"),
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "dst_port", "53"),
				),
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_icmp(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigWithICMP,
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_multiple_address_groups(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config:      testAccFirewallRuleConfigMultipleAddressGroups,
				ExpectError: regexp.MustCompile("firewall rule groups must be of different group types"),
			},
		},
	})
}

func TestAccFirewallRule_date_and_time(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigWithDateTime,
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_multiple_port_groups(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config:      testAccFirewallRuleConfigMultiplePortGroups,
				ExpectError: regexp.MustCompile("firewall rule groups must be of different group types"),
			},
		},
	})
}

func TestAccFirewallRule_address_and_port_group(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigAddressAndPortGroup,
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckFirewallGroupExists(t, "name"),
				// ),
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_IPv6_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigIPv6,
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_IPv6_dst_port(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigIPv6WithPort,
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

// func TestAccFirewallRule_firewall_group(t *testing.T) {
// func TestAccFirewallRule_network(t *testing.T) {

const testAccFirewallRuleConfigWithDateTime = `
resource "unifi_firewall_group" "test" {
	name = "tf acc rule date time"
	type = "address-group"

	members = ["192.168.1.1", "192.168.1.2"]
}

resource "unifi_firewall_rule" "test" {
	name    = "tf acc rule date time"
	action  = "accept"
	ruleset = "LAN_IN"

	rule_index = 2021

	protocol = "all"

	src_firewall_group_ids = [unifi_firewall_group.test.id]
	startdate = "1/1/2021"
	stopdate  = "12/31/2021"
	starttime = "00:00:00"
	endtime   = "23:59:59"
	weekdays  = "Mon,Tue,Wed,Thu,Fri"
	utc       = false
}
`
const testAccFirewallRuleConfig = `
resource "unifi_firewall_group" "test" {
	name = "tf acc rule"
	type = "address-group"

	members = ["192.168.1.1", "192.168.1.2"]
}

resource "unifi_firewall_rule" "test" {
	name    = "tf acc"
	action  = "accept"
	ruleset = "LAN_IN"

	rule_index = 2010

	protocol = "all"

	src_firewall_group_ids = [unifi_firewall_group.test.id]

	dst_address = "192.168.1.1"
}
`

const testAccFirewallRuleConfigWithPort = `
resource "unifi_firewall_rule" "test" {
	name    = "tf acc"
	action  = "accept"
	ruleset = "LAN_IN"

	rule_index = 2011

	protocol = "tcp"

	src_address = "192.168.3.3"
  src_port    = 123
	dst_address = "192.168.1.1"
	dst_port    = 53
}
`

const testAccFirewallRuleConfigWithICMP = `
resource "unifi_firewall_rule" "test" {
	name    = "tf acc"
	action  = "accept"
	ruleset = "LAN_LOCAL"

	rule_index = 2012

	protocol      = "icmp"
	icmp_typename = "echo-request"
}
`

const testAccFirewallRuleConfigMultipleAddressGroups = `
resource "unifi_firewall_group" "test_a" {
	name = "tf acc rule multiple address groups a"
	type = "address-group"

	members = ["192.168.1.1", "192.168.1.2"]
}

resource "unifi_firewall_group" "test_b" {
	name = "tf acc rule multiple address groups b"
	type = "address-group"

	members = ["192.168.1.3"]
}

resource "unifi_firewall_rule" "test" {
	name    = "tf acc"
	action  = "accept"
	ruleset = "LAN_IN"

	rule_index = 2013

	protocol = "all"

	src_firewall_group_ids = [
		unifi_firewall_group.test_a.id,
		unifi_firewall_group.test_b.id,
	]

	dst_address = "192.168.1.1"
}
`

const testAccFirewallRuleConfigMultiplePortGroups = `
resource "unifi_firewall_group" "test_a" {
	name = "tf acc rule multiple port groups a"
	type = "port-group"

	members = ["53"]
}

resource "unifi_firewall_group" "test_b" {
	name = "tf acc rule multiple port groups b"
	type = "port-group"

	members = ["80", "443"]
}

resource "unifi_firewall_rule" "test" {
	name    = "tf acc"
	action  = "accept"
	ruleset = "LAN_IN"

	rule_index = 2014

	protocol = "all"

	src_firewall_group_ids = [
		unifi_firewall_group.test_a.id,
		unifi_firewall_group.test_b.id,
	]

	dst_address = "192.168.1.1"
}
`

const testAccFirewallRuleConfigAddressAndPortGroup = `
resource "unifi_firewall_group" "test_a" {
	name = "tf acc rule address and port group a"
	type = "address-group"

	members = ["192.168.1.1", "192.168.1.2"]
}

resource "unifi_firewall_group" "test_b" {
	name = "tf acc rule address and port group b"
	type = "port-group"

	members = ["80", "443"]
}

resource "unifi_firewall_rule" "test" {
	name    = "tf acc"
	action  = "accept"
	ruleset = "LAN_IN"

	rule_index = 2015

	protocol = "all"

	src_firewall_group_ids = [
		unifi_firewall_group.test_a.id,
		unifi_firewall_group.test_b.id,
	]

	dst_address = "192.168.1.1"
}
`

const testAccFirewallRuleConfigIPv6 = `
resource "unifi_firewall_group" "test_a" {
	name = "tf acc rule IPv6 group a"
	type = "ipv6-address-group"

	members = ["fd6a:37be:e364::/64", "fd6a:37be:e365::/64",]
}

resource "unifi_firewall_group" "test_b" {
	name = "tf acc rule IPv6 group b"
	type = "ipv6-address-group"

	members = ["2001:4860:4860::8888", "2001:4860:4860::8844"]
}

resource "unifi_firewall_rule" "test" {
	name    = "tf acc"
	action  = "drop"
	ruleset = "LANv6_IN"

	rule_index = 2510

	protocol_v6 = "all"

	src_firewall_group_ids = [unifi_firewall_group.test_a.id]

	dst_firewall_group_ids = [unifi_firewall_group.test_b.id]
}
`

const testAccFirewallRuleConfigIPv6WithPort = `
resource "unifi_firewall_rule" "test" {
	name    = "tf acc"
	action  = "accept"
	ruleset = "LANv6_IN"

	rule_index = 2511

	protocol = "tcp"

	src_address_ipv6 = "fd6a:37be:e364::1/64"
	dst_address_ipv6 = "fd6a:37be:e364::2/64"
	dst_port    = 53
}
`
