package unifi

import "context"

func (c *Client) ListDynamicDNS(ctx context.Context, site string) ([]DynamicDNS, error) {
	return c.listDynamicDNS(ctx, site)
}

func (c *Client) GetDynamicDNS(ctx context.Context, site, id string) (*DynamicDNS, error) {
	return c.getDynamicDNS(ctx, site, id)
}

func (c *Client) DeleteDynamicDNS(ctx context.Context, site, id string) error {
	return c.deleteDynamicDNS(ctx, site, id)
}

func (c *Client) CreateDynamicDNS(ctx context.Context, site string, d *DynamicDNS) (*DynamicDNS, error) {
	return c.createDynamicDNS(ctx, site, d)
}

func (c *Client) UpdateDynamicDNS(ctx context.Context, site string, d *DynamicDNS) (*DynamicDNS, error) {
	return c.updateDynamicDNS(ctx, site, d)
}
