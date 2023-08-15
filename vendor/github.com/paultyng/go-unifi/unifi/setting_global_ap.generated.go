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

type SettingGlobalAp struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	ApExclusions    []string `json:"ap_exclusions,omitempty"`    // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	NaChannelSize   int      `json:"na_channel_size,omitempty"`  // 20|40|80|160
	NaTxPower       int      `json:"na_tx_power,omitempty"`      // [0-9]|[1-4][0-9]
	NaTxPowerMode   string   `json:"na_tx_power_mode,omitempty"` // auto|medium|high|low|custom
	NgChannelSize   int      `json:"ng_channel_size,omitempty"`  // 20|40
	NgTxPower       int      `json:"ng_tx_power,omitempty"`      // [0-9]|[1-4][0-9]
	NgTxPowerMode   string   `json:"ng_tx_power_mode,omitempty"` // auto|medium|high|low|custom
	SixEChannelSize int      `json:"6e_channel_size,omitempty"`  // 20|40|80|160
	SixETxPower     int      `json:"6e_tx_power,omitempty"`      // [0-9]|[1-4][0-9]
	SixETxPowerMode string   `json:"6e_tx_power_mode,omitempty"` // auto|medium|high|low|custom
}

func (dst *SettingGlobalAp) UnmarshalJSON(b []byte) error {
	type Alias SettingGlobalAp
	aux := &struct {
		NaChannelSize   emptyStringInt `json:"na_channel_size"`
		NaTxPower       emptyStringInt `json:"na_tx_power"`
		NgChannelSize   emptyStringInt `json:"ng_channel_size"`
		NgTxPower       emptyStringInt `json:"ng_tx_power"`
		SixEChannelSize emptyStringInt `json:"6e_channel_size"`
		SixETxPower     emptyStringInt `json:"6e_tx_power"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.NaChannelSize = int(aux.NaChannelSize)
	dst.NaTxPower = int(aux.NaTxPower)
	dst.NgChannelSize = int(aux.NgChannelSize)
	dst.NgTxPower = int(aux.NgTxPower)
	dst.SixEChannelSize = int(aux.SixEChannelSize)
	dst.SixETxPower = int(aux.SixETxPower)

	return nil
}

func (c *Client) getSettingGlobalAp(ctx context.Context, site string) (*SettingGlobalAp, error) {
	var respBody struct {
		Meta meta              `json:"meta"`
		Data []SettingGlobalAp `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/global_ap", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingGlobalAp(ctx context.Context, site string, d *SettingGlobalAp) (*SettingGlobalAp, error) {
	var respBody struct {
		Meta meta              `json:"meta"`
		Data []SettingGlobalAp `json:"data"`
	}

	d.Key = "global_ap"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/global_ap", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
