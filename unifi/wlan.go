package unifi

import (
	"fmt"
)

/*
{
  "meta": {
    "rc": "ok"
  },
  "data": [
    {
      "_id": "5deeabfc439adf048407dcf3",
      "enabled": true,
      "name": "fred",
      "security": "wpapsk",
      "wpa_enc": "ccmp",
      "wpa_mode": "wpa2",
      "x_passphrase": "Liebeskind",
      "wlangroup_id": "5d6d8b0c439adf048407dce9",
      "name_combine_enabled": true,
      "site_id": "5d6d8b07439adf048407dcd9",
      "x_iapp_key": "a199e6f225c01127e4135211236f9767",
      "minrate_ng_enabled": true,
      "minrate_ng_beacon_rate_kbps": 6000,
      "minrate_ng_data_rate_kbps": 6000,
      "no2ghz_oui": true,
      "wep_idx": 1,
      "usergroup_id": "5d6d8b0c439adf048407dce8",
      "dtim_mode": "default",
      "dtim_ng": 1,
      "dtim_na": 1,
      "minrate_ng_advertising_rates": false,
      "minrate_ng_cck_rates_enabled": true,
      "minrate_na_enabled": false,
      "minrate_na_advertising_rates": false,
      "minrate_na_data_rate_kbps": 6000,
      "mac_filter_enabled": false,
      "mac_filter_policy": "allow",
      "mac_filter_list": [],
      "bc_filter_enabled": false,
      "bc_filter_list": [],
      "group_rekey": 3600,
      "vlan_enabled": true,
      "vlan": "20",
      "radius_das_enabled": false,
      "schedule": [],
      "minrate_ng_mgmt_rate_kbps": 6000,
      "minrate_na_mgmt_rate_kbps": 6000,
      "minrate_na_beacon_rate_kbps": 6000
    },
    {
      "_id": "5deecfa31e801c052a1a5f5a",
      "enabled": true,
      "is_guest": true,
      "name": "fred-guest",
      "security": "wpapsk",
      "usergroup_id": "5d6d8b0c439adf048407dce8",
      "vlan_enabled": true,
      "wlangroup_id": "5d6d8b0c439adf048407dce9",
      "x_passphrase": "Enzo and Alice",
      "site_id": "5d6d8b07439adf048407dcd9",
      "x_iapp_key": "3e6b24dda0d11ce3c097107fc429596e",
      "minrate_ng_enabled": true,
      "minrate_ng_beacon_rate_kbps": 6000,
      "minrate_ng_data_rate_kbps": 6000,
      "vlan": "30",
      "wep_idx": 1,
      "wpa_mode": "wpa2",
      "wpa_enc": "ccmp",
      "dtim_mode": "default",
      "dtim_ng": 1,
      "dtim_na": 1,
      "minrate_ng_advertising_rates": false,
      "minrate_ng_cck_rates_enabled": true,
      "minrate_na_enabled": false,
      "minrate_na_advertising_rates": false,
      "minrate_na_data_rate_kbps": 6000,
      "mac_filter_enabled": false,
      "mac_filter_policy": "allow",
      "mac_filter_list": [],
      "name_combine_enabled": true,
      "bc_filter_enabled": false,
      "bc_filter_list": [],
      "group_rekey": 3600,
      "radius_das_enabled": false,
      "schedule": [],
      "minrate_ng_mgmt_rate_kbps": 6000,
      "minrate_na_mgmt_rate_kbps": 6000,
      "minrate_na_beacon_rate_kbps": 6000
    },
    {
      "_id": "5deed0d51e801c052a1a5f64",
      "enabled": true,
      "name": "patpat",
      "security": "wpapsk",
      "usergroup_id": "5d6d8b0c439adf048407dce8",
      "wlangroup_id": "5d6d8b0c439adf048407dce9",
      "x_passphrase": "forever home 23",
      "site_id": "5d6d8b07439adf048407dcd9",
      "x_iapp_key": "5aa68d6ab2ddedf45230f7f87bf11921",
      "minrate_ng_enabled": true,
      "minrate_ng_beacon_rate_kbps": 6000,
      "minrate_ng_data_rate_kbps": 6000,
      "vlan": "40",
      "vlan_enabled": true,
      "wep_idx": 1,
      "wpa_mode": "wpa2",
      "wpa_enc": "ccmp",
      "dtim_mode": "default",
      "dtim_ng": 1,
      "dtim_na": 1,
      "minrate_ng_advertising_rates": false,
      "minrate_ng_cck_rates_enabled": true,
      "minrate_na_enabled": false,
      "minrate_na_advertising_rates": false,
      "minrate_na_data_rate_kbps": 6000,
      "mac_filter_enabled": false,
      "mac_filter_policy": "allow",
      "mac_filter_list": [],
      "name_combine_enabled": true,
      "bc_filter_enabled": false,
      "bc_filter_list": [],
      "group_rekey": 3600,
      "radius_das_enabled": false,
      "schedule": [],
      "minrate_ng_mgmt_rate_kbps": 6000,
      "minrate_na_mgmt_rate_kbps": 6000,
      "minrate_na_beacon_rate_kbps": 6000
    }
  ]
}
*/

type WLAN struct {
	ID      string `json:"_id,omitempty"`
	SiteID  string `json:"site_id,omitempty"`
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`

	Security    string `json:"security"` // "wpapsk", "wpaeap", "open"
	WPAEnc      string `json:"wpa_enc"`  // "ccmp", "tkip"?
	WPAMode     string `json:"wpa_mode"` // "wpa2"
	XPassphrase string `json:"x_passphrase"`

	// create only?
	FastRoamingEnabled      bool `json:"fast_roaming_enabled,omitempty"`
	HideSSID                bool `json:"hide_ssid,omitempty"`
	IsGuest                 bool `json:"is_guest,omitempty"`
	MulticastEnhanceEnabled bool `json:"mcastenhance_enabled,omitempty"`

	RADIUSDasEnabled bool `json:"radius_das_enabled"`

	WLANGroupID string `json:"wlangroup_id"`

	NameCombineEnabled bool   `json:"name_combine_enabled"`
	NameCombineSuffix  string `json:"name_combine_suffix"`

	XIappKey string `json:"x_iapp_key,omitempty"`

	No2GhzOui   bool   `json:"no2ghz_oui"`
	WEPIdx      int    `json:"wep_idx,omitempty"`
	UserGroupID string `json:"usergroup_id"`
	DTIMMode    string `json:"dtim_mode"`
	DTIMNg      int    `json:"dtim_ng,omitempty"`
	DTIMNa      int    `json:"dtim_na,omitempty"`

	MinrateNgEnabled          bool `json:"minrate_ng_enabled"`
	MinrateNgBeaconRateKbps   int  `json:"minrate_ng_beacon_rate_kbps"`
	MinrateNgDataRateKbps     int  `json:"minrate_ng_data_rate_kbps"`
	MinrateNgAdvertisingRates bool `json:"minrate_ng_advertising_rates"`
	MinrateNgCckRatesEnabled  bool `json:"minrate_ng_cck_rates_enabled"`
	MinrateNaEnabled          bool `json:"minrate_na_enabled"`
	MinrateNaAdvertisingRates bool `json:"minrate_na_advertising_rates"`
	MinrateNaDataRateKbps     int  `json:"minrate_na_data_rate_kbps"`
	MinrateNgMgmtRateKbps     int  `json:"minrate_ng_mgmt_rate_kbps"`
	MinrateNaMgmtRateKbps     int  `json:"minrate_na_mgmt_rate_kbps"`
	MinrateNaBeaconRateKbps   int  `json:"minrate_na_beacon_rate_kbps"`

	MACFilterEnabled bool     `json:"mac_filter_enabled"`
	MACFilterPolicy  string   `json:"mac_filter_policy,omitempty"`
	MACFilterList    []string `json:"mac_filter_list,omitempty"`

	BroadcastFilterEnabled bool     `json:"bc_filter_enabled"`
	BroadcastFilterList    []string `json:"bc_filter_list,omitempty"`

	GroupRekey int `json:"group_rekey"`

	VLANEnabled bool   `json:"vlan_enabled"`
	VLAN        string `json:"vlan"`

	Schedule []string `json:"schedule"`
}

func (c *Client) ListWLAN(site string) ([]WLAN, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/wlanconf", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) GetWLAN(site, id string) (*WLAN, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/wlanconf/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) DeleteWLAN(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/wlanconf/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateWLAN(site string, d *WLAN) (*WLAN, error) {
	if d.Schedule == nil {
		d.Schedule = []string{}
	}

	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/wlanconf", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
