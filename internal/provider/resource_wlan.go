package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceWLAN() *schema.Resource {
	return &schema.Resource{
		Description: `
unifi_wlan manages a WiFi network / SSID.
`,
		Create: resourceWLANCreate,
		Read:   resourceWLANRead,
		Update: resourceWLANUpdate,
		Delete: resourceWLANDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The SSID of the network.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vlan_id": {
				Description: "VLAN ID for the network.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
			},
			"wlan_group_id": {
				Description: "ID of the WLAN group to use for this network.",
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
			"passphrase": {
				Description: "The passphrase for the network, this is only required if `security` is not set to `open`.",
				Type:        schema.TypeString,
				// only required if security != open
				Optional:  true,
				Sensitive: true,
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
						"block_start": {
							Description:      "Time of day to start the block.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validation.StringMatch(timeOfDayRegexp, "Time of day is invalid"),
							DiffSuppressFunc: timeOfDayDiffSuppress,
						},
						"block_end": {
							Description:      "Time of day to end the block.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validation.StringMatch(timeOfDayRegexp, "Time of day is invalid"),
							DiffSuppressFunc: timeOfDayDiffSuppress,
						},
					},
				},
			},
		},
	}
}

func resourceWLANGetResourceData(d *schema.ResourceData) (*unifi.WLAN, error) {
	vlan := d.Get("vlan_id").(int)

	security := d.Get("security").(string)
	passphrase := d.Get("passphrase").(string)
	switch security {
	case "open":
		passphrase = ""
	}

	macFilterEnabled := d.Get("mac_filter_enabled").(bool)
	macFilterList, err := setToStringSlice(d.Get("mac_filter_list").(*schema.Set))
	if err != nil {
		return nil, err
	}
	if !macFilterEnabled {
		macFilterList = nil
	}

	schedule, err := listToScheduleStrings(d.Get("schedule").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to process schedule block: %w", err)
	}

	log.Printf("[TRACE] TF Schedule: %#v", schedule)

	return &unifi.WLAN{
		Name:                    d.Get("name").(string),
		VLAN:                    vlan,
		XPassphrase:             passphrase,
		HideSSID:                d.Get("hide_ssid").(bool),
		IsGuest:                 d.Get("is_guest").(bool),
		WLANGroupID:             d.Get("wlan_group_id").(string),
		UserGroupID:             d.Get("user_group_id").(string),
		Security:                security,
		MulticastEnhanceEnabled: d.Get("multicast_enhance").(bool),
		MACFilterEnabled:        macFilterEnabled,
		MACFilterList:           macFilterList,
		MACFilterPolicy:         d.Get("mac_filter_policy").(string),
		RADIUSProfileID:         d.Get("radius_profile_id").(string),
		Schedule:                schedule,
		ScheduleEnabled:         len(schedule) > 0,

		VLANEnabled: vlan != 0 && vlan != 1,

		// TODO: add to schema
		WPAEnc:             "ccmp",
		WPAMode:            "wpa2",
		Enabled:            true,
		NameCombineEnabled: true,

		GroupRekey:               3600,
		DTIMMode:                 "default",
		No2GhzOui:                true,
		MinrateNgCckRatesEnabled: true,
	}, nil
}

func resourceWLANCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceWLANGetResourceData(d)
	if err != nil {
		return err
	}

	resp, err := c.c.CreateWLAN(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceWLANSetResourceData(resp, d)
}

func resourceWLANSetResourceData(resp *unifi.WLAN, d *schema.ResourceData) error {
	vlan := 0
	if resp.VLANEnabled {
		vlan = resp.VLAN
	}

	security := resp.Security
	passphrase := resp.XPassphrase
	switch security {
	case "open":
		passphrase = ""
	}

	macFilterEnabled := resp.MACFilterEnabled
	var macFilterList *schema.Set
	macFilterPolicy := "deny"
	if macFilterEnabled {
		macFilterList = stringSliceToSet(resp.MACFilterList)
		macFilterPolicy = resp.MACFilterPolicy
	}

	log.Printf("[TRACE] API Schedule: %#v", resp.Schedule)

	schedule, err := listFromScheduleStrings(resp.Schedule)
	if err != nil {
		return fmt.Errorf("unable to parse schedule: %w", err)
	}

	d.Set("name", resp.Name)
	d.Set("vlan_id", vlan)
	d.Set("passphrase", passphrase)
	d.Set("hide_ssid", resp.HideSSID)
	d.Set("is_guest", resp.IsGuest)
	d.Set("wlan_group_id", resp.WLANGroupID)
	d.Set("user_group_id", resp.UserGroupID)
	d.Set("security", security)
	d.Set("multicast_enhance", resp.MulticastEnhanceEnabled)
	d.Set("mac_filter_enabled", macFilterEnabled)
	d.Set("mac_filter_list", macFilterList)
	d.Set("mac_filter_policy", macFilterPolicy)
	d.Set("radius_profile_id", resp.RADIUSProfileID)
	d.Set("schedule", schedule)

	return nil
}

func resourceWLANRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	resp, err := c.c.GetWLAN(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceWLANSetResourceData(resp, d)
}

func resourceWLANUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceWLANGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = c.site

	resp, err := c.c.UpdateWLAN(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	return resourceWLANSetResourceData(resp, d)
}

func resourceWLANDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	err := c.c.DeleteWLAN(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}

func listToScheduleStrings(list []interface{}) ([]string, error) {
	schedStrings := make([]string, 0, len(list))
	for _, item := range list {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data in block")
		}
		ss, err := toScheduleString(data)
		if err != nil {
			return nil, fmt.Errorf("unable to create schedule string: %w", err)
		}
		schedStrings = append(schedStrings, ss)
	}
	return schedStrings, nil
}

func toScheduleString(data map[string]interface{}) (string, error) {
	// TODO: error check these?
	dow := data["day_of_week"].(string)
	start := timeFromConfig(data["block_start"].(string))
	end := timeFromConfig(data["block_end"].(string))

	return fmt.Sprintf("%s|%s-%s", dow, start, end), nil
}

func fromScheduleString(s string) (map[string]interface{}, error) {
	parts := strings.Split(s, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed schedule string %q", s)
	}
	dow, times := parts[0], parts[1]
	timeParts := strings.Split(times, "-")
	if len(timeParts) != 2 {
		return nil, fmt.Errorf("malformed schedule times %q", s)
	}

	start, end := timeFromUnifi(timeParts[0]), timeFromUnifi(timeParts[1])

	return map[string]interface{}{
		"day_of_week": dow,
		"block_start": start,
		"block_end":   end,
	}, nil
}

func listFromScheduleStrings(ss []string) ([]interface{}, error) {
	list := make([]interface{}, 0, len(ss))
	for _, s := range ss {
		v, err := fromScheduleString(s)
		if err != nil {
			return nil, fmt.Errorf("unable to parse schedule string %q: %w", s, err)
		}
		list = append(list, v)
	}
	return list, nil
}
