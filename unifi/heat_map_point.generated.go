// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"fmt"
)

// just to fix compile issues with the import
var _ fmt.Formatter

type HeatMapPoint struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	DownloadSpeed float64 `json:"download_speed,omitempty"`
	HeatmapID     string  `json:"heatmap_id"`
	UploadSpeed   float64 `json:"upload_speed,omitempty"`
	X             float64 `json:"x,omitempty"`
	Y             float64 `json:"y,omitempty"`
}

func (c *Client) listHeatMapPoint(site string) ([]HeatMapPoint, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []HeatMapPoint `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/heatmappoint", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getHeatMapPoint(site, id string) (*HeatMapPoint, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []HeatMapPoint `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/heatmappoint/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteHeatMapPoint(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/heatmappoint/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createHeatMapPoint(site string, d *HeatMapPoint) (*HeatMapPoint, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []HeatMapPoint `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/heatmappoint", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateHeatMapPoint(site string, d *HeatMapPoint) (*HeatMapPoint, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []HeatMapPoint `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/heatmappoint/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
