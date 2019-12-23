package unifi

type Site struct {
	ID       string `json:"_id,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`

	Name        string `json:"name"`
	Description string `json:"desc"`

	//Role string `json:"role"`
}

func (c *Client) ListSites() ([]Site, error) {
	var respBody struct {
		Meta struct {
			RC string `json:"rc"`
		} `json:"meta"`
		Data []Site `json:"data"`
	}

	err := c.do("GET", "self/sites", nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}
