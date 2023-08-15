package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

type Setting struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`
	Key    string `json:"key"`
}

func (s *Setting) newFields() (interface{}, error) {
	switch s.Key {
	case "auto_speedtest":
		return &SettingAutoSpeedtest{}, nil
	case "baresip":
		return &SettingBaresip{}, nil
	case "broadcast":
		return &SettingBroadcast{}, nil
	case "connectivity":
		return &SettingConnectivity{}, nil
	case "country":
		return &SettingCountry{}, nil
	case "dpi":
		return &SettingDpi{}, nil
	case "element_adopt":
		return &SettingElementAdopt{}, nil
	case "guest_access":
		return &SettingGuestAccess{}, nil
	// case "ips":
	// 	return &SettingI
	case "lcm":
		return &SettingLcm{}, nil
	case "locale":
		return &SettingLocale{}, nil
	case "mgmt":
		return &SettingMgmt{}, nil
	case "network_optimization":
		return &SettingNetworkOptimization{}, nil
	case "ntp":
		return &SettingNtp{}, nil
	case "porta":
		return &SettingPorta{}, nil
	case "provider_capabilities":
		return &SettingProviderCapabilities{}, nil
	case "radio_ai":
		return &SettingRadioAi{}, nil
	case "radius":
		return &SettingRadius{}, nil
	case "rsyslogd":
		return &SettingRsyslogd{}, nil
	case "snmp":
		return &SettingSnmp{}, nil
	case "super_cloudaccess":
		return &SettingSuperCloudaccess{}, nil
	case "super_events":
		return &SettingSuperEvents{}, nil
	case "super_fwupdate":
		return &SettingSuperFwupdate{}, nil
	case "super_identity":
		return &SettingSuperIdentity{}, nil
	case "super_mail":
		return &SettingSuperMail{}, nil
	case "super_mgmt":
		return &SettingSuperMgmt{}, nil
	case "super_sdn":
		return &SettingSuperSdn{}, nil
	case "super_smtp":
		return &SettingSuperSmtp{}, nil
	case "usg":
		return &SettingUsg{}, nil
	case "usw":
		return &SettingUsw{}, nil
	}

	return nil, fmt.Errorf("unexpected key %q", s.Key)
}

func (c *Client) GetSetting(ctx context.Context, site, key string) (*Setting, interface{}, error) {
	var respBody struct {
		Meta meta              `json:"meta"`
		Data []json.RawMessage `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting", site), nil, &respBody)
	if err != nil {
		return nil, nil, err
	}

	var raw json.RawMessage
	var setting *Setting
	for _, d := range respBody.Data {
		err = json.Unmarshal(d, &setting)
		if err != nil {
			return nil, nil, err
		}
		if setting.Key == key {
			raw = d
			break
		}
	}
	if setting == nil {
		return nil, nil, &NotFoundError{}
	}

	fields, err := setting.newFields()
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(raw, &fields)
	if err != nil {
		return nil, nil, err
	}

	return setting, fields, nil
}
