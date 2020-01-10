// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

type SettingMgmt struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	AdvancedFeatureEnabled  bool     `json:"advanced_feature_enabled"`
	AlertEnabled            bool     `json:"alert_enabled"`
	AutoUpgrade             bool     `json:"auto_upgrade"`
	LedEnabled              bool     `json:"led_enabled"`
	OutdoorModeEnabled      bool     `json:"outdoor_mode_enabled"`
	UnifiIDpEnabled         bool     `json:"unifi_idp_enabled"`
	XMgmtKey                string   `json:"x_mgmt_key,omitempty"` // [0-9a-f]{32}
	XSshAuthPasswordEnabled bool     `json:"x_ssh_auth_password_enabled"`
	XSshBindWildcard        bool     `json:"x_ssh_bind_wildcard"`
	XSshEnabled             bool     `json:"x_ssh_enabled"`
	XSshKeys                []string `json:"x_ssh_keys,omitempty"`
	XSshMd5Passwd           string   `json:"x_ssh_md5passwd,omitempty"`
	XSshPassword            string   `json:"x_ssh_password,omitempty"` // .{1,128}
	XSshSha512Passwd        string   `json:"x_ssh_sha512passwd,omitempty"`
	XSshUsername            string   `json:"x_ssh_username,omitempty"` // ^[_A-Za-z0-9][-_.A-Za-z0-9]{0,29}$
}
