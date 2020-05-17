package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_PASSWORD", ""),
			},
			"api_url": {
				Description: "URL of the controller API.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_API", ""),
			},
			"site": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_SITE", "default"),
			},
			"allow_insecure": {
				Description: "Skip verification of TLS certificates of API requests. " +
					"You may need to set this to `true` if you are using your local API without " +
					"setting up a signed certificate.",
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_INSECURE", false),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"unifi_radius_profile": dataRADIUSProfile(),
			"unifi_user_group":     dataUserGroup(),
			"unifi_wlan_group":     dataWLANGroup(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"unifi_firewall_group": resourceFirewallGroup(),
			"unifi_firewall_rule":  resourceFirewallRule(),
			"unifi_network":        resourceNetwork(),
			"unifi_port_forward":   resourcePortForward(),
			"unifi_user_group":     resourceUserGroup(),
			"unifi_user":           resourceUser(),
			"unifi_wlan":           resourceWLAN(),
		},
	}
	p.ConfigureFunc = configure(p)
	return p
}

func configure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		user := d.Get("username").(string)
		pass := d.Get("password").(string)
		baseURL := d.Get("api_url").(string)
		site := d.Get("site").(string)
		insecure := d.Get("allow_insecure").(bool)

		c := &client{
			c: &lazyClient{
				user:     user,
				pass:     pass,
				baseURL:  baseURL,
				insecure: insecure,
			},
			site: site,
		}

		return c, nil
	}
}

type unifiClient interface {
	ListUserGroup(ctx context.Context, site string) ([]unifi.UserGroup, error)
	DeleteUserGroup(ctx context.Context, site, id string) error
	CreateUserGroup(ctx context.Context, site string, d *unifi.UserGroup) (*unifi.UserGroup, error)
	GetUserGroup(ctx context.Context, site, id string) (*unifi.UserGroup, error)
	UpdateUserGroup(ctx context.Context, site string, d *unifi.UserGroup) (*unifi.UserGroup, error)

	ListFirewallGroup(ctx context.Context, site string) ([]unifi.FirewallGroup, error)
	DeleteFirewallGroup(ctx context.Context, site, id string) error
	CreateFirewallGroup(ctx context.Context, site string, d *unifi.FirewallGroup) (*unifi.FirewallGroup, error)
	GetFirewallGroup(ctx context.Context, site, id string) (*unifi.FirewallGroup, error)
	UpdateFirewallGroup(ctx context.Context, site string, d *unifi.FirewallGroup) (*unifi.FirewallGroup, error)

	ListFirewallRule(ctx context.Context, site string) ([]unifi.FirewallRule, error)
	DeleteFirewallRule(ctx context.Context, site, id string) error
	CreateFirewallRule(ctx context.Context, site string, d *unifi.FirewallRule) (*unifi.FirewallRule, error)
	GetFirewallRule(ctx context.Context, site, id string) (*unifi.FirewallRule, error)
	UpdateFirewallRule(ctx context.Context, site string, d *unifi.FirewallRule) (*unifi.FirewallRule, error)

	ListWLANGroup(ctx context.Context, site string) ([]unifi.WLANGroup, error)

	DeleteNetwork(ctx context.Context, site, id, name string) error
	CreateNetwork(ctx context.Context, site string, d *unifi.Network) (*unifi.Network, error)
	GetNetwork(ctx context.Context, site, id string) (*unifi.Network, error)
	UpdateNetwork(ctx context.Context, site string, d *unifi.Network) (*unifi.Network, error)

	DeleteWLAN(ctx context.Context, site, id string) error
	CreateWLAN(ctx context.Context, site string, d *unifi.WLAN) (*unifi.WLAN, error)
	GetWLAN(ctx context.Context, site, id string) (*unifi.WLAN, error)
	UpdateWLAN(ctx context.Context, site string, d *unifi.WLAN) (*unifi.WLAN, error)

	GetUser(ctx context.Context, site, id string) (*unifi.User, error)
	GetUserByMAC(ctx context.Context, site, mac string) (*unifi.User, error)
	CreateUser(ctx context.Context, site string, d *unifi.User) (*unifi.User, error)
	BlockUserByMAC(ctx context.Context, site, mac string) error
	UnblockUserByMAC(ctx context.Context, site, mac string) error
	UpdateUser(ctx context.Context, site string, d *unifi.User) (*unifi.User, error)
	DeleteUserByMAC(ctx context.Context, site, mac string) error

	GetPortForward(ctx context.Context, site, id string) (*unifi.PortForward, error)
	DeletePortForward(ctx context.Context, site, id string) error
	CreatePortForward(ctx context.Context, site string, d *unifi.PortForward) (*unifi.PortForward, error)
	UpdatePortForward(ctx context.Context, site string, d *unifi.PortForward) (*unifi.PortForward, error)

	ListRADIUSProfile(ctx context.Context, site string) ([]unifi.RADIUSProfile, error)
	GetRADIUSProfile(ctx context.Context, site, id string) (*unifi.RADIUSProfile, error)
	DeleteRADIUSProfile(ctx context.Context, site, id string) error
	CreateRADIUSProfile(ctx context.Context, site string, d *unifi.RADIUSProfile) (*unifi.RADIUSProfile, error)
	UpdateRADIUSProfile(ctx context.Context, site string, d *unifi.RADIUSProfile) (*unifi.RADIUSProfile, error)
}

type client struct {
	c    unifiClient
	site string
}
