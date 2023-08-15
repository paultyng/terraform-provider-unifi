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

type Account struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	IP               string `json:"ip,omitempty"`   // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	Name             string `json:"name,omitempty"` // ^[^"' ]+$
	NetworkID        string `json:"networkconf_id,omitempty"`
	TunnelConfigType string `json:"tunnel_config_type,omitempty"` // vpn|802.1x|custom
	TunnelMediumType int    `json:"tunnel_medium_type,omitempty"` // [1-9]|1[0-5]|^$
	TunnelType       int    `json:"tunnel_type,omitempty"`        // [1-9]|1[0-3]|^$
	VLAN             int    `json:"vlan,omitempty"`               // [2-9]|[1-9][0-9]{1,2}|[1-3][0-9]{3}|400[0-9]|^$
	XPassword        string `json:"x_password,omitempty"`
}

func (dst *Account) UnmarshalJSON(b []byte) error {
	type Alias Account
	aux := &struct {
		TunnelMediumType emptyStringInt `json:"tunnel_medium_type"`
		TunnelType       emptyStringInt `json:"tunnel_type"`
		VLAN             emptyStringInt `json:"vlan"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.TunnelMediumType = int(aux.TunnelMediumType)
	dst.TunnelType = int(aux.TunnelType)
	dst.VLAN = int(aux.VLAN)

	return nil
}

func (c *Client) listAccount(ctx context.Context, site string) ([]Account, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Account `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/account", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getAccount(ctx context.Context, site, id string) (*Account, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Account `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/account/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteAccount(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/account/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createAccount(ctx context.Context, site string, d *Account) (*Account, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Account `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/account", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateAccount(ctx context.Context, site string, d *Account) (*Account, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Account `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/account/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
