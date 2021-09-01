package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

var firewallRuleProtocolRegexp = regexp.MustCompile("^$|all|([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])|tcp_udp|ah|ax.25|dccp|ddp|egp|eigrp|encap|esp|etherip|fc|ggp|gre|hip|hmp|icmp|idpr-cmtp|idrp|igmp|igp|ip|ipcomp|ipencap|ipip|ipv6|ipv6-frag|ipv6-icmp|ipv6-nonxt|ipv6-opts|ipv6-route|isis|iso-tp4|l2tp|manet|mobility-header|mpls-in-ip|ospf|pim|pup|rdp|rohc|rspf|rsvp|sctp|shim6|skip|st|tcp|udp|udplite|vmtp|vrrp|wesp|xns-idp|xtp")

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_firewall_rule` manages an individual firewall rule on the gateway.",

		CreateContext: resourceFirewallRuleCreate,
		ReadContext:   resourceFirewallRuleRead,
		UpdateContext: resourceFirewallRuleUpdate,
		DeleteContext: resourceFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the firewall rule.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the firewall rule with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the firewall rule.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description:  "The action of the firewall rule. Must be one of `drop`, `accept`, or `reject`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"drop", "accept", "reject"}, false),
			},
			"ruleset": {
				Description: "The ruleset for the rule. This is from the perspective of the security gateway. " +
					"Must be one of `WAN_IN`, `WAN_OUT`, `WAN_LOCAL`, `LAN_IN`, `LAN_OUT`, `LAN_LOCAL`, `GUEST_IN`, " +
					"`GUEST_OUT`, `GUEST_LOCAL`, `WANv6_IN`, `WANv6_OUT`, `WANv6_LOCAL`, `LANv6_IN`, `LANv6_OUT`, " +
					"`LANv6_LOCAL`, `GUESTv6_IN`, `GUESTv6_OUT`, or `GUESTv6_LOCAL`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"WAN_IN", "WAN_OUT", "WAN_LOCAL", "LAN_IN", "LAN_OUT", "LAN_LOCAL", "GUEST_IN", "GUEST_OUT", "GUEST_LOCAL", "WANv6_IN", "WANv6_OUT", "WANv6_LOCAL", "LANv6_IN", "LANv6_OUT", "LANv6_LOCAL", "GUESTv6_IN", "GUESTv6_OUT", "GUESTv6_LOCAL"}, false),
			},
			"rule_index": {
				Description: "The index of the rule. Must be >= 2000 < 3000 or >= 4000 < 5000.",
				Type:        schema.TypeInt,
				Required:    true,
				// 2[0-9]{3}|4[0-9]{3}
			},
			"protocol": {
				Description:  "The protocol of the rule.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(firewallRuleProtocolRegexp, "must be a valid protocol"),
			},
			"icmp_typename": {
				Description: "ICMP type name.",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// sources
			"src_network_id": {
				Description: "The source network ID for the firewall rule.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"src_network_type": {
				Description:  "The source network type of the firewall rule. Can be one of `ADDRv4` or `NETv4`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NETv4",
				ValidateFunc: validation.StringInSlice([]string{"ADDRv4", "NETv4"}, false),
			},
			"src_firewall_group_ids": {
				Description: "The source firewall group IDs for the firewall rule.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"src_address": {
				Description: "The source address for the firewall rule.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"src_mac": {
				Description: "The source MAC address of the firewall rule.",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// destinations
			"dst_network_id": {
				Description: "The destination network ID of the firewall rule.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dst_network_type": {
				Description:  "The destination network type of the firewall rule. Can be one of `ADDRv4` or `NETv4`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NETv4",
				ValidateFunc: validation.StringInSlice([]string{"ADDRv4", "NETv4"}, false),
			},
			"dst_firewall_group_ids": {
				Description: "The destination firewall group IDs of the firewall rule.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dst_address": {
				Description: "The destination address of the firewall rule.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dst_port": {
				Description:  "The destination port of the firewall rule.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePortRange,
			},

			// advanced
			"logging": {
				Description: "Enable logging for the firewall rule.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"state_established": {
				Description: "Match where the state is established.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"state_invalid": {
				Description: "Match where the state is invalid.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"state_new": {
				Description: "Match where the state is new.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"state_related": {
				Description: "Match where the state is related.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"ip_sec": {
				Description:  "Specify whether the rule matches on IPsec packets. Can be one of `match-ipset` or `match-none`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"match-ipsec", "match-none"}, false),
			},
			"weekdays": {
				Description: "Days of the week the rule is enforced.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"weekdays_negate": {
				Description: "Negate (inverse) the days of the week the rule is enforced.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"monthdays": {
				Description: "Days of the month the rule is enforced.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"monthdays_negate": {
				Description: "Negate (inverse) the days of the month the rule is enforced.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"startdate": {
				Description: "Date that the rule starts being enforced.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"starttime": {
				Description: "Time that the rule starts being enforced.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"stopdate": {
				Description: "Date that the rule stops being enforced.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"stoptime": {
				Description: "Time that the rule stops being enforced.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"utc": {
				Description: "Controls whether time calculations for firewall rule is based on UTC time.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}

func resourceFirewallRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceFirewallRuleGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateFirewallRule(ctx, site, req)
	if err != nil {
		var apiErr *unifi.APIError
		if errors.As(err, &apiErr) && apiErr.Message == "api.err.FirewallGroupTypeExists" {
			return diag.Errorf("firewall rule groups must be of different group types (ie. a port group and address group): %s", err)
		}

		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourceFirewallRuleSetResourceData(resp, d, site)
}

func resourceFirewallRuleGetResourceData(d *schema.ResourceData) (*unifi.FirewallRule, error) {
	srcFirewallGroupIDs, err := setToStringSlice(d.Get("src_firewall_group_ids").(*schema.Set))
	if err != nil {
		return nil, err
	}

	dstFirewallGroupIDs, err := setToStringSlice(d.Get("dst_firewall_group_ids").(*schema.Set))
	if err != nil {
		return nil, err
	}

	return &unifi.FirewallRule{
		Enabled:          true,
		Name:             d.Get("name").(string),
		Action:           d.Get("action").(string),
		Ruleset:          d.Get("ruleset").(string),
		RuleIndex:        d.Get("rule_index").(int),
		Protocol:         d.Get("protocol").(string),
		ICMPTypename:     d.Get("icmp_typename").(string),
		Logging:          d.Get("logging").(bool),
		IPSec:            d.Get("ip_sec").(string),
		StateEstablished: d.Get("state_established").(bool),
		StateInvalid:     d.Get("state_invalid").(bool),
		StateNew:         d.Get("state_new").(bool),
		StateRelated:     d.Get("state_related").(bool),

		SrcNetworkType:      d.Get("src_network_type").(string),
		SrcMACAddress:       d.Get("src_mac").(string),
		SrcAddress:          d.Get("src_address").(string),
		SrcNetworkID:        d.Get("src_network_id").(string),
		SrcFirewallGroupIDs: srcFirewallGroupIDs,

		DstNetworkType:      d.Get("dst_network_type").(string),
		DstAddress:          d.Get("dst_address").(string),
		DstPort:             d.Get("dst_port").(string),
		DstNetworkID:        d.Get("dst_network_id").(string),
		DstFirewallGroupIDs: dstFirewallGroupIDs,

		Weekdays:        d.Get("weekdays").(string),
		WeekdaysNegate:  d.Get("weekdays_negate").(bool),
		MonthDays:       d.Get("monthdays").(string),
		MonthDaysNegate: d.Get("monthdays_negate").(bool),
		StartDate:       d.Get("startdate").(string),
		StartTime:       d.Get("starttime").(string),
		StopDate:        d.Get("stopdate").(string),
		StopTime:        d.Get("stoptime").(string),
		UTC:             d.Get("utc").(bool),
	}, nil
}

func resourceFirewallRuleSetResourceData(resp *unifi.FirewallRule, d *schema.ResourceData, site string) diag.Diagnostics {
	d.Set("site", site)
	d.Set("name", resp.Name)
	d.Set("action", resp.Action)
	d.Set("ruleset", resp.Ruleset)
	d.Set("rule_index", resp.RuleIndex)
	d.Set("protocol", resp.Protocol)
	d.Set("icmp_typename", resp.ICMPTypename)
	d.Set("logging", resp.Logging)
	d.Set("ip_sec", resp.IPSec)
	d.Set("state_established", resp.StateEstablished)
	d.Set("state_invalid", resp.StateInvalid)
	d.Set("state_new", resp.StateNew)
	d.Set("state_related", resp.StateRelated)

	// TODO: handle IPv6
	d.Set("src_network_type", resp.SrcNetworkType)
	d.Set("src_firewall_group_ids", stringSliceToSet(resp.SrcFirewallGroupIDs))
	d.Set("src_mac", resp.SrcMACAddress)
	d.Set("src_address", resp.SrcAddress)
	d.Set("src_network_id", resp.SrcNetworkID)

	// TODO: handle IPv6
	d.Set("dst_network_type", resp.DstNetworkType)
	d.Set("dst_firewall_group_ids", stringSliceToSet(resp.DstFirewallGroupIDs))
	d.Set("dst_address", resp.DstAddress)
	d.Set("dst_network_id", resp.DstNetworkID)
	d.Set("dst_port", resp.DstPort)

	d.Set("weekdays", resp.Weekdays)
	d.Set("weekdays_negate", resp.WeekdaysNegate)
	d.Set("monthdays", resp.MonthDays)
	d.Set("monthdays_negate", resp.MonthDaysNegate)
	d.Set("startdate", resp.StartDate)
	d.Set("starttime", resp.StartTime)
	d.Set("stopdate", resp.StopDate)
	d.Set("stoptime", resp.StopTime)
	d.Set("utc", resp.UTC)

	return nil
}

func resourceFirewallRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetFirewallRule(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFirewallRuleSetResourceData(resp, d, site)
}

func resourceFirewallRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceFirewallRuleGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateFirewallRule(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFirewallRuleSetResourceData(resp, d, site)
}

func resourceFirewallRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	err := c.c.DeleteFirewallRule(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return diag.FromErr(err)
}
