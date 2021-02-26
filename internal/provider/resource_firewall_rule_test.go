package provider

import (
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

func TestAccFirewallRule_dst_port(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigWithPort,
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

// func TestAccFirewallRule_firewall_group(t *testing.T) {
// func TestAccFirewallRule_network(t *testing.T) {

const testAccFirewallRuleConfig = `
resource "unifi_firewall_group" "test" {
	name = "tf acc"
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
	dst_address = "192.168.1.1"
	dst_port    = 53
}
`

// resource "unifi_firewall_rule" "can_print_drop" {
// 	name    = "[tf] can-print (drop all)"
// 	action  = "drop"
// 	ruleset = "LAN_IN"

// 	rule_index = 2011

// 	protocol = "all"

// 	dst_address = "192.168.1.1"
// }
// `
