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

type SettingRadioAi struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	AutoAdjustChannelsToCountry bool     `json:"auto_adjust_channels_to_country"`
	ChannelsNa                  []int    `json:"channels_na,omitempty"` // 34|36|38|40|42|44|46|48|52|56|60|64|100|104|108|112|116|120|124|128|132|136|140|144|149|153|157|161|165|169
	ChannelsNg                  []int    `json:"channels_ng,omitempty"` // 1|2|3|4|5|6|7|8|9|10|11|12|13|14
	CronExpr                    string   `json:"cron_expr,omitempty"`
	Default                     bool     `json:"default"`
	Enabled                     bool     `json:"enabled"`
	ExcludeDevices              []string `json:"exclude_devices,omitempty"`    // ([0-9a-z]{2}:){5}[0-9a-z]{2}
	HtModesNa                   []int    `json:"ht_modes_na,omitempty"`        // ^(20|40|80|160)$
	HtModesNg                   []int    `json:"ht_modes_ng,omitempty"`        // ^(20|40)$
	Optimize                    []string `json:"optimize,omitempty"`           // channel|power
	Radios                      []string `json:"radios,omitempty"`             // na|ng
	SettingPreference           string   `json:"setting_preference,omitempty"` // auto|manual
	UseXY                       bool     `json:"useXY"`
}

func (dst *SettingRadioAi) UnmarshalJSON(b []byte) error {
	type Alias SettingRadioAi
	aux := &struct {
		ChannelsNa []emptyStringInt `json:"channels_na"`
		ChannelsNg []emptyStringInt `json:"channels_ng"`
		HtModesNa  []emptyStringInt `json:"ht_modes_na"`
		HtModesNg  []emptyStringInt `json:"ht_modes_ng"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.ChannelsNa = make([]int, len(aux.ChannelsNa))
	for i, v := range aux.ChannelsNa {
		dst.ChannelsNa[i] = int(v)
	}
	dst.ChannelsNg = make([]int, len(aux.ChannelsNg))
	for i, v := range aux.ChannelsNg {
		dst.ChannelsNg[i] = int(v)
	}
	dst.HtModesNa = make([]int, len(aux.HtModesNa))
	for i, v := range aux.HtModesNa {
		dst.HtModesNa[i] = int(v)
	}
	dst.HtModesNg = make([]int, len(aux.HtModesNg))
	for i, v := range aux.HtModesNg {
		dst.HtModesNg[i] = int(v)
	}

	return nil
}

func (c *Client) getSettingRadioAi(ctx context.Context, site string) (*SettingRadioAi, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []SettingRadioAi `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/radio_ai", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingRadioAi(ctx context.Context, site string, d *SettingRadioAi) (*SettingRadioAi, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []SettingRadioAi `json:"data"`
	}

	d.Key = "radio_ai"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/radio_ai", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
