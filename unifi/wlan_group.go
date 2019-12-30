package unifi

import (
	"fmt"
)

type WLANGroup struct {
	ID     string `json:"_id"`
	SiteID string `json:"site_id"`
	Name   string `json:"name"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	BSupported         bool   `json:"b_supported"`
	LoadBalanceEnabled bool   `json:"loadbalance_enabled"`
	PMFMode            string `json:"pmf_mode"` // "disabled", "optional", "required"
	// roam_channel_na: 36
	// roam_channel_ng: 1
	// roam_enabled: false
	// roam_radio: "ng"
}

func (c *Client) ListWLANGroup(site string) ([]WLANGroup, error) {
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

func (c *Client) GetWLANGroup(site, id string) (*WLANGroup, error) {
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

func (c *Client) DeleteWLANGroup(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/wlangroup/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateWLANGroup(site string, d *WLANGroup) (*WLANGroup, error) {
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

func (c *Client) UpdateWLANGroup(site string, d *WLANGroup) (*WLANGroup, error) {
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
