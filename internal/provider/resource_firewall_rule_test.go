package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/paultyng/go-unifi/unifi"
)

func testAccCheckFirewallRuleDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "unifi_firewall_rule" {
			continue
		}

		rule, err := testClient.GetFirewallRule(ctx, rs.Primary.Attributes["site"], rs.Primary.ID)
		if rule != nil {
			return fmt.Errorf("Firewall rule still exists with ID %v", rs.Primary.ID)
		}
		if _, ok := err.(*unifi.NotFoundError); !ok {
			return err
		}
	}

	return nil
}

func TestAccFirewallRule_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfig(name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "name", name),
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "enabled", "true"),
				),
			},
			importStep("unifi_firewall_rule.test"),
			{
				Config: testAccFirewallRuleConfig(name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "enabled", "false"),
				),
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_port(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigWithPort(name),
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
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigWithICMP(name),
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_multiple_address_groups(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccFirewallRuleConfigMultipleAddressGroups(name),
				ExpectError: regexp.MustCompile("firewall rule groups must be of different group types"),
			},
		},
	})
}

func TestAccFirewallRule_multiple_port_groups(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccFirewallRuleConfigMultiplePortGroups(name),
				ExpectError: regexp.MustCompile("firewall rule groups must be of different group types"),
			},
		},
	})
}

func TestAccFirewallRule_address_and_port_group(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigAddressAndPortGroup(name),
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckFirewallGroupExists(t, "name"),
				// ),
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_IPv6_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigIPv6(name),
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_IPv6_dst_port(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigIPv6WithPort(name),
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func TestAccFirewallRule_computed_index(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfigNoIndex(name),
			},
			importStep("unifi_firewall_rule.test"),
		},
	})
}

func testAccFirewallRuleConfig(name string, enabled bool) string {
	return fmt.Sprintf(`
resource "unifi_firewall_group" "test" {
	name = "%[1]s"
	type = "address-group"

	members = ["192.168.1.1", "192.168.1.2"]
}

resource "unifi_firewall_rule" "test" {
	name    = "%[1]s"
	action  = "accept"
	ruleset = "LAN_IN"
  enabled = %[2]t

	rule_index = 2010

	protocol = "all"

	src_firewall_group_ids = [unifi_firewall_group.test.id]

	dst_address = "192.168.1.1"
}
`, name, enabled)
}

func testAccFirewallRuleConfigWithPort(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_rule" "test" {
	name    = "%s"
	action  = "accept"
	ruleset = "LAN_IN"

	rule_index = 2011

	protocol = "tcp"

	src_address = "192.168.3.3"
  src_port    = 123
	dst_address = "192.168.1.1"
	dst_port    = 53
}
`, name)
}

func testAccFirewallRuleConfigWithICMP(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_rule" "test" {
	name    = "%s"
	action  = "accept"
	ruleset = "LAN_LOCAL"

	rule_index = 2012

	protocol      = "icmp"
	icmp_typename = "echo-request"
}
`, name)
}

func testAccFirewallRuleConfigMultipleAddressGroups(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_group" "test_a" {
	name = "%[1]s-a"
	type = "address-group"

	members = ["192.168.1.1", "192.168.1.2"]
}

resource "unifi_firewall_group" "test_b" {
	name = "%[1]s-b"
	type = "address-group"

	members = ["192.168.1.3"]
}

resource "unifi_firewall_rule" "test" {
	name    = "%[1]s"
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
`, name)
}

func testAccFirewallRuleConfigMultiplePortGroups(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_group" "test_a" {
	name = "%[1]s-a"
	type = "port-group"

	members = ["53"]
}

resource "unifi_firewall_group" "test_b" {
	name = "%[1]s-b"
	type = "port-group"

	members = ["80", "443"]
}

resource "unifi_firewall_rule" "test" {
	name    = "%[1]s"
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
`, name)
}

func testAccFirewallRuleConfigAddressAndPortGroup(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_group" "test_a" {
	name = "%[1]s-a"
	type = "address-group"

	members = ["192.168.1.1", "192.168.1.2"]
}

resource "unifi_firewall_group" "test_b" {
	name = "%[1]s-b"
	type = "port-group"

	members = ["80", "443"]
}

resource "unifi_firewall_rule" "test" {
	name    = "%[1]s"
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
`, name)
}

func testAccFirewallRuleConfigIPv6(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_group" "test_a" {
	name = "%[1]s-a"
	type = "ipv6-address-group"

	members = ["fd6a:37be:e364::/64", "fd6a:37be:e365::/64"]
}

resource "unifi_firewall_group" "test_b" {
	name = "%[1]s-b"
	type = "ipv6-address-group"

	members = ["2001:4860:4860::8888", "2001:4860:4860::8844"]
}

resource "unifi_firewall_rule" "test" {
	name    = "%[1]s"
	action  = "drop"
	ruleset = "LANv6_IN"

	rule_index = 2510

	protocol_v6 = "all"

	src_firewall_group_ids = [unifi_firewall_group.test_a.id]

	dst_firewall_group_ids = [unifi_firewall_group.test_b.id]
}
`, name)
}

func testAccFirewallRuleConfigIPv6WithPort(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_rule" "test" {
	name    = "%s"
	action  = "accept"
	ruleset = "LANv6_IN"

	rule_index = 2511

	protocol = "tcp"

	src_address_ipv6 = "fd6a:37be:e364::1/64"
	dst_address_ipv6 = "fd6a:37be:e364::2/64"
	dst_port    = 53
}
`, name)
}

func testAccFirewallRuleConfigNoIndex(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_group" "test" {
	name = "%[1]s"
	type = "address-group"

	members = ["192.168.1.1", "192.168.1.2"]
}

resource "unifi_firewall_rule" "test" {
	name    = "%[1]s"
	action  = "accept"
	ruleset = "LAN_IN"

	protocol = "all"

	src_firewall_group_ids = [unifi_firewall_group.test.id]

	dst_address = "192.168.1.1"
}
`, name)
}
