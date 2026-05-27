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

type SettingMagicSiteToSiteVpn struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	Enabled bool `json:"enabled"`
}

func (dst *SettingMagicSiteToSiteVpn) UnmarshalJSON(b []byte) error {
	type Alias SettingMagicSiteToSiteVpn
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

func (c *Client) getSettingMagicSiteToSiteVpn(ctx context.Context, site string) (*SettingMagicSiteToSiteVpn, error) {
	var respBody struct {
		Meta meta                        `json:"meta"`
		Data []SettingMagicSiteToSiteVpn `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/magic_site_to_site_vpn", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingMagicSiteToSiteVpn(ctx context.Context, site string, d *SettingMagicSiteToSiteVpn) (*SettingMagicSiteToSiteVpn, error) {
	var respBody struct {
		Meta meta                        `json:"meta"`
		Data []SettingMagicSiteToSiteVpn `json:"data"`
	}

	d.Key = "magic_site_to_site_vpn"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/magic_site_to_site_vpn", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
