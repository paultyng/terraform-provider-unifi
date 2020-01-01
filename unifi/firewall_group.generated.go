// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type FirewallGroup struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	GroupMembers []string `json:"group_members,omitempty"`
	GroupType    string   `json:"group_type,omitempty"` // address-group|port-group|ipv6-address-group
	Name         string   `json:"name,omitempty"`       // .{1,64}
}
