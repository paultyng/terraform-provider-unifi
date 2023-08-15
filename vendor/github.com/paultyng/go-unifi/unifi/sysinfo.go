package unifi

import (
	"context"
	"fmt"
)

type sysInfo struct {
	Timezone        string `json:"timezone"`
	Version         string `json:"version"`
	PreviousVersion string `json:"previous_version"`
	UBNTDeviceType  string `json:"ubnt_device_type"`
	UDMVersion      string `json:"udm_version"`

	/*

	   {
	       "meta": {
	           "rc": "ok"
	       },
	       "data": [
	           {
	               "timezone": "America/New_York",
	               "autobackup": false,
	               "build": "atag_6.0.43_14348",
	               "version": "6.0.43",
	               "previous_version": "5.12.60",
	               "debug_mgmt": "warn",
	               "debug_system": "warn",
	               "debug_device": "warn",
	               "debug_sdn": "warn",
	               "data_retention_days": 90,
	               "data_retention_time_in_hours_for_5minutes_scale": 24,
	               "data_retention_time_in_hours_for_hourly_scale": 720,
	               "data_retention_time_in_hours_for_daily_scale": 2160,
	               "data_retention_time_in_hours_for_monthly_scale": 8760,
	               "data_retention_time_in_hours_for_others": 2160,
	               "update_available": false,
	               "update_downloaded": false,
	               "live_chat": "super-only",
	               "store_enabled": "super-only",
	               "hostname": "example-domain.ui.com",
	               "name": "Dream Machine",
	               "ip_addrs": [
	                   "1.2.3.4"
	               ],
	               "inform_port": 8080,
	               "https_port": 8443,
	               "override_inform_host": false,
	               "image_maps_use_google_engine": false,
	               "radius_disconnect_running": false,
	               "facebook_wifi_registered": false,
	               "sso_app_id": "",
	               "sso_app_sec": "",
	               "uptime": 2541796,
	               "anonymous_controller_id": "",
	               "ubnt_device_type": "UDMB",
	               "udm_version": "1.8.6.2969",
	               "unsupported_device_count": 0,
	               "unsupported_device_list": [],
	               "unifi_go_enabled": false
	           }
	       ]
	   }

	*/
}

func (c *Client) sysinfo(ctx context.Context, id string) (*sysInfo, error) {
	var respBody struct {
		Meta meta      `json:"meta"`
		Data []sysInfo `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/stat/sysinfo", id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	return &respBody.Data[0], nil
}
