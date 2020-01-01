package unifi

import (
	"encoding/json"
	"fmt"
)

func (n *Network) UnmarshalJSON(b []byte) error {
	type Alias Network
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

func (c *Client) ListNetwork(site string) ([]Network, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/networkconf", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) GetNetwork(site, id string) (*Network, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/networkconf/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) DeleteNetwork(site, id, name string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/networkconf/%s", site, id), struct {
		Name string `json:"name"`
	}{
		Name: name,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateNetwork(site string, d *Network) (*Network, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/networkconf", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) UpdateNetwork(site string, d *Network) (*Network, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/networkconf/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
