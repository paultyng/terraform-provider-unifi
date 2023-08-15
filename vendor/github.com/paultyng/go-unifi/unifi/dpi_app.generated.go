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

type DpiApp struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Apps           []int  `json:"apps,omitempty"`
	Blocked        bool   `json:"blocked"`
	Cats           []int  `json:"cats,omitempty"`
	Enabled        bool   `json:"enabled"`
	Log            bool   `json:"log"`
	Name           string `json:"name,omitempty"`              // .{1,128}
	QOSRateMaxDown int    `json:"qos_rate_max_down,omitempty"` // -1|[2-9]|[1-9][0-9]{1,4}|100000|10[0-1][0-9]{3}|102[0-3][0-9]{2}|102400
	QOSRateMaxUp   int    `json:"qos_rate_max_up,omitempty"`   // -1|[2-9]|[1-9][0-9]{1,4}|100000|10[0-1][0-9]{3}|102[0-3][0-9]{2}|102400
}

func (dst *DpiApp) UnmarshalJSON(b []byte) error {
	type Alias DpiApp
	aux := &struct {
		Apps           []emptyStringInt `json:"apps"`
		Cats           []emptyStringInt `json:"cats"`
		QOSRateMaxDown emptyStringInt   `json:"qos_rate_max_down"`
		QOSRateMaxUp   emptyStringInt   `json:"qos_rate_max_up"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Apps = make([]int, len(aux.Apps))
	for i, v := range aux.Apps {
		dst.Apps[i] = int(v)
	}
	dst.Cats = make([]int, len(aux.Cats))
	for i, v := range aux.Cats {
		dst.Cats[i] = int(v)
	}
	dst.QOSRateMaxDown = int(aux.QOSRateMaxDown)
	dst.QOSRateMaxUp = int(aux.QOSRateMaxUp)

	return nil
}

func (c *Client) listDpiApp(ctx context.Context, site string) ([]DpiApp, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []DpiApp `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/dpiapp", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getDpiApp(ctx context.Context, site, id string) (*DpiApp, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []DpiApp `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/dpiapp/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteDpiApp(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/dpiapp/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createDpiApp(ctx context.Context, site string, d *DpiApp) (*DpiApp, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []DpiApp `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/dpiapp", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateDpiApp(ctx context.Context, site string, d *DpiApp) (*DpiApp, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []DpiApp `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/dpiapp/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
