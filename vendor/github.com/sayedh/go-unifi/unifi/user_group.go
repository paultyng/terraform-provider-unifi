package unifi

import "context"

func (c *Client) ListUserGroup(ctx context.Context, site string) ([]UserGroup, error) {
	return c.listUserGroup(ctx, site)
}

func (c *Client) GetUserGroup(ctx context.Context, site, id string) (*UserGroup, error) {
	return c.getUserGroup(ctx, site, id)
}

func (c *Client) DeleteUserGroup(ctx context.Context, site, id string) error {
	return c.deleteUserGroup(ctx, site, id)
}

func (c *Client) CreateUserGroup(ctx context.Context, site string, d *UserGroup) (*UserGroup, error) {
	return c.createUserGroup(ctx, site, d)
}

func (c *Client) UpdateUserGroup(ctx context.Context, site string, d *UserGroup) (*UserGroup, error) {
	return c.updateUserGroup(ctx, site, d)
}
