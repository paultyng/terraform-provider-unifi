package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sayedh/go-unifi/unifi"
	// "github.com/paultyng/go-unifi/unifi" (old go-unifi version)
)

func resourcePortProfile() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_port_profile` manages a port profile for use on network switches.",

		CreateContext: resourcePortProfileCreate,
		ReadContext:   resourcePortProfileRead,
		UpdateContext: resourcePortProfileUpdate,
		DeleteContext: resourcePortProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
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
				Description: "Enable link auto negotiation for the port profile. When set to `true` this overrides `speed`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"dot1x_ctrl": {
				Description:  "The type of 802.1X control to use. Can be `auto`, `force_authorized`, `force_unauthorized`, `mac_based` or `multi_host`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "force_authorized",
				ValidateFunc: validation.StringInSlice([]string{"auto", "force_authorized", "force_unauthorized", "mac_based", "multi_host"}, false),
			},
			"dot1x_idle_timeout": {
				Description:  "The timeout, in seconds, to use when using the MAC Based 802.1X control. Can be between 0 and 65535",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"egress_rate_limit_kbps": {
				Description:  "The egress rate limit, in kpbs, for the port profile. Can be between `64` and `9999999`.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(64, 9999999),
			},
			"egress_rate_limit_kbps_enabled": {
				Description: "Enable egress rate limiting for the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"excluded_network_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true, // Add Computed: true
				Description: "The IDs of networks to exclude from this port profile.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// Adding this new attribute
			"included_network_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true, // Add Computed: true
				Description:   "The IDs of networks to include in this port profile.",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"excluded_network_ids"},
			},
			"forward": {
				Description:  "The type forwarding to use for the port profile. Can be `all`, `native`, `customize` or `disabled`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "customize",
				ValidateFunc: validation.StringInSlice([]string{"all", "native", "customize", "disabled"}, false),
			},
			"full_duplex": {
				Description: "Enable full duplex for the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"isolation": {
				Description: "Enable port isolation for the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"lldpmed_enabled": {
				Description: "Enable LLDP-MED for the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"lldpmed_notify_enabled": {
				Description: "Enable LLDP-MED topology change notifications for the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				//ValidateFunc: ,
			},
			// TODO: rename to native_network_id
			"native_networkconf_id": {
				Description: "The ID of network to use as the main network on the port profile.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "The name of the port profile.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"op_mode": {
				Description:  "The operation mode for the port profile. Can only be `switch`",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "switch",
				ValidateFunc: validation.StringInSlice([]string{"switch"}, false),
			},
			"poe_mode": {
				Description:  "The POE mode for the port profile. Can be one of `auto`, `passv24`, `passthrough` or `off`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"auto", "passv24", "passthrough", "off"}, false),
			},
			"port_security_enabled": {
				Description: "Enable port security for the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"port_security_mac_address": {
				Description: "The MAC addresses associated with the port security for the port profile.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"priority_queue1_level": {
				Description:  "The priority queue 1 level for the port profile. Can be between 0 and 100.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"priority_queue2_level": {
				Description:  "The priority queue 2 level for the port profile. Can be between 0 and 100.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"priority_queue3_level": {
				Description:  "The priority queue 3 level for the port profile. Can be between 0 and 100.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"priority_queue4_level": {
				Description:  "The priority queue 4 level for the port profile. Can be between 0 and 100.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"speed": {
				Description:  "The link speed to set for the port profile. Can be one of `10`, `100`, `1000`, `2500`, `5000`, `10000`, `20000`, `25000`, `40000`, `50000` or `100000`",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{10, 100, 1000, 2500, 5000, 10000, 20000, 25000, 40000, 50000, 100000}),
			},
			"stormctrl_bcast_enabled": {
				Description: "Enable broadcast Storm Control for the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"stormctrl_bcast_level": {
				Description:   "The broadcast Storm Control level for the port profile. Can be between 0 and 100.",
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(0, 100),
				ConflictsWith: []string{"stormctrl_bcast_rate"},
			},
			"stormctrl_bcast_rate": {
				Description:  "The broadcast Storm Control rate for the port profile. Can be between 0 and 14880000.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 14880000),
			},
			"stormctrl_mcast_enabled": {
				Description: "Enable multicast Storm Control for the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"stormctrl_mcast_level": {
				Description:   "The multicast Storm Control level for the port profile. Can be between 0 and 100.",
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(0, 100),
				ConflictsWith: []string{"stormctrl_mcast_rate"},
			},
			"stormctrl_mcast_rate": {
				Description:  "The multicast Storm Control rate for the port profile. Can be between 0 and 14880000.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 14880000),
			},
			"stormctrl_type": {
				Description:  "The type of Storm Control to use for the port profile. Can be one of `level` or `rate`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"level", "rate"}, false),
			},
			"stormctrl_ucast_enabled": {
				Description: "Enable unknown unicast Storm Control for the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"stormctrl_ucast_level": {
				Description:   "The unknown unicast Storm Control level for the port profile. Can be between 0 and 100.",
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(0, 100),
				ConflictsWith: []string{"stormctrl_ucast_rate"},
			},
			"stormctrl_ucast_rate": {
				Description:  "The unknown unicast Storm Control rate for the port profile. Can be between 0 and 14880000.",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 14880000),
			},
			"stp_port_mode": {
				Description: "Enable spanning tree protocol on the port profile.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},

			// Deprecated from paultyng/go-unifi - UnifiVersion = "7.4.162"
			// "tagged_networkconf_ids": {
			// 	Description: "The IDs of networks to tag traffic with for the port profile.",
			// 	Type:        schema.TypeSet,
			// 	Optional:    true,
			// 	Elem:        &schema.Schema{Type: schema.TypeString},
			// },

			// New update from sayedh/go-unifi - UnifiVersion = "8.3.32"
			"tagged_vlan_mgmt": {
				Description: "The VLAN management type for the port profile. Can be one of 'auto', 'block_all', or 'custom'.",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// TODO: rename to voice_network_id
			"voice_networkconf_id": {
				Description: "The ID of network to use as the voice network on the port profile.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourcePortProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	// Compute excludedNetworkIDs based on includedNetworkIDs
	excludedNetworkIDs, diags := computeExcludedNetworkIDs(ctx, c, d)
	if diags.HasError() {
		return diags
	}

	// Prepare the PortProfile request
	req, err := resourcePortProfileGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the computed excludedNetworkIDs
	req.ExcludedNetworkIDs = excludedNetworkIDs

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	resp, err := c.c.CreatePortProfile(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	// Fetch networks to compute included_network_ids
	networks, err := c.c.ListNetwork(ctx, site)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePortProfileSetResourceData(resp, d, site, networks)
}

func resourcePortProfileGetResourceData(d *schema.ResourceData) (*unifi.PortProfile, error) {
	portSecurityMacAddress, err := setToStringSlice(d.Get("port_security_mac_address").(*schema.Set))
	if err != nil {
		return nil, err
	}

	// Deprecated from paultyng/go-unifi - UnifiVersion = "7.4.162"
	// taggedNetworkconfIds, err := setToStringSlice(d.Get("tagged_networkconf_ids").(*schema.Set))
	// if err != nil {
	// 	return nil, err
	// }

	excludedNetworkIDs, err := setToStringSlice(d.Get("excluded_network_ids").(*schema.Set))
	if err != nil {
		return nil, err
	}

	return &unifi.PortProfile{
		Autoneg:                      d.Get("autoneg").(bool),
		Dot1XCtrl:                    d.Get("dot1x_ctrl").(string),
		Dot1XIDleTimeout:             d.Get("dot1x_idle_timeout").(int),
		EgressRateLimitKbps:          d.Get("egress_rate_limit_kbps").(int),
		EgressRateLimitKbpsEnabled:   d.Get("egress_rate_limit_kbps_enabled").(bool),
		ExcludedNetworkIDs:           excludedNetworkIDs,
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

		// Deprecated from paultyng/go-unifi - UnifiVersion = "7.4.162"
		// TaggedNetworkIDs:             taggedNetworkconfIds,

		// New update from sayedh/go-unifi - UnifiVersion = "8.3.32"
		TaggedVLANMgmt: d.Get("tagged_vlan_mgmt").(string),

		VoiceNetworkID: d.Get("voice_networkconf_id").(string),
	}, nil
}

func resourcePortProfileSetResourceData(resp *unifi.PortProfile, d *schema.ResourceData, site string, networks []unifi.Network) diag.Diagnostics {
	d.Set("site", site)
	d.Set("autoneg", resp.Autoneg)
	d.Set("dot1x_ctrl", resp.Dot1XCtrl)
	d.Set("dot1x_idle_timeout", resp.Dot1XIDleTimeout)
	d.Set("egress_rate_limit_kbps", resp.EgressRateLimitKbps)
	d.Set("egress_rate_limit_kbps_enabled", resp.EgressRateLimitKbpsEnabled)
	d.Set("excluded_network_ids", stringSliceToSet(resp.ExcludedNetworkIDs))
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

	// Deprecated from paultyng/go-unifi - UnifiVersion = "7.4.162"
	// d.Set("tagged_networkconf_ids", stringSliceToSet(resp.TaggedNetworkIDs))

	// New update from sayedh/go-unifi - UnifiVersion = "8.3.32"
	d.Set("tagged_vlan_mgmt", resp.TaggedVLANMgmt)

	d.Set("voice_networkconf_id", resp.VoiceNetworkID)

	// Set included_network_ids if it was used in configuration
	if _, ok := d.GetOk("included_network_ids"); ok {
		// Create a set of all network IDs
		allNetworkIDs := make(map[string]struct{})
		for _, network := range networks {
			allNetworkIDs[network.ID] = struct{}{}
		}

		// Create a set of excluded network IDs from the response
		excludedSet := make(map[string]struct{})
		for _, id := range resp.ExcludedNetworkIDs {
			excludedSet[id] = struct{}{}
		}

		// Compute included network IDs
		var includedNetworkIDs []string
		for id := range allNetworkIDs {
			if _, found := excludedSet[id]; !found {
				includedNetworkIDs = append(includedNetworkIDs, id)
			}
		}

		d.Set("included_network_ids", stringSliceToSet(includedNetworkIDs))
	} else {
		// Ensure included_network_ids is not set in the state
		d.Set("included_network_ids", nil)
	}

	return nil
}

// Modify resourcePortProfileRead to fetch networks and pass them to SetResourceData
func resourcePortProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	resp, err := c.c.GetPortProfile(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	// Fetch networks to compute included_network_ids
	networks, err := c.c.ListNetwork(ctx, site)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePortProfileSetResourceData(resp, d, site, networks)
}

func resourcePortProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	// Compute excludedNetworkIDs based on includedNetworkIDs
	excludedNetworkIDs, diags := computeExcludedNetworkIDs(ctx, c, d)
	if diags.HasError() {
		return diags
	}

	// Prepare the PortProfile request
	req, err := resourcePortProfileGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()
	req.ExcludedNetworkIDs = excludedNetworkIDs

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdatePortProfile(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fetch networks to compute included_network_ids
	networks, err := c.c.ListNetwork(ctx, site)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePortProfileSetResourceData(resp, d, site, networks)
}

func resourcePortProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	err := c.c.DeletePortProfile(ctx, site, id)
	return diag.FromErr(err)
}

func computeExcludedNetworkIDs(ctx context.Context, c *client, d *schema.ResourceData) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var excludedNetworkIDs []string

	if v, ok := d.GetOk("included_network_ids"); ok {
		includedNetworkIDs, err := setToStringSlice(v.(*schema.Set))
		if err != nil {
			return nil, diag.FromErr(err)
		}

		site := d.Get("site").(string)
		if site == "" {
			site = c.site
		}

		// Fetch all networks from the site
		networks, err := c.c.ListNetwork(ctx, site)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		// Get the native network ID
		nativeNetworkID := d.Get("native_networkconf_id").(string)

		// Create a set of included network IDs for quick lookup
		includedSet := make(map[string]struct{})
		for _, id := range includedNetworkIDs {
			includedSet[id] = struct{}{}
		}

		// Compute the excluded network IDs
		for _, network := range networks {
			// Skip WAN networks
			if network.Purpose == "wan" || network.Purpose == "WAN" {
				continue
			}
			// Skip included networks
			if _, found := includedSet[network.ID]; found {
				continue
			}
			// Skip the native VLAN
			if network.ID == nativeNetworkID {
				continue
			}
			// Add to excluded networks
			excludedNetworkIDs = append(excludedNetworkIDs, network.ID)
		}

		// Add debug statements
		tflog.Debug(ctx, "Included Network IDs", map[string]interface{}{
			"included_network_ids": includedNetworkIDs,
		})
		tflog.Debug(ctx, "Excluded Network IDs", map[string]interface{}{
			"excluded_network_ids": excludedNetworkIDs,
		})

	} else if v, ok := d.GetOk("excluded_network_ids"); ok {
		var err error
		excludedNetworkIDs, err = setToStringSlice(v.(*schema.Set))
		if err != nil {
			return nil, diag.FromErr(err)
		}

		// Add debug statements
		tflog.Debug(ctx, "Excluded Network IDs (from config)", map[string]interface{}{
			"excluded_network_ids": excludedNetworkIDs,
		})
	}

	return excludedNetworkIDs, diags
}
