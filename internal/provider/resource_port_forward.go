package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

//	GetPortForward(site string, id string) (*unifi.PortForward, error)
//	DeletePortForward(site string, id string) (*unifi.PortForward, error)
//	CreatePortForward(site string, d *unifi.PortForward) (*unifi.PortForward, error)
//	UpdatePortForward(site string, d *unifi.PortForward) (*unifi.PortForward, error)

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
				Required:     false,
				Optional:     true,
				ValidateFunc: validation.StringMatch(portForwardDstPortRegexp, "dst_port is invalid"),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"fwd": {
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				ValidateFunc: validation.StringMatch(portForwardFwdRegexp, "fwd is invalid"),
			},
			"fwd_port": {
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				ValidateFunc: validation.StringMatch(portForwardFwdPortRegexp, "fwd_port is invalid"),
			},
			"log": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"port_forward_interface": {
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"wan", "wan2", "both", ""}, false),
			},
			"proto": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "tcp_udp",
				ValidateFunc: validation.StringInSlice([]string{"tcp_udp", "tcp", "udp"}, false),
			},
			"src": {
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				Default:      "any",
				ValidateFunc: validation.StringMatch(portForwardSrcRegexp, "src is invalid"),
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

	resp, err := c.c.CreatePortForward(c.site, req)
	if err != nil {
		apiErr, ok := err.(*unifi.APIError)
		if !ok {

			return fmt.Errorf("ResourcePortForwardCreate api error %w. With %w", apiErr, err)
		}

	}

	d.SetId(resp.ID)

	return resourcePortForwardSetResourceData(resp, d)
}

func resourcePortForwardGetResourceData(d *schema.ResourceData) (*unifi.PortForward, error) {
	return &unifi.PortForward{
		DstPort:       d.Get("dst_port").(string),
		Enabled:       d.Get("enabled").(bool),
		Fwd:           d.Get("fwd").(string),
		FwdPort:       d.Get("fwd_port").(string),
		Log:           d.Get("log").(bool),
		Name:          d.Get("name").(string),
		PfwdInterface: d.Get("port_forward_interface").(string),
		Proto:         d.Get("proto").(string),
		Src:           d.Get("src").(string),
	}, nil
}

func resourcePortForwardSetResourceData(resp *unifi.PortForward, d *schema.ResourceData) error {
	log := false
	if resp.Log {
		log = resp.Log
	}
	d.Set("dst_port", resp.DstPort)
	d.Set("enabled", resp.Enabled)
	d.Set("fwd", resp.Fwd)
	d.Set("fwd_port", resp.FwdPort)
	d.Set("log", log)
	d.Set("name", resp.Name)
	d.Set("port_forward_interface", resp.PfwdInterface)
	d.Set("proto", resp.Proto)
	d.Set("src", resp.Src)

	return nil
}

func resourcePortForwardRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	resp, err := c.c.GetPortForward(c.site, id)
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

	resp, err := c.c.UpdatePortForward(c.site, req)
	if err != nil {
		return err
	}

	return resourcePortForwardSetResourceData(resp, d)
}

func resourcePortForwardDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	err := c.c.DeletePortForward(c.site, id)
	return err
}
