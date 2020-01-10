// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

type VirtualDevice struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	HeightInMeters float64 `json:"heightInMeters,omitempty"`
	Locked         bool    `json:"locked"`
	MapID          string  `json:"map_id"`
	Type           string  `json:"type,omitempty"` // uap|usg|usw
	X              string  `json:"x,omitempty"`
	Y              string  `json:"y,omitempty"`
}

func (c *Client) listVirtualDevice(site string) ([]VirtualDevice, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []VirtualDevice `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/virtualdevice", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getVirtualDevice(site, id string) (*VirtualDevice, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []VirtualDevice `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/virtualdevice/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteVirtualDevice(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/virtualdevice/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createVirtualDevice(site string, d *VirtualDevice) (*VirtualDevice, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []VirtualDevice `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/virtualdevice", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateVirtualDevice(site string, d *VirtualDevice) (*VirtualDevice, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []VirtualDevice `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/virtualdevice/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
