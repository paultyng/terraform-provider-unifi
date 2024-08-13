package unifi

import (
	"context"
)

func (c *Client) ListRADIUSProfile(ctx context.Context, site string) ([]RADIUSProfile, error) {
	return c.listRADIUSProfile(ctx, site)
}

func (c *Client) GetRADIUSProfile(ctx context.Context, site, id string) (*RADIUSProfile, error) {
	return c.getRADIUSProfile(ctx, site, id)
}

func (c *Client) DeleteRADIUSProfile(ctx context.Context, site, id string) error {
	return c.deleteRADIUSProfile(ctx, site, id)
}

func (c *Client) CreateRADIUSProfile(ctx context.Context, site string, d *RADIUSProfile) (*RADIUSProfile, error) {
	return c.createRADIUSProfile(ctx, site, d)
}

func (c *Client) UpdateRADIUSProfile(ctx context.Context, site string, d *RADIUSProfile) (*RADIUSProfile, error) {
	return c.updateRADIUSProfile(ctx, site, d)
}
