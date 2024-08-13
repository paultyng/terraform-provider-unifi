package unifi

import (
	"context"
	"fmt"
)

type FirewallRuleIndexUpdate struct {
	ID        string `json:"_id"`
	RuleIndex int    `json:"rule_index,string"`
}

func (c *Client) ListFirewallRule(ctx context.Context, site string) ([]FirewallRule, error) {
	return c.listFirewallRule(ctx, site)
}

func (c *Client) GetFirewallRule(ctx context.Context, site, id string) (*FirewallRule, error) {
	return c.getFirewallRule(ctx, site, id)
}

func (c *Client) DeleteFirewallRule(ctx context.Context, site, id string) error {
	return c.deleteFirewallRule(ctx, site, id)
}

func (c *Client) CreateFirewallRule(ctx context.Context, site string, d *FirewallRule) (*FirewallRule, error) {
	return c.createFirewallRule(ctx, site, d)
}

func (c *Client) UpdateFirewallRule(ctx context.Context, site string, d *FirewallRule) (*FirewallRule, error) {
	return c.updateFirewallRule(ctx, site, d)
}

func (c *Client) ReorderFirewallRules(ctx context.Context, site, ruleset string, reorder []FirewallRuleIndexUpdate) error {
	reqBody := struct {
		Cmd     string                    `json:"cmd"`
		Ruleset string                    `json:"ruleset"`
		Rules   []FirewallRuleIndexUpdate `json:"rules"`
	}{
		Cmd:     "reorder",
		Ruleset: ruleset,
		Rules:   reorder,
	}
	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/cmd/firewall", site), reqBody, nil)
	if err != nil {
		return err
	}

	return nil
}
