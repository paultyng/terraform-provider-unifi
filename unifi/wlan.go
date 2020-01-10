package unifi

import (
	"encoding/json"
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

func (c *Client) CreateWLAN(site string, d *WLAN) (*WLAN, error) {
	if d.Schedule == nil {
		d.Schedule = []string{}
	}

	return c.createWLAN(site, d)
}

func (c *Client) ListWLAN(site string) ([]WLAN, error) {
	return c.listWLAN(site)
}

func (c *Client) GetWLAN(site, id string) (*WLAN, error) {
	return c.getWLAN(site, id)
}

func (c *Client) DeleteWLAN(site, id string) error {
	return c.deleteWLAN(site, id)
}

func (c *Client) UpdateWLAN(site string, d *WLAN) (*WLAN, error) {
	return c.updateWLAN(site, d)
}
