package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: `
unifi_network manages LAN/VLAN networks.
`,

		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"purpose": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"corporate", "guest", "vlan-only"}, false),
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"subnet": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: cidrDiffSuppress,
			},
			"network_group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LAN",
			},
			"dhcp_start": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"dhcp_stop": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
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
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"igmp_snooping": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceNetworkGetResourceData(d)
	if err != nil {
		return err
	}

	resp, err := c.c.CreateNetwork(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceNetworkSetResourceData(resp, d)
}

func resourceNetworkGetResourceData(d *schema.ResourceData) (*unifi.Network, error) {
	vlan := d.Get("vlan_id").(int)

	return &unifi.Network{
		Name:           d.Get("name").(string),
		Purpose:        d.Get("purpose").(string),
		VLAN:           vlan,
		IPSubnet:       cidrOneBased(d.Get("subnet").(string)),
		NetworkGroup:   d.Get("network_group").(string),
		DHCPDStart:     d.Get("dhcp_start").(string),
		DHCPDStop:      d.Get("dhcp_stop").(string),
		DHCPDEnabled:   d.Get("dhcp_enabled").(bool),
		DHCPDLeaseTime: d.Get("dhcp_lease").(int),
		DomainName:     d.Get("domain_name").(string),
		IGMPSnooping:   d.Get("igmp_snooping").(bool),

		VLANEnabled: vlan != 0 && vlan != 1,

		Enabled:           true,
		IPV6InterfaceType: "none",
		// IPV6InterfaceType string `json:"ipv6_interface_type"` // "none"
		// IPV6PDStart       string `json:"ipv6_pd_start"`       // "::2"
		// IPV6PDStop        string `json:"ipv6_pd_stop"`        // "::7d1"
	}, nil
}

func resourceNetworkSetResourceData(resp *unifi.Network, d *schema.ResourceData) error {
	vlan := 0
	if resp.VLANEnabled {
		vlan = resp.VLAN
	}

	dhcpLease := resp.DHCPDLeaseTime
	if resp.DHCPDEnabled && dhcpLease == 0 {
		dhcpLease = 86400
	}

	d.Set("name", resp.Name)
	d.Set("purpose", resp.Purpose)
	d.Set("vlan_id", vlan)
	d.Set("subnet", cidrZeroBased(resp.IPSubnet))
	d.Set("network_group", resp.NetworkGroup)
	d.Set("dhcp_start", resp.DHCPDStart)
	d.Set("dhcp_stop", resp.DHCPDStop)
	d.Set("dhcp_enabled", resp.DHCPDEnabled)
	d.Set("dhcp_lease", dhcpLease)
	d.Set("domain_name", resp.DomainName)
	d.Set("igmp_snooping", resp.IGMPSnooping)

	return nil
}

func resourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	resp, err := c.c.GetNetwork(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceNetworkSetResourceData(resp, d)
}

func resourceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceNetworkGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = c.site

	resp, err := c.c.UpdateNetwork(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	return resourceNetworkSetResourceData(resp, d)
}

func resourceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	name := d.Get("name").(string)
	id := d.Id()

	err := c.c.DeleteNetwork(context.TODO(), c.site, id, name)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
