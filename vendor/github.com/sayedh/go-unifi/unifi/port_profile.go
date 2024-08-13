package unifi

import (
	"context"
)

func (c *Client) ListPortProfile(ctx context.Context, site string) ([]PortProfile, error) {
	return c.listPortProfile(ctx, site)
}

func (c *Client) GetPortProfile(ctx context.Context, site, id string) (*PortProfile, error) {
	return c.getPortProfile(ctx, site, id)
}

func (c *Client) DeletePortProfile(ctx context.Context, site, id string) error {
	return c.deletePortProfile(ctx, site, id)
}

func (c *Client) CreatePortProfile(ctx context.Context, site string, d *PortProfile) (*PortProfile, error) {
	return c.createPortProfile(ctx, site, d)
}

func (c *Client) UpdatePortProfile(ctx context.Context, site string, d *PortProfile) (*PortProfile, error) {
	return c.updatePortProfile(ctx, site, d)
}
