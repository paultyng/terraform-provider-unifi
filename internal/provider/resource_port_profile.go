package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourcePortProfile() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_port_profile` manages a port profile for use on network switches.",

		Create: resourcePortProfileCreate,
		Read:   resourcePortProfileRead,
		Update: resourcePortProfileUpdate,
		Delete: resourcePortProfileDelete,
		Importer: &schema.ResourceImporter{
			State: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the port profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the port profile with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"autoneg": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"dot1x_ctrl": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
				//ValidateFunc: ,
			},
			"dot1x_idle_timeout": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"egress_rate_limit_kbps": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"egress_rate_limit_kbps_enabled": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"forward": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
				//ValidateFunc: ,
			},
			"full_duplex": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"isolation": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"lldpmed_enabled": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"lldpmed_notify_enabled": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				//ValidateFunc: ,
			},
			"native_networkconf_id": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "The name of the port profile.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"op_mode": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
				//ValidateFunc: ,
			},
			"poe_mode": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
				//ValidateFunc: ,
			},
			"port_security_enabled": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"port_security_mac_address": {
				Description: "",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				//ValidateFunc: ,
			},
			"priority_queue1_level": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"priority_queue2_level": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"priority_queue3_level": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"priority_queue4_level": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"speed": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"stormctrl_bcast_enabled": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"stormctrl_bcast_level": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"stormctrl_bcast_rate": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"stormctrl_mcast_enabled": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"stormctrl_mcast_level": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"stormctrl_mcast_rate": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"stormctrl_type": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
				//ValidateFunc: ,
			},
			"stormctrl_ucast_enabled": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"stormctrl_ucast_level": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"stormctrl_ucast_rate": {
				Description: "",
				Type:        schema.TypeInt,
				Optional:    true,
				//ValidateFunc: ,
			},
			"stp_port_mode": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"tagged_networkconf_ids": {
				Description: "",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"voice_networkconf_id": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
				//ValidateFunc: ,
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

func resourcePortProfileCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourcePortProfileGetResourceData(d)
	if err != nil {
		return err
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	resp, err := c.c.CreatePortProfile(context.TODO(), site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourcePortProfileSetResourceData(resp, d, site)
}

func resourcePortProfileGetResourceData(d *schema.ResourceData) (*unifi.PortProfile, error) {
	portSecurityMacAddress, err := setToStringSlice(d.Get("port_security_mac_address").(*schema.Set))
	if err != nil {
		return nil, err
	}

	taggedNetworkconfIds, err := setToStringSlice(d.Get("tagged_networkconf_ids").(*schema.Set))
	if err != nil {
		return nil, err
	}

	return &unifi.PortProfile{
		Autoneg:                      d.Get("autoneg").(bool),
		Dot1XCtrl:                    d.Get("dot1x_ctrl").(string),
		Dot1XIDleTimeout:             d.Get("dot1x_idle_timeout").(int),
		EgressRateLimitKbps:          d.Get("egress_rate_limit_kbps").(int),
		EgressRateLimitKbpsEnabled:   d.Get("egress_rate_limit_kbps_enabled").(bool),
		Forward:                      d.Get("forward").(string),
		FullDuplex:                   d.Get("full_duplex").(bool),
		Isolation:                    d.Get("isolation").(bool),
		LldpmedEnabled:               d.Get("lldpmed_enabled").(bool),
		LldpmedNotifyEnabled:         d.Get("lldpmed_notify_enabled").(bool),
		NATiveNetworkID:              d.Get("native_networkconf_id").(string),
		Name:                         d.Get("name").(string),
		OpMode:                       d.Get("op_mode").(string),
		PoeMode:                      d.Get("poe_mode").(string),
		PortSecurityEnabled:          d.Get("port_security_enabled").(bool),
		PortSecurityMACAddress:       portSecurityMacAddress,
		PriorityQueue1Level:          d.Get("priority_queue1_level").(int),
		PriorityQueue2Level:          d.Get("priority_queue2_level").(int),
		PriorityQueue3Level:          d.Get("priority_queue3_level").(int),
		PriorityQueue4Level:          d.Get("priority_queue4_level").(int),
		Speed:                        d.Get("speed").(int),
		StormctrlBroadcastastEnabled: d.Get("stormctrl_bcast_enabled").(bool),
		StormctrlBroadcastastLevel:   d.Get("stormctrl_bcast_level").(int),
		StormctrlBroadcastastRate:    d.Get("stormctrl_bcast_rate").(int),
		StormctrlMcastEnabled:        d.Get("stormctrl_mcast_enabled").(bool),
		StormctrlMcastLevel:          d.Get("stormctrl_mcast_level").(int),
		StormctrlMcastRate:           d.Get("stormctrl_mcast_rate").(int),
		StormctrlType:                d.Get("stormctrl_type").(string),
		StormctrlUcastEnabled:        d.Get("stormctrl_ucast_enabled").(bool),
		StormctrlUcastLevel:          d.Get("stormctrl_ucast_level").(int),
		StormctrlUcastRate:           d.Get("stormctrl_ucast_rate").(int),
		StpPortMode:                  d.Get("stp_port_mode").(bool),
		TaggedNetworkIDs:             taggedNetworkconfIds,
		VoiceNetworkID:               d.Get("voice_networkconf_id").(string),
	}, nil
}

func resourcePortProfileSetResourceData(resp *unifi.PortProfile, d *schema.ResourceData, site string) error {
	d.Set("site", site)
	d.Set("autoneg", resp.Autoneg)
	d.Set("dot1x_ctrl", resp.Dot1XCtrl)
	d.Set("dot1x_idle_timeout", resp.Dot1XIDleTimeout)
	d.Set("egress_rate_limit_kbps", resp.EgressRateLimitKbps)
	d.Set("egress_rate_limit_kbps_enabled", resp.EgressRateLimitKbpsEnabled)
	d.Set("forward", resp.Forward)
	d.Set("full_duplex", resp.FullDuplex)
	d.Set("isolation", resp.Isolation)
	d.Set("lldpmed_enabled", resp.LldpmedEnabled)
	d.Set("lldpmed_notify_enabled", resp.LldpmedNotifyEnabled)
	d.Set("native_networkconf_id", resp.NATiveNetworkID)
	d.Set("name", resp.Name)
	d.Set("op_mode", resp.OpMode)
	d.Set("poe_mode", resp.PoeMode)
	d.Set("port_security_enabled", resp.PortSecurityEnabled)
	d.Set("port_security_mac_address", stringSliceToSet(resp.PortSecurityMACAddress))
	d.Set("priority_queue1_level", resp.PriorityQueue1Level)
	d.Set("priority_queue2_level", resp.PriorityQueue2Level)
	d.Set("priority_queue3_level", resp.PriorityQueue3Level)
	d.Set("priority_queue4_level", resp.PriorityQueue4Level)
	d.Set("speed", resp.Speed)
	d.Set("stormctrl_bcast_enabled", resp.StormctrlBroadcastastEnabled)
	d.Set("stormctrl_bcast_level", resp.StormctrlBroadcastastLevel)
	d.Set("stormctrl_bcast_rate", resp.StormctrlBroadcastastRate)
	d.Set("stormctrl_mcast_enabled", resp.StormctrlMcastEnabled)
	d.Set("stormctrl_mcast_level", resp.StormctrlMcastLevel)
	d.Set("stormctrl_mcast_rate", resp.StormctrlMcastRate)
	d.Set("stormctrl_type", resp.StormctrlType)
	d.Set("stormctrl_ucast_enabled", resp.StormctrlUcastEnabled)
	d.Set("stormctrl_ucast_level", resp.StormctrlUcastLevel)
	d.Set("stormctrl_ucast_rate", resp.StormctrlUcastRate)
	d.Set("stp_port_mode", resp.StpPortMode)
	d.Set("tagged_networkconf_ids", stringSliceToSet(resp.TaggedNetworkIDs))
	d.Set("voice_networkconf_id", resp.VoiceNetworkID)

	return nil
}

func resourcePortProfileRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	resp, err := c.c.GetPortProfile(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourcePortProfileSetResourceData(resp, d, site)
}

func resourcePortProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourcePortProfileGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdatePortProfile(context.TODO(), site, req)
	if err != nil {
		return err
	}

	return resourcePortProfileSetResourceData(resp, d, site)
}

func resourcePortProfileDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	err := c.c.DeletePortProfile(context.TODO(), site, id)
	return err
}
