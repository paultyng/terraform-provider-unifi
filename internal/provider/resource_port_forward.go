package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourcePortForward() *schema.Resource {
	return &schema.Resource{
		Create: resourcePortForwardCreate,
		Read:   resourcePortForwardRead,
		Update: resourcePortForwardUpdate,
		Delete: resourcePortForwardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"dst_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePortRange,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"fwd_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"fwd_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePortRange,
			},
			"log": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_forward_interface": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"wan", "wan2", "both"}, false),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "tcp_udp",
				ValidateFunc: validation.StringInSlice([]string{"tcp_udp", "tcp", "udp"}, false),
			},
			"src_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "any",
				ValidateFunc: validation.Any(
					validation.StringInSlice([]string{"any"}, false),
					validation.IsIPv4Address,
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

	resp, err := c.c.CreatePortForward(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePortForwardSetResourceData(resp, d)
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

func resourcePortForwardSetResourceData(resp *unifi.PortForward, d *schema.ResourceData) error {
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

	resp, err := c.c.GetPortForward(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourcePortForwardSetResourceData(resp, d)
}

func resourcePortForwardUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourcePortForwardGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = c.site

	resp, err := c.c.UpdatePortForward(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	return resourcePortForwardSetResourceData(resp, d)
}

func resourcePortForwardDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	err := c.c.DeletePortForward(context.TODO(), c.site, id)
	return err
}
