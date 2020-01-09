package unifi

import (
	"encoding/json"
	"fmt"
)

func (n *WLAN) UnmarshalJSON(b []byte) error {
	type Alias WLAN
	aux := &struct {
		VLAN json.Number `json:"vlan"`
		*Alias
	}{
		Alias: (*Alias)(n),
	}
	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}
	n.VLAN = 0
	if aux.VLAN.String() != "" {
		vlan, err := aux.VLAN.Int64()
		if err != nil {
			return err
		}
		n.VLAN = int(vlan)
	}
	return nil
}

func (c *Client) ListWLAN(site string) ([]WLAN, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/wlanconf", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) GetWLAN(site, id string) (*WLAN, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/wlanconf/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) DeleteWLAN(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/wlanconf/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateWLAN(site string, d *WLAN) (*WLAN, error) {
	if d.Schedule == nil {
		d.Schedule = []string{}
	}

	var respBody struct {
		Meta meta   `json:"meta"`
		Data []WLAN `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/wlanconf", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
