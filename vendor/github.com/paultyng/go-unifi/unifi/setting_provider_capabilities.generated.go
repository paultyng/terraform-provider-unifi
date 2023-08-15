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

type SettingProviderCapabilities struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	Download int  `json:"download,omitempty"` // ^[1-9][0-9]*$
	Enabled  bool `json:"enabled"`
	Upload   int  `json:"upload,omitempty"` // ^[1-9][0-9]*$
}

func (dst *SettingProviderCapabilities) UnmarshalJSON(b []byte) error {
	type Alias SettingProviderCapabilities
	aux := &struct {
		Download emptyStringInt `json:"download"`
		Upload   emptyStringInt `json:"upload"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Download = int(aux.Download)
	dst.Upload = int(aux.Upload)

	return nil
}

func (c *Client) getSettingProviderCapabilities(ctx context.Context, site string) (*SettingProviderCapabilities, error) {
	var respBody struct {
		Meta meta                          `json:"meta"`
		Data []SettingProviderCapabilities `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/provider_capabilities", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingProviderCapabilities(ctx context.Context, site string, d *SettingProviderCapabilities) (*SettingProviderCapabilities, error) {
	var respBody struct {
		Meta meta                          `json:"meta"`
		Data []SettingProviderCapabilities `json:"data"`
	}

	d.Key = "provider_capabilities"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/provider_capabilities", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
