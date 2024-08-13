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

type SettingDashboard struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	LayoutPreference string                    `json:"layout_preference,omitempty"` // auto|custom
	Widgets          []SettingDashboardWidgets `json:"widgets,omitempty"`
}

func (dst *SettingDashboard) UnmarshalJSON(b []byte) error {
	type Alias SettingDashboard
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

type SettingDashboardWidgets struct {
	Name string `json:"name,omitempty"` // traffic_identification|connection_types|wifi_technology|most_active_clients|most_active_aps|meshing|network_activity|wireless_experience|internet|wifi_activity|wifi_channels|wifi_client_experience|wifi_tx_retries|admin_activity|device_client_count|server_ip
}

func (dst *SettingDashboardWidgets) UnmarshalJSON(b []byte) error {
	type Alias SettingDashboardWidgets
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

func (c *Client) getSettingDashboard(ctx context.Context, site string) (*SettingDashboard, error) {
	var respBody struct {
		Meta meta               `json:"meta"`
		Data []SettingDashboard `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/dashboard", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingDashboard(ctx context.Context, site string, d *SettingDashboard) (*SettingDashboard, error) {
	var respBody struct {
		Meta meta               `json:"meta"`
		Data []SettingDashboard `json:"data"`
	}

	d.Key = "dashboard"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/dashboard", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
