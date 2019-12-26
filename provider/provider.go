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
			// "scaffolding_data_source": dataSourceScaffolding(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"unifi_network": resourceNetwork(),
			"unifi_wlan":    resourceWLAN(),
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
			c:    &unifi.Client{},
			site: site,
		}

		c.c.SetBaseURL(baseURL)

		// TODO: defer this to first call?
		err := c.c.Login(user, pass)
		if err != nil {
			return nil, err
		}

		return c, nil
	}
}

type client struct {
	c    *unifi.Client
	site string
}
