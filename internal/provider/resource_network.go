package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

var (
	wanUsernameRegexp   = regexp.MustCompile("[^\"' ]+")
	validateWANUsername = validation.StringMatch(wanUsernameRegexp, "invalid WAN username")

	wanTypeRegexp   = regexp.MustCompile("disabled|dhcp|static|pppoe")
	validateWANType = validation.StringMatch(wanTypeRegexp, "invalid WAN connection type")

	wanTypeV6Regexp   = regexp.MustCompile("disabled|dhcpv6|static")
	validateWANTypeV6 = validation.StringMatch(wanTypeV6Regexp, "invalid WANv6 connection type")

	wanPasswordRegexp   = regexp.MustCompile("[^\"' ]+")
	validateWANPassword = validation.StringMatch(wanPasswordRegexp, "invalid WAN password")

	wanNetworkGroupRegexp   = regexp.MustCompile("WAN[2]?|WAN_LTE_FAILOVER")
	validateWANNetworkGroup = validation.StringMatch(wanNetworkGroupRegexp, "invalid WAN network group")

	wanV6NetworkGroupRegexp   = regexp.MustCompile("wan[2]?")
	validateWANV6NetworkGroup = validation.StringMatch(wanV6NetworkGroupRegexp, "invalid WANv6 network group")

	ipV6InterfaceTypeRegexp   = regexp.MustCompile("none|pd|static")
	validateIpV6InterfaceType = validation.StringMatch(ipV6InterfaceTypeRegexp, "invalid IPv6 interface type")

	// This is a slightly larger range than the UI, it includes some reserved ones, so could be tightened up.
	validateVLANID = validation.IntBetween(0, 4096)

	ipV6RAPriorityRegexp   = regexp.MustCompile("high|medium|low")
	validateIpV6RAPriority = validation.StringMatch(ipV6RAPriorityRegexp, "invalid IPv6 RA priority")
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_network` manages WAN/LAN/VLAN networks.",

		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importNetwork,
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
				Description:  "The VLAN ID of the network.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateVLANID,
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
				Description: "Specifies the lease time for DHCP addresses in seconds.",
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
			"dhcpd_boot_enabled": {
				Description: "Toggles on the DHCP boot options. Should be set to true when you want to have dhcpd_boot_filename, and dhcpd_boot_server to take effect.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"dhcpd_boot_server": {
				Description: "Specifies the IPv4 address of a TFTP server to network boot from.",
				Type:        schema.TypeString,
				// TODO: IPv4 validation?
				Optional: true,
			},
			"dhcpd_boot_filename": {
				Description: "Specifies the file to PXE boot from on the dhcpd_boot_server.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dhcp_relay_enabled": {
				Description: "Specifies whether DHCP relay is enabled or not on this network.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"dhcp_v6_dns": {
				Description: "Specifies the IPv6 addresses for the DNS server to be returned from the DHCP " +
					"server. Used if `dhcp_v6_dns_auto` is set to `false`.",
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 4,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv6Address,
					// TODO: should this ensure blank can't get through?
				},
			},
			"dhcp_v6_dns_auto": {
				Description: "Specifies DNS source to propagate. If set `false` the entries in `dhcp_v6_dns` are used, the upstream entries otherwise",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"dhcp_v6_enabled": {
				Description: "Enable stateful DHCPv6 for static configuration.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"dhcp_v6_lease": {
				Description: "Specifies the lease time for DHCPv6 addresses in seconds.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     86400,
			},
			"dhcp_v6_start": {
				Description:  "Start address of the DHCPv6 range. Used in static DHCPv6 configuration.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv6Address,
			},
			"dhcp_v6_stop": {
				Description:  "End address of the DHCPv6 range. Used in static DHCPv6 configuration.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv6Address,
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
				Description:  "Specifies which type of IPv6 connection to use. Must be one of either `static`, `pd`, or `none`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validateIpV6InterfaceType,
			},
			"ipv6_static_subnet": {
				Description: "Specifies the static IPv6 subnet when `ipv6_interface_type` is 'static'.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"ipv6_pd_interface": {
				Description:  "Specifies which WAN interface to use for IPv6 PD. Must be one of either `wan` or `wan2`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateWANV6NetworkGroup,
			},
			"ipv6_pd_prefixid": {
				Description: "Specifies the IPv6 Prefix ID.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ipv6_pd_start": {
				Description:  "Start address of the DHCPv6 range. Used if `ipv6_interface_type` is set to `pd`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv6Address,
			},
			"ipv6_pd_stop": {
				Description:  "End address of the DHCPv6 range. Used if `ipv6_interface_type` is set to `pd`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv6Address,
			},
			"ipv6_ra_enable": {
				Description: "Specifies whether to enable router advertisements or not.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"internet_access_enabled": {
				Description: "Specifies whether this network should be allowed to access the internet or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"intra_network_access_enabled": {
				Description: "Specifies whether this network should be allowed to access other local networks or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"ipv6_ra_preferred_lifetime": {
				Description: "Lifetime in which the address can be used. Address becomes deprecated afterwards. Must be lower than or equal to `ipv6_ra_valid_lifetime`",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     14400,
			},
			"ipv6_ra_priority": {
				Description:  "IPv6 router advertisement priority. Must be one of either `high`, `medium`, or `low`",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIpV6RAPriority,
			},
			"ipv6_ra_valid_lifetime": {
				Description: "Total lifetime in which the address can be used. Must be equal to or greater than `ipv6_ra_preferred_lifetime`.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     86400,
			},
			"multicast_dns": {
				Description: "Specifies whether Multicast DNS (mDNS) is enabled or not on the network (Controller >=v7).",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"wan_ip": {
				Description:  "The IPv4 address of the WAN.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"wan_netmask": {
				Description:  "The IPv4 netmask of the WAN.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"wan_gateway": {
				Description:  "The IPv4 gateway of the WAN.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"wan_dns": {
				Description: "DNS servers IPs of the WAN.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    4,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv4Address,
				},
			},
			"wan_type": {
				Description:  "Specifies the IPV4 WAN connection type. Must be one of either `disabled`, `static`, `dhcp`, or `pppoe`.",
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
			"wan_type_v6": {
				Description:  "Specifies the IPV6 WAN connection type. Must be one of either `disabled`, `static`, or `dhcpv6`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateWANTypeV6,
			},
			"wan_dhcp_v6_pd_size": {
				Description:  "Specifies the IPv6 prefix size to request from ISP. Must be between 48 and 64.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(48, 64),
			},
			"wan_ipv6": {
				Description:  "The IPv6 address of the WAN.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv6Address,
			},
			"wan_gateway_v6": {
				Description:  "The IPv6 gateway of the WAN.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv6Address,
			},
			"wan_prefixlen": {
				Description:  "The IPv6 prefix length of the WAN. Must be between 1 and 128.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 128),
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceNetworkGetResourceData(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateNetwork(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourceNetworkSetResourceData(resp, d, site)
}

func resourceNetworkGetResourceData(d *schema.ResourceData, meta interface{}) (*unifi.Network, error) {
	// c := meta.(*client)

	vlan := d.Get("vlan_id").(int)
	dhcpDNS, err := listToStringSlice(d.Get("dhcp_dns").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to convert dhcp_dns to string slice: %w", err)
	}
	dhcpV6DNS, err := listToStringSlice(d.Get("dhcp_v6_dns").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to convert dhcp_v6_dns to string slice: %w", err)
	}
	wanDNS, err := listToStringSlice(d.Get("wan_dns").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to convert wan_dns to string slice: %w", err)
	}

	return &unifi.Network{
		Name:              d.Get("name").(string),
		Purpose:           d.Get("purpose").(string),
		VLAN:              vlan,
		IPSubnet:          cidrOneBased(d.Get("subnet").(string)),
		NetworkGroup:      d.Get("network_group").(string),
		DHCPDStart:        d.Get("dhcp_start").(string),
		DHCPDStop:         d.Get("dhcp_stop").(string),
		DHCPDEnabled:      d.Get("dhcp_enabled").(bool),
		DHCPDLeaseTime:    d.Get("dhcp_lease").(int),
		DHCPDBootEnabled:  d.Get("dhcpd_boot_enabled").(bool),
		DHCPDBootServer:   d.Get("dhcpd_boot_server").(string),
		DHCPDBootFilename: d.Get("dhcpd_boot_filename").(string),
		DHCPRelayEnabled:  d.Get("dhcp_relay_enabled").(bool),
		DomainName:        d.Get("domain_name").(string),
		IGMPSnooping:      d.Get("igmp_snooping").(bool),
		MdnsEnabled:       d.Get("multicast_dns").(bool),

		DHCPDDNSEnabled: len(dhcpDNS) > 0,
		// this is kinda hacky but ¯\_(ツ)_/¯
		DHCPDDNS1: append(dhcpDNS, "")[0],
		DHCPDDNS2: append(dhcpDNS, "", "")[1],
		DHCPDDNS3: append(dhcpDNS, "", "", "")[2],
		DHCPDDNS4: append(dhcpDNS, "", "", "", "")[3],

		VLANEnabled: vlan != 0 && vlan != 1,

		Enabled: true,

		// Same hackish code as for DHCPv4 ¯\_(ツ)_/¯
		DHCPDV6DNS1: append(dhcpV6DNS, "")[0],
		DHCPDV6DNS2: append(dhcpV6DNS, "", "")[1],
		DHCPDV6DNS3: append(dhcpV6DNS, "", "", "")[2],
		DHCPDV6DNS4: append(dhcpV6DNS, "", "", "", "")[3],

		DHCPDV6DNSAuto:   d.Get("dhcp_v6_dns_auto").(bool),
		DHCPDV6Enabled:   d.Get("dhcp_v6_enabled").(bool),
		DHCPDV6LeaseTime: d.Get("dhcp_v6_lease").(int),
		DHCPDV6Start:     d.Get("dhcp_v6_start").(string),
		DHCPDV6Stop:      d.Get("dhcp_v6_stop").(string),

		IPV6InterfaceType:       d.Get("ipv6_interface_type").(string),
		IPV6Subnet:              d.Get("ipv6_static_subnet").(string),
		IPV6PDInterface:         d.Get("ipv6_pd_interface").(string),
		IPV6PDPrefixid:          d.Get("ipv6_pd_prefixid").(string),
		IPV6PDStart:             d.Get("ipv6_pd_start").(string),
		IPV6PDStop:              d.Get("ipv6_pd_stop").(string),
		IPV6RaEnabled:           d.Get("ipv6_ra_enable").(bool),
		IPV6RaPreferredLifetime: d.Get("ipv6_ra_preferred_lifetime").(int),
		IPV6RaPriority:          d.Get("ipv6_ra_priority").(string),
		IPV6RaValidLifetime:     d.Get("ipv6_ra_valid_lifetime").(int),

		InternetAccessEnabled:     d.Get("internet_access_enabled").(bool),
		IntraNetworkAccessEnabled: d.Get("intra_network_access_enabled").(bool),

		WANIP:           d.Get("wan_ip").(string),
		WANType:         d.Get("wan_type").(string),
		WANNetmask:      d.Get("wan_netmask").(string),
		WANGateway:      d.Get("wan_gateway").(string),
		WANNetworkGroup: d.Get("wan_networkgroup").(string),
		WANEgressQOS:    d.Get("wan_egress_qos").(int),
		WANUsername:     d.Get("wan_username").(string),
		XWANPassword:    d.Get("x_wan_password").(string),

		WANTypeV6:       d.Get("wan_type_v6").(string),
		WANDHCPv6PDSize: d.Get("wan_dhcp_v6_pd_size").(int),
		WANIPV6:         d.Get("wan_ipv6").(string),
		WANGatewayV6:    d.Get("wan_gateway_v6").(string),
		WANPrefixlen:    d.Get("wan_prefixlen").(int),

		// this is kinda hacky but ¯\_(ツ)_/¯
		WANDNS1: append(wanDNS, "")[0],
		WANDNS2: append(wanDNS, "", "")[1],
		WANDNS3: append(wanDNS, "", "", "")[2],
		WANDNS4: append(wanDNS, "", "", "", "")[3],
	}, nil
}

func resourceNetworkSetResourceData(resp *unifi.Network, d *schema.ResourceData, site string) diag.Diagnostics {
	wanType := ""
	wanDNS := []string{}
	wanIP := ""
	wanNetmask := ""
	wanGateway := ""

	if resp.Purpose == "wan" {
		wanType = resp.WANType

		for _, dns := range []string{
			resp.WANDNS1,
			resp.WANDNS2,
			resp.WANDNS3,
			resp.WANDNS4,
		} {
			if dns == "" {
				continue
			}
			wanDNS = append(wanDNS, dns)
		}

		if wanType != "dhcp" {
			wanIP = resp.WANIP
			wanNetmask = resp.WANNetmask
			wanGateway = resp.WANGateway
		}

		// TODO: set other wan only fields here?
	}

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

	dhcpV6DNS := []string{}
	for _, dns := range []string{
		resp.DHCPDV6DNS1,
		resp.DHCPDV6DNS2,
		resp.DHCPDV6DNS3,
		resp.DHCPDV6DNS4,
	} {
		if dns == "" {
			continue
		}
		dhcpV6DNS = append(dhcpV6DNS, dns)
	}

	d.Set("site", site)
	d.Set("name", resp.Name)
	d.Set("purpose", resp.Purpose)
	d.Set("vlan_id", vlan)
	d.Set("subnet", cidrZeroBased(resp.IPSubnet))
	d.Set("network_group", resp.NetworkGroup)

	d.Set("dhcp_dns", dhcpDNS)
	d.Set("dhcp_enabled", resp.DHCPDEnabled)
	d.Set("dhcp_lease", dhcpLease)
	d.Set("dhcp_relay_enabled", resp.DHCPRelayEnabled)
	d.Set("dhcp_start", resp.DHCPDStart)
	d.Set("dhcp_stop", resp.DHCPDStop)
	d.Set("dhcp_v6_dns_auto", resp.DHCPDV6DNSAuto)
	d.Set("dhcp_v6_dns", dhcpV6DNS)
	d.Set("dhcp_v6_enabled", resp.DHCPDV6Enabled)
	d.Set("dhcp_v6_lease", resp.DHCPDV6LeaseTime)
	d.Set("dhcp_v6_start", resp.DHCPDV6Start)
	d.Set("dhcp_v6_stop", resp.DHCPDV6Stop)
	d.Set("dhcpd_boot_enabled", resp.DHCPDBootEnabled)
	d.Set("dhcpd_boot_filename", resp.DHCPDBootFilename)
	d.Set("dhcpd_boot_server", resp.DHCPDBootServer)
	d.Set("domain_name", resp.DomainName)
	d.Set("igmp_snooping", resp.IGMPSnooping)
	d.Set("internet_access_enabled", resp.InternetAccessEnabled)
	d.Set("intra_network_access_enabled", resp.IntraNetworkAccessEnabled)
	d.Set("ipv6_interface_type", resp.IPV6InterfaceType)
	d.Set("ipv6_pd_interface", resp.IPV6PDInterface)
	d.Set("ipv6_pd_prefixid", resp.IPV6PDPrefixid)
	d.Set("ipv6_pd_start", resp.IPV6PDStart)
	d.Set("ipv6_pd_stop", resp.IPV6PDStop)
	d.Set("ipv6_ra_enable", resp.IPV6RaEnabled)
	d.Set("ipv6_ra_preferred_lifetime", resp.IPV6RaPreferredLifetime)
	d.Set("ipv6_ra_priority", resp.IPV6RaPriority)
	d.Set("ipv6_ra_valid_lifetime", resp.IPV6RaValidLifetime)
	d.Set("ipv6_static_subnet", resp.IPV6Subnet)
	d.Set("multicast_dns", resp.MdnsEnabled)
	d.Set("wan_dhcp_v6_pd_size", resp.WANDHCPv6PDSize)
	d.Set("wan_dns", wanDNS)
	d.Set("wan_egress_qos", resp.WANEgressQOS)
	d.Set("wan_gateway_v6", resp.WANGatewayV6)
	d.Set("wan_gateway", wanGateway)
	d.Set("wan_ip", wanIP)
	d.Set("wan_ipv6", resp.WANIPV6)
	d.Set("wan_netmask", wanNetmask)
	d.Set("wan_networkgroup", resp.WANNetworkGroup)
	d.Set("wan_prefixlen", resp.WANPrefixlen)
	d.Set("wan_type_v6", resp.WANTypeV6)
	d.Set("wan_type", wanType)
	d.Set("wan_username", resp.WANUsername)
	d.Set("x_wan_password", resp.XWANPassword)

	return nil
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetNetwork(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetworkSetResourceData(resp, d, site)
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceNetworkGetResourceData(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateNetwork(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetworkSetResourceData(resp, d, site)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	name := d.Get("name").(string)
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	id := d.Id()

	err := c.c.DeleteNetwork(ctx, site, id, name)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return diag.FromErr(err)
}

func importNetwork(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(*client)
	id := d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	if strings.Contains(id, ":") {
		importParts := strings.SplitN(id, ":", 2)
		site = importParts[0]
		id = importParts[1]
	}

	if strings.HasPrefix(id, "name=") {
		targetName := strings.TrimPrefix(id, "name=")
		var err error
		if id, err = getNetworkIDByName(ctx, c.c, targetName, site); err != nil {
			return nil, err
		}
	}

	if id != "" {
		d.SetId(id)
	}
	if site != "" {
		d.Set("site", site)
	}

	return []*schema.ResourceData{d}, nil
}

func getNetworkIDByName(ctx context.Context, client unifiClient, networkName, site string) (string, error) {
	networks, err := client.ListNetwork(ctx, site)
	if err != nil {
		return "", err
	}

	idMatchingName := ""
	allNames := []string{}
	for _, network := range networks {
		allNames = append(allNames, network.Name)
		if network.Name != networkName {
			continue
		}
		if idMatchingName != "" {
			return "", fmt.Errorf("Found multiple networks with name '%s'", networkName)
		}
		idMatchingName = network.ID
	}
	if idMatchingName == "" {
		return "", fmt.Errorf("Found no networks with name '%s', found: %s", networkName, strings.Join(allNames, ", "))
	}
	return idMatchingName, nil
}
