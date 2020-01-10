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
		err = c.inner.SetBaseURL(c.baseURL)
		if err != nil {
			return
		}

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
func (c *lazyClient) UpdateNetwork(site string, d *unifi.Network) (*unifi.Network, error) {
	c.init()
	return c.inner.UpdateNetwork(site, d)
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
func (c *lazyClient) UpdateWLAN(site string, d *unifi.WLAN) (*unifi.WLAN, error) {
	c.init()
	return c.inner.UpdateWLAN(site, d)
}
func (c *lazyClient) DeleteUserGroup(site, id string) error {
	c.init()
	return c.inner.DeleteUserGroup(site, id)
}
func (c *lazyClient) CreateUserGroup(site string, d *unifi.UserGroup) (*unifi.UserGroup, error) {
	c.init()
	return c.inner.CreateUserGroup(site, d)
}
func (c *lazyClient) GetUserGroup(site, id string) (*unifi.UserGroup, error) {
	c.init()
	return c.inner.GetUserGroup(site, id)
}
func (c *lazyClient) UpdateUserGroup(site string, d *unifi.UserGroup) (*unifi.UserGroup, error) {
	c.init()
	return c.inner.UpdateUserGroup(site, d)
}
func (c *lazyClient) GetUser(site, id string) (*unifi.User, error) {
	c.init()
	return c.inner.GetUser(site, id)
}
func (c *lazyClient) GetUserByMAC(site, mac string) (*unifi.User, error) {
	c.init()
	return c.inner.GetUserByMAC(site, mac)
}
func (c *lazyClient) CreateUser(site string, d *unifi.User) (*unifi.User, error) {
	c.init()
	return c.inner.CreateUser(site, d)
}
func (c *lazyClient) UpdateUser(site string, d *unifi.User) (*unifi.User, error) {
	c.init()
	return c.inner.UpdateUser(site, d)
}
func (c *lazyClient) DeleteUserByMAC(site, mac string) error {
	c.init()
	return c.inner.DeleteUserByMAC(site, mac)
}
func (c *lazyClient) BlockUserByMAC(site, mac string) error {
	c.init()
	return c.inner.BlockUserByMAC(site, mac)
}
func (c *lazyClient) UnblockUserByMAC(site, mac string) error {
	c.init()
	return c.inner.UnblockUserByMAC(site, mac)
}
