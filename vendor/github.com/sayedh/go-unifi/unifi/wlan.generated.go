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

type WLAN struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	ApGroupIDs                  []string                   `json:"ap_group_ids,omitempty"`
	ApGroupMode                 string                     `json:"ap_group_mode,omitempty"` // all|groups|devices
	AuthCache                   bool                       `json:"auth_cache"`
	BSupported                  bool                       `json:"b_supported"`
	BroadcastFilterEnabled      bool                       `json:"bc_filter_enabled"`
	BroadcastFilterList         []string                   `json:"bc_filter_list,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	BssTransition               bool                       `json:"bss_transition"`
	CountryBeacon               bool                       `json:"country_beacon"`
	DPIEnabled                  bool                       `json:"dpi_enabled"`
	DPIgroupID                  string                     `json:"dpigroup_id"`         // [\d\w]+|^$
	DTIM6E                      int                        `json:"dtim_6e,omitempty"`   // ^([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DTIMMode                    string                     `json:"dtim_mode,omitempty"` // default|custom
	DTIMNa                      int                        `json:"dtim_na,omitempty"`   // ^([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DTIMNg                      int                        `json:"dtim_ng,omitempty"`   // ^([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	ElementAdopt                bool                       `json:"element_adopt"`
	Enabled                     bool                       `json:"enabled"`
	FastRoamingEnabled          bool                       `json:"fast_roaming_enabled"`
	GroupRekey                  int                        `json:"group_rekey,omitempty"` // ^(0|[6-9][0-9]|[1-9][0-9]{2,3}|[1-7][0-9]{4}|8[0-5][0-9]{3}|86[0-3][0-9][0-9]|86400)$
	HideSSID                    bool                       `json:"hide_ssid"`
	Hotspot2ConfEnabled         bool                       `json:"hotspot2conf_enabled"`
	Hotspot2ConfID              string                     `json:"hotspot2conf_id"`
	IappEnabled                 bool                       `json:"iapp_enabled"`
	IsGuest                     bool                       `json:"is_guest"`
	L2Isolation                 bool                       `json:"l2_isolation"`
	LogLevel                    string                     `json:"log_level,omitempty"`
	MACFilterEnabled            bool                       `json:"mac_filter_enabled"`
	MACFilterList               []string                   `json:"mac_filter_list,omitempty"`   // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	MACFilterPolicy             string                     `json:"mac_filter_policy,omitempty"` // allow|deny
	MinrateNaAdvertisingRates   bool                       `json:"minrate_na_advertising_rates"`
	MinrateNaDataRateKbps       int                        `json:"minrate_na_data_rate_kbps,omitempty"`
	MinrateNaEnabled            bool                       `json:"minrate_na_enabled"`
	MinrateNgAdvertisingRates   bool                       `json:"minrate_ng_advertising_rates"`
	MinrateNgDataRateKbps       int                        `json:"minrate_ng_data_rate_kbps,omitempty"`
	MinrateNgEnabled            bool                       `json:"minrate_ng_enabled"`
	MinrateSettingPreference    string                     `json:"minrate_setting_preference,omitempty"` // auto|manual
	MloEnabled                  bool                       `json:"mlo_enabled"`
	MulticastEnhanceEnabled     bool                       `json:"mcastenhance_enabled"`
	Name                        string                     `json:"name,omitempty"` // .{1,32}
	NameCombineEnabled          bool                       `json:"name_combine_enabled"`
	NameCombineSuffix           string                     `json:"name_combine_suffix,omitempty"` // .{0,8}
	NetworkID                   string                     `json:"networkconf_id"`
	No2GhzOui                   bool                       `json:"no2ghz_oui"`
	OptimizeIotWifiConnectivity bool                       `json:"optimize_iot_wifi_connectivity"`
	P2P                         bool                       `json:"p2p"`
	P2PCrossConnect             bool                       `json:"p2p_cross_connect"`
	PMFCipher                   string                     `json:"pmf_cipher,omitempty"` // auto|aes-128-cmac|bip-gmac-256
	PMFMode                     string                     `json:"pmf_mode,omitempty"`   // disabled|optional|required
	Priority                    string                     `json:"priority,omitempty"`   // medium|high|low
	PrivatePresharedKeys        []WLANPrivatePresharedKeys `json:"private_preshared_keys,omitempty"`
	PrivatePresharedKeysEnabled bool                       `json:"private_preshared_keys_enabled"`
	ProxyArp                    bool                       `json:"proxy_arp"`
	RADIUSDasEnabled            bool                       `json:"radius_das_enabled"`
	RADIUSMACAuthEnabled        bool                       `json:"radius_mac_auth_enabled"`
	RADIUSMACaclEmptyPassword   bool                       `json:"radius_macacl_empty_password"`
	RADIUSMACaclFormat          string                     `json:"radius_macacl_format,omitempty"` // none_lower|hyphen_lower|colon_lower|none_upper|hyphen_upper|colon_upper
	RADIUSProfileID             string                     `json:"radiusprofile_id"`
	RoamClusterID               int                        `json:"roam_cluster_id,omitempty"` // [0-9]|[1-2][0-9]|[3][0-1]|^$
	RrmEnabled                  bool                       `json:"rrm_enabled"`
	SaeAntiClogging             int                        `json:"sae_anti_clogging,omitempty"`
	SaeGroups                   []int                      `json:"sae_groups,omitempty"`
	SaePsk                      []WLANSaePsk               `json:"sae_psk,omitempty"`
	SaePskVLANRequired          bool                       `json:"sae_psk_vlan_required"`
	SaeSync                     int                        `json:"sae_sync,omitempty"`
	Schedule                    []string                   `json:"schedule,omitempty"` // (sun|mon|tue|wed|thu|fri|sat)(\-(sun|mon|tue|wed|thu|fri|sat))?\|([0-2][0-9][0-5][0-9])\-([0-2][0-9][0-5][0-9])
	ScheduleEnabled             bool                       `json:"schedule_enabled"`
	ScheduleReversed            bool                       `json:"schedule_reversed"`
	ScheduleWithDuration        []WLANScheduleWithDuration `json:"schedule_with_duration"`
	Security                    string                     `json:"security,omitempty"`           // open|wpapsk|wep|wpaeap|osen
	SettingPreference           string                     `json:"setting_preference,omitempty"` // auto|manual
	TdlsProhibit                bool                       `json:"tdls_prohibit"`
	UapsdEnabled                bool                       `json:"uapsd_enabled"`
	UidWorkspaceUrl             string                     `json:"uid_workspace_url,omitempty"`
	UserGroupID                 string                     `json:"usergroup_id"`
	VLAN                        int                        `json:"vlan,omitempty"` // [2-9]|[1-9][0-9]{1,2}|[1-3][0-9]{3}|40[0-8][0-9]|409[0-5]|^$
	VLANEnabled                 bool                       `json:"vlan_enabled"`
	WEPIDX                      int                        `json:"wep_idx,omitempty"`    // [1-4]
	WLANBand                    string                     `json:"wlan_band,omitempty"`  // 2g|5g|both
	WLANBands                   []string                   `json:"wlan_bands,omitempty"` // 2g|5g|6g
	WLANGroupID                 string                     `json:"wlangroup_id"`
	WPA3Enhanced192             bool                       `json:"wpa3_enhanced_192"`
	WPA3FastRoaming             bool                       `json:"wpa3_fast_roaming"`
	WPA3Support                 bool                       `json:"wpa3_support"`
	WPA3Transition              bool                       `json:"wpa3_transition"`
	WPAEnc                      string                     `json:"wpa_enc,omitempty"`        // auto|ccmp|gcmp|ccmp-256|gcmp-256
	WPAMode                     string                     `json:"wpa_mode,omitempty"`       // auto|wpa1|wpa2
	WPAPskRADIUS                string                     `json:"wpa_psk_radius,omitempty"` // disabled|optional|required
	XIappKey                    string                     `json:"x_iapp_key,omitempty"`     // [0-9A-Fa-f]{32}
	XPassphrase                 string                     `json:"x_passphrase,omitempty"`   // [\x20-\x7E]{8,255}|[0-9a-fA-F]{64}
	XWEP                        string                     `json:"x_wep,omitempty"`
}

func (dst *WLAN) UnmarshalJSON(b []byte) error {
	type Alias WLAN
	aux := &struct {
		DTIM6E                emptyStringInt   `json:"dtim_6e"`
		DTIMNa                emptyStringInt   `json:"dtim_na"`
		DTIMNg                emptyStringInt   `json:"dtim_ng"`
		GroupRekey            emptyStringInt   `json:"group_rekey"`
		MinrateNaDataRateKbps emptyStringInt   `json:"minrate_na_data_rate_kbps"`
		MinrateNgDataRateKbps emptyStringInt   `json:"minrate_ng_data_rate_kbps"`
		RoamClusterID         emptyStringInt   `json:"roam_cluster_id"`
		SaeAntiClogging       emptyStringInt   `json:"sae_anti_clogging"`
		SaeGroups             []emptyStringInt `json:"sae_groups"`
		SaeSync               emptyStringInt   `json:"sae_sync"`
		VLAN                  emptyStringInt   `json:"vlan"`
		WEPIDX                emptyStringInt   `json:"wep_idx"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.DTIM6E = int(aux.DTIM6E)
	dst.DTIMNa = int(aux.DTIMNa)
	dst.DTIMNg = int(aux.DTIMNg)
	dst.GroupRekey = int(aux.GroupRekey)
	dst.MinrateNaDataRateKbps = int(aux.MinrateNaDataRateKbps)
	dst.MinrateNgDataRateKbps = int(aux.MinrateNgDataRateKbps)
	dst.RoamClusterID = int(aux.RoamClusterID)
	dst.SaeAntiClogging = int(aux.SaeAntiClogging)
	dst.SaeGroups = make([]int, len(aux.SaeGroups))
	for i, v := range aux.SaeGroups {
		dst.SaeGroups[i] = int(v)
	}
	dst.SaeSync = int(aux.SaeSync)
	dst.VLAN = int(aux.VLAN)
	dst.WEPIDX = int(aux.WEPIDX)

	return nil
}

type WLANPrivatePresharedKeys struct {
	NetworkID string `json:"networkconf_id"`
	Password  string `json:"password,omitempty"` // [\x20-\x7E]{8,255}
}

func (dst *WLANPrivatePresharedKeys) UnmarshalJSON(b []byte) error {
	type Alias WLANPrivatePresharedKeys
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

type WLANSaePsk struct {
	ID   string `json:"id"`             // .{0,128}
	MAC  string `json:"mac,omitempty"`  // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	Psk  string `json:"psk,omitempty"`  // [\x20-\x7E]{8,255}
	VLAN int    `json:"vlan,omitempty"` // [0-9]|[1-9][0-9]{1,2}|[1-3][0-9]{3}|40[0-8][0-9]|409[0-5]|^$
}

func (dst *WLANSaePsk) UnmarshalJSON(b []byte) error {
	type Alias WLANSaePsk
	aux := &struct {
		VLAN emptyStringInt `json:"vlan"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.VLAN = int(aux.VLAN)

	return nil
}

type WLANScheduleWithDuration struct {
	DurationMinutes int      `json:"duration_minutes,omitempty"`   // ^[1-9][0-9]*$
	Name            string   `json:"name,omitempty"`               // .*
	StartDaysOfWeek []string `json:"start_days_of_week,omitempty"` // ^(sun|mon|tue|wed|thu|fri|sat)$
	StartHour       int      `json:"start_hour,omitempty"`         // ^(1?[0-9])|(2[0-3])$
	StartMinute     int      `json:"start_minute,omitempty"`       // ^[0-5]?[0-9]$
}

func (dst *WLANScheduleWithDuration) UnmarshalJSON(b []byte) error {
	type Alias WLANScheduleWithDuration
	aux := &struct {
		DurationMinutes emptyStringInt `json:"duration_minutes"`
		StartHour       emptyStringInt `json:"start_hour"`
		StartMinute     emptyStringInt `json:"start_minute"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.DurationMinutes = int(aux.DurationMinutes)
	dst.StartHour = int(aux.StartHour)
	dst.StartMinute = int(aux.StartMinute)

	return nil
}

func (c *Client) listWLAN(ctx context.Context, site string) ([]WLAN, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/wlanconf", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getWLAN(ctx context.Context, site, id string) (*WLAN, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/wlanconf/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteWLAN(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/wlanconf/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createWLAN(ctx context.Context, site string, d *WLAN) (*WLAN, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/wlanconf", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateWLAN(ctx context.Context, site string, d *WLAN) (*WLAN, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/wlanconf/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
