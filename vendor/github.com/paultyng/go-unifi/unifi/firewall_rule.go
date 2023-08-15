package unifi

import (
	"context"
)

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
