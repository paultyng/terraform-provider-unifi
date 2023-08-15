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

type Dashboard struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	ControllerVersion string             `json:"controller_version,omitempty"`
	Desc              string             `json:"desc,omitempty"`
	IsPublic          bool               `json:"is_public"`
	Modules           []DashboardModules `json:"modules,omitempty"`
	Name              string             `json:"name,omitempty"`
}

func (dst *Dashboard) UnmarshalJSON(b []byte) error {
	type Alias Dashboard
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

type DashboardModules struct {
	Config       string `json:"config,omitempty"`
	ID           string `json:"id"`
	ModuleID     string `json:"module_id"`
	Restrictions string `json:"restrictions,omitempty"`
}

func (dst *DashboardModules) UnmarshalJSON(b []byte) error {
	type Alias DashboardModules
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

func (c *Client) listDashboard(ctx context.Context, site string) ([]Dashboard, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []Dashboard `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/dashboard", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getDashboard(ctx context.Context, site, id string) (*Dashboard, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []Dashboard `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/dashboard/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteDashboard(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/dashboard/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createDashboard(ctx context.Context, site string, d *Dashboard) (*Dashboard, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []Dashboard `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/dashboard", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateDashboard(ctx context.Context, site string, d *Dashboard) (*Dashboard, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []Dashboard `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/dashboard/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
