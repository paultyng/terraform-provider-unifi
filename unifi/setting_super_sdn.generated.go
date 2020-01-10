// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

type SettingSuperSdn struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	AuthToken         string   `json:"auth_token,omitempty"`
	DeviceID          string   `json:"device_id"`
	Enabled           bool     `json:"enabled"`
	Migrated          bool     `json:"migrated"`
	OauthAppID        string   `json:"oauth_app_id"`
	OauthEnabled      bool     `json:"oauth_enabled"`
	OauthRedirectUris []string `json:"oauth_redirect_uris,omitempty"`
	SsoLoginEnabled   string   `json:"sso_login_enabled,omitempty"`
	UbicUuid          string   `json:"ubic_uuid,omitempty"`
	XOauthAppSecret   string   `json:"x_oauth_app_secret,omitempty"`
}
