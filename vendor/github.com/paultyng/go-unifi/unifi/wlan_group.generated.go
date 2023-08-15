// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type WLANGroup struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Name string `json:"name,omitempty"` // .{1,128}
}

func (dst *WLANGroup) UnmarshalJSON(b []byte) error {
	type Alias WLANGroup
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

func (c *Client) listWLANGroup(ctx context.Context, site string) ([]WLANGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []WLANGroup `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/wlangroup", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getWLANGroup(ctx context.Context, site, id string) (*WLANGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []WLANGroup `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteWLANGroup(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createWLANGroup(ctx context.Context, site string, d *WLANGroup) (*WLANGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []WLANGroup `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/wlangroup", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateWLANGroup(ctx context.Context, site string, d *WLANGroup) (*WLANGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []WLANGroup `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
