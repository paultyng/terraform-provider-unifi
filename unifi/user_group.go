package unifi

func (c *Client) ListUserGroup(site string) ([]UserGroup, error) {
	return c.listUserGroup(site)
}

func (c *Client) GetUserGroup(site, id string) (*UserGroup, error) {
	return c.getUserGroup(site, id)
}

func (c *Client) DeleteUserGroup(site, id string) error {
	return c.deleteUserGroup(site, id)
}

func (c *Client) CreateUserGroup(site string, d *UserGroup) (*UserGroup, error) {
	return c.createUserGroup(site, d)
}

func (c *Client) UpdateUserGroup(site string, d *UserGroup) (*UserGroup, error) {
	return c.updateUserGroup(site, d)
}
