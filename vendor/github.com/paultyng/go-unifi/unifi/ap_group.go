package unifi

import (
	"context"
	"fmt"
)

// just to fix compile issues with the import
var (
	_ fmt.Formatter
	_ context.Context
)

// This is a v2 API object, so manually coded for now, need to figure out generation...

type APGroup struct {
	ID string `json:"_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Name       string   `json:"name"`
	DeviceMACs []string `json:"device_macs"`
}

func (c *Client) ListAPGroup(ctx context.Context, site string) ([]APGroup, error) {
	var respBody []APGroup

	err := c.do(ctx, "GET", fmt.Sprintf("%s/site/%s/apgroups", c.apiV2Path, site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// func (c *Client) getWLANGroup(ctx context.Context, site, id string) (*WLANGroup, error) {
// 	var respBody struct {
// 		Meta meta        `json:"meta"`
// 		Data []WLANGroup `json:"data"`
// 	}

// 	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, id), nil, &respBody)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(respBody.Data) != 1 {
// 		return nil, &NotFoundError{}
// 	}

// 	d := respBody.Data[0]
// 	return &d, nil
// }

// func (c *Client) deleteWLANGroup(ctx context.Context, site, id string) error {
// 	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, id), struct{}{}, nil)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (c *Client) CreateAPGroup(ctx context.Context, site string, d *APGroup) (*APGroup, error) {
	var respBody APGroup

	err := c.do(ctx, "POST", fmt.Sprintf("%s/site/%s/apgroups", c.apiV2Path, site), d, &respBody)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}

// func (c *Client) updateWLANGroup(ctx context.Context, site string, d *WLANGroup) (*WLANGroup, error) {
// 	var respBody struct {
// 		Meta meta        `json:"meta"`
// 		Data []WLANGroup `json:"data"`
// 	}

// 	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, d.ID), d, &respBody)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(respBody.Data) != 1 {
// 		return nil, &NotFoundError{}
// 	}

// 	new := respBody.Data[0]

// 	return &new, nil
// }
