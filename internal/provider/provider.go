package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown

	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Deprecated != "" {
			desc += " " + s.Deprecated
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"username": {
					Description: "Local user name for the Unifi controller API. Can be specified with the `UNIFI_USERNAME` " +
						"environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("UNIFI_USERNAME", ""),
				},
				"password": {
					Description: "Password for the user accessing the API. Can be specified with the `UNIFI_PASSWORD` " +
						"environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("UNIFI_PASSWORD", ""),
				},
				"api_url": {
					Description: "URL of the controller API. Can be specified with the `UNIFI_API` environment variable. " +
						"You should **NOT** supply the path (`/api`), the SDK will discover the appropriate paths. This is " +
						"to support UDM Pro style API paths as well as more standard controller paths.",

					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("UNIFI_API", ""),
				},
				"site": {
					Description: "The site in the Unifi controller this provider will manage. Can be specified with " +
						"the `UNIFI_SITE` environment variable. Default: `default`",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("UNIFI_SITE", "default"),
				},
				"allow_insecure": {
					Description: "Skip verification of TLS certificates of API requests. You may need to set this to `true` " +
						"if you are using your local API without setting up a signed certificate. Can be specified with the " +
						"`UNIFI_INSECURE` environment variable.",
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("UNIFI_INSECURE", false),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"unifi_ap_group":       dataAPGroup(),
				"unifi_network":        dataNetwork(),
				"unifi_port_profile":   dataPortProfile(),
				"unifi_radius_profile": dataRADIUSProfile(),
				"unifi_user":           dataUser(),
				"unifi_user_group":     dataUserGroup(),
				"unifi_wlan_group":     dataWLANGroup(),
			},
			ResourcesMap: map[string]*schema.Resource{
				// TODO: "unifi_ap_group"
				"unifi_device":         resourceDevice(),
				"unifi_dynamic_dns":    resourceDynamicDNS(),
				"unifi_firewall_group": resourceFirewallGroup(),
				"unifi_firewall_rule":  resourceFirewallRule(),
				"unifi_network":        resourceNetwork(),
				"unifi_port_forward":   resourcePortForward(),
				"unifi_port_profile":   resourcePortProfile(),
				"unifi_site":           resourceSite(),
				"unifi_static_route":   resourceStaticRoute(),
				"unifi_user_group":     resourceUserGroup(),
				"unifi_user":           resourceUser(),
				"unifi_wlan":           resourceWLAN(),

				"unifi_setting_mgmt": resourceSettingMgmt(),
				"unifi_setting_usg":  resourceSettingUsg(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

func configure(version string, p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
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
	Version() string

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

	ListAPGroup(ctx context.Context, site string) ([]unifi.APGroup, error)

	DeleteNetwork(ctx context.Context, site, id, name string) error
	CreateNetwork(ctx context.Context, site string, d *unifi.Network) (*unifi.Network, error)
	GetNetwork(ctx context.Context, site, id string) (*unifi.Network, error)
	ListNetwork(ctx context.Context, site string) ([]unifi.Network, error)
	UpdateNetwork(ctx context.Context, site string, d *unifi.Network) (*unifi.Network, error)

	DeleteWLAN(ctx context.Context, site, id string) error
	CreateWLAN(ctx context.Context, site string, d *unifi.WLAN) (*unifi.WLAN, error)
	GetWLAN(ctx context.Context, site, id string) (*unifi.WLAN, error)
	UpdateWLAN(ctx context.Context, site string, d *unifi.WLAN) (*unifi.WLAN, error)

	GetDevice(ctx context.Context, site, id string) (*unifi.Device, error)
	CreateDevice(ctx context.Context, site string, d *unifi.Device) (*unifi.Device, error)
	UpdateDevice(ctx context.Context, site string, d *unifi.Device) (*unifi.Device, error)
	DeleteDevice(ctx context.Context, site, id string) error
	ListDevice(ctx context.Context, site string) ([]unifi.Device, error)

	GetUser(ctx context.Context, site, id string) (*unifi.User, error)
	GetUserByMAC(ctx context.Context, site, mac string) (*unifi.User, error)
	CreateUser(ctx context.Context, site string, d *unifi.User) (*unifi.User, error)
	BlockUserByMAC(ctx context.Context, site, mac string) error
	UnblockUserByMAC(ctx context.Context, site, mac string) error
	OverrideUserFingerprint(ctx context.Context, site, mac string, devIdOveride int) error
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

	GetSite(ctx context.Context, id string) (*unifi.Site, error)
	ListSites(ctx context.Context) ([]unifi.Site, error)
	CreateSite(ctx context.Context, Description string) ([]unifi.Site, error)
	UpdateSite(ctx context.Context, Name, Description string) ([]unifi.Site, error)
	DeleteSite(ctx context.Context, ID string) ([]unifi.Site, error)

	ListPortProfile(ctx context.Context, site string) ([]unifi.PortProfile, error)
	GetPortProfile(ctx context.Context, site, id string) (*unifi.PortProfile, error)
	DeletePortProfile(ctx context.Context, site, id string) error
	CreatePortProfile(ctx context.Context, site string, d *unifi.PortProfile) (*unifi.PortProfile, error)
	UpdatePortProfile(ctx context.Context, site string, d *unifi.PortProfile) (*unifi.PortProfile, error)

	ListRouting(ctx context.Context, site string) ([]unifi.Routing, error)
	GetRouting(ctx context.Context, site, id string) (*unifi.Routing, error)
	DeleteRouting(ctx context.Context, site, id string) error
	CreateRouting(ctx context.Context, site string, d *unifi.Routing) (*unifi.Routing, error)
	UpdateRouting(ctx context.Context, site string, d *unifi.Routing) (*unifi.Routing, error)

	ListDynamicDNS(ctx context.Context, site string) ([]unifi.DynamicDNS, error)
	GetDynamicDNS(ctx context.Context, site, id string) (*unifi.DynamicDNS, error)
	DeleteDynamicDNS(ctx context.Context, site, id string) error
	CreateDynamicDNS(ctx context.Context, site string, d *unifi.DynamicDNS) (*unifi.DynamicDNS, error)
	UpdateDynamicDNS(ctx context.Context, site string, d *unifi.DynamicDNS) (*unifi.DynamicDNS, error)

	GetSettingMgmt(ctx context.Context, id string) (*unifi.SettingMgmt, error)
	GetSettingUsg(ctx context.Context, id string) (*unifi.SettingUsg, error)
	UpdateSettingMgmt(ctx context.Context, site string, d *unifi.SettingMgmt) (*unifi.SettingMgmt, error)
	UpdateSettingUsg(ctx context.Context, site string, d *unifi.SettingUsg) (*unifi.SettingUsg, error)
}

type client struct {
	c    unifiClient
	site string
}
