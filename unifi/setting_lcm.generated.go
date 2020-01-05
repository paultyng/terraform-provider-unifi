// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type SettingLcm struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Brightness  int  `json:"brightness,omitempty"` // [1-9]|[1-9][0-9]|100
	Enabled     bool `json:"enabled"`
	IDleTimeout int  `json:"idle_timeout,omitempty"` // [1-9][0-9]|[1-9][0-9][0-9]|[1-2][0-9][0-9][0-9]|3[0-5][0-9][0-9]|3600
	Sync        bool `json:"sync"`
	TouchEvent  bool `json:"touch_event"`
}
