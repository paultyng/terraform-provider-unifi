// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

type BroadcastGroup struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	MemberTable []string `json:"member_table,omitempty"`
	Name        string   `json:"name,omitempty"`
}

func (c *Client) listBroadcastGroup(site string) ([]BroadcastGroup, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []BroadcastGroup `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/broadcastgroup", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getBroadcastGroup(site, id string) (*BroadcastGroup, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []BroadcastGroup `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/broadcastgroup/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteBroadcastGroup(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/broadcastgroup/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createBroadcastGroup(site string, d *BroadcastGroup) (*BroadcastGroup, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []BroadcastGroup `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/broadcastgroup", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateBroadcastGroup(site string, d *BroadcastGroup) (*BroadcastGroup, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []BroadcastGroup `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/broadcastgroup/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
