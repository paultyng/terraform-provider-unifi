package provider

import (
	"sync"

	"github.com/paultyng/terraform-provider-unifi/unifi"
)

type lazyClient struct {
	baseURL string
	user    string
	pass    string

	once  sync.Once
	inner *unifi.Client
}

func (c *lazyClient) init() error {
	var err error
	c.once.Do(func() {
		c.inner = &unifi.Client{}
		c.inner.SetBaseURL(c.baseURL)

		err = c.inner.Login(c.user, c.pass)
	})
	return err
}

func (c *lazyClient) ListUserGroup(site string) ([]unifi.UserGroup, error) {
	c.init()
	return c.inner.ListUserGroup(site)
}
func (c *lazyClient) ListWLANGroup(site string) ([]unifi.WLANGroup, error) {
	c.init()
	return c.inner.ListWLANGroup(site)
}
func (c *lazyClient) DeleteNetwork(site, id, name string) error {
	c.init()
	return c.inner.DeleteNetwork(site, id, name)
}
func (c *lazyClient) CreateNetwork(site string, d *unifi.Network) (*unifi.Network, error) {
	c.init()
	return c.inner.CreateNetwork(site, d)
}
func (c *lazyClient) GetNetwork(site, id string) (*unifi.Network, error) {
	c.init()
	return c.inner.GetNetwork(site, id)
}
func (c *lazyClient) DeleteWLAN(site, id string) error {
	c.init()
	return c.inner.DeleteWLAN(site, id)
}
func (c *lazyClient) CreateWLAN(site string, d *unifi.WLAN) (*unifi.WLAN, error) {
	c.init()
	return c.inner.CreateWLAN(site, d)
}
func (c *lazyClient) GetWLAN(site, id string) (*unifi.WLAN, error) {
	c.init()
	return c.inner.GetWLAN(site, id)
}
