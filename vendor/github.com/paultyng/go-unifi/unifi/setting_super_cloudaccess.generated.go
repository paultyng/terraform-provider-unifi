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

type SettingSuperCloudaccess struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	DeviceAuth      string `json:"device_auth,omitempty"`
	DeviceID        string `json:"device_id"`
	Enabled         bool   `json:"enabled"`
	UbicUuid        string `json:"ubic_uuid,omitempty"`
	XCertificateArn string `json:"x_certificate_arn,omitempty"`
	XCertificatePem string `json:"x_certificate_pem,omitempty"`
	XPrivateKey     string `json:"x_private_key,omitempty"`
}

func (dst *SettingSuperCloudaccess) UnmarshalJSON(b []byte) error {
	type Alias SettingSuperCloudaccess
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

func (c *Client) getSettingSuperCloudaccess(ctx context.Context, site string) (*SettingSuperCloudaccess, error) {
	var respBody struct {
		Meta meta                      `json:"meta"`
		Data []SettingSuperCloudaccess `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/super_cloudaccess", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingSuperCloudaccess(ctx context.Context, site string, d *SettingSuperCloudaccess) (*SettingSuperCloudaccess, error) {
	var respBody struct {
		Meta meta                      `json:"meta"`
		Data []SettingSuperCloudaccess `json:"data"`
	}

	d.Key = "super_cloudaccess"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/super_cloudaccess", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
