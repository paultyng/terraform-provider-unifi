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

type SettingRadius struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	AccountingEnabled     bool   `json:"accounting_enabled"`
	AcctPort              int    `json:"acct_port,omitempty"` // [1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5]
	AuthPort              int    `json:"auth_port,omitempty"` // [1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5]
	ConfigureWholeNetwork bool   `json:"configure_whole_network"`
	Enabled               bool   `json:"enabled"`
	InterimUpdateInterval int    `json:"interim_update_interval,omitempty"` // ^([6-9][0-9]|[1-9][0-9]{2,3}|[1-7][0-9]{4}|8[0-5][0-9]{3}|86[0-3][0-9][0-9]|86400)$
	TunneledReply         bool   `json:"tunneled_reply"`
	XSecret               string `json:"x_secret,omitempty"` // ^[^\\"' ]{1,48}$
}

func (dst *SettingRadius) UnmarshalJSON(b []byte) error {
	type Alias SettingRadius
	aux := &struct {
		AcctPort              emptyStringInt `json:"acct_port"`
		AuthPort              emptyStringInt `json:"auth_port"`
		InterimUpdateInterval emptyStringInt `json:"interim_update_interval"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.AcctPort = int(aux.AcctPort)
	dst.AuthPort = int(aux.AuthPort)
	dst.InterimUpdateInterval = int(aux.InterimUpdateInterval)

	return nil
}

func (c *Client) getSettingRadius(ctx context.Context, site string) (*SettingRadius, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []SettingRadius `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/radius", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingRadius(ctx context.Context, site string, d *SettingRadius) (*SettingRadius, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []SettingRadius `json:"data"`
	}

	d.Key = "radius"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/radius", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
