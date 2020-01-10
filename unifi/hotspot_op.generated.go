// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

type HotspotOp struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Name      string `json:"name,omitempty"` // .{1,256}
	Note      string `json:"note,omitempty"`
	XPassword string `json:"x_password,omitempty"` // .{1,256}
}

func (c *Client) listHotspotOp(site string) ([]HotspotOp, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []HotspotOp `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/hotspotop", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getHotspotOp(site, id string) (*HotspotOp, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []HotspotOp `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/hotspotop/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteHotspotOp(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/hotspotop/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createHotspotOp(site string, d *HotspotOp) (*HotspotOp, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []HotspotOp `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/hotspotop", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateHotspotOp(site string, d *HotspotOp) (*HotspotOp, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []HotspotOp `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/hotspotop/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
