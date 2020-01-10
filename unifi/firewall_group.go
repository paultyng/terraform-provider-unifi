package unifi

func (c *Client) ListFirewallGroup(site string) ([]FirewallGroup, error) {
	return c.listFirewallGroup(site)
}

func (c *Client) GetFirewallGroup(site, id string) (*FirewallGroup, error) {
	return c.getFirewallGroup(site, id)
}

func (c *Client) DeleteFirewallGroup(site, id string) error {
	return c.deleteFirewallGroup(site, id)
}

func (c *Client) CreateFirewallGroup(site string, d *FirewallGroup) (*FirewallGroup, error) {
	return c.createFirewallGroup(site, d)
}

func (c *Client) UpdateFirewallGroup(site string, d *FirewallGroup) (*FirewallGroup, error) {
	return c.updateFirewallGroup(site, d)
}
