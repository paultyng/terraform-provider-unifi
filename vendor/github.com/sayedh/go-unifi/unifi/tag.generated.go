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

type Tag struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	MemberTable []string `json:"member_table,omitempty"`
	Name        string   `json:"name,omitempty"`
}

func (dst *Tag) UnmarshalJSON(b []byte) error {
	type Alias Tag
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

func (c *Client) listTag(ctx context.Context, site string) ([]Tag, error) {
	var respBody struct {
		Meta meta  `json:"meta"`
		Data []Tag `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/tag", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getTag(ctx context.Context, site, id string) (*Tag, error) {
	var respBody struct {
		Meta meta  `json:"meta"`
		Data []Tag `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/tag/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteTag(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/tag/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createTag(ctx context.Context, site string, d *Tag) (*Tag, error) {
	var respBody struct {
		Meta meta  `json:"meta"`
		Data []Tag `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/tag", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateTag(ctx context.Context, site string, d *Tag) (*Tag, error) {
	var respBody struct {
		Meta meta  `json:"meta"`
		Data []Tag `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/tag/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
