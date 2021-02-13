package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/paultyng/go-unifi/unifi"
)

type lazyClient struct {
	baseURL  string
	user     string
	pass     string
	insecure bool

	once  sync.Once
	inner *unifi.Client
}

func setHTTPClient(c *unifi.Client, insecure bool) {
	httpClient := &http.Client{}
	httpClient.Transport = &http.Transport{
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

		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
	}

	httpClient.Transport = logging.NewTransport("Unifi", httpClient.Transport)

	jar, _ := cookiejar.New(nil)
	httpClient.Jar = jar

	c.SetHTTPClient(httpClient)
}

var initErr error

func (c *lazyClient) init(ctx context.Context) error {
	c.once.Do(func() {
		c.inner = &unifi.Client{}
		setHTTPClient(c.inner, c.insecure)

		initErr = c.inner.SetBaseURL(c.baseURL)
		if initErr != nil {
			return
		}

		initErr = c.inner.Login(ctx, c.user, c.pass)

		log.Printf("[TRACE] Unifi controller version: %q", c.inner.Version())
	})
	return initErr
}

func (c *lazyClient) Version() string {
	if err := c.init(context.Background()); err != nil {
		panic(fmt.Sprintf("client not initialized: %s", err))
	}
	return c.inner.Version()
}
func (c *lazyClient) ListUserGroup(ctx context.Context, site string) ([]unifi.UserGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListUserGroup(ctx, site)
}
func (c *lazyClient) ListWLANGroup(ctx context.Context, site string) ([]unifi.WLANGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListWLANGroup(ctx, site)
}
func (c *lazyClient) ListAPGroup(ctx context.Context, site string) ([]unifi.APGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListAPGroup(ctx, site)
}
func (c *lazyClient) DeleteNetwork(ctx context.Context, site, id, name string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.DeleteNetwork(ctx, site, id, name)
}
func (c *lazyClient) CreateNetwork(ctx context.Context, site string, d *unifi.Network) (*unifi.Network, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreateNetwork(ctx, site, d)
}
func (c *lazyClient) GetNetwork(ctx context.Context, site, id string) (*unifi.Network, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetNetwork(ctx, site, id)
}
func (c *lazyClient) UpdateNetwork(ctx context.Context, site string, d *unifi.Network) (*unifi.Network, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdateNetwork(ctx, site, d)
}
func (c *lazyClient) DeleteWLAN(ctx context.Context, site, id string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.DeleteWLAN(ctx, site, id)
}
func (c *lazyClient) CreateWLAN(ctx context.Context, site string, d *unifi.WLAN) (*unifi.WLAN, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreateWLAN(ctx, site, d)
}
func (c *lazyClient) GetWLAN(ctx context.Context, site, id string) (*unifi.WLAN, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetWLAN(ctx, site, id)
}
func (c *lazyClient) UpdateWLAN(ctx context.Context, site string, d *unifi.WLAN) (*unifi.WLAN, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdateWLAN(ctx, site, d)
}
func (c *lazyClient) DeleteUserGroup(ctx context.Context, site, id string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.DeleteUserGroup(ctx, site, id)
}
func (c *lazyClient) CreateUserGroup(ctx context.Context, site string, d *unifi.UserGroup) (*unifi.UserGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreateUserGroup(ctx, site, d)
}
func (c *lazyClient) GetUserGroup(ctx context.Context, site, id string) (*unifi.UserGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetUserGroup(ctx, site, id)
}
func (c *lazyClient) UpdateUserGroup(ctx context.Context, site string, d *unifi.UserGroup) (*unifi.UserGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdateUserGroup(ctx, site, d)
}
func (c *lazyClient) GetUser(ctx context.Context, site, id string) (*unifi.User, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetUser(ctx, site, id)
}
func (c *lazyClient) GetUserByMAC(ctx context.Context, site, mac string) (*unifi.User, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetUserByMAC(ctx, site, mac)
}
func (c *lazyClient) CreateUser(ctx context.Context, site string, d *unifi.User) (*unifi.User, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreateUser(ctx, site, d)
}
func (c *lazyClient) UpdateUser(ctx context.Context, site string, d *unifi.User) (*unifi.User, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdateUser(ctx, site, d)
}
func (c *lazyClient) DeleteUserByMAC(ctx context.Context, site, mac string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.DeleteUserByMAC(ctx, site, mac)
}
func (c *lazyClient) BlockUserByMAC(ctx context.Context, site, mac string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.BlockUserByMAC(ctx, site, mac)
}
func (c *lazyClient) UnblockUserByMAC(ctx context.Context, site, mac string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.UnblockUserByMAC(ctx, site, mac)
}
func (c *lazyClient) ListFirewallGroup(ctx context.Context, site string) ([]unifi.FirewallGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListFirewallGroup(ctx, site)
}
func (c *lazyClient) DeleteFirewallGroup(ctx context.Context, site, id string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.DeleteFirewallGroup(ctx, site, id)
}
func (c *lazyClient) CreateFirewallGroup(ctx context.Context, site string, d *unifi.FirewallGroup) (*unifi.FirewallGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreateFirewallGroup(ctx, site, d)
}
func (c *lazyClient) GetFirewallGroup(ctx context.Context, site, id string) (*unifi.FirewallGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetFirewallGroup(ctx, site, id)
}
func (c *lazyClient) UpdateFirewallGroup(ctx context.Context, site string, d *unifi.FirewallGroup) (*unifi.FirewallGroup, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdateFirewallGroup(ctx, site, d)
}
func (c *lazyClient) ListFirewallRule(ctx context.Context, site string) ([]unifi.FirewallRule, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListFirewallRule(ctx, site)
}
func (c *lazyClient) DeleteFirewallRule(ctx context.Context, site, id string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.DeleteFirewallRule(ctx, site, id)
}
func (c *lazyClient) CreateFirewallRule(ctx context.Context, site string, d *unifi.FirewallRule) (*unifi.FirewallRule, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreateFirewallRule(ctx, site, d)
}
func (c *lazyClient) GetFirewallRule(ctx context.Context, site, id string) (*unifi.FirewallRule, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetFirewallRule(ctx, site, id)
}
func (c *lazyClient) UpdateFirewallRule(ctx context.Context, site string, d *unifi.FirewallRule) (*unifi.FirewallRule, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdateFirewallRule(ctx, site, d)
}
func (c *lazyClient) GetPortForward(ctx context.Context, site, id string) (*unifi.PortForward, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetPortForward(ctx, site, id)
}
func (c *lazyClient) DeletePortForward(ctx context.Context, site, id string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.DeletePortForward(ctx, site, id)
}
func (c *lazyClient) CreatePortForward(ctx context.Context, site string, d *unifi.PortForward) (*unifi.PortForward, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreatePortForward(ctx, site, d)
}
func (c *lazyClient) UpdatePortForward(ctx context.Context, site string, d *unifi.PortForward) (*unifi.PortForward, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdatePortForward(ctx, site, d)
}
func (c *lazyClient) ListRADIUSProfile(ctx context.Context, site string) ([]unifi.RADIUSProfile, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListRADIUSProfile(ctx, site)
}
func (c *lazyClient) GetRADIUSProfile(ctx context.Context, site, id string) (*unifi.RADIUSProfile, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetRADIUSProfile(ctx, site, id)
}
func (c *lazyClient) DeleteRADIUSProfile(ctx context.Context, site, id string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.DeleteRADIUSProfile(ctx, site, id)
}
func (c *lazyClient) CreateRADIUSProfile(ctx context.Context, site string, d *unifi.RADIUSProfile) (*unifi.RADIUSProfile, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreateRADIUSProfile(ctx, site, d)
}
func (c *lazyClient) UpdateRADIUSProfile(ctx context.Context, site string, d *unifi.RADIUSProfile) (*unifi.RADIUSProfile, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdateRADIUSProfile(ctx, site, d)
}
func (c *lazyClient) GetSite(ctx context.Context, id string) (*unifi.Site, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetSite(ctx, id)
}
func (c *lazyClient) ListSites(ctx context.Context) ([]unifi.Site, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListSites(ctx)
}
func (c *lazyClient) CreateSite(ctx context.Context, description string) ([]unifi.Site, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreateSite(ctx, description)
}
func (c *lazyClient) DeleteSite(ctx context.Context, id string) ([]unifi.Site, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.DeleteSite(ctx, id)
}
func (c *lazyClient) UpdateSite(ctx context.Context, name, description string) ([]unifi.Site, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdateSite(ctx, name, description)
}

func (c *lazyClient) ListPortProfile(ctx context.Context, site string) ([]unifi.PortProfile, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListPortProfile(ctx, site)
}

func (c *lazyClient) GetPortProfile(ctx context.Context, site, id string) (*unifi.PortProfile, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetPortProfile(ctx, site, id)
}

func (c *lazyClient) DeletePortProfile(ctx context.Context, site, id string) error {
	if err := c.init(ctx); err != nil {
		return err
	}
	return c.inner.DeletePortProfile(ctx, site, id)
}

func (c *lazyClient) CreatePortProfile(ctx context.Context, site string, d *unifi.PortProfile) (*unifi.PortProfile, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.CreatePortProfile(ctx, site, d)
}

func (c *lazyClient) UpdatePortProfile(ctx context.Context, site string, d *unifi.PortProfile) (*unifi.PortProfile, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.inner.UpdatePortProfile(ctx, site, d)
}
