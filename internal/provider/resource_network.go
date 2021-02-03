package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

var (
	wanUsernameRegexp   = regexp.MustCompile("[^\"' ]+")
	validateWANUsername = validation.StringMatch(wanUsernameRegexp, "invalid WAN username")

	wanTypeRegexp   = regexp.MustCompile("disabled|dhcp|static|pppoe")
	validateWANType = validation.StringMatch(wanTypeRegexp, "invalid WAN connection type")

	wanPasswordRegexp   = regexp.MustCompile("[^\"' ]+")
	validateWANPassword = validation.StringMatch(wanPasswordRegexp, "invalid WAN password")

	wanNetworkGroupRegexp   = regexp.MustCompile("WAN[2]?|WAN_LTE_FAILOVER")
	validateWANNetworkGroup = validation.StringMatch(wanNetworkGroupRegexp, "invalid WAN network group")
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_network` manages WAN/LAN/VLAN networks.",

		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the network.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the network with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the network.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"purpose": {
				Description:  "The purpose of the network. Must be one of `corporate`, `guest`, `wan`, or `vlan-only`.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"corporate", "guest", "wan", "vlan-only"}, false),
			},
			"vlan_id": {
				Description: "The VLAN ID of the network.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"subnet": {
				Description:      "The subnet of the network. Must be a valid CIDR address.",
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: cidrDiffSuppress,
				ValidateFunc:     cidrValidate,
			},
			"network_group": {
				Description: "The group of the network.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "LAN",
			},
			"dhcp_start": {
				Description:  "The IPv4 address where the DHCP range of addresses starts.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"dhcp_stop": {
				Description:  "The IPv4 address where the DHCP range of addresses stops.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"dhcp_enabled": {
				Description: "Specifies whether DHCP is enabled or not on this network.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"dhcp_lease": {
				Description: "Specifies the lease time for DHCP addresses.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     86400,
			},
			"dhcp_dns": {
				Description: "Specifies the IPv4 addresses for the DNS server to be returned from the DHCP " +
					"server. Leave blank to disable this feature.",
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 4,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						validation.IsIPv4Address,
						// this doesn't let blank through
						validation.StringLenBetween(1, 50),
					),
				},
			},
			"domain_name": {
				Description: "The domain name of this network.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"igmp_snooping": {
				Description: "Specifies whether IGMP snooping is enabled or not.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"ipv6_interface_type": {
				Description: "Specifies which type of IPv6 connection to use.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "none",
			},
			"ipv6_static_subnet": {
				Description: "Specifies the static IPv6 subnet when ipv6_interface_type is 'static'.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ipv6_pd_interface": {
				Description: "Specifies which WAN interface to use for IPv6 PD.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ipv6_pd_prefixid": {
				Description: "Specifies the IPv6 Prefix ID.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ipv6_ra_enable": {
				Description: "Specifies whether to enable router advertisements or not.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"wan_ip": {
				Description:  "The IPv4 address of the WAN.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"wan_type": {
				Description:  "Specifies the IPV4 WAN connection type. Must be one of either `disabled` or `pppoe`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateWANType,
			},
			"wan_networkgroup": {
				Description:  "Specifies the WAN network group. Must be one of either `WAN`, `WAN2` or `WAN_LTE_FAILOVER`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateWANNetworkGroup,
			},
			"wan_egress_qos": {
				Description: "Specifies the WAN egress quality of service.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"wan_username": {
				Description:  "Specifies the IPV4 WAN username.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateWANUsername,
			},
			"x_wan_password": {
				Description:  "Specifies the IPV4 WAN password.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateWANPassword,
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

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateNetwork(context.TODO(), site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceNetworkSetResourceData(resp, d, site)
}

func resourceNetworkGetResourceData(d *schema.ResourceData) (*unifi.Network, error) {
	vlan := d.Get("vlan_id").(int)
	dhcpDNS, err := listToStringSlice(d.Get("dhcp_dns").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to convert dhcp_dns to string slice: %w", err)
	}

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

		DHCPDDNSEnabled: len(dhcpDNS) > 0,
		// this is kinda hacky but ¯\_(ツ)_/¯
		DHCPDDNS1: append(dhcpDNS, "")[0],
		DHCPDDNS2: append(dhcpDNS, "", "")[1],
		DHCPDDNS3: append(dhcpDNS, "", "", "")[2],
		DHCPDDNS4: append(dhcpDNS, "", "", "", "")[3],

		VLANEnabled: vlan != 0 && vlan != 1,

		Enabled: true,

		IPV6InterfaceType: d.Get("ipv6_interface_type").(string),
		IPV6Subnet:        d.Get("ipv6_static_subnet").(string),
		IPV6PDInterface:   d.Get("ipv6_pd_interface").(string),
		IPV6PDPrefixid:    d.Get("ipv6_pd_prefixid").(string),
		IPV6RaEnabled:     d.Get("ipv6_ra_enable").(bool),

		WANIP:           d.Get("wan_ip").(string),
		WANType:         d.Get("wan_type").(string),
		WANNetworkGroup: d.Get("wan_networkgroup").(string),
		WANEgressQOS:    d.Get("wan_egress_qos").(int),
		WANUsername:     d.Get("wan_username").(string),
		XWANPassword:    d.Get("x_wan_password").(string),
	}, nil
}

func resourceNetworkSetResourceData(resp *unifi.Network, d *schema.ResourceData, site string) error {
	vlan := 0
	if resp.VLANEnabled {
		vlan = resp.VLAN
	}

	dhcpLease := resp.DHCPDLeaseTime
	if resp.DHCPDEnabled && dhcpLease == 0 {
		dhcpLease = 86400
	}

	dhcpDNS := []string{}
	if resp.DHCPDDNSEnabled {
		for _, dns := range []string{
			resp.DHCPDDNS1,
			resp.DHCPDDNS2,
			resp.DHCPDDNS3,
			resp.DHCPDDNS4,
		} {
			if dns == "" {
				continue
			}
			dhcpDNS = append(dhcpDNS, dns)
		}
	}

	d.Set("site", site)
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
	d.Set("dhcp_dns", dhcpDNS)
	d.Set("ipv6_interface_type", resp.IPV6InterfaceType)
	d.Set("ipv6_static_subnet", resp.IPV6Subnet)
	d.Set("ipv6_pd_interface", resp.IPV6PDInterface)
	d.Set("ipv6_pd_prefixid", resp.IPV6PDPrefixid)
	d.Set("ipv6_ra_enable", resp.IPV6RaEnabled)
	d.Set("wan_ip", resp.WANIP)
	d.Set("wan_type", resp.WANType)
	d.Set("wan_networkgroup", resp.WANNetworkGroup)
	d.Set("wan_egress_qos", resp.WANEgressQOS)
	d.Set("wan_username", resp.WANUsername)
	d.Set("x_wan_password", resp.XWANPassword)

	return nil
}

func resourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetNetwork(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceNetworkSetResourceData(resp, d, site)
}

func resourceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceNetworkGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateNetwork(context.TODO(), site, req)
	if err != nil {
		return err
	}

	return resourceNetworkSetResourceData(resp, d, site)
}

func resourceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	name := d.Get("name").(string)
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	id := d.Id()

	err := c.c.DeleteNetwork(context.TODO(), site, id, name)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
