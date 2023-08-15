package unifi

import (
	"context"
)

func (c *Client) CreateWLAN(ctx context.Context, site string, d *WLAN) (*WLAN, error) {
	if d.Schedule == nil {
		d.Schedule = []string{}
	}

	return c.createWLAN(ctx, site, d)
}

func (c *Client) ListWLAN(ctx context.Context, site string) ([]WLAN, error) {
	return c.listWLAN(ctx, site)
}

func (c *Client) GetWLAN(ctx context.Context, site, id string) (*WLAN, error) {
	return c.getWLAN(ctx, site, id)
}

func (c *Client) DeleteWLAN(ctx context.Context, site, id string) error {
	return c.deleteWLAN(ctx, site, id)
}

func (c *Client) UpdateWLAN(ctx context.Context, site string, d *WLAN) (*WLAN, error) {
	return c.updateWLAN(ctx, site, d)
}
