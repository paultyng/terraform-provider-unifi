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

type HeatMap struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Description string `json:"description,omitempty"`
	MapID       string `json:"map_id"`
	Name        string `json:"name,omitempty"` // .*[^\s]+.*
	Type        string `json:"type,omitempty"` // download|upload
}

func (dst *HeatMap) UnmarshalJSON(b []byte) error {
	type Alias HeatMap
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

func (c *Client) listHeatMap(ctx context.Context, site string) ([]HeatMap, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []HeatMap `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/heatmap", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getHeatMap(ctx context.Context, site, id string) (*HeatMap, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []HeatMap `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/heatmap/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteHeatMap(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/heatmap/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createHeatMap(ctx context.Context, site string, d *HeatMap) (*HeatMap, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []HeatMap `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/heatmap", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateHeatMap(ctx context.Context, site string, d *HeatMap) (*HeatMap, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []HeatMap `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/heatmap/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
