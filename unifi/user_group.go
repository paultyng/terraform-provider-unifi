package unifi

import (
	"fmt"
)

type UserGroup struct {
	ID     string `json:"_id"`
	SiteID string `json:"site_id"`
	Name   string `json:"name"`

	//Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	//NoEdit   bool   `json:"attr_no_edit,omitempty"`

	QOSRateMaxDown int `json:"qos_rate_max_down"`
	QOSRateMaxUp   int `json:"qos_rate_max_up"`
}

func (c *Client) ListUserGroup(site string) ([]UserGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []UserGroup `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/usergroup", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}
