// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

type DHCPOption struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Code   string `json:"code,omitempty"` // ^(?!(?:15|42|43|44|51|66|67|252)$)([7-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-4])$
	Name   string `json:"name,omitempty"` // ^[A-Za-z0-9-_]{1,25}$
	Signed bool   `json:"signed"`
	Type   string `json:"type,omitempty"`  // ^(boolean|hexarray|integer|ipaddress|macaddress|text)$
	Width  int    `json:"width,omitempty"` // ^(8|16|32)$
}

func (c *Client) listDHCPOption(site string) ([]DHCPOption, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []DHCPOption `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/dhcpoption", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getDHCPOption(site, id string) (*DHCPOption, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []DHCPOption `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/dhcpoption/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteDHCPOption(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/dhcpoption/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createDHCPOption(site string, d *DHCPOption) (*DHCPOption, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []DHCPOption `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/dhcpoption", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateDHCPOption(site string, d *DHCPOption) (*DHCPOption, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []DHCPOption `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/dhcpoption/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
