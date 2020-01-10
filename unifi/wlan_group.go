package unifi

func (c *Client) ListWLANGroup(site string) ([]WLANGroup, error) {
	return c.listWLANGroup(site)
}

func (c *Client) GetWLANGroup(site, id string) (*WLANGroup, error) {
	return c.getWLANGroup(site, id)
}

func (c *Client) DeleteWLANGroup(site, id string) error {
	return c.deleteWLANGroup(site, id)
}

func (c *Client) CreateWLANGroup(site string, d *WLANGroup) (*WLANGroup, error) {
	return c.createWLANGroup(site, d)
}

func (c *Client) UpdateWLANGroup(site string, d *WLANGroup) (*WLANGroup, error) {
	return c.updateWLANGroup(site, d)
}
