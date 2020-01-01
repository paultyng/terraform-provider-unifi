// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type Account struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	IP               string `json:"ip"`                           // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	Name             string `json:"name,omitempty"`               // ^[^"' ]+$
	TunnelConfigType string `json:"tunnel_config_type,omitempty"` // vpn|802.1x|custom
	TunnelMediumType int    `json:"tunnel_medium_type,omitempty"` // [1-9]|1[0-5]|^$
	TunnelType       int    `json:"tunnel_type,omitempty"`        // [1-9]|1[0-3]|^$
	VLAN             int    `json:"vlan,omitempty"`               // [2-9]|[1-9][0-9]{1,2}|[1-3][0-9]{3}|400[0-9]|^$
	XPassword        string `json:"x_password,omitempty"`
}
