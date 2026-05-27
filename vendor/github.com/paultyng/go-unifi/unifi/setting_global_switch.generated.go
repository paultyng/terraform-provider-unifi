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

type SettingGlobalSwitch struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	AclDeviceIsolation     []string                            `json:"acl_device_isolation,omitempty"`
	AclL3Isolation         []SettingGlobalSwitchAclL3Isolation `json:"acl_l3_isolation,omitempty"`
	DHCPSnoop              bool                                `json:"dhcp_snoop"`
	Dot1XFallbackNetworkID string                              `json:"dot1x_fallback_networkconf_id"` // [\d\w]+|
	Dot1XPortctrlEnabled   bool                                `json:"dot1x_portctrl_enabled"`
	FlowctrlEnabled        bool                                `json:"flowctrl_enabled"`
	JumboframeEnabled      bool                                `json:"jumboframe_enabled"`
	RADIUSProfileID        string                              `json:"radiusprofile_id"`
	StpVersion             string                              `json:"stp_version,omitempty"`       // stp|rstp|disabled
	SwitchExclusions       []string                            `json:"switch_exclusions,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
}

func (dst *SettingGlobalSwitch) UnmarshalJSON(b []byte) error {
	type Alias SettingGlobalSwitch
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

type SettingGlobalSwitchAclL3Isolation struct {
	DestinationNetworks []string `json:"destination_networks,omitempty"`
	SourceNetwork       string   `json:"source_network,omitempty"`
}

func (dst *SettingGlobalSwitchAclL3Isolation) UnmarshalJSON(b []byte) error {
	type Alias SettingGlobalSwitchAclL3Isolation
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

func (c *Client) getSettingGlobalSwitch(ctx context.Context, site string) (*SettingGlobalSwitch, error) {
	var respBody struct {
		Meta meta                  `json:"meta"`
		Data []SettingGlobalSwitch `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/global_switch", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingGlobalSwitch(ctx context.Context, site string, d *SettingGlobalSwitch) (*SettingGlobalSwitch, error) {
	var respBody struct {
		Meta meta                  `json:"meta"`
		Data []SettingGlobalSwitch `json:"data"`
	}

	d.Key = "global_switch"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/global_switch", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
