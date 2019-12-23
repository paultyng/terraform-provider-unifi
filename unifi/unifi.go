package unifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	// "io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/hashicorp/terraform/helper/logging"
)

type NotFoundError struct{}

func (err *NotFoundError) Error() string {
	return "not found"
}

type Client struct {
	c       *http.Client
	baseURL *url.URL
}

func (c *Client) SetBaseURL(base string) {
	c.baseURL, _ = url.Parse(base)
}

func (c *Client) Login(user, pass string) error {
	if c.c == nil {
		c.c = &http.Client{}
		c.c.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,

			// TODO: make this opt-in
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		c.c.Transport = logging.NewTransport("Unifi", c.c.Transport)

		jar, _ := cookiejar.New(nil)
		c.c.Jar = jar
	}

	err := c.do("POST", "login", &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: user,
		Password: pass,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) do(method, relativeURL string, reqBody interface{}, respBody interface{}) error {
	var reqReader io.Reader
	if reqBody != nil {
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}
		reqReader = bytes.NewReader(reqBytes)
	}

	reqURL, _ := url.Parse(relativeURL)
	url := c.baseURL.ResolveReference(reqURL)

	req, err := http.NewRequest(method, url.String(), reqReader)
	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{}
	}

	if resp.StatusCode != 200 {
		// body, _ := ioutil.ReadAll(resp.Body)
		//TODO: debug logging?
		// fmt.Printf("%s %s\nStatus: %s\n%s", method, url.String(), resp.Status, string(body))

		return fmt.Errorf("error from API %s", resp.Status)
	}

	if respBody == nil || resp.ContentLength == 0 {
		return nil
	}

	// TODO: check rc?

	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		return err
	}

	return nil
}
