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

type SettingNtp struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	NtpServer1        string `json:"ntp_server_1,omitempty"`
	NtpServer2        string `json:"ntp_server_2,omitempty"`
	NtpServer3        string `json:"ntp_server_3,omitempty"`
	NtpServer4        string `json:"ntp_server_4,omitempty"`
	SettingPreference string `json:"setting_preference,omitempty"` // auto|manual
}

func (dst *SettingNtp) UnmarshalJSON(b []byte) error {
	type Alias SettingNtp
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

func (c *Client) getSettingNtp(ctx context.Context, site string) (*SettingNtp, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []SettingNtp `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/ntp", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingNtp(ctx context.Context, site string, d *SettingNtp) (*SettingNtp, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []SettingNtp `json:"data"`
	}

	d.Key = "ntp"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/ntp", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
