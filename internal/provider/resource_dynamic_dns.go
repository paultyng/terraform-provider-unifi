package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceDynamicDNS() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_dynamic_dns` manages dynamic DNS settings for different providers.",

		Create: resourceDynamicDNSCreate,
		Read:   resourceDynamicDNSRead,
		Update: resourceDynamicDNSUpdate,
		Delete: resourceDynamicDNSDelete,
		Importer: &schema.ResourceImporter{
			State: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the dynamic DNS.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the dynamic DNS with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"interface": {
				Description: "The interface for the dynamic DNS. Can be `wan` or `wan2`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "wan",
				ForceNew:    true,
			},
			"service": {
				Description: "The Dynamic DNS service provider, various values are supported (for example `dyndns`, etc.).",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"host_name": {
				Description: "The host name to update in the dynamic DNS service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"server": {
				Description: "The server for the dynamic DNS service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"login": {
				Description: "The server for the dynamic DNS service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "The server for the dynamic DNS service.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},

			//TODO: options support?
		},
	}
}

func resourceDynamicDNSCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceDynamicDNSGetResourceData(d)
	if err != nil {
		return err
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateDynamicDNS(context.TODO(), site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceDynamicDNSSetResourceData(resp, d, site)
}

func resourceDynamicDNSGetResourceData(d *schema.ResourceData) (*unifi.DynamicDNS, error) {
	r := &unifi.DynamicDNS{
		Interface: d.Get("interface").(string),
		Service:   d.Get("service").(string),

		HostName: d.Get("host_name").(string),

		Server:    d.Get("server").(string),
		Login:     d.Get("login").(string),
		XPassword: d.Get("password").(string),
	}

	return r, nil
}

func resourceDynamicDNSSetResourceData(resp *unifi.DynamicDNS, d *schema.ResourceData, site string) error {
	d.Set("interface", resp.Interface)
	d.Set("service", resp.Service)

	d.Set("host_name", resp.HostName)

	d.Set("server", resp.Server)
	d.Set("login", resp.Login)
	d.Set("password", resp.XPassword)

	return nil
}

func resourceDynamicDNSRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetDynamicDNS(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceDynamicDNSSetResourceData(resp, d, site)
}

func resourceDynamicDNSUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceDynamicDNSGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateDynamicDNS(context.TODO(), site, req)
	if err != nil {
		return err
	}

	return resourceDynamicDNSSetResourceData(resp, d, site)
}

func resourceDynamicDNSDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	err := c.c.DeleteDynamicDNS(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
