package unifi

import "fmt"

// GET https://73.212.25.176:8443/api/s/default/stat/user/e4:f0:42:bf:bd:11
// {"meta":{"rc":"ok"},"data":[{"_id":"5deeac76439adf048407dd01","mac":"e4:f0:42:bf:bd:11","site_id":"5d6d8b07439adf048407dcd9","oui":"Google","is_guest":false,"first_seen":1575922805,"last_seen":1577907316,"is_wired":true,"fingerprint_engine":"tdts","dev_cat":6,"dev_family":9,"os_class":16,"os_name":56,"dev_vendor":7,"dev_id":2822,"priority":101,"fingerprint_source":0,"name":"Google WiFi","usergroup_id":"","noted":true,"fixed_ip":"10.0.6.11","note":"","confidence":100,"network_id":"5d6d8b0c439adf048407dce7","use_fixedip":true,"duration":1966599,"tx_bytes":2391811124,"tx_packets":2227771,"rx_bytes":1431694370,"rx_packets":2606611,"wifi_tx_attempts":0,"tx_retries":0,"assoc_time":1577889972,"latest_assoc_time":1577889972,"user_id":"5deeac76439adf048407dd01","_uptime_by_ugw":1984511,"_last_seen_by_ugw":1577907316,"_is_guest_by_ugw":false,"gw_mac":"74:83:c2:d6:ff:83","network":"LAN","ip":"192.168.1.135","uptime":17344,"tx_bytes-r":0,"rx_bytes-r":0,"authorized":true,"qos_policy_applied":true,"_uptime_by_usw":1984511,"_last_seen_by_usw":1577907316,"_is_guest_by_usw":false,"sw_mac":"74:83:c2:d6:ff:83","sw_depth":0,"sw_port":4,"wired-tx_bytes":2176183895,"wired-rx_bytes":1406877963,"wired-tx_packets":2156208,"wired-rx_packets":2444036,"wired-tx_bytes-r":1968908,"wired-rx_bytes-r":105359}]}

// PUT https://73.212.25.176:8443/api/s/default/rest/user/5deeac76439adf048407dd01
// { use_fixedip: true, network_id: "5df7f70f1e801c052a1ab032", fixed_ip: "10.0.6.11" }
// { note: "my note", usergroup_id: "", name: "Google WiFi alias"}
// {"meta":{"rc":"ok"},"data":[{"_id":"5deeac76439adf048407dd01","mac":"e4:f0:42:bf:bd:11","site_id":"5d6d8b07439adf048407dcd9","oui":"Google","is_guest":false,"first_seen":1575922805,"last_seen":1577889639,"is_wired":true,"fingerprint_engine":"tdts","dev_cat":6,"dev_family":9,"os_class":16,"os_name":56,"dev_vendor":7,"dev_id":2822,"priority":101,"fingerprint_source":0,"name":"Google WiFi","usergroup_id":"","noted":true,"fixed_ip":"10.0.6.11","note":"","confidence":100,"network_id":"5df7f70f1e801c052a1ab032","use_fixedip":true}]}

type User struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`
	Name   string `json:"name"`
	MAC    string `json:"mac"`

	UserGroupID string `json:"user_group_id"`
	Note        string `json:"note"`
	UseFixedIP  bool   `json:"use_fixedip"`
	FixedIP     string `json:"fixed_ip,omitempty"`
	NetworkID   string `json:"network_id"`

	// not sure if you can end this for create/update, etc, only
	// observed modifying via stamgr
	Blocked bool `json:"blocked,omitempty"`
}

func (c *Client) ListUser(site string) ([]User, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/user", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) GetUser(site, id string) (*User, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/user/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) GetUserByMAC(site, mac string) (*User, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/stat/user/%s", site, mac), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) CreateUser(site string, d *User) (*User, error) {
	reqBody := struct {
		Objects []struct {
			Data *User `json:"data"`
		} `json:"objects"`
	}{
		Objects: []struct {
			Data *User `json:"data"`
		}{
			{Data: d},
		},
	}

	var respBody struct {
		Meta meta `json:"meta"`
		Data []struct {
			Meta meta   `json:"meta"`
			Data []User `json:"data"`
		} `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/group/user", site), reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, fmt.Errorf("malformed group response")
	}

	if err := respBody.Data[0].Meta.error(); err != nil {
		return nil, err
	}

	if len(respBody.Data[0].Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0].Data[0]

	return &new, nil
}

func (c *Client) stamgr(site, cmd string, data map[string]interface{}) ([]User, error) {
	reqBody := map[string]interface{}{}

	for k, v := range data {
		reqBody[k] = v
	}

	reqBody["cmd"] = cmd

	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/cmd/stamgr", site), reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) BlockUserByMAC(site, mac string) error {
	users, err := c.stamgr(site, "block-sta", map[string]interface{}{
		"mac": mac,
	})
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return &NotFoundError{}
	}
	return nil
}

func (c *Client) UnblockUserByMAC(site, mac string) error {
	users, err := c.stamgr(site, "unblock-sta", map[string]interface{}{
		"mac": mac,
	})
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return &NotFoundError{}
	}
	return nil
}

func (c *Client) UpdateUser(site string, d *User) (*User, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/user/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) DeleteUserByMAC(site, mac string) error {
	users, err := c.stamgr(site, "forget-sta", map[string]interface{}{
		"macs": []string{mac},
	})
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return &NotFoundError{}
	}
	return nil
}
