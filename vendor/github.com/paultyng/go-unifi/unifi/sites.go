package unifi

import (
	"context"
	"fmt"
)

type Site struct {
	ID string `json:"_id,omitempty"`

	// Hidden   bool   `json:"attr_hidden,omitempty"`
	// HiddenID string `json:"attr_hidden_id,omitempty"`
	// NoDelete bool   `json:"attr_no_delete,omitempty"`
	// NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Name        string `json:"name"`
	Description string `json:"desc"`

	//Role string `json:"role"`
}

func (c *Client) ListSites(ctx context.Context) ([]Site, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []Site `json:"data"`
	}

	err := c.do(ctx, "GET", "self/sites", nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) GetSite(ctx context.Context, id string) (*Site, error) {
	sites, err := c.ListSites(ctx)

	if err != nil {
		return nil, err
	}

	for _, s := range sites {
		if s.ID == id {
			return &s, nil
		}
	}

	return nil, &NotFoundError{}
}

func (c *Client) CreateSite(ctx context.Context, description string) ([]Site, error) {
	reqBody := struct {
		Cmd  string `json:"cmd"`
		Desc string `json:"desc"`
	}{
		Cmd:  "add-site",
		Desc: description,
	}

	var respBody struct {
		Meta meta   `json:"meta"`
		Data []Site `json:"data"`
	}

	err := c.do(ctx, "POST", "s/default/cmd/sitemgr", reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) DeleteSite(ctx context.Context, id string) ([]Site, error) {
	reqBody := struct {
		Cmd  string `json:"cmd"`
		Site string `json:"site"`
	}{
		Cmd:  "delete-site",
		Site: id,
	}

	var respBody struct {
		Meta meta   `json:"meta"`
		Data []Site `json:"data"`
	}

	err := c.do(ctx, "POST", "s/default/cmd/sitemgr", reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) UpdateSite(ctx context.Context, name, description string) ([]Site, error) {
	reqBody := struct {
		Cmd  string `json:"cmd"`
		Desc string `json:"desc"`
	}{
		Cmd:  "update-site",
		Desc: description,
	}

	var respBody struct {
		Meta meta   `json:"meta"`
		Data []Site `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/cmd/sitemgr", name), reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}
