package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/paultyng/terraform-provider-unifi/unifi"
)

func resourceWLAN() *schema.Resource {
	return &schema.Resource{
		Create: resourceWLANCreate,
		Read:   resourceWLANRead,
		Update: resourceWLANUpdate,
		Delete: resourceWLANDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"wlan_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"wpapsk", "wpaeap", "open"}, false),
			},
			"passphrase": {
				Type: schema.TypeString,
				// only required if security != open
				Optional:  true,
				Sensitive: true,
			},
			"hide_ssid": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_guest": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceWLANCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	vlan := d.Get("vlan_id").(int)

	req := &unifi.WLAN{
		Name:        d.Get("name").(string),
		VLAN:        fmt.Sprintf("%d", vlan),
		XPassphrase: d.Get("passphrase").(string),
		HideSSID:    d.Get("hide_ssid").(bool),
		IsGuest:     d.Get("is_guest").(bool),
		WLANGroupID: d.Get("wlan_group_id").(string),
		UserGroupID: d.Get("user_group_id").(string),
		Security:    d.Get("security").(string),

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
	}

	resp, err := c.c.CreateWLAN(c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceWLANSetResourceData(resp, d)
}

func resourceWLANSetResourceData(resp *unifi.WLAN, d *schema.ResourceData) error {
	var err error
	vlan := 0
	if resp.VLANEnabled {
		vlan, err = strconv.Atoi(resp.VLAN)
		if err != nil {
			return err
		}
	}

	d.Set("name", resp.Name)
	d.Set("vlan_id", vlan)
	d.Set("passphrase", resp.XPassphrase)
	d.Set("hide_ssid", resp.HideSSID)
	d.Set("is_guest", resp.IsGuest)
	d.Set("wlan_group_id", resp.WLANGroupID)
	d.Set("user_group_id", resp.UserGroupID)
	d.Set("security", resp.Security)

	return nil
}

func resourceWLANRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	resp, err := c.c.GetWLAN(c.site, id)
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
	panic("not implemented")
}

func resourceWLANDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	err := c.c.DeleteWLAN(c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
