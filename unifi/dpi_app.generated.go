// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

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

func (c *Client) listDpiApp(site string) ([]DpiApp, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []DpiApp `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/dpiapp", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getDpiApp(site, id string) (*DpiApp, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []DpiApp `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/dpiapp/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteDpiApp(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/dpiapp/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createDpiApp(site string, d *DpiApp) (*DpiApp, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []DpiApp `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/dpiapp", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateDpiApp(site string, d *DpiApp) (*DpiApp, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []DpiApp `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/dpiapp/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
