package unifi

func (c *Client) ListPortForward(site string) ([]PortForward, error) {
	return c.listPortForward(site)
}

func (c *Client) GetPortForward(site, id string) (*PortForward, error) {
	return c.getPortForward(site, id)
}

func (c *Client) DeletePortForward(site, id string) error {
	return c.deletePortForward(site, id)
}

func (c *Client) CreatePortForward(site string, d *PortForward) (*PortForward, error) {
	return c.createPortForward(site, d)
}

func (c *Client) UpdatePortForward(site string, d *PortForward) (*PortForward, error) {
	return c.updatePortForward(site, d)
}
