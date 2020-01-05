// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type SettingBaresip struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Enabled       bool   `json:"enabled"`
	OutboundProxy string `json:"outbound_proxy,omitempty"`
	PackageUrl    string `json:"package_url,omitempty"`
	Server        string `json:"server,omitempty"`
}
