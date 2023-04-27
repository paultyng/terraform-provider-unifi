package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFirewallRuleset_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	rulesMap := make(map[int]int, 6)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRulesetConfig(name, [3]int{1, 2, 3}, [3]int{4, 5, 6}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{1, 2, 3}, []int{4, 5, 6}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfig(name, [3]int{3, 1, 2}, [3]int{6, 4, 5}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{3, 1, 2}, []int{6, 4, 5}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfig(name, [3]int{1, 3, 2}, [3]int{4, 6, 5}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{1, 3, 2}, []int{4, 6, 5}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfig(name, [3]int{1, 3, 5}, [3]int{4, 6, 2}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{1, 3, 5}, []int{4, 6, 2}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfig(name, [3]int{6, 3, 5}, [3]int{4, 1, 2}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{6, 3, 5}, []int{4, 1, 2}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfig(name, [3]int{1, 2, 3}, [3]int{4, 5, 6}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{1, 2, 3}, []int{4, 5, 6}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfigRemove(name, [2]int{1, 2}, [2]int{4, 5}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 4, 5}, &rulesMap),
					testAccCheckRulesOrder([]int{1, 2}, []int{4, 5}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfig(name, [3]int{1, 2, 3}, [3]int{4, 5, 6}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{1, 2, 3}, []int{4, 5, 6}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfigRemove(name, [2]int{2, 1}, [2]int{5, 4}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 4, 5}, &rulesMap),
					testAccCheckRulesOrder([]int{2, 1}, []int{5, 4}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfig(name, [3]int{1, 2, 3}, [3]int{4, 5, 6}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{1, 2, 3}, []int{4, 5, 6}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.lan_in"),
			{
				Config: testAccFirewallRulesetConfigRemove(name, [2]int{4, 1}, [2]int{2, 5}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 4, 5}, &rulesMap),
					testAccCheckRulesOrder([]int{4, 1}, []int{2, 5}, &rulesMap),
				),
			},
		},
	})
}

func TestAccFirewallRuleset_unrelated_unmanaged_rules(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	rulesMap := make(map[int]int, 6)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRulesetConfigUnrelatedUnmanaged(name, [3]int{1, 2, 3}, [3]int{4, 5, 6}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{1, 2, 3}, []int{4, 5, 6}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.guest_in"),
			{
				Config: testAccFirewallRulesetConfigUnrelatedUnmanaged(name, [3]int{4, 1, 6}, [3]int{3, 5, 2}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulesExist([]int{1, 2, 3, 4, 5, 6}, &rulesMap),
					testAccCheckRulesOrder([]int{4, 1, 6}, []int{3, 5, 2}, &rulesMap),
				),
			},
			importStep("unifi_firewall_ruleset.guest_in"),
		},
	})
}

func TestAccFirewallRuleset_unmanaged_rules(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccFirewallRulesetConfigMissingRule(name),
				ExpectError: regexp.MustCompile("The set of existing rule IDs in ruleset LAN_OUT of site default does not match the set of managed rule IDs"),
			},
			{
				Config:      testAccFirewallRulesetConfigDuplicateRule(name),
				ExpectError: regexp.MustCompile("The set of existing rule IDs in ruleset LAN_OUT of site default does not match the set of managed rule IDs"),
			},
			{
				Config:      testAccFirewallRulesetConfigMissingDuplicateRule(name),
				ExpectError: regexp.MustCompile("The set of existing rule IDs in ruleset LAN_OUT of site default does not match the set of managed rule IDs"),
			},
		},
	})
}

func testAccFirewallRulesetConfig(name string, preOrder [3]int, postOrder [3]int) string {
	return fmt.Sprintf(`
resource "unifi_firewall_rule" "test1" {
	name    = "%[1]s-1"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.1"
}

resource "unifi_firewall_rule" "test2" {
	name    = "%[1]s-2"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.2"
}

resource "unifi_firewall_rule" "test3" {
	name    = "%[1]s-3"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.3"
}

resource "unifi_firewall_rule" "test4" {
	name    = "%[1]s-4"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.4"
}

resource "unifi_firewall_rule" "test5" {
	name    = "%[1]s-5"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.5"
}

resource "unifi_firewall_rule" "test6" {
	name    = "%[1]s-6"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.6"
}

resource "unifi_firewall_ruleset" "lan_in" {
	ruleset = "LAN_IN"
	before_predefined = [
		unifi_firewall_rule.test%[2]d.id,
		unifi_firewall_rule.test%[3]d.id,
		unifi_firewall_rule.test%[4]d.id,
	]
	after_predefined = [
		unifi_firewall_rule.test%[5]d.id,
		unifi_firewall_rule.test%[6]d.id,
		unifi_firewall_rule.test%[7]d.id,
	]
}
`, name, preOrder[0], preOrder[1], preOrder[2], postOrder[0], postOrder[1], postOrder[2])
}

func testAccFirewallRulesetConfigRemove(name string, preOrder [2]int, postOrder [2]int) string {
	return fmt.Sprintf(`
resource "unifi_firewall_rule" "test1" {
	name    = "%[1]s-1"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.1"
}

resource "unifi_firewall_rule" "test2" {
	name    = "%[1]s-2"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.2"
}

resource "unifi_firewall_rule" "test4" {
	name    = "%[1]s-4"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.4"
}

resource "unifi_firewall_rule" "test5" {
	name    = "%[1]s-5"
	action  = "accept"
	ruleset = "LAN_IN"
	dst_address = "192.168.1.5"
}

resource "unifi_firewall_ruleset" "lan_in" {
	ruleset = "LAN_IN"
	before_predefined = [
		unifi_firewall_rule.test%[2]d.id,
		unifi_firewall_rule.test%[3]d.id,
	]
	after_predefined = [
		unifi_firewall_rule.test%[4]d.id,
		unifi_firewall_rule.test%[5]d.id,
	]
}
`, name, preOrder[0], preOrder[1], postOrder[0], postOrder[1])
}

func testAccFirewallRulesetConfigUnrelatedUnmanaged(name string, preOrder [3]int, postOrder [3]int) string {
	return fmt.Sprintf(`
resource "unifi_firewall_rule" "testunrelated" {
	name    = "%[1]s-1"
	action  = "accept"
	ruleset = "GUEST_LOCAL"
	src_address = "192.168.1.6"
}

resource "unifi_firewall_rule" "test1" {
	name    = "%[1]s-1"
	action  = "accept"
	ruleset = "GUEST_IN"
	dst_address = "192.168.1.1"
}

resource "unifi_firewall_rule" "test2" {
	name    = "%[1]s-2"
	action  = "accept"
	ruleset = "GUEST_IN"
	dst_address = "192.168.1.2"
}

resource "unifi_firewall_rule" "test3" {
	name    = "%[1]s-3"
	action  = "accept"
	ruleset = "GUEST_IN"
	dst_address = "192.168.1.3"
}

resource "unifi_firewall_rule" "test4" {
	name    = "%[1]s-4"
	action  = "accept"
	ruleset = "GUEST_IN"
	dst_address = "192.168.1.4"
}

resource "unifi_firewall_rule" "test5" {
	name    = "%[1]s-5"
	action  = "accept"
	ruleset = "GUEST_IN"
	dst_address = "192.168.1.5"
}

resource "unifi_firewall_rule" "test6" {
	name    = "%[1]s-6"
	action  = "accept"
	ruleset = "GUEST_IN"
	dst_address = "192.168.1.6"
}

resource "unifi_firewall_ruleset" "guest_in" {
	ruleset = "GUEST_IN"
	before_predefined = [
		unifi_firewall_rule.test%[2]d.id,
		unifi_firewall_rule.test%[3]d.id,
		unifi_firewall_rule.test%[4]d.id,
	]
	after_predefined = [
		unifi_firewall_rule.test%[5]d.id,
		unifi_firewall_rule.test%[6]d.id,
		unifi_firewall_rule.test%[7]d.id,
	]
}
`, name, preOrder[0], preOrder[1], preOrder[2], postOrder[0], postOrder[1], postOrder[2])
}

func testAccFirewallRulesetConfigMissingRule(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_rule" "test1" {
	name    = "%[1]s-1"
	action  = "accept"
	ruleset = "LAN_OUT"
	dst_address = "192.168.1.1"
}

resource "unifi_firewall_rule" "test2" {
	name    = "%[1]s-2"
	action  = "accept"
	ruleset = "LAN_OUT"
	dst_address = "192.168.1.2"
}

resource "unifi_firewall_rule" "test3" {
	name    = "%[1]s-3"
	action  = "accept"
	ruleset = "LAN_OUT"
	dst_address = "192.168.1.3"
}

resource "unifi_firewall_ruleset" "lan_out" {
	# Ensure the test always passes, but creating a new rule and not having
	# it listed in the same run can actually cause an API error sometimes.
	depends_on = [unifi_firewall_rule.test3]
	ruleset = "LAN_OUT"
	before_predefined = [
		unifi_firewall_rule.test2.id,
		unifi_firewall_rule.test1.id,
	]
}
`, name)
}

func testAccFirewallRulesetConfigDuplicateRule(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_rule" "test1" {
	name    = "%[1]s-1"
	action  = "accept"
	ruleset = "LAN_OUT"
	dst_address = "192.168.1.1"
}

resource "unifi_firewall_rule" "test2" {
	name    = "%[1]s-2"
	action  = "accept"
	ruleset = "LAN_OUT"
	dst_address = "192.168.1.2"
}

resource "unifi_firewall_rule" "test3" {
	name    = "%[1]s-3"
	action  = "accept"
	ruleset = "LAN_OUT"
	dst_address = "192.168.1.3"
}

resource "unifi_firewall_ruleset" "lan_out" {
	ruleset = "LAN_OUT"
	before_predefined = [
		unifi_firewall_rule.test1.id,
		unifi_firewall_rule.test1.id,
		unifi_firewall_rule.test2.id,
		unifi_firewall_rule.test3.id,
	]
}
`, name)
}

func testAccFirewallRulesetConfigMissingDuplicateRule(name string) string {
	return fmt.Sprintf(`
resource "unifi_firewall_rule" "test1" {
	name    = "%[1]s-1"
	action  = "accept"
	ruleset = "LAN_OUT"
	dst_address = "192.168.1.1"
}

resource "unifi_firewall_rule" "test2" {
	name    = "%[1]s-2"
	action  = "accept"
	ruleset = "LAN_OUT"
	dst_address = "192.168.1.2"
}

resource "unifi_firewall_rule" "test3" {
	name    = "%[1]s-3"
	action  = "accept"
	ruleset = "LAN_OUT"
	dst_address = "192.168.1.3"
}

resource "unifi_firewall_ruleset" "lan_out" {
	ruleset = "LAN_OUT"
	before_predefined = [
		unifi_firewall_rule.test1.id,
		unifi_firewall_rule.test1.id,
		unifi_firewall_rule.test2.id,
	]
}
`, name)
}

func testAccCheckRulesExist(rules []int, rulesMap *map[int]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		site := "default"
		for _, ir := range rules {
			resourceName := "unifi_firewall_rule.test" + strconv.Itoa(ir)
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("Rule not found: %s", resourceName)
			}
			if rs.Primary.ID == "" {
				return fmt.Errorf("Rule ID is not set")
			}
			rule, err := testClient.GetFirewallRule(ctx, site, rs.Primary.ID)
			if err != nil {
				return err
			}
			(*rulesMap)[ir] = rule.RuleIndex
		}
		return nil
	}
}

func testAccCheckRulesOrder(preRules []int, postRules []int, rulesMap *map[int]int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		for i, ir := range preRules {
			index, ok := (*rulesMap)[ir]
			if !ok {
				return fmt.Errorf("Internal error: could not access rule information for rule test%d", ir)
			}
			expected := 2000 + i
			if index != expected {
				return fmt.Errorf("Rule test%d: expected %d, actual %d", ir, index, expected)
			}
		}
		for i, ir := range postRules {
			index, ok := (*rulesMap)[ir]
			if !ok {
				return fmt.Errorf("Internal error: could not access rule information for rule test%d", ir)
			}
			expected := 4000 + i
			if index != expected {
				return fmt.Errorf("Rule test%d: expected %d, actual %d", ir, index, expected)
			}
		}
		return nil
	}
}
