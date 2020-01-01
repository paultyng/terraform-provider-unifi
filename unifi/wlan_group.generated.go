// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type WLANGroup struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	BSupported         bool   `json:"b_supported"`
	LoadbalanceEnabled bool   `json:"loadbalance_enabled"`
	Maxsta             int    `json:"maxsta,omitempty"` // [1-9]|[1-9][0-9]|1[0-9]{2}|200|^$
	MinRSSI            string `json:"minrssi,omitempty"`
	MinRSSIEnabled     bool   `json:"minrssi_enabled"`
	Name               string `json:"name,omitempty"`     // .{1,128}
	PMFMode            string `json:"pmf_mode,omitempty"` // disabled|optional|required
}
