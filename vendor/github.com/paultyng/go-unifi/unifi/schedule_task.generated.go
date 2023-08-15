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

type ScheduleTask struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Action                  string                       `json:"action,omitempty"` // stream|upgrade
	AdditionalSoundsEnabled bool                         `json:"additional_sounds_enabled"`
	BroadcastgroupID        string                       `json:"broadcastgroup_id"`
	CronExpr                string                       `json:"cron_expr,omitempty"`
	ExecuteOnlyOnce         bool                         `json:"execute_only_once"`
	MediafileID             string                       `json:"mediafile_id"`
	Name                    string                       `json:"name,omitempty"`
	SampleFilename          string                       `json:"sample_filename,omitempty"`
	StreamType              string                       `json:"stream_type,omitempty"` // media|sample
	UpgradeTargets          []ScheduleTaskUpgradeTargets `json:"upgrade_targets,omitempty"`
}

func (dst *ScheduleTask) UnmarshalJSON(b []byte) error {
	type Alias ScheduleTask
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

type ScheduleTaskUpgradeTargets struct {
	MAC string `json:"mac,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
}

func (dst *ScheduleTaskUpgradeTargets) UnmarshalJSON(b []byte) error {
	type Alias ScheduleTaskUpgradeTargets
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

func (c *Client) listScheduleTask(ctx context.Context, site string) ([]ScheduleTask, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []ScheduleTask `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/scheduletask", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getScheduleTask(ctx context.Context, site, id string) (*ScheduleTask, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []ScheduleTask `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/scheduletask/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteScheduleTask(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/scheduletask/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createScheduleTask(ctx context.Context, site string, d *ScheduleTask) (*ScheduleTask, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []ScheduleTask `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/scheduletask", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateScheduleTask(ctx context.Context, site string, d *ScheduleTask) (*ScheduleTask, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []ScheduleTask `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/scheduletask/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
