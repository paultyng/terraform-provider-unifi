// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type SettingProviderCapabilities struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Download int  `json:"download,omitempty"` // ^[1-9][0-9]*$
	Enabled  bool `json:"enabled"`
	Upload   int  `json:"upload,omitempty"` // ^[1-9][0-9]*$
}
