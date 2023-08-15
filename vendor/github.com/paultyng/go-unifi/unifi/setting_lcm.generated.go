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

type SettingLcm struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	Brightness  int  `json:"brightness,omitempty"` // [1-9]|[1-9][0-9]|100
	Enabled     bool `json:"enabled"`
	IDleTimeout int  `json:"idle_timeout,omitempty"` // [1-9][0-9]|[1-9][0-9][0-9]|[1-2][0-9][0-9][0-9]|3[0-5][0-9][0-9]|3600
	Sync        bool `json:"sync"`
	TouchEvent  bool `json:"touch_event"`
}

func (dst *SettingLcm) UnmarshalJSON(b []byte) error {
	type Alias SettingLcm
	aux := &struct {
		Brightness  emptyStringInt `json:"brightness"`
		IDleTimeout emptyStringInt `json:"idle_timeout"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Brightness = int(aux.Brightness)
	dst.IDleTimeout = int(aux.IDleTimeout)

	return nil
}

func (c *Client) getSettingLcm(ctx context.Context, site string) (*SettingLcm, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []SettingLcm `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/lcm", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingLcm(ctx context.Context, site string, d *SettingLcm) (*SettingLcm, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []SettingLcm `json:"data"`
	}

	d.Key = "lcm"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/lcm", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
