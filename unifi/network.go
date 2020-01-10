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

func (c *Client) ListNetwork(site string) ([]Network, error) {
	return c.listNetwork(site)
}

func (c *Client) GetNetwork(site, id string) (*Network, error) {
	return c.getNetwork(site, id)
}

func (c *Client) CreateNetwork(site string, d *Network) (*Network, error) {
	return c.createNetwork(site, d)
}

func (c *Client) UpdateNetwork(site string, d *Network) (*Network, error) {
	return c.updateNetwork(site, d)
}
