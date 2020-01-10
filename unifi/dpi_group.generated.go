// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

type DpiGroup struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	DPIappIDs []string `json:"dpiapp_ids,omitempty"` // [\d\w]+
	Enabled   bool     `json:"enabled"`
	Name      string   `json:"name,omitempty"` // .{1,128}
}

func (c *Client) listDpiGroup(site string) ([]DpiGroup, error) {
	var respBody struct {
		Meta meta       `json:"meta"`
		Data []DpiGroup `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/dpigroup", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getDpiGroup(site, id string) (*DpiGroup, error) {
	var respBody struct {
		Meta meta       `json:"meta"`
		Data []DpiGroup `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/dpigroup/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteDpiGroup(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/dpigroup/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createDpiGroup(site string, d *DpiGroup) (*DpiGroup, error) {
	var respBody struct {
		Meta meta       `json:"meta"`
		Data []DpiGroup `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/dpigroup", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateDpiGroup(site string, d *DpiGroup) (*DpiGroup, error) {
	var respBody struct {
		Meta meta       `json:"meta"`
		Data []DpiGroup `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/dpigroup/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
