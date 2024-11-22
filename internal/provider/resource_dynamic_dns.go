package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceDynamicDNS() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_dynamic_dns` manages dynamic DNS settings for different providers.",

		CreateContext: resourceDynamicDNSCreate,
		ReadContext:   resourceDynamicDNSRead,
		UpdateContext: resourceDynamicDNSUpdate,
		DeleteContext: resourceDynamicDNSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
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

			// TODO: options support?
		},
	}
}

func resourceDynamicDNSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceDynamicDNSGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateDynamicDNS(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
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

func resourceDynamicDNSSetResourceData(resp *unifi.DynamicDNS, d *schema.ResourceData, site string) diag.Diagnostics {
	d.Set("interface", resp.Interface)
	d.Set("service", resp.Service)

	d.Set("host_name", resp.HostName)

	d.Set("server", resp.Server)
	d.Set("login", resp.Login)
	d.Set("password", resp.XPassword)

	return nil
}

func resourceDynamicDNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetDynamicDNS(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDynamicDNSSetResourceData(resp, d, site)
}

func resourceDynamicDNSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceDynamicDNSGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateDynamicDNS(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDynamicDNSSetResourceData(resp, d, site)
}

func resourceDynamicDNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	err := c.c.DeleteDynamicDNS(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return diag.FromErr(err)
}
