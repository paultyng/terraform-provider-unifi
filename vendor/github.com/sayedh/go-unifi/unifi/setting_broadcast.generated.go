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

type SettingBroadcast struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	SoundAfterEnabled   bool   `json:"sound_after_enabled"`
	SoundAfterResource  string `json:"sound_after_resource,omitempty"`
	SoundAfterType      string `json:"sound_after_type,omitempty"` // sample|media
	SoundBeforeEnabled  bool   `json:"sound_before_enabled"`
	SoundBeforeResource string `json:"sound_before_resource,omitempty"`
	SoundBeforeType     string `json:"sound_before_type,omitempty"` // sample|media
}

func (dst *SettingBroadcast) UnmarshalJSON(b []byte) error {
	type Alias SettingBroadcast
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

func (c *Client) getSettingBroadcast(ctx context.Context, site string) (*SettingBroadcast, error) {
	var respBody struct {
		Meta meta               `json:"meta"`
		Data []SettingBroadcast `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/broadcast", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingBroadcast(ctx context.Context, site string, d *SettingBroadcast) (*SettingBroadcast, error) {
	var respBody struct {
		Meta meta               `json:"meta"`
		Data []SettingBroadcast `json:"data"`
	}

	d.Key = "broadcast"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/broadcast", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
