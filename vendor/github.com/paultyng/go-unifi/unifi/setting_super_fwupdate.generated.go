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

type SettingSuperFwupdate struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	ControllerChannel string `json:"controller_channel,omitempty"` // internal|alpha|beta|release-candidate|release
	FirmwareChannel   string `json:"firmware_channel,omitempty"`   // internal|alpha|beta|release-candidate|release
	SsoEnabled        bool   `json:"sso_enabled"`
}

func (dst *SettingSuperFwupdate) UnmarshalJSON(b []byte) error {
	type Alias SettingSuperFwupdate
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

func (c *Client) getSettingSuperFwupdate(ctx context.Context, site string) (*SettingSuperFwupdate, error) {
	var respBody struct {
		Meta meta                   `json:"meta"`
		Data []SettingSuperFwupdate `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/super_fwupdate", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingSuperFwupdate(ctx context.Context, site string, d *SettingSuperFwupdate) (*SettingSuperFwupdate, error) {
	var respBody struct {
		Meta meta                   `json:"meta"`
		Data []SettingSuperFwupdate `json:"data"`
	}

	d.Key = "super_fwupdate"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/super_fwupdate", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
