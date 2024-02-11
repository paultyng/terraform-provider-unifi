// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type User struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	DevIdOverride int    `json:"dev_id_override,omitempty"` // non-generated field
	IP            string `json:"ip,omitempty"`              // non-generated field

	Blocked                       bool   `json:"blocked,omitempty"`
	FixedApEnabled                bool   `json:"fixed_ap_enabled"`
	FixedApMAC                    string `json:"fixed_ap_mac,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	FixedIP                       string `json:"fixed_ip,omitempty"`
	Hostname                      string `json:"hostname,omitempty"`
	LastSeen                      int    `json:"last_seen,omitempty"`
	LocalDNSRecord                string `json:"local_dns_record,omitempty"`
	LocalDNSRecordEnabled         bool   `json:"local_dns_record_enabled"`
	MAC                           string `json:"mac,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	Name                          string `json:"name,omitempty"`
	NetworkID                     string `json:"network_id"`
	Note                          string `json:"note,omitempty"`
	UseFixedIP                    bool   `json:"use_fixedip"`
	UserGroupID                   string `json:"usergroup_id"`
	VirtualNetworkOverrideEnabled bool   `json:"virtual_network_override_enabled"`
	VirtualNetworkOverrideID      string `json:"virtual_network_override_id"`
}

func (dst *User) UnmarshalJSON(b []byte) error {
	type Alias User
	aux := &struct {
		LastSeen emptyStringInt `json:"last_seen"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.LastSeen = int(aux.LastSeen)

	return nil
}

func (c *Client) listUser(ctx context.Context, site string) ([]User, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/user", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getUser(ctx context.Context, site, id string) (*User, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/user/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteUser(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/user/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createUser(ctx context.Context, site string, d *User) (*User, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/user", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateUser(ctx context.Context, site string, d *User) (*User, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/user/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
