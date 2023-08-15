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

func (dst *FirewallGroup) UnmarshalJSON(b []byte) error {
	type Alias FirewallGroup
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

func (c *Client) listFirewallGroup(ctx context.Context, site string) ([]FirewallGroup, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []FirewallGroup `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/firewallgroup", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getFirewallGroup(ctx context.Context, site, id string) (*FirewallGroup, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []FirewallGroup `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/firewallgroup/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteFirewallGroup(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/firewallgroup/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createFirewallGroup(ctx context.Context, site string, d *FirewallGroup) (*FirewallGroup, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []FirewallGroup `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/firewallgroup", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateFirewallGroup(ctx context.Context, site string, d *FirewallGroup) (*FirewallGroup, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []FirewallGroup `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/firewallgroup/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
