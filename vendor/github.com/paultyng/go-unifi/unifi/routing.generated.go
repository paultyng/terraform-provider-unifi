// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type Routing struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Enabled              bool   `json:"enabled"`
	GatewayDevice        string `json:"gateway_device,omitempty"`        // ^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$
	GatewayType          string `json:"gateway_type,omitempty"`          // default|switch
	Name                 string `json:"name,omitempty"`                  // .{1,128}
	StaticRouteDistance  int    `json:"static-route_distance,omitempty"` // ^[1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]$|^$
	StaticRouteInterface string `json:"static-route_interface"`          // WAN1|WAN2|[\d\w]+|^$
	StaticRouteNetwork   string `json:"static-route_network,omitempty"`  // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2][0-9]|3[0-2])$|^([a-fA-F0-9:]+\/(([1-9]|[1-8][0-9]|9[0-9]|1[01][0-9]|12[0-8])))$
	StaticRouteNexthop   string `json:"static-route_nexthop"`            // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^([a-fA-F0-9:]+)$|^$
	StaticRouteType      string `json:"static-route_type,omitempty"`     // nexthop-route|interface-route|blackhole
	Type                 string `json:"type,omitempty"`                  // static-route
}

func (dst *Routing) UnmarshalJSON(b []byte) error {
	type Alias Routing
	aux := &struct {
		StaticRouteDistance emptyStringInt `json:"static-route_distance"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.StaticRouteDistance = int(aux.StaticRouteDistance)

	return nil
}

func (c *Client) listRouting(ctx context.Context, site string) ([]Routing, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Routing `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/routing", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getRouting(ctx context.Context, site, id string) (*Routing, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Routing `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/routing/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteRouting(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/routing/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createRouting(ctx context.Context, site string, d *Routing) (*Routing, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Routing `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/routing", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateRouting(ctx context.Context, site string, d *Routing) (*Routing, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []Routing `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/routing/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
