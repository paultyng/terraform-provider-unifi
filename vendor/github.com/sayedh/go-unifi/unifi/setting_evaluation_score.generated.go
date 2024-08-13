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

type SettingEvaluationScore struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	DismissedIDs []string `json:"dismissed_ids,omitempty"` // ^[a-zA-Z]{2}[0-9]{2,3}$|^$
}

func (dst *SettingEvaluationScore) UnmarshalJSON(b []byte) error {
	type Alias SettingEvaluationScore
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

func (c *Client) getSettingEvaluationScore(ctx context.Context, site string) (*SettingEvaluationScore, error) {
	var respBody struct {
		Meta meta                     `json:"meta"`
		Data []SettingEvaluationScore `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/evaluation_score", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingEvaluationScore(ctx context.Context, site string, d *SettingEvaluationScore) (*SettingEvaluationScore, error) {
	var respBody struct {
		Meta meta                     `json:"meta"`
		Data []SettingEvaluationScore `json:"data"`
	}

	d.Key = "evaluation_score"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/evaluation_score", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
