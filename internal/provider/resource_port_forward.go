package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourcePortForward() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_port_forward` manages a port forwarding rule on the gateway.",

		Create: resourcePortForwardCreate,
		Read:   resourcePortForwardRead,
		Update: resourcePortForwardUpdate,
		Delete: resourcePortForwardDelete,
		Importer: &schema.ResourceImporter{
			State: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the port forwarding rule.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the port forwarding rule with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"dst_port": {
				Description:  "The destination port for the forwarding.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePortRange,
			},
			// TODO: remove this, disabled rules should just be deleted.
			"enabled": {
				Description: "Specifies whether the port forwarding rule is enabled or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Deprecated: "This will attribute will be removed in a future release. Instead of disabling a " +
					"port forwarding rule you can remove it from your configuration.",
			},
			"fwd_ip": {
				Description:  "The IPv4 address to forward traffic to.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"fwd_port": {
				Description:  "The port to forward traffic to.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePortRange,
			},
			"log": {
				Description: "Specifies whether to log forwarded traffic or not.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"name": {
				Description: "The name of the port forwarding rule.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"port_forward_interface": {
				Description:  "The port forwarding interface. Can be `wan`, `wan2`, or `both`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"wan", "wan2", "both"}, false),
			},
			"protocol": {
				Description:  "The protocol for the port forwarding rule. Can be `tcp`, `udp`, or `tcp_udp`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "tcp_udp",
				ValidateFunc: validation.StringInSlice([]string{"tcp_udp", "tcp", "udp"}, false),
			},
			"src_ip": {
				Description: "The source IPv4 address (or CIDR) of the port forwarding rule. For all traffic, specify `any`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "any",
				ValidateFunc: validation.Any(
					validation.StringInSlice([]string{"any"}, false),
					validation.IsIPv4Address,
					cidrValidate,
				),
			},
		},
	}
}

func resourcePortForwardCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourcePortForwardGetResourceData(d)
	if err != nil {
		return err
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	resp, err := c.c.CreatePortForward(context.TODO(), site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePortForwardSetResourceData(resp, d, site)
}

func resourcePortForwardGetResourceData(d *schema.ResourceData) (*unifi.PortForward, error) {
	return &unifi.PortForward{
		DstPort:       d.Get("dst_port").(string),
		Enabled:       d.Get("enabled").(bool),
		Fwd:           d.Get("fwd_ip").(string),
		FwdPort:       d.Get("fwd_port").(string),
		Log:           d.Get("log").(bool),
		Name:          d.Get("name").(string),
		PfwdInterface: d.Get("port_forward_interface").(string),
		Proto:         d.Get("protocol").(string),
		Src:           d.Get("src_ip").(string),
	}, nil
}

func resourcePortForwardSetResourceData(resp *unifi.PortForward, d *schema.ResourceData, site string) error {
	d.Set("site", site)
	d.Set("dst_port", resp.DstPort)
	d.Set("enabled", resp.Enabled)
	d.Set("fwd_ip", resp.Fwd)
	d.Set("fwd_port", resp.FwdPort)
	d.Set("log", resp.Log)
	d.Set("name", resp.Name)
	d.Set("port_forward_interface", resp.PfwdInterface)
	d.Set("protocol", resp.Proto)
	d.Set("src_ip", resp.Src)

	return nil
}

func resourcePortForwardRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	resp, err := c.c.GetPortForward(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourcePortForwardSetResourceData(resp, d, site)
}

func resourcePortForwardUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourcePortForwardGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdatePortForward(context.TODO(), site, req)
	if err != nil {
		return err
	}

	return resourcePortForwardSetResourceData(resp, d, site)
}

func resourcePortForwardDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	err := c.c.DeletePortForward(context.TODO(), site, id)
	return err
}
