package unifi

import "fmt"

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
