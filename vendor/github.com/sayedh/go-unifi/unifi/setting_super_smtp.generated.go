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

type SettingSuperSmtp struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	Enabled   bool   `json:"enabled"`
	Host      string `json:"host,omitempty"`
	Port      int    `json:"port,omitempty"` // [1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5]|^$
	Sender    string `json:"sender,omitempty"`
	UseAuth   bool   `json:"use_auth"`
	UseSender bool   `json:"use_sender"`
	UseSsl    bool   `json:"use_ssl"`
	Username  string `json:"username,omitempty"`
	XPassword string `json:"x_password,omitempty"`
}

func (dst *SettingSuperSmtp) UnmarshalJSON(b []byte) error {
	type Alias SettingSuperSmtp
	aux := &struct {
		Port emptyStringInt `json:"port"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Port = int(aux.Port)

	return nil
}

func (c *Client) getSettingSuperSmtp(ctx context.Context, site string) (*SettingSuperSmtp, error) {
	var respBody struct {
		Meta meta               `json:"meta"`
		Data []SettingSuperSmtp `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/super_smtp", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingSuperSmtp(ctx context.Context, site string, d *SettingSuperSmtp) (*SettingSuperSmtp, error) {
	var respBody struct {
		Meta meta               `json:"meta"`
		Data []SettingSuperSmtp `json:"data"`
	}

	d.Key = "super_smtp"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/super_smtp", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
