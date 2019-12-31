package provider

import (
	// "fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/paultyng/terraform-provider-unifi/unifi"
)

func Provider() terraform.ResourceProvider {
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
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_API", ""),
			},
			"site": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_SITE", "default"),
			},
			// "allow_insecure": {
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// 	Default:  false,
			// },
		},
		DataSourcesMap: map[string]*schema.Resource{
			"unifi_user_group": dataUserGroup(),
			"unifi_wlan_group": dataWLANGroup(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"unifi_network":    resourceNetwork(),
			"unifi_user_group": resourceUserGroup(),
			"unifi_user":       resourceUser(),
			"unifi_wlan":       resourceWLAN(),
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
		//insecure := d.Get("allow_insecure").(bool)

		c := &client{
			c: &lazyClient{
				user:    user,
				pass:    pass,
				baseURL: baseURL,
			},
			site: site,
		}

		return c, nil
	}
}

type unifiClient interface {
	ListUserGroup(site string) ([]unifi.UserGroup, error)
	DeleteUserGroup(site, id string) error
	CreateUserGroup(site string, d *unifi.UserGroup) (*unifi.UserGroup, error)
	GetUserGroup(site, id string) (*unifi.UserGroup, error)
	UpdateUserGroup(site string, d *unifi.UserGroup) (*unifi.UserGroup, error)

	ListWLANGroup(site string) ([]unifi.WLANGroup, error)

	DeleteNetwork(site, id, name string) error
	CreateNetwork(site string, d *unifi.Network) (*unifi.Network, error)
	GetNetwork(site, id string) (*unifi.Network, error)
	UpdateNetwork(site string, d *unifi.Network) (*unifi.Network, error)

	DeleteWLAN(site, id string) error
	CreateWLAN(site string, d *unifi.WLAN) (*unifi.WLAN, error)
	GetWLAN(site, id string) (*unifi.WLAN, error)

	GetUser(site, id string) (*unifi.User, error)
	GetUserByMAC(site, mac string) (*unifi.User, error)
	CreateUser(site string, d *unifi.User) (*unifi.User, error)
	BlockUserByMAC(site, mac string) error
	UnblockUserByMAC(site, mac string) error
	UpdateUser(site string, d *unifi.User) (*unifi.User, error)
	DeleteUserByMAC(site, mac string) error
}

type client struct {
	c    unifiClient
	site string
}
