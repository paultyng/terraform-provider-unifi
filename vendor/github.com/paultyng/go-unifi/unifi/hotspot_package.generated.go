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

type HotspotPackage struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Amount                         float64 `json:"amount,omitempty"`
	ChargedAs                      string  `json:"charged_as,omitempty"`
	Currency                       string  `json:"currency,omitempty"` // [A-Z]{3}
	CustomPaymentFieldsEnabled     bool    `json:"custom_payment_fields_enabled"`
	Hours                          int     `json:"hours,omitempty"`
	Index                          int     `json:"index,omitempty"`
	LimitDown                      int     `json:"limit_down,omitempty"`
	LimitOverwrite                 bool    `json:"limit_overwrite"`
	LimitQuota                     int     `json:"limit_quota,omitempty"`
	LimitUp                        int     `json:"limit_up,omitempty"`
	Name                           string  `json:"name,omitempty"`
	PaymentFieldsAddressEnabled    bool    `json:"payment_fields_address_enabled"`
	PaymentFieldsAddressRequired   bool    `json:"payment_fields_address_required"`
	PaymentFieldsCityEnabled       bool    `json:"payment_fields_city_enabled"`
	PaymentFieldsCityRequired      bool    `json:"payment_fields_city_required"`
	PaymentFieldsCountryEnabled    bool    `json:"payment_fields_country_enabled"`
	PaymentFieldsCountryRequired   bool    `json:"payment_fields_country_required"`
	PaymentFieldsEmailEnabled      bool    `json:"payment_fields_email_enabled"`
	PaymentFieldsEmailRequired     bool    `json:"payment_fields_email_required"`
	PaymentFieldsFirstNameEnabled  bool    `json:"payment_fields_first_name_enabled"`
	PaymentFieldsFirstNameRequired bool    `json:"payment_fields_first_name_required"`
	PaymentFieldsLastNameEnabled   bool    `json:"payment_fields_last_name_enabled"`
	PaymentFieldsLastNameRequired  bool    `json:"payment_fields_last_name_required"`
	PaymentFieldsStateEnabled      bool    `json:"payment_fields_state_enabled"`
	PaymentFieldsStateRequired     bool    `json:"payment_fields_state_required"`
	PaymentFieldsZipEnabled        bool    `json:"payment_fields_zip_enabled"`
	PaymentFieldsZipRequired       bool    `json:"payment_fields_zip_required"`
	TrialDurationMinutes           int     `json:"trial_duration_minutes,omitempty"`
	TrialReset                     float64 `json:"trial_reset,omitempty"`
}

func (dst *HotspotPackage) UnmarshalJSON(b []byte) error {
	type Alias HotspotPackage
	aux := &struct {
		Hours                emptyStringInt `json:"hours"`
		Index                emptyStringInt `json:"index"`
		LimitDown            emptyStringInt `json:"limit_down"`
		LimitQuota           emptyStringInt `json:"limit_quota"`
		LimitUp              emptyStringInt `json:"limit_up"`
		TrialDurationMinutes emptyStringInt `json:"trial_duration_minutes"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Hours = int(aux.Hours)
	dst.Index = int(aux.Index)
	dst.LimitDown = int(aux.LimitDown)
	dst.LimitQuota = int(aux.LimitQuota)
	dst.LimitUp = int(aux.LimitUp)
	dst.TrialDurationMinutes = int(aux.TrialDurationMinutes)

	return nil
}

func (c *Client) listHotspotPackage(ctx context.Context, site string) ([]HotspotPackage, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []HotspotPackage `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/hotspotpackage", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getHotspotPackage(ctx context.Context, site, id string) (*HotspotPackage, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []HotspotPackage `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/hotspotpackage/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteHotspotPackage(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/hotspotpackage/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createHotspotPackage(ctx context.Context, site string, d *HotspotPackage) (*HotspotPackage, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []HotspotPackage `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/hotspotpackage", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateHotspotPackage(ctx context.Context, site string, d *HotspotPackage) (*HotspotPackage, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []HotspotPackage `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/hotspotpackage/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
