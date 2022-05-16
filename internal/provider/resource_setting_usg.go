package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

var resourceSettingUsgLock = sync.Mutex{}

func resourceSettingUsgLocker(f func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		resourceSettingUsgLock.Lock()
		defer resourceSettingUsgLock.Unlock()
		return f(ctx, d, meta)
	}
}

func resourceSettingUsg() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_setting_usg` manages settings for a Unifi Security Gateway.",

		CreateContext: resourceSettingUsgLocker(resourceSettingUsgUpsert),
		ReadContext:   resourceSettingUsgLocker(resourceSettingUsgRead),
		UpdateContext: resourceSettingUsgLocker(resourceSettingUsgUpsert),
		DeleteContext: schema.NoopContext,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the settings.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the settings with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"multicast_dns_enabled": {
				Description: "Whether multicast DNS is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"firewall_guest_default_log": {
				Description: "Whether the guest firewall log is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"firewall_lan_default_log": {
				Description: "Whether the LAN firewall log is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"firewall_wan_default_log": {
				Description: "Whether the WAN firewall log is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"dhcp_relay_servers": {
				Description: "The DHCP relay servers.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    5,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						validation.IsIPv4Address,
						// this doesn't let blank through
						validation.StringLenBetween(1, 50),
					),
				},
			},
		},
	}
}

func resourceSettingUsgUpdateResourceData(d *schema.ResourceData, meta interface{}, setting *unifi.SettingUsg) error {
	c := meta.(*client)

	if mdns, hasMdns := d.GetOkExists("multicast_dns_enabled"); hasMdns {
		if v := c.ControllerVersion(); v.GreaterThanOrEqual(controllerV7) {
			return fmt.Errorf("multicast_dns_enabled is not supported on controller version %v", c.ControllerVersion())
		}

		setting.MdnsEnabled = mdns.(bool)
	}

	setting.FirewallGuestDefaultLog = d.Get("firewall_guest_default_log").(bool)
	setting.FirewallLanDefaultLog = d.Get("firewall_lan_default_log").(bool)
	setting.FirewallWANDefaultLog = d.Get("firewall_wan_default_log").(bool)

	dhcpRelay, err := listToStringSlice(d.Get("dhcp_relay_servers").([]interface{}))
	if err != nil {
		return fmt.Errorf("unable to convert dhcp_relay_servers to string slice: %w", err)
	}
	setting.DHCPRelayServer1 = append(dhcpRelay, "")[0]
	setting.DHCPRelayServer2 = append(dhcpRelay, "", "")[1]
	setting.DHCPRelayServer3 = append(dhcpRelay, "", "", "")[2]
	setting.DHCPRelayServer4 = append(dhcpRelay, "", "", "", "")[3]
	setting.DHCPRelayServer5 = append(dhcpRelay, "", "", "", "", "")[4]

	return nil
}

func resourceSettingUsgUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	req, err := c.c.GetSettingUsg(ctx, c.site)
	if err != nil {
		return diag.FromErr(err)
	}

	err = resourceSettingUsgUpdateResourceData(d, meta, req)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.c.UpdateSettingUsg(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	return resourceSettingUsgSetResourceData(resp, d, meta, site)
}

func resourceSettingUsgSetResourceData(resp *unifi.SettingUsg, d *schema.ResourceData, meta interface{}, site string) diag.Diagnostics {
	d.Set("site", site)
	d.Set("multicast_dns_enabled", resp.MdnsEnabled)
	d.Set("firewall_guest_default_log", resp.FirewallGuestDefaultLog)
	d.Set("firewall_lan_default_log", resp.FirewallLanDefaultLog)
	d.Set("firewall_wan_default_log", resp.FirewallWANDefaultLog)

	dhcpRelay := []string{}
	for _, s := range []string{
		resp.DHCPRelayServer1,
		resp.DHCPRelayServer2,
		resp.DHCPRelayServer3,
		resp.DHCPRelayServer4,
		resp.DHCPRelayServer5,
	} {
		if s == "" {
			continue
		}
		dhcpRelay = append(dhcpRelay, s)
	}
	d.Set("dhcp_relay_servers", dhcpRelay)

	return nil
}

func resourceSettingUsgRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetSettingUsg(ctx, site)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSettingUsgSetResourceData(resp, d, meta, site)
}
