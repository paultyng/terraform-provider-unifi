// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type WLAN struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	BroadcastFilterEnabled    bool     `json:"bc_filter_enabled"`
	BroadcastFilterList       []string `json:"bc_filter_list,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	CountryBeacon             bool     `json:"country_beacon"`
	DPIEnabled                bool     `json:"dpi_enabled"`
	DPIgroupID                string   `json:"dpigroup_id"`         // [\d\w]+|^$
	DTIMMode                  string   `json:"dtim_mode,omitempty"` // default|custom
	DTIMNa                    int      `json:"dtim_na,omitempty"`   // ^([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DTIMNg                    int      `json:"dtim_ng,omitempty"`   // ^([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	ElementAdopt              bool     `json:"element_adopt"`
	Enabled                   bool     `json:"enabled"`
	FastRoamingEnabled        bool     `json:"fast_roaming_enabled"`
	GroupRekey                int      `json:"group_rekey,omitempty"` // ^(0|[6-9][0-9]|[1-9][0-9]{2,3}|[1-7][0-9]{4}|8[0-5][0-9]{3}|86[0-3][0-9][0-9]|86400)$
	HideSSID                  bool     `json:"hide_ssid"`
	Hotspot2ConfEnabled       bool     `json:"hotspot2conf_enabled"`
	Hotspot2ConfID            string   `json:"hotspot2conf_id"`
	IappEnabled               bool     `json:"iapp_enabled"`
	IsGuest                   bool     `json:"is_guest"`
	MACFilterEnabled          bool     `json:"mac_filter_enabled"`
	MACFilterList             []string `json:"mac_filter_list,omitempty"`   // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	MACFilterPolicy           string   `json:"mac_filter_policy,omitempty"` // allow|deny
	MulticastEnhanceEnabled   bool     `json:"mcastenhance_enabled"`
	MinrateNaAdvertisingRates bool     `json:"minrate_na_advertising_rates"`
	MinrateNaBeaconRateKbps   int      `json:"minrate_na_beacon_rate_kbps,omitempty"`
	MinrateNaDataRateKbps     int      `json:"minrate_na_data_rate_kbps,omitempty"`
	MinrateNaEnabled          bool     `json:"minrate_na_enabled"`
	MinrateNaMgmtRateKbps     int      `json:"minrate_na_mgmt_rate_kbps,omitempty"`
	MinrateNgAdvertisingRates bool     `json:"minrate_ng_advertising_rates"`
	MinrateNgBeaconRateKbps   int      `json:"minrate_ng_beacon_rate_kbps,omitempty"`
	MinrateNgCckRatesEnabled  bool     `json:"minrate_ng_cck_rates_enabled"`
	MinrateNgDataRateKbps     int      `json:"minrate_ng_data_rate_kbps,omitempty"`
	MinrateNgEnabled          bool     `json:"minrate_ng_enabled"`
	MinrateNgMgmtRateKbps     int      `json:"minrate_ng_mgmt_rate_kbps,omitempty"`
	Name                      string   `json:"name,omitempty"` // .{1,32}
	NameCombineEnabled        bool     `json:"name_combine_enabled"`
	NameCombineSuffix         string   `json:"name_combine_suffix,omitempty"` // .{0,8}
	No2GhzOui                 bool     `json:"no2ghz_oui"`
	Priority                  string   `json:"priority,omitempty"` // medium|high|low
	RADIUSDasEnabled          bool     `json:"radius_das_enabled"`
	RADIUSMACAuthEnabled      bool     `json:"radius_mac_auth_enabled"`
	RADIUSMACaclEmptyPassword bool     `json:"radius_macacl_empty_password"`
	RADIUSMACaclFormat        string   `json:"radius_macacl_format,omitempty"` // none_lower|hyphen_lower|colon_lower|none_upper|hyphen_upper|colon_upper
	RADIUSprofileID           string   `json:"radiusprofile_id"`
	RoamClusterID             int      `json:"roam_cluster_id,omitempty"` // [0-9]|[1-2][0-9]|[3][0-1]|^$
	RrmEnabled                bool     `json:"rrm_enabled"`
	Schedule                  []string `json:"schedule,omitempty"` // (sun|mon|tue|wed|thu|fri|sat)(\-(sun|mon|tue|wed|thu|fri|sat))?\|([0-2][0-9][0-5][0-9])\-([0-2][0-9][0-5][0-9])
	ScheduleEnabled           bool     `json:"schedule_enabled"`
	ScheduleReversed          bool     `json:"schedule_reversed"`
	Security                  string   `json:"security,omitempty"` // open|wpapsk|wep|wpaeap
	UapsdEnabled              bool     `json:"uapsd_enabled"`
	UserGroupID               string   `json:"usergroup_id"`
	VLAN                      int      `json:"vlan,omitempty"` // [2-9]|[1-9][0-9]{1,2}|[1-3][0-9]{3}|40[0-8][0-9]|409[0-5]|^$
	VLANEnabled               bool     `json:"vlan_enabled"`
	WEPIDX                    int      `json:"wep_idx,omitempty"` // [1-4]
	WLANGroupID               string   `json:"wlangroup_id"`
	WPAEnc                    string   `json:"wpa_enc,omitempty"`      // auto|ccmp
	WPAMode                   string   `json:"wpa_mode,omitempty"`     // auto|wpa1|wpa2
	XIappKey                  string   `json:"x_iapp_key,omitempty"`   // [0-9A-Fa-f]{32}
	XPassphrase               string   `json:"x_passphrase,omitempty"` // [\x20-\x7E]{8,63}|[0-9a-fA-F]{64}
	XWEP                      string   `json:"x_wep,omitempty"`
}
