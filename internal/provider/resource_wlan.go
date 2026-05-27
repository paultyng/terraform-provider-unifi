package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

var (
	wlanValidMinimumDataRate2g = []int{1000, 2000, 5500, 6000, 9000, 11000, 12000, 18000, 24000, 36000, 48000, 54000}
	wlanValidMinimumDataRate5g = []int{6000, 9000, 12000, 18000, 24000, 36000, 48000, 54000}
)

func resourceWLAN() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_wlan` manages a WiFi network / SSID.",

		CreateContext: resourceWLANCreate,
		ReadContext:   resourceWLANRead,
		UpdateContext: resourceWLANUpdate,
		DeleteContext: resourceWLANDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the network.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the wlan with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The SSID of the network.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user_group_id": {
				Description: "ID of the user group to use for this network.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"security": {
				Description:  "The type of WiFi security for this network. Valid values are: `wpapsk`, `wpaeap`, and `open`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"wpapsk", "wpaeap", "open"}, false),
			},
			"wpa3_support": {
				Description: "Enable WPA 3 support (security must be `wpapsk` and PMF must be turned on).",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"wpa3_transition": {
				Description: "Enable WPA 3 and WPA 2 support (security must be `wpapsk` and `wpa3_support` must be true).",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"pmf_mode": {
				Description:  "Enable Protected Management Frames. This cannot be disabled if using WPA 3. Valid values are `required`, `optional` and `disabled`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"required", "optional", "disabled"}, false),
				Default:      "disabled",
			},
			"passphrase": {
				Description: "The passphrase for the network. Not used if security is open or if private_preshared_keys_enabled is true.",
				Type:        schema.TypeString,
				// only required if security != open
				Optional:  true,
				Sensitive: true,
				ConflictsWith: []string{"private_preshared_keys_enabled", "private_preshared_key"},
			},
			"private_preshared_keys_enabled": {
				Description: "Enable Private Pre-Shared Keys (PPSK) for this WLAN. If true, `passphrase` should not be set.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"private_preshared_key": {
				Description: "A list of private pre-shared keys. Required if `private_preshared_keys_enabled` is true.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Description: "The pre-shared key passphrase.",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},
						"network_id": {
							Description: "The ID of the network (VLAN) to assign to clients using this PSK.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				RequiredWith: []string{"private_preshared_keys_enabled"},
			},
			"hide_ssid": {
				Description: "Indicates whether or not to hide the SSID from broadcast.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"is_guest": {
				Description: "Indicates that this is a guest WLAN and should use guest behaviors.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"multicast_enhance": {
				Description: "Indicates whether or not Multicast Enhance is turned of for the network.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"mac_filter_enabled": {
				Description: "Indicates whether or not the MAC filter is turned of for the network.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"mac_filter_list": {
				Description: "List of MAC addresses to filter (only valid if `mac_filter_enabled` is `true`).",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateFunc:     validation.StringMatch(macAddressRegexp, "Mac address is invalid"),
					DiffSuppressFunc: macDiffSuppressFunc,
				},
			},
			"mac_filter_policy": {
				Description:  "MAC address filter policy (only valid if `mac_filter_enabled` is `true`).",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "deny",
				ValidateFunc: validation.StringInSlice([]string{"allow", "deny"}, false),
			},
			"radius_profile_id": {
				Description: "ID of the RADIUS profile to use when security `wpaeap`. You can query this via the " +
					"`unifi_radius_profile` data source.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule": {
				Description: "Start and stop schedules for the WLAN",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_of_week": {
							Description:  "Day of week for the block. Valid values are `sun`, `mon`, `tue`, `wed`, `thu`, `fri`, `sat`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"sun", "mon", "tue", "wed", "thu", "fri", "sat", "sun"}, false),
						},
						"start_hour": {
							Description:  "Start hour for the block (0-23).",
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},
						"start_minute": {
							Description:  "Start minute for the block (0-59).",
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 59),
						},
						"duration": {
							Description:  "Length of the block in minutes.",
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"name": {
							Description: "Name of the block.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"no2ghz_oui": {
				Description: "Connect high performance clients to 5 GHz only.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"l2_isolation": {
				Description: "Isolates stations on layer 2 (ethernet) level.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"proxy_arp": {
				Description: "Reduces airtime usage by allowing APs to \"proxy\" common broadcast frames as unicast.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"bss_transition": {
				Description: "Improves client transitions between APs when they have a weak signal.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"uapsd": {
				Description: "Enable Unscheduled Automatic Power Save Delivery.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"fast_roaming_enabled": {
				Description: "Enables 802.11r fast roaming.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"minimum_data_rate_2g_kbps": {
				Description: "Set minimum data rate control for 2G devices, in Kbps. " +
					"Use `0` to disable minimum data rates. " +
					"Valid values are: " + markdownValueListInt(wlanValidMinimumDataRate2g) + ".",
				Type:     schema.TypeInt,
				Optional: true,
				// TODO: this validation is from the UI, if other values work, perhaps remove this is set it to a range instead?
				ValidateFunc: validation.IntInSlice(append([]int{0}, wlanValidMinimumDataRate2g...)),
			},
			"minimum_data_rate_5g_kbps": {
				Description: "Set minimum data rate control for 5G devices, in Kbps. " +
					"Use `0` to disable minimum data rates. " +
					"Valid values are: " + markdownValueListInt(wlanValidMinimumDataRate5g) + ".",
				Type:     schema.TypeInt,
				Optional: true,
				// TODO: this validation is from the UI, if other values work, perhaps remove this is set it to a range instead?
				ValidateFunc: validation.IntInSlice(append([]int{0}, wlanValidMinimumDataRate5g...)),
			},
			"wlan_band": {
				Description:  "Radio band your WiFi network will use.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"2g", "5g", "both"}, false),
				Default:      "both",
			},
			"network_id": {
				Description: "ID of the network for this SSID. Not used and must not be set if `private_preshared_keys_enabled` is true.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{"private_preshared_keys_enabled"},
			},
			"ap_group_ids": {
				Description: "IDs of the AP groups to use for this network.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceWLANGetResourceData(d *schema.ResourceData, meta interface{}) (*unifi.WLAN, error) {
	c := meta.(*client)

	security := d.Get("security").(string)
	mainPassphrase := d.Get("passphrase").(string)
	switch security {
	case "open":
		mainPassphrase = ""
	}

	pmf := d.Get("pmf_mode").(string)
	wpa3 := d.Get("wpa3_support").(bool)
	wpa3Transition := d.Get("wpa3_transition").(bool)
	switch security {
	case "wpapsk":
		// nothing
	default:
		if wpa3 || wpa3Transition {
			return nil, fmt.Errorf("wpa3_support and wpa3_transition are only valid for security type wpapsk")
		}
	}
	if v := c.ControllerVersion(); v.LessThanOrEqual(controllerVersionWPA3) {
		if wpa3 || wpa3Transition {
			return nil, fmt.Errorf("WPA 3 support is not available on controller version %q, you must be on %q or higher", v, controllerVersionWPA3)
		}
	}

	if wpa3Transition && pmf == "disabled" {
		return nil, fmt.Errorf("WPA 3 transition mode requires pmf_mode to be turned on.")
	} else if wpa3 && !wpa3Transition && pmf != "required" {
		return nil, fmt.Errorf("For WPA 3 you must set pmf_mode to required.")
	}

	macFilterEnabled := d.Get("mac_filter_enabled").(bool)
	macFilterList, err := setToStringSlice(d.Get("mac_filter_list").(*schema.Set))
	if err != nil {
		return nil, err
	}
	if !macFilterEnabled {
		macFilterList = nil
	}

	ppskEnabled := d.Get("private_preshared_keys_enabled").(bool)
	var ppskEntries []unifi.WLANPrivatePresharedKeys

	if ppskEnabled {
		// Schema `ConflictsWith` should ensure `d.Get("passphrase").(string)` is empty here.
		if v, ok := d.GetOk("private_preshared_key"); ok {
			tfPPSKList := v.([]interface{})
			if len(tfPPSKList) == 0 && ppskEnabled {
				return nil, fmt.Errorf("`private_preshared_key` block cannot be empty when `private_preshared_keys_enabled` is true")
			}
			ppskEntries = make([]unifi.WLANPrivatePresharedKeys, len(tfPPSKList))
			for i, item := range tfPPSKList {
				entryMap := item.(map[string]interface{})
				ppskEntries[i] = unifi.WLANPrivatePresharedKeys{
					Password:  entryMap["password"].(string),
					NetworkID: entryMap["network_id"].(string),
				}
			}
		} else {
			return nil, fmt.Errorf("`private_preshared_key` attribute is required when `private_preshared_keys_enabled` is true")
		}
		mainPassphrase = "" // Main passphrase is not used
	} else if security == "wpapsk" && mainPassphrase == "" {
		return nil, fmt.Errorf("`passphrase` is required when security is 'wpapsk' and `private_preshared_keys_enabled` is false")
	}

	// version specific fields and validation
	networkID := d.Get("network_id").(string)
	apGroupIDs, err := setToStringSlice(d.Get("ap_group_ids").(*schema.Set))
	if err != nil {
		return nil, err
	}
	wlanBand := d.Get("wlan_band").(string)

	schedule, err := listToSchedules(d.Get("schedule").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to process schedule block: %w", err)
	}

	minrateSettingPreference := "auto"
	if d.Get("minimum_data_rate_2g_kbps").(int) != 0 || d.Get("minimum_data_rate_5g_kbps").(int) != 0 {
		if d.Get("minimum_data_rate_2g_kbps").(int) == 0 || d.Get("minimum_data_rate_5g_kbps").(int) == 0 {
			// this is really only true I think in >= 7.2, but easier to just apply this in general
			return nil, fmt.Errorf("you must set minimum data rates on both 2g and 5g if setting either")
		}
		minrateSettingPreference = "manual"
	}

	return &unifi.WLAN{
		Name:                    d.Get("name").(string),
		XPassphrase:             mainPassphrase,
		HideSSID:                d.Get("hide_ssid").(bool),
		IsGuest:                 d.Get("is_guest").(bool),
		NetworkID:               networkID,
		ApGroupIDs:              apGroupIDs,
		UserGroupID:             d.Get("user_group_id").(string),
		Security:                security,
		WPA3Support:             wpa3,
		WPA3Transition:          wpa3Transition,
		MulticastEnhanceEnabled: d.Get("multicast_enhance").(bool),
		MACFilterEnabled:        macFilterEnabled,
		MACFilterList:           macFilterList,
		MACFilterPolicy:         d.Get("mac_filter_policy").(string),
		RADIUSProfileID:         d.Get("radius_profile_id").(string),
		ScheduleWithDuration:    schedule,
		ScheduleEnabled:         len(schedule) > 0,
		WLANBand:                wlanBand,
		PMFMode:                 pmf,

		// TODO: add to schema
		WPAEnc:             "ccmp",
		WPAMode:            "wpa2",
		Enabled:            true,
		NameCombineEnabled: true,

		GroupRekey:         3600,
		DTIMMode:           "default",
		No2GhzOui:          d.Get("no2ghz_oui").(bool),
		L2Isolation:        d.Get("l2_isolation").(bool),
		ProxyArp:           d.Get("proxy_arp").(bool),
		BssTransition:      d.Get("bss_transition").(bool),
		UapsdEnabled:       d.Get("uapsd").(bool),
		FastRoamingEnabled: d.Get("fast_roaming_enabled").(bool),

		MinrateSettingPreference: minrateSettingPreference,

		MinrateNgEnabled:      d.Get("minimum_data_rate_2g_kbps").(int) != 0,
		MinrateNgDataRateKbps: d.Get("minimum_data_rate_2g_kbps").(int),

		MinrateNaEnabled:      d.Get("minimum_data_rate_5g_kbps").(int) != 0,
		MinrateNaDataRateKbps: d.Get("minimum_data_rate_5g_kbps").(int),

		PrivatePresharedKeysEnabled: ppskEnabled,
		PrivatePresharedKeys:        ppskEntries,
	}, nil
}

func resourceWLANCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceWLANGetResourceData(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateWLAN(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourceWLANSetResourceData(resp, d, meta, site)
}

func resourceWLANSetResourceData(resp *unifi.WLAN, d *schema.ResourceData, meta interface{}, site string) diag.Diagnostics {
	// c := meta.(*client)
	security := resp.Security
	wpa3 := false
	wpa3Transition := false
	switch security {
	case "open":
	case "wpapsk":
		wpa3 = resp.WPA3Support
		wpa3Transition = resp.WPA3Transition
	}

	macFilterEnabled := resp.MACFilterEnabled
	var macFilterList *schema.Set
	macFilterPolicy := "deny"
	if macFilterEnabled {
		macFilterList = stringSliceToSet(resp.MACFilterList)
		macFilterPolicy = resp.MACFilterPolicy
	}

	apGroupIDs := stringSliceToSet(resp.ApGroupIDs)

	schedule := listFromSchedules(resp.ScheduleWithDuration)

	d.Set("site", site)
	d.Set("name", resp.Name)
	d.Set("user_group_id", resp.UserGroupID)
	d.Set("hide_ssid", resp.HideSSID)
	d.Set("is_guest", resp.IsGuest)
	d.Set("security", security)
	d.Set("wpa3_support", wpa3)
	d.Set("wpa3_transition", wpa3Transition)
	d.Set("multicast_enhance", resp.MulticastEnhanceEnabled)
	d.Set("mac_filter_enabled", macFilterEnabled)
	d.Set("mac_filter_list", macFilterList)
	d.Set("mac_filter_policy", macFilterPolicy)
	d.Set("radius_profile_id", resp.RADIUSProfileID)
	d.Set("schedule", schedule)
	d.Set("wlan_band", resp.WLANBand)
	d.Set("no2ghz_oui", resp.No2GhzOui)
	d.Set("l2_isolation", resp.L2Isolation)
	d.Set("proxy_arp", resp.ProxyArp)
	d.Set("bss_transition", resp.BssTransition)
	d.Set("uapsd", resp.UapsdEnabled)
	d.Set("fast_roaming_enabled", resp.FastRoamingEnabled)
	d.Set("ap_group_ids", apGroupIDs)
	d.Set("pmf_mode", resp.PMFMode)
	if resp.MinrateSettingPreference != "auto" && resp.MinrateNgEnabled {
		d.Set("minimum_data_rate_2g_kbps", resp.MinrateNgDataRateKbps)
	} else {
		d.Set("minimum_data_rate_2g_kbps", 0)
	}
	if resp.MinrateSettingPreference != "auto" && resp.MinrateNaEnabled {
		d.Set("minimum_data_rate_5g_kbps", resp.MinrateNaDataRateKbps)
	} else {
		d.Set("minimum_data_rate_5g_kbps", 0)
	}
	d.Set("private_preshared_keys_enabled", resp.PrivatePresharedKeysEnabled)
	if resp.PrivatePresharedKeysEnabled {
		// PPSK is ENABLED.
		d.Set("passphrase", "") // Override main passphrase state to empty
		d.Set("network_id", "") // Set main network_id state to empty (to prevent drift)

		tfPPSKList := make([]interface{}, len(resp.PrivatePresharedKeys))
		for i, apiEntry := range resp.PrivatePresharedKeys {
			entryMap := map[string]interface{}{
				"password":   apiEntry.Password,
				"network_id": apiEntry.NetworkID,
			}
			tfPPSKList[i] = entryMap
		}
		d.Set("private_preshared_key", tfPPSKList)
	} else {
		d.Set("private_preshared_key", nil)
		passphraseToSet := resp.XPassphrase
		if resp.Security == "open" {
			passphraseToSet = ""
		}
		d.Set("passphrase", passphraseToSet)

		// Set the main network_id from the API response
		d.Set("network_id", resp.NetworkID)
	}

	return nil
}

func resourceWLANRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetWLAN(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWLANSetResourceData(resp, d, meta, site)
}

func resourceWLANUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceWLANGetResourceData(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateWLAN(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWLANSetResourceData(resp, d, meta, site)
}

func resourceWLANDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	err := c.c.DeleteWLAN(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return diag.FromErr(err)
}

func listToSchedules(list []interface{}) ([]unifi.WLANScheduleWithDuration, error) {
	schedules := make([]unifi.WLANScheduleWithDuration, 0, len(list))
	for _, item := range list {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data in block")
		}
		ss := toSchedule(data)
		schedules = append(schedules, ss)
	}
	return schedules, nil
}

func toSchedule(data map[string]interface{}) unifi.WLANScheduleWithDuration {
	// TODO: error check these?
	dow := data["day_of_week"].(string)
	startHour := data["start_hour"].(int)
	startMinute := data["start_minute"].(int)
	duration := data["duration"].(int)
	name := data["name"].(string)

	return unifi.WLANScheduleWithDuration{
		StartDaysOfWeek: []string{dow},
		StartHour:       startHour,
		StartMinute:     startMinute,
		DurationMinutes: duration,
		Name:            name,
	}
}

func fromSchedule(dow string, s unifi.WLANScheduleWithDuration) map[string]interface{} {
	return map[string]interface{}{
		"day_of_week":  dow,
		"start_hour":   s.StartHour,
		"start_minute": s.StartMinute,
		"duration":     s.DurationMinutes,
		"name":         s.Name,
	}
}

func listFromSchedules(ss []unifi.WLANScheduleWithDuration) []interface{} {
	// this explodes days of week lists in to individual schedules
	list := make([]interface{}, 0, len(ss))
	for _, s := range ss {
		for _, dow := range s.StartDaysOfWeek {
			v := fromSchedule(dow, s)
			list = append(list, v)
		}
	}
	return list
}
