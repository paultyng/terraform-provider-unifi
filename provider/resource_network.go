package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/paultyng/terraform-provider-unifi/unifi"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,

		// TODO: handle site + ID (or name)
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },

		Schema: map[string]*schema.Schema{
			"site": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"purpose": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				//"corporate", "guest", "vlan-only"
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"subnet": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LAN",
			},
			"dhcp_start": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp_stop": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dhcp_lease": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  86400,
			},
		},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)

	req := &unifi.Network{
		Name:           d.Get("name").(string),
		Purpose:        d.Get("purpose").(string),
		VLAN:           fmt.Sprintf("%d", d.Get("vlan_id").(int)),
		IPSubnet:       d.Get("subnet").(string),
		NetworkGroup:   d.Get("network_group").(string),
		DHCPDStart:     d.Get("dhcp_start").(string),
		DHCPDStop:      d.Get("dhcp_stop").(string),
		DHCPDEnabled:   d.Get("dhcp_enabled").(bool),
		DHCPDLeaseTime: d.Get("dhcp_lease").(int),

		Enabled:           true,
		VLANEnabled:       true,
		IPV6InterfaceType: "none",
		// IPV6InterfaceType string `json:"ipv6_interface_type"` // "none"
		// IPV6PDStart       string `json:"ipv6_pd_start"`       // "::2"
		// IPV6PDStop        string `json:"ipv6_pd_stop"`        // "::7d1"
	}

	resp, err := c.c.CreateNetwork(site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return nil
}

func resourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)
	id := d.Id()

	_, err := c.c.GetNetwork(site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	panic("not implemented")
}

func resourceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)
	name := d.Get("name").(string)
	id := d.Id()

	err := c.c.DeleteNetwork(site, id, name)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
