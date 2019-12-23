package unifi

import (
	"fmt"
)

type Network struct {
	ID       string `json:"_id,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`

	Purpose      string `json:"purpose"`      // "corporate"
	NetworkGroup string `json:"networkgroup"` // "LAN"
	Name         string `json:"name"`
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

func (c *Client) ListNetworks(site string) ([]Network, error) {
	var respBody struct {
		Meta struct {
			RC string `json:"rc"`
		} `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/networkconf", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) GetNetwork(site, id string) (*Network, error) {
	list, err := c.ListNetworks(site)
	if err != nil {
		return nil, err
	}
	for _, net := range list {
		if net.ID == id {
			return &net, nil
		}
	}
	return nil, &NotFoundError{}

	// 	var respBody struct {
	// 		Meta struct {
	// 			RC string `json:"rc"`
	// 		} `json:"meta"`
	// 		Data []Network `json:"data"`
	// 	}

	// 	err := c.do("GET", fmt.Sprintf("s/%s/rest/networkconf/%s", site, id), nil, &respBody)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if len(respBody.Data) != 1 {
	// 		return nil, &NotFoundError{}
	// 	}

	// 	net := respBody.Data[0]
	// 	return &net, nil
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

func (c *Client) CreateNetwork(site string, network *Network) (*Network, error) {
	var respBody struct {
		Meta struct {
			RC string `json:"rc"`
		} `json:"meta"`
		Data []Network `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/networkconf", site), network, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	newNetwork := respBody.Data[0]

	return &newNetwork, nil
}
