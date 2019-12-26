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
	PMFMode            string `json:"pmf_mode"`
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
