// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

type WLANGroup struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	BSupported         bool   `json:"b_supported"`
	LoadbalanceEnabled bool   `json:"loadbalance_enabled"`
	Maxsta             int    `json:"maxsta,omitempty"` // [1-9]|[1-9][0-9]|1[0-9]{2}|200|^$
	MinRSSI            string `json:"minrssi,omitempty"`
	MinRSSIEnabled     bool   `json:"minrssi_enabled"`
	Name               string `json:"name,omitempty"`     // .{1,128}
	PMFMode            string `json:"pmf_mode,omitempty"` // disabled|optional|required
}

func (c *Client) listWLANGroup(site string) ([]WLANGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []WLANGroup `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/wlangroup", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getWLANGroup(site, id string) (*WLANGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []WLANGroup `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteWLANGroup(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createWLANGroup(site string, d *WLANGroup) (*WLANGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []WLANGroup `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/wlangroup", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateWLANGroup(site string, d *WLANGroup) (*WLANGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []WLANGroup `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
