package unifi

func (c *Client) ListFirewallRule(site string) ([]FirewallRule, error) {
	return c.listFirewallRule(site)
}

func (c *Client) GetFirewallRule(site, id string) (*FirewallRule, error) {
	return c.getFirewallRule(site, id)
}

func (c *Client) DeleteFirewallRule(site, id string) error {
	return c.deleteFirewallRule(site, id)
}

func (c *Client) CreateFirewallRule(site string, d *FirewallRule) (*FirewallRule, error) {
	return c.createFirewallRule(site, d)
}

func (c *Client) UpdateFirewallRule(site string, d *FirewallRule) (*FirewallRule, error) {
	return c.updateFirewallRule(site, d)
}
