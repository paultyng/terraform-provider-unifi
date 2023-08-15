package unifi

import (
	"context"
	"fmt"
)

//go:generate go run golang.org/x/tools/cmd/stringer -trimprefix DeviceState -type DeviceState
type DeviceState int

const (
	DeviceStateUnknown          DeviceState = 0
	DeviceStateConnected        DeviceState = 1
	DeviceStatePending          DeviceState = 2
	DeviceStateFirmwareMismatch DeviceState = 3
	DeviceStateUpgrading        DeviceState = 4
	DeviceStateProvisioning     DeviceState = 5
	DeviceStateHeartbeatMissed  DeviceState = 6
	DeviceStateAdopting         DeviceState = 7
	DeviceStateDeleting         DeviceState = 8
	DeviceStateInformError      DeviceState = 9
	DeviceStateAdoptFailed      DeviceState = 10
	DeviceStateIsolated         DeviceState = 11
)

func (c *Client) ListDevice(ctx context.Context, site string) ([]Device, error) {
	return c.listDevice(ctx, site)
}

func (c *Client) GetDeviceByMAC(ctx context.Context, site, mac string) (*Device, error) {
	return c.getDevice(ctx, site, mac)
}

func (c *Client) DeleteDevice(ctx context.Context, site, id string) error {
	return c.deleteDevice(ctx, site, id)
}

func (c *Client) CreateDevice(ctx context.Context, site string, d *Device) (*Device, error) {
	return c.createDevice(ctx, site, d)
}

func (c *Client) UpdateDevice(ctx context.Context, site string, d *Device) (*Device, error) {
	return c.updateDevice(ctx, site, d)
}

func (c *Client) GetDevice(ctx context.Context, site, id string) (*Device, error) {
	devices, err := c.ListDevice(ctx, site)

	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		if d.ID == id {
			return &d, nil
		}
	}

	return nil, &NotFoundError{}
}

func (c *Client) AdoptDevice(ctx context.Context, site, mac string) error {
	reqBody := struct {
		Cmd string `json:"cmd"`
		MAC string `json:"mac"`
	}{
		Cmd: "adopt",
		MAC: mac,
	}

	var respBody struct {
		Meta meta `json:"meta"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/cmd/devmgr", site), reqBody, &respBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ForgetDevice(ctx context.Context, site, mac string) error {
	reqBody := struct {
		Cmd  string   `json:"cmd"`
		MACs []string `json:"macs"`
	}{
		Cmd:  "delete-device",
		MACs: []string{mac},
	}

	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Device `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/cmd/sitemgr", site), reqBody, &respBody)
	if err != nil {
		return err
	}

	return nil
}
