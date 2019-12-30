package unifi

import (
	"fmt"
)

type Network struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`
	Name   string `json:"name"`

	// Hidden   bool   `json:"attr_hidden,omitempty"`
	// HiddenID string `json:"attr_hidden_id,omitempty"`
	// NoDelete bool   `json:"attr_no_delete,omitempty"`
	// NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Purpose      string `json:"purpose"`      // "corporate"
	NetworkGroup string `json:"networkgroup"` // "LAN"
	VLAN         string `json:"vlan"`
	VLANEnabled  bool   `json:"vlan_enabled"`
	IPSubnet     string `json:"ip_subnet"`
	Enabled      bool   `json:"enabled"`
	IsNAT        bool   `json:"is_nat"`
	DomainName   string `json:"domain_name"`

	DHCPRelayEnabled bool `json:"dhcp_relay_enabled"`

	DHCPDStart             string `json:"dhcpd_start"`     // "10.0.0.6"
	DHCPDStop              string `json:"dhcpd_stop"`      // "10.0.0.254"
	DHCPDEnabled           bool   `json:"dhcpd_enabled"`   // true
	DHCPDLeaseTime         int    `json:"dhcpd_leasetime"` // 86400
	DHCPDDNSEnabled        bool   `json:"dhcpd_dns_enabled"`
	DHCPDGatewayEnabled    bool   `json:"dhcpd_gateway_enabled"`
	DHCPDTimeOffsetEnabled bool   `json:"dhcpd_time_offset_enabled"`

	IPV6InterfaceType string `json:"ipv6_interface_type"` // "none"
	IPV6PDStart       string `json:"ipv6_pd_start"`       // "::2"
	IPV6PDStop        string `json:"ipv6_pd_stop"`        // "::7d1"
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
