package unifi

import (
	"fmt"
)

func (c *Client) ListUserGroup(site string) ([]UserGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []UserGroup `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/usergroup", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) GetUserGroup(site, id string) (*UserGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []UserGroup `json:"data"`
	}

	err := c.do("GET", fmt.Sprintf("s/%s/rest/usergroup/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) DeleteUserGroup(site, id string) error {
	err := c.do("DELETE", fmt.Sprintf("s/%s/rest/usergroup/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateUserGroup(site string, d *UserGroup) (*UserGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []UserGroup `json:"data"`
	}

	err := c.do("POST", fmt.Sprintf("s/%s/rest/usergroup", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) UpdateUserGroup(site string, d *UserGroup) (*UserGroup, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []UserGroup `json:"data"`
	}

	err := c.do("PUT", fmt.Sprintf("s/%s/rest/usergroup/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
