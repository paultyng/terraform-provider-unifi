// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type User struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Blocked     bool   `json:"blocked,omitempty"`
	FixedIP     string `json:"fixed_ip,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	LastSeen    string `json:"last_seen,omitempty"`
	MAC         string `json:"mac,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	Name        string `json:"name,omitempty"`
	NetworkID   string `json:"network_id"`
	Note        string `json:"note,omitempty"`
	UseFixedIP  bool   `json:"use_fixedip"`
	UserGroupID string `json:"usergroup_id"`
}
