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

type SettingPorta struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	Ugw3WAN2Enabled bool `json:"ugw3_wan2_enabled"`
}

func (dst *SettingPorta) UnmarshalJSON(b []byte) error {
	type Alias SettingPorta
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

func (c *Client) getSettingPorta(ctx context.Context, site string) (*SettingPorta, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []SettingPorta `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/porta", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingPorta(ctx context.Context, site string, d *SettingPorta) (*SettingPorta, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []SettingPorta `json:"data"`
	}

	d.Key = "porta"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/porta", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
