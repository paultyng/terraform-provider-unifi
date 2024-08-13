package unifi

import "context"

func (c *Client) ListPortForward(ctx context.Context, site string) ([]PortForward, error) {
	return c.listPortForward(ctx, site)
}

func (c *Client) GetPortForward(ctx context.Context, site, id string) (*PortForward, error) {
	return c.getPortForward(ctx, site, id)
}

func (c *Client) DeletePortForward(ctx context.Context, site, id string) error {
	return c.deletePortForward(ctx, site, id)
}

func (c *Client) CreatePortForward(ctx context.Context, site string, d *PortForward) (*PortForward, error) {
	return c.createPortForward(ctx, site, d)
}

func (c *Client) UpdatePortForward(ctx context.Context, site string, d *PortForward) (*PortForward, error) {
	return c.updatePortForward(ctx, site, d)
}
